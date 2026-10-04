//go:build windows

package main

import (
	"errors"
	"fmt"
	"strings"
)

func cleans(kind CleanKind, paths []string) []Stats {
	var out []Stats
	for _, p := range paths {
		if strings.TrimSpace(p) == "" {
			out = append(out, Stats{Missing: true})
			continue
		}
		out = append(out, Clean(kind, p))
	}
	return out
}

func foldStats(list []Stats) (removed int, freed int64, skipped int, present int, errs []error) {
	for _, s := range list {
		if s.Err != nil {
			errs = append(errs, s.Err)
			continue
		}
		if s.Missing {
			continue
		}
		present++
		removed += s.Removed
		freed += s.Freed
		skipped += s.Skipped
	}
	return
}

func recordClean(rep *Report, ui *UI, lang Lang, name tr, list []Stats) {
	removed, freed, skipped, present, errs := foldStats(list)
	if present > 0 {
		m := msg{"done_clean", []any{name, removed, size(freed), skipped}}
		rep.addDone(m)
		ui.ok(m.Text(lang))
	} else if len(errs) == 0 {
		m := msg{"done_missing", []any{name}}
		rep.addDone(m)
		ui.warn(m.Text(lang))
	}
	for _, e := range errs {
		m := msg{"fail_clean", []any{name, errDetail(e)}}
		rep.addFail(m)
		ui.err(m.Text(lang))
	}
}

func recordCommand(rep *Report, ui *UI, lang Lang, name tr, code int, tail string, rebootOK bool, line *scanLine, doneLine string) {
	if code == 0 || (rebootOK && code == 3010) {
		key := "done_cmd"
		if code == 3010 {
			key = "done_cmd_reboot"
		}
		m := msg{key, []any{name, code}}
		rep.addDone(m)
		final := doneLine
		if code == 3010 {
			final = m.Text(lang)
		}
		line.paint = func(s string) string { return ui.paint(ansiGreen, s) }
		line.finish(fitConsole(ui, final))
		return
	}
	var detail any
	if strings.TrimSpace(tail) == "" {
		detail = tr("err_exit")
	} else {
		detail = clip(tail, 220)
	}
	m := msg{"fail_cmd", []any{name, code, detail}}
	rep.addFail(m)
	line.paint = func(s string) string { return ui.paint(ansiRed, s) }
	line.finish(fitConsole(ui, m.Text(lang)))
}

func fitConsole(ui *UI, s string) string {
	text := "  " + s
	w := ui.termWidth() - 1
	if w < 1 {
		w = 1
	}
	if runeLen(text) > w {
		return clip(text, w)
	}
	return text
}

func scanPaint(ui *UI) func(string) string {
	return func(s string) string { return ui.paint(ansiBoldCyan, s) }
}

func execute(lang Lang, acts []Action, ui *UI, rep *Report) {
	fresh := rep.markSnap()
	if fresh {
		collectInfo(rep)
	}
	if fresh || Contains(acts, ActInfo) {
		ui.section(lang, "rep_found", rep.FoundCopy())
	}
	for _, a := range acts {
		switch a {
		case ActInfo:
		case ActRestore:
			ui.step(T(lang, "running", T(lang, "name_restore")))
			if err := createRestorePoint(); err != nil {
				m := msg{"fail_restore", []any{errDetail(err)}}
				rep.addFail(m)
				ui.err(m.Text(lang))
			} else {
				m := msg{"done_restore", nil}
				rep.addDone(m)
				ui.ok(m.Text(lang))
			}
		case ActTemp:
			ui.step(T(lang, "running", T(lang, "name_temp")))
			recordClean(rep, ui, lang, tr("name_user_temp"), cleans(CleanTemp, userTempTargets()))
			recordClean(rep, ui, lang, tr("name_win_temp"), cleans(CleanTemp, []string{windowsTemp()}))
		case ActWU:
			ui.step(T(lang, "running", T(lang, "name_wu")))
			recordClean(rep, ui, lang, tr("name_wu"), cleans(CleanWU, []string{wuDownload()}))
		case ActThumb:
			ui.step(T(lang, "running", T(lang, "name_thumb")))
			dir := thumbDir()
			var st Stats
			if dir == "" {
				st = Stats{Missing: true}
			} else {
				st = CleanMatching(dir, thumbPatterns)
			}
			recordClean(rep, ui, lang, tr("name_thumb"), []Stats{st})
		case ActDO:
			ui.step(T(lang, "running", T(lang, "name_do")))
			recordClean(rep, ui, lang, tr("name_do"), cleans(CleanDO, []string{doCache()}))
		case ActWER:
			ui.step(T(lang, "running", T(lang, "name_wer")))
			recordClean(rep, ui, lang, tr("name_wer"), cleans(CleanWER, werTargets()))
		case ActPower:
			ui.step(T(lang, "running", T(lang, "name_power")))
			guid, err := setHighPerformance()
			if err != nil {
				m := msg{"fail_power", []any{errDetail(err)}}
				rep.addFail(m)
				ui.err(m.Text(lang))
			} else {
				m := msg{"done_power", []any{guid}}
				rep.addDone(m)
				ui.ok(m.Text(lang))
			}
		case ActSFC:
			ui.step(T(lang, "running", T(lang, "name_sfc")))
			ui.warn(T(lang, "long_running"))
			code, tail, line := runSFC(ui.out, ui.termWidth()-1, scanPaint(ui))
			recordCommand(rep, ui, lang, tr("name_sfc"), code, tail, false, line, T(lang, "done_sfc_line"))
		case ActDISM:
			ui.step(T(lang, "running", T(lang, "name_dism")))
			ui.warn(T(lang, "long_running"))
			code, tail, line := runDISM(ui.out, ui.termWidth()-1, scanPaint(ui))
			recordCommand(rep, ui, lang, tr("name_dism"), code, tail, true, line, T(lang, "done_dism_line"))
		case ActStartup:
			runStartupUI(lang, ui, rep)
		}
		rep.setLang(lang)
		_, _ = rep.Write()
	}
}

