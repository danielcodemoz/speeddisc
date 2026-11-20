//go:build linux

// Package linuxdesk runs local, read-only Linux diagnostics.
package linuxdesk

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/danielcodemoz/speeddisc/internal/report"
	"golang.org/x/sys/unix"
)

type Task struct{ ID, Category, Title, Detail, Command string }

func en(lang string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(lang)), "en")
}

func text(lang, pt, english string) string {
	if en(lang) {
		return english
	}
	return pt
}

func Tasks(lang string) []Task {
	return []Task{
		{"disk", "disk", text(lang, "Espaço em disco", "Disk space"), text(lang, "Capacidade e espaço disponível na unidade que aloja a raiz do sistema.", "Capacity and available space on the filesystem that hosts the system root."), "df -h /"},
		{"memory", "diagnostics", text(lang, "Memória e swap", "Memory and swap"), text(lang, "Memória disponível, uso estimado e utilização de swap segundo o kernel.", "Available memory, estimated use and swap use according to the kernel."), "free -h"},
		{"processes", "diagnostics", text(lang, "Processos com mais memória", "Top memory processes"), text(lang, "Doze processos visíveis por memória residente, sem ler argumentos ou dados privados dos processos.", "Twelve visible processes ranked by resident memory, without reading their arguments or process data."), "ps -eo pid,comm,rss --sort=-rss"},
		{"network", "diagnostics", text(lang, "Interfaces de rede", "Network interfaces"), text(lang, "Tráfego e erros por interface contados localmente pelo kernel. Não testa a Internet nem comunica com outros sistemas.", "Traffic and errors per interface as counted locally by the kernel. Does not test the Internet or contact other systems."), "ip -s link"},
		{"large-files", "disk", text(lang, "Ficheiros grandes", "Large files"), text(lang, "Ficheiros grandes na pasta pessoal. Não segue ligações simbólicas nem atravessa dispositivos. Não apaga nada.", "Large files in the home folder. Does not follow symbolic links or cross device boundaries. Does not delete anything."), ""},
		{"temp-estimate", "disk", text(lang, "Estimativa de temporários", "Temporary files estimate"), text(lang, "Estima o tamanho dos ficheiros temporários acessíveis do utilizador. Não apaga nada.", "Estimates the size of the current user's accessible temporary files. Does not delete anything."), ""},
	}
}

type environment struct {
	proc, root, home, temp string
	uid                    int
	pageSize               int
}

func Run(ctx context.Context, id, lang string) report.Result {
	home, _ := os.UserHomeDir()
	// Use the conventional system temp root rather than a caller-controlled
	// TMPDIR; the estimate is restricted to the current user's files there.
	return run(ctx, id, lang, environment{"/proc", "/", home, "/tmp", os.Getuid(), os.Getpagesize()})
}

