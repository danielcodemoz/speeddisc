//go:build windows && speeddiscgui

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/danielcodemoz/speeddisc/internal/report"
)

type guiItem struct {
	ID, Category, Title, Detail string
	NeedsAdmin, Mutates, Guided bool
}

type guiOptions struct {
	Winsock    bool
	StartupIDs []string
	Confirmed  bool
}

type guiStartupItem struct {
	Name, Detail, ID string
	CanDisable       bool
}

func guiText(lang Lang, pt, en string) string {
	if lang == LangEN {
		return en
	}
	return pt
}

func guiItems(lang Lang) []guiItem {
	items := make([]guiItem, 0, 19)
	for _, item := range menuCatalog() {
		id := strconv.Itoa(item.num)
		category := "advanced"
		readOnly := false
		guided := false
		switch item.num {
		case 1, 2, 4:
			category = "packages"
		case 3:
			category, readOnly = "packages", true
		case 5, 6, 7, 8, 10, 16, 18:
			category = "cleanup"
		case 14:
			guided, readOnly = true, true
		case 15, 17:
			category, readOnly = "diagnostics", true
		case 19:
			category, guided = "network", true
		}
		needsAdmin := !readOnly && item.num != 7 && item.num != 18
		detail := T(lang, item.desc)
		if item.num == 14 {
			detail = guiText(lang,
				"Lista aplicações de arranque. A listagem não altera nada. Só desactiva entradas seleccionadas e confirmadas separadamente.",
				"List startup apps. Listing changes nothing. Only explicitly selected and separately confirmed entries can be disabled.")
		}
		items = append(items, guiItem{
			ID: id, Category: category, Title: T(lang, item.title), Detail: detail,
			NeedsAdmin: needsAdmin, Mutates: !readOnly, Guided: guided,
		})
	}
	return items
}

func guiRequiresAdmin(item guiItem, options guiOptions) bool {
	return item.NeedsAdmin || (item.ID == "14" && len(options.StartupIDs) > 0)
}

func guiRun(id string, lang Lang, options guiOptions) report.Result {
	start := time.Now()
	res := report.Result{Started: start, Lang: guiText(lang, "pt", "en")}
	if len(options.StartupIDs) > 0 {
		options.StartupIDs = append([]string(nil), options.StartupIDs...)
	}
	outLabel := guiText(lang, "Observação", "Note")
	finishError := func(pt, en string) report.Result {
		res.Status = guiText(lang, "Não executado", "Not run")
		res.Duration = time.Since(start)
		res.Rows = []report.Row{{Label: outLabel, Value: guiText(lang, pt, en)}}
		return res
	}
	var item guiItem
	for _, it := range guiItems(lang) {
		if it.ID == id {
			item = it
			break
		}
	}
	if item.ID == "" {
		res.Title = "SpeedDisc"
		return finishError("Opção desconhecida.", "Unknown operation.")
	}
	res.Title = item.Title
	if (options.Winsock && id != "19") || (len(options.StartupIDs) > 0 && id != "14") || len(options.StartupIDs) > 512 {
		return finishError("Opções incompatíveis com esta operação.", "Options do not match this operation.")
	}
	if (item.Mutates || len(options.StartupIDs) > 0) && !options.Confirmed {
		return finishError("É necessária confirmação explícita antes de alterar o sistema.", "Explicit confirmation is required before changing the system.")
	}
	if guiRequiresAdmin(item, options) && !elevated() {
		return finishError("Esta operação requer permissões de administrador.", "This operation requires administrator privileges.")
	}
	acts, err := ParseSelection(item.ID)
	if err != nil {
		return finishError("Não foi possível preparar a operação.", "Unable to prepare the operation.")
	}
	rep := newReport()
	rep.setLang(lang)
	ui := NewUI(io.Discard, strings.NewReader(""))
	ui.cols = 100
	hooks := runHooks{
		noAutoWrite: true, selectedInfoOnly: true,
		startup: func(l Lang, u *UI, r *Report) { guiRunStartup(l, u, r, options.StartupIDs) },
		winsock: func() bool { return options.Winsock },
	}
	executeWithHooks(lang, acts, ui, rep, hooks)

	res.Duration = time.Since(start)
	res.Status = guiText(lang, "Concluído", "Completed")
	if len(rep.Failed) != 0 {
		if len(rep.Done) > 0 || len(rep.Found) > 0 {
			res.Status = guiText(lang, "Concluído com ressalvas", "Completed with issues")
		} else {
			res.Status = guiText(lang, "Falhou", "Failed")
		}
	}
	for _, lines := range []struct {
		label string
		items []msg
	}{
		{T(lang, "rep_found"), rep.Found},
		{T(lang, "rep_done"), rep.Done},
		{T(lang, "rep_failed"), rep.Failed},
	} {
		for _, line := range lines.items {
			res.Rows = append(res.Rows, report.Row{Label: lines.label, Value: line.Text(lang)})
		}
	}
	if item.Mutates || len(options.StartupIDs) > 0 {
		for _, line := range rep.Done {
			if line.key == "disk_summary" {
				res.Rows = append(res.Rows, report.Row{Label: outLabel, Value: guiText(lang,
					"A variação do espaço livre resulta de duas medições da unidade do sistema e pode incluir actividade de outros programas.",
					"The change in free space is based on two system-drive readings and may include activity from other programs.")})
				break
			}
		}
	}
	if id == "17" {
		res.Rows = append(res.Rows, report.Row{Label: outLabel, Value: guiText(lang,
			"Mostra até 20 ficheiros de pelo menos 100 MiB em locais acessíveis. A pesquisa é limitada a 30 segundos e pode ficar incompleta.",
			"Shows up to 20 files of at least 100 MiB in accessible locations. Search is limited to 30 seconds and may be incomplete.")})
	}
	return res
}