func runStartupUI(lang Lang, ui *UI, rep *Report) {
	ui.step(T(lang, "running", T(lang, "name_startup")))
	entries, warns := listStartup()
	for _, w := range warns {
		m := msg{"fail_info", []any{tr("name_startup"), clip(w, ui.textWidth())}}
		rep.addFail(m)
		ui.warn(m.Text(lang))
	}
	if len(entries) == 0 {
		ui.info(T(lang, "startup_empty"))
		m := msg{"done_startup_none", nil}
		rep.addDone(m)
		return
	}
	fmt.Fprintln(ui.out)
	for i, e := range entries {
		mark := ""
		if !e.CanDisable {
			mark = "  " + T(lang, "not_text_value")
		}
		ui.rawWrap(fmt.Sprintf("  %3d  %s%s", i+1, e.Name, mark))
		ui.rawWrap("       " + e.Where)
		cmd := e.Command
		if cmd == "" {
			cmd = T(lang, "not_text_value")
		}
		ui.rawWrap("       " + cmd)
	}
	fmt.Fprintln(ui.out)
	ui.info(T(lang, "startup_pick"))
	var idxs []int
	var picked bool
	for attempt := 0; attempt < 3; attempt++ {
		line := ui.Ask("> ")
		if line == "" || isQuit(line) {
			m := msg{"done_startup_none", nil}
			rep.addDone(m)
			ui.info(m.Text(lang))
			return
		}
		got, err := parseIndexes(line, len(entries))
		if err != nil {
			ui.warn(T(lang, "startup_bad"))
			continue
		}
		idxs = got
		picked = true
		break
	}
	if !picked {
		m := msg{"done_startup_none", nil}
		rep.addDone(m)
		return
	}
	fmt.Fprintln(ui.out)
	for _, i := range idxs {
		e := entries[i]
		ui.rawWrap("  " + e.Name)
		ui.rawWrap("  " + e.Where)
		ui.rawWrap("  " + e.Command)
		fmt.Fprintln(ui.out)
	}
	if !isYes(ui.Ask("  " + T(lang, "startup_confirm"))) {
		ui.warn(T(lang, "cancelled"))
		m := msg{"done_startup_none", nil}
		rep.addDone(m)
		return
	}
	for _, i := range idxs {
		e := entries[i]
		if !e.CanDisable {
			m := msg{"fail_startup_type", []any{e.Name}}
			rep.addFail(m)
			ui.err(m.Text(lang))
			continue
		}
		if e.File != "" {
			dest, err := disableStartupFile(e.File)
			if err != nil {
				m := msg{"fail_startup", []any{e.Name, errDetail(err)}}
				rep.addFail(m)
				ui.err(m.Text(lang))
				continue
			}
			m := msg{"done_startup_file", []any{e.Name, dest}}
			rep.addDone(m)
			ui.ok(m.Text(lang))
			continue
		}
		err := disableRunValue(e)
		if err != nil {
			var m msg
			switch {
			case errors.Is(err, ErrNotText):
				m = msg{"fail_startup_type", []any{e.Name}}
			case errors.Is(err, ErrValueChanged):
				m = msg{"fail_startup_changed", []any{e.Name}}
			default:
				m = msg{"fail_startup", []any{e.Name, errDetail(err)}}
			}
			rep.addFail(m)
			ui.err(m.Text(lang))
			continue
		}
		m := msg{"done_startup_reg", []any{e.Name, e.Where, clip(e.Command, 300)}}
		rep.addDone(m)
		ui.ok(m.Text(lang))
	}
}