func run(ctx context.Context, id, lang string, fs environment) (r report.Result) {
	start := time.Now()
	r = report.Result{Title: "SpeedDisc Linux", Lang: text(lang, "pt", "en"), Started: start, Status: text(lang, "Concluído", "Completed")}
	defer func() { r.Duration = time.Since(start) }()
	for _, t := range Tasks(lang) {
		if t.ID == id {
			r.Title = t.Title
			break
		}
	}
	add := func(label, value string) { r.Rows = append(r.Rows, report.Row{Label: label, Value: value}) }
	fail := func(msg string) {
		r.Status = text(lang, "Indisponível", "Unavailable")
		add(text(lang, "Estado", "Status"), msg)
	}
	partial := func(msg string) {
		r.Status = text(lang, "Concluído com ressalvas", "Completed with limits")
		add(text(lang, "Limite", "Limit"), msg)
	}
	if ctx.Err() != nil {
		fail(text(lang, "Operação cancelada.", "Operation cancelled."))
		return
	}
	switch id {
	case "disk":
		var stat unix.Statfs_t
		if err := unix.Statfs(fs.root, &stat); err != nil {
			fail(text(lang, "Não foi possível consultar o sistema de ficheiros da raiz.", "Unable to inspect the root filesystem."))
			break
		}
		size := uint64(stat.Bsize)
		if stat.Frsize > 0 {
			size = uint64(stat.Frsize)
		}
		add(text(lang, "Sistema de ficheiros", "Filesystem"), "/")
		add(text(lang, "Capacidade", "Capacity"), byteSize(stat.Blocks*size))
		add(text(lang, "Disponível para o utilizador", "Available to current user"), byteSize(stat.Bavail*size))
		if stat.Blocks >= stat.Bfree {
			add(text(lang, "Ocupado", "Used"), byteSize((stat.Blocks-stat.Bfree)*size))
		}
	case "memory":
		m, err := memInfo(filepath.Join(fs.proc, "meminfo"))
		if err != nil || m["MemTotal"] == 0 {
			fail(text(lang, "Não foi possível ler as métricas de memória do kernel.", "Unable to read kernel memory metrics."))
			break
		}
		add(text(lang, "Memória total", "Total memory"), byteSize(m["MemTotal"]*1024))
		if a, ok := m["MemAvailable"]; ok && a <= m["MemTotal"] {
			add(text(lang, "Memória disponível", "Available memory"), byteSize(a*1024))
			add(text(lang, "Memória ocupada (estimativa)", "Memory in use (estimate)"), byteSize((m["MemTotal"]-a)*1024))
		} else {
			partial(text(lang, "O kernel não forneceu uma leitura válida de memória disponível.", "The kernel did not supply a valid available-memory reading."))
		}
		if s, ok := m["SwapTotal"]; ok {
			add(text(lang, "Swap total", "Total swap"), byteSize(s*1024))
			if f, ok := m["SwapFree"]; ok && f <= s {
				add(text(lang, "Swap em uso", "Swap in use"), byteSize((s-f)*1024))
			}
		}
	case "network":
		rows, limited, err := networkRows(filepath.Join(fs.proc, "net", "dev"), lang)
		if err != nil || len(rows) == 0 {
			fail(text(lang, "Não foi possível ler os contadores das interfaces.", "Unable to read interface counters."))
			break
		}
		r.Rows = rows
		add(text(lang, "Âmbito", "Scope"), text(lang, "Contadores cumulativos do kernel desde o início de cada interface; inclui interfaces locais.", "Cumulative kernel counters since each interface started; includes local interfaces."))
		if limited {
			partial(text(lang, "A listagem foi limitada a 64 interfaces.", "The list was limited to 64 interfaces."))
		}
	case "processes":
		p, limited, err := processes(ctx, fs.proc, fs.pageSize)
		if err != nil {
			fail(text(lang, "Não foi possível ler a lista de processos do kernel.", "Unable to read the kernel process list."))
			break
		}
		for _, x := range p {
			add("PID "+strconv.Itoa(x.pid), fmt.Sprintf("%s · %s RSS", display(x.name), byteSize(x.rss)))
		}
		add(text(lang, "Âmbito", "Scope"), text(lang, "A memória residente pode incluir páginas partilhadas. Processos não visíveis ou terminados entretanto não são listados.", "Resident memory can include shared pages. Hidden or already exited processes are not listed."))
		if limited {
			partial(text(lang, "A lista foi limitada para manter a análise responsiva.", "The list was bounded to keep the analysis responsive."))
		}
	case "large-files", "temp-estimate":
		temp := id == "temp-estimate"
		base := fs.home
		if temp {
			base = fs.temp
		}
		v, err := scan(ctx, base, fs.uid, temp, scanLimits{depth: 8, entries: 10000, time: 2 * time.Second})
		if err != nil {
			fail(text(lang, "A pasta não está disponível, não é local ou contém uma ligação simbólica.", "The folder is unavailable, non-local, or contains a symbolic link."))
			break
		}
		if temp {
			add(text(lang, "Ficheiros temporários visíveis", "Visible temporary files"), strconv.Itoa(v.count))
			add(text(lang, "Tamanho aparente", "Apparent size"), byteSize(v.bytes))
			add(text(lang, "Âmbito", "Scope"), text(lang, "Só ficheiros regulares que pertencem ao utilizador em /tmp, com tamanho lógico. Não é uma estimativa de espaço recuperável.", "Only regular files owned by the current user in /tmp, measured by logical size. This is not a reclaimable-space estimate."))
		} else {
			for i, f := range v.top {
				add(fmt.Sprintf("%02d · %s", i+1, byteSize(f.bytes)), display(f.name))
			}
			add(text(lang, "Âmbito", "Scope"), fmt.Sprintf(text(lang, "%d ficheiros examinados; mostra até 15 com pelo menos 50 MiB. Caminhos relativos à pasta pessoal.", "%d files examined; shows up to 15 of at least 50 MiB. Paths are relative to the home folder."), v.count))
		}
		if v.limited || v.skipped > 0 {
			partial(fmt.Sprintf(text(lang, "Pesquisa parcial: até 2 s, 10 000 entradas, profundidade 8; %d entradas inacessíveis ou excluídas.", "Partial scan: up to 2 s, 10,000 entries, depth 8; %d inaccessible or excluded entries."), v.skipped))
		}
	default:
		fail(text(lang, "Operação desconhecida.", "Unknown operation."))
	}
	return
}

