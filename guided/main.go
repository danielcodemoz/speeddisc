package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const version = "0.5.0-preview"

type Report struct {
	Version     string
	Flow        string
	Created     string
	Findings    any
	Limitations string
	NextSteps   string
	Partial     bool
	Error       string
}
type Disk struct {
	Files   int
	Bytes   int64
	Denied  int
	Skipped int
	Stopped bool
	Largest []Large
}
type Large struct {
	RelativePath string
	Bytes        int64
}

func tr(lang, pt, en string) string {
	if lang == "en" {
		return en
	}
	return pt
}
func excluded(p string) bool {
	for _, s := range strings.FieldsFunc(p, func(r rune) bool { return r == '/' || r == '\\' }) {
		if strings.EqualFold(s, "CrossDevice") {
			return true
		}
	}
	return false
}

// Only metadata is inspected. Root confines resolution; links/reparse points are skipped.
// Scan a static local directory; concurrent filesystem changes can make counts incomplete.
func scan(ctx context.Context, path string, limit int) (Disk, error) {
	d := Disk{}
	abs, e := filepath.Abs(path)
	if e != nil {
		return d, e
	}
	if excluded(abs) {
		return d, fmt.Errorf("CrossDevice excluded")
	}
	if e = localRoot(abs); e != nil {
		return d, e
	}
	root, e := os.OpenRoot(abs)
	if e != nil {
		return d, e
	}
	defer root.Close()
	var walk func(string)
	walk = func(rel string) {
		if ctx.Err() != nil || d.Files >= limit {
			d.Stopped = true
			return
		}
		f, e := root.Open(rel)
		if e != nil {
			d.Denied++
			return
		}
		defer f.Close()
		for {
			entries, err := f.ReadDir(128)
			for _, x := range entries {
				if ctx.Err() != nil || d.Files >= limit {
					d.Stopped = true
					return
				}
				p := filepath.Join(rel, x.Name())
				if excluded(p) || x.Type()&os.ModeSymlink != 0 {
					d.Skipped++
					continue
				}
				if unsafeEntry(filepath.Join(abs, p)) {
					d.Skipped++
					continue
				}
				info, e := root.Lstat(p)
				if e != nil {
					d.Denied++
					continue
				}
				if info.Mode()&os.ModeSymlink != 0 {
					d.Skipped++
					continue
				}
				if info.IsDir() {
					walk(p)
				} else if info.Mode().IsRegular() {
					d.Files++
					d.Bytes += info.Size()
					d.Largest = append(d.Largest, Large{p, info.Size()})
					sort.Slice(d.Largest, func(i, j int) bool { return d.Largest[i].Bytes > d.Largest[j].Bytes })
					if len(d.Largest) > 10 {
						d.Largest = d.Largest[:10]
					}
				}
			}
			if err == io.EOF {
				return
			}
			if err != nil {
				d.Denied++
				return
			}
		}
	}
	walk(".")
	return d, nil
}

type Adapter struct {
	Name     string
	Up       bool
	Loopback bool
	MTU      int
}

func network() ([]Adapter, error) {
	xs, e := net.Interfaces()
	if e != nil {
		return nil, e
	}
	out := []Adapter{}
	for _, x := range xs {
		out = append(out, Adapter{x.Name, x.Flags&net.FlagUp != 0, x.Flags&net.FlagLoopback != 0, x.MTU})
	}
	return out, nil
}