func guiStartupID(entry StartupEntry) (string, bool) {
	if !entry.CanDisable {
		return "", false
	}
	var fileState string
	if entry.File != "" {
		info, err := os.Lstat(entry.File)
		if err != nil || !info.Mode().IsRegular() || unsafePathHasReparse(entry.File, false) {
			return "", false
		}
		if unsafePathHasReparse(filepath.Join(filepath.Dir(entry.File), "_SpeedDisc_disabled"), true) {
			return "", false
		}
		fileState = fmt.Sprintf("%d/%d/%d", info.Size(), info.ModTime().UnixNano(), info.Mode())
	}
	raw, _ := json.Marshal(struct {
		Entry     StartupEntry
		FileState string
	}{entry, fileState})
	return fmt.Sprintf("%x", sha256.Sum256(raw)), true
}

func guiStartup(lang Lang) ([]guiStartupItem, []string) {
	entries, warnings := listStartup()
	items := make([]guiStartupItem, 0, len(entries))
	for _, entry := range entries {
		id, canDisable := guiStartupID(entry)
		detail := entry.Where
		if entry.Command != "" {
			detail += "\r\n" + entry.Command
		} else {
			detail += "\r\n" + T(lang, "not_text_value")
		}
		items = append(items, guiStartupItem{Name: entry.Name, Detail: detail, ID: id, CanDisable: canDisable})
	}
	return items, warnings
}

func guiRunStartup(lang Lang, ui *UI, rep *Report, ids []string) {
	list, warnings := listStartup()
	for _, w := range warnings {
		rep.addFail(msg{"fail_info", []any{tr("name_startup"), clip(w, 180)}})
	}
	if len(ids) == 0 {
		rep.addFound(msg{"found_startup_count", []any{len(list)}})
		for _, entry := range list {
			rep.addFound(msg{"top_file_entry", []any{entry.Name, entry.Where}})
		}
		rep.addDone(msg{"done_startup_none", nil})
		return
	}
	byID := make(map[string]StartupEntry)
	for _, entry := range list {
		id, ok := guiStartupID(entry)
		if ok {
			byID[id] = entry
		}
	}
	// Validate the full selection before the first write. Preview always
	// presents the exact current entry, never a mutable list index.
	var chosen []StartupEntry
	seen := make(map[string]bool)
	for _, id := range ids {
		entry, ok := byID[id]
		if !ok {
			rep.addFail(msg{"fail_startup_changed", []any{T(lang, "name_startup")}})
			return
		}
		if !seen[id] {
			seen[id] = true
			chosen = append(chosen, entry)
		}
	}
	for _, entry := range chosen {
		disableStartupUI(lang, ui, rep, entry)
	}
}