func byteSize(n uint64) string {
	size := float64(n)
	for _, u := range []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"} {
		if size < 1024 || u == "PiB" {
			if u == "B" {
				return fmt.Sprintf("%.0f B", size)
			}
			return fmt.Sprintf("%.1f %s", size, u)
		}
		size /= 1024
	}
	return "0 B"
}

func limitscan(path string, parse func(*bufio.Scanner)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(io.LimitReader(f, 65536))
	parse(s)
	return s.Err()
}

func memInfo(path string) (map[string]uint64, error) {
	rows := map[string]uint64{}
	err := limitscan(path, func(s *bufio.Scanner) {
		for s.Scan() {
			f := strings.Fields(s.Text())
			if len(f) < 3 || f[2] != "kB" {
				continue
			}
			k := strings.TrimSuffix(f[0], ":")
			switch k {
			case "MemTotal", "MemAvailable", "SwapTotal", "SwapFree":
				if v, e := strconv.ParseUint(f[1], 10, 53); e == nil {
					rows[k] = v
				}
			}
		}
	})
	return rows, err
}

func networkRows(path, lang string) (rows []report.Row, limited bool, err error) {
	err = limitscan(path, func(s *bufio.Scanner) {
		for s.Scan() {
			line := s.Text()
			i := strings.LastIndexByte(line, ':')
			if i < 0 {
				continue
			}
			name := strings.TrimSpace(line[:i])
			f := strings.Fields(line[i+1:])
			if name == "" || len(f) < 16 {
				continue
			}
			var stat [4]uint64
			valid := true
			for i, pos := range []int{0, 2, 8, 10} {
				v, e := strconv.ParseUint(f[pos], 10, 64)
				if e != nil {
					valid = false
					break
				}
				stat[i] = v
			}
			if !valid {
				continue
			}
			if len(rows) >= 64 {
				limited = true
				return
			}
			rows = append(rows, report.Row{Label: display(name), Value: fmt.Sprintf(text(lang, "Recebido: %s · Enviado: %s · Erros RX/TX: %d/%d", "Received: %s · Sent: %s · RX/TX errors: %d/%d"), byteSize(stat[0]), byteSize(stat[2]), stat[1], stat[3])})
		}
	})
	return
}

func display(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 32 || r == 127 || r >= 0x202a && r <= 0x202e || r >= 0x2066 && r <= 0x2069 {
			fmt.Fprintf(&b, "<U+%04X>", r)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

type procRow struct {
	pid  int
	name string
	rss  uint64
}

func processes(ctx context.Context, base string, pageSize int) (items []procRow, limited bool, err error) {
	f, err := os.Open(base)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	count := 0
	deadline := time.Now().Add(2 * time.Second)
	for count < 15000 && ctx.Err() == nil && time.Now().Before(deadline) {
		batch, e := f.ReadDir(128)
		for _, d := range batch {
			pid, e := strconv.Atoi(d.Name())
			if e != nil || pid <= 0 || !d.IsDir() {
				continue
			}
			count++
			path := filepath.Join(base, d.Name())
			var resident uint64
			e = limitscan(filepath.Join(path, "statm"), func(s *bufio.Scanner) {
				if s.Scan() {
					f := strings.Fields(s.Text())
					if len(f) > 1 {
						resident, _ = strconv.ParseUint(f[1], 10, 52)
					}
				}
			})
			if e != nil || resident == 0 {
				continue
			}
			var name string
			e = limitscan(filepath.Join(path, "comm"), func(s *bufio.Scanner) {
				if s.Scan() {
					name = strings.TrimSpace(s.Text())
				}
			})
			if e != nil || name == "" {
				continue
			}
			items = append(items, procRow{pid, name, resident * uint64(pageSize)})
		}
		if e == io.EOF {
			break
		}
		if e != nil {
			return items, true, e
		}
	}
	limited = count >= 15000 || ctx.Err() != nil || !time.Now().Before(deadline)
	sort.Slice(items, func(i, j int) bool {
		if items[i].rss == items[j].rss {
			return items[i].pid < items[j].pid
		}
		return items[i].rss > items[j].rss
	})
	if len(items) > 12 {
		items = items[:12]
	}
	return items, limited, nil
}