// Report creation is the only write. Never overwrite an existing report.
func saveReport(path string, r Report) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if excluded(abs) {
		return fmt.Errorf("CrossDevice report destination excluded")
	}
	if err := localRoot(filepath.Dir(abs)); err != nil {
		return err
	}
	path = abs
	b, e := json.MarshalIndent(r, "", "  ")
	if e != nil {
		return e
	}
	parent, e := os.OpenRoot(filepath.Dir(path))
	if e != nil {
		return e
	}
	defer parent.Close()
	f, e := parent.OpenFile(filepath.Base(path), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(append(b, '\n'))
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
func showReport(out io.Writer, lang string, r Report) {
	fmt.Fprintln(out, "\n"+tr(lang, "RESULTADO • ", "RESULT • ")+r.Flow)
	if r.Partial {
		fmt.Fprintln(out, tr(lang, "Observação parcial: não representa uma verificação completa.", "Partial observation: this is not a complete check."))
	}
	if r.Error != "" {
		fmt.Fprintln(out, tr(lang, "Limite encontrado: ", "Issue encountered: ")+r.Error)
	}
	switch x := r.Findings.(type) {
	case Disk:
		fmt.Fprintf(out, tr(lang, "Visitados: %d ficheiros | tamanho lógico: %.2f MiB\nPastas negadas: %d | itens ignorados: %d | interrompido/limite: %t\n", "Visited: %d files | logical size: %.2f MiB\nDenied folders: %d | skipped items: %d | stopped/limit: %t\n"), x.Files, float64(x.Bytes)/1048576, x.Denied, x.Skipped, x.Stopped)
		fmt.Fprintln(out, tr(lang, "Maiores ficheiros visitados (não são recomendações de eliminação):", "Largest visited files (not deletion recommendations):"))
		for i, v := range x.Largest {
			fmt.Fprintf(out, "%d. %.2f MiB  %s\n", i+1, float64(v.Bytes)/1048576, v.RelativePath)
		}
	case []Adapter:
		if len(x) == 0 {
			fmt.Fprintln(out, tr(lang, "Nenhum adaptador disponível no inventário.", "No adapters available in the inventory."))
		}
		for _, v := range x {
			state := tr(lang, "inactivo", "down")
			if v.Up {
				state = tr(lang, "activo", "up")
			}
			fmt.Fprintf(out, "• %s: %s | MTU %d", v.Name, state, v.MTU)
			if v.Loopback {
				fmt.Fprint(out, tr(lang, " | interface interna", " | loopback"))
			}
			fmt.Fprintln(out)
		}
	default:
		fmt.Fprintln(out, tr(lang, "Nomes de aplicações de arranque observados:", "Observed startup application names:"))
		switch names := x.(type) {
		case []any:
			for _, n := range names {
				fmt.Fprintln(out, "•", n)
			}
		case []string:
			for _, n := range names {
				fmt.Fprintln(out, "•", n)
			}
		case string:
			fmt.Fprintln(out, "•", names)
		default:
			fmt.Fprintln(out, tr(lang, "Inventário indisponível.", "Inventory unavailable."))
		}
	}
	fmt.Fprintln(out, "\n"+tr(lang, "Limitações: ", "Limitations: ")+r.Limitations)
	fmt.Fprintln(out, "\n"+tr(lang, "Próximo passo: ", "Next step: ")+r.NextSteps)
}

func run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	ask := func(prompt string) (string, bool) {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			return "", false
		}
		return strings.TrimSpace(scanner.Text()), true
	}
	lang := "pt"
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	for {
		fmt.Fprintln(out, "\nSPEEDDISC | GUIDED / GUIADO | "+version)
		fmt.Fprintln(out, tr(lang, "Diagnóstico sem reparações • conta normal • pré-visualização não validada no Windows", "Read-only diagnostics • standard account • preview not validated on Windows"))
		fmt.Fprintln(out, tr(lang, "1 Disco cheio   2 Sem internet   3 Arranque lento\nL PT/EN   Q Sair", "1 Disk full   2 No internet   3 Slow startup\nL PT/EN   Q Quit"))
		choice, ok := ask("> ")
		if !ok || strings.EqualFold(choice, "q") {
			return
		}
		if strings.EqualFold(choice, "l") {
			if lang == "pt" {
				lang = "en"
			} else {
				lang = "pt"
			}
			continue
		}
		r := Report{Version: version, Created: time.Now().UTC().Format(time.RFC3339)}
		var e error
		switch choice {
		case "1":
			r.Flow = "disk"
			fmt.Fprintln(out, tr(lang, "Escolha uma pasta local. Só metadados; sem apagar. CrossDevice, ligações e pontos de reanálise são ignorados. Caminhos relativos aparecem no relatório. B para voltar.", "Choose a local folder. Metadata only; no deletion. CrossDevice, links and reparse points are skipped. Relative paths appear in the report. B to go back."))
			p, ok := ask(tr(lang, "Pasta: ", "Folder: "))
			if !ok {
				return
			}
			if strings.EqualFold(p, "b") {
				continue
			}
			scanCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			var d Disk
			d, e = scan(scanCtx, p, 100000)
			cancel()
			r.Findings = d
			r.Partial = d.Denied > 0 || d.Skipped > 0 || d.Stopped || e != nil
			r.Limitations = tr(lang, "Tamanho lógico dos ficheiros visitados, não espaço recuperável nem uso total do disco. Máximo 30 s / 100 000 ficheiros; pastas negadas/ignoradas tornam a contagem parcial.", "Logical size of visited files, not reclaimable space or total disk usage. Maximum 30 s / 100,000 files; denied/skipped folders make counts partial.")
			r.NextSteps = tr(lang, "Reveja os maiores ficheiros e confirme cópias de segurança. Use Definições > Sistema > Armazenamento para opções oficiais; nada foi limpo.", "Review the largest files and verify backups. Use Settings > System > Storage for official options; nothing was cleaned.")
		case "2":
			r.Flow = "internet"
			r.Findings, e = network()
			r.Limitations = tr(lang, "Só estado local dos adaptadores. Sem recolher IP/MAC/SSID ou palavras-passe, sem tráfego de teste. Adaptador activo não comprova acesso à internet.", "Local adapter state only. No IP/MAC/SSID/password collection and no test traffic. An active adapter does not prove internet access.")
			r.NextSteps = tr(lang, "Confirme modo de avião e ligação Wi-Fi/cabo. Compare um website habitual noutro dispositivo; não foram alteradas definições.", "Check airplane mode and Wi-Fi/cable connection. Compare a familiar website on another device; no settings were changed.")
		case "3":
			r.Flow = "startup"
			r.Findings, e = startup(ctx)
			r.Limitations = tr(lang, "Inventário parcial dos nomes de arranque disponíveis. Sem comandos, argumentos ou destinos. Não mede tempo de arranque e não inclui todos os serviços/tarefas.", "Partial inventory of available startup names. No commands, arguments or targets. Does not measure boot time or include every service/task.")
			r.Partial = true
			r.NextSteps = tr(lang, "Compare com Gestor de Tarefas > Aplicações de arranque. Confirme a finalidade antes de decidir; o SpeedDisc não desactivou nada.", "Compare with Task Manager > Startup apps. Confirm each item's purpose before deciding; SpeedDisc disabled nothing.")
		default:
			fmt.Fprintln(out, tr(lang, "Escolha 1, 2, 3, L ou Q.", "Choose 1, 2, 3, L or Q."))
			continue
		}
		if e != nil {
			r.Partial = true
			r.Error = e.Error()
			fmt.Fprintln(out, tr(lang, "Não foi possível concluir: ", "Could not complete: ")+e.Error())
		}
		showReport(out, lang, r)
		fmt.Fprintln(out, tr(lang, "Relatório local pode conter nomes de ficheiros/aplicações. Reveja antes de partilhar.", "Local report may contain file/application names. Review before sharing."))
		p, ok := ask(tr(lang, "Guardar JSON (caminho novo), ou Enter para não guardar: ", "Save JSON (new path), or Enter to skip: "))
		if !ok {
			return
		}
		if p != "" {
			if err := saveReport(p, r); err != nil {
				fmt.Fprintln(out, err)
			} else {
				fmt.Fprintln(out, tr(lang, "Relatório guardado.", "Report saved."))
			}
		}
		if ctx.Err() != nil {
			return
		}
	}
}
func main() { run(os.Stdin, os.Stdout) }
