//go:build linux

// SpeedDisc Linux is a read-only, single-user toolkit. Its local workbench
// starts without root and is never served on a routable interface.
package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/danielcodemoz/speeddisc/internal/linuxdesk"
	"github.com/danielcodemoz/speeddisc/internal/report"
)

//go:embed web
var web embed.FS

const version = "0.1.0-preview"

type taskView struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Command  string `json:"command"`
}
type rowView struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type resultView struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Status   string    `json:"status"`
	Duration string    `json:"duration"`
	Rows     []rowView `json:"rows"`
}
type desk struct {
	host, prefix string
	assets       fs.FS
	stop         func()
	demo         bool
	runTask      func(context.Context, string, string) report.Result
	last         atomic.Int64
	busy         chan struct{}
	mu           sync.Mutex
	results      map[string]report.Result
	order        []string
}

func token(n int) (string, error) {
	b := make([]byte, n)
	_, e := rand.Read(b)
	return hex.EncodeToString(b), e
}
func newDesk(host, secret string, assets fs.FS, stop func(), demo bool) *desk {
	return &desk{host: host, prefix: "/" + secret + "/", assets: assets, stop: stop, demo: demo, runTask: linuxdesk.Run, busy: make(chan struct{}, 1), results: make(map[string]report.Result)}
}
func jsonTo(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
func problem(w http.ResponseWriter, code int, msg string) {
	jsonTo(w, code, map[string]string{"error": msg})
}

func (d *desk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; font-src 'self'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	peer, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(peer)
	if r.Host != d.host || ip == nil || !ip.IsLoopback() || !strings.HasPrefix(r.URL.Path, d.prefix) {
		http.NotFound(w, r)
		return
	}
	// Reject cross-site fetches even if a local link was accidentally shared.
	if origin := r.Header.Get("Origin"); origin != "" && origin != "http://"+d.host {
		problem(w, 403, "Origin not allowed")
		return
	}
	if s := r.Header.Get("Sec-Fetch-Site"); s == "cross-site" || s == "same-site" {
		problem(w, 403, "Cross-origin request not allowed")
		return
	}
	d.last.Store(time.Now().UnixNano())
	endpoint := strings.TrimPrefix(r.URL.Path, d.prefix)
	if strings.HasPrefix(endpoint, "api/") {
		d.api(w, r, endpoint)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		problem(w, 405, "Method not allowed")
		return
	}
	if endpoint == "" {
		endpoint = "index.html"
	}
	if !fs.ValidPath(endpoint) || strings.Contains(endpoint, "/") {
		http.NotFound(w, r)
		return
	}
	b, err := fs.ReadFile(d.assets, endpoint)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	m := mime.TypeByExtension(path.Ext(endpoint))
	if m == "" {
		m = "application/octet-stream"
	}
	w.Header().Set("Content-Type", m)
	w.Header().Set("Content-Length", fmt.Sprint(len(b)))
	if r.Method != http.MethodHead {
		_, _ = w.Write(b)
	}
}

func (d *desk) api(w http.ResponseWriter, r *http.Request, endpoint string) {
	switch {
	case endpoint == "api/tasks" && r.Method == http.MethodGet:
		from := linuxdesk.Tasks(language(r.URL.Query().Get("lang")))
		out := make([]taskView, 0, len(from))
		for _, t := range from {
			out = append(out, taskView{t.ID, t.Category, t.Title, t.Detail, t.Command})
		}
		jsonTo(w, 200, out)
	case endpoint == "api/run" && r.Method == http.MethodPost:
		ct, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if ct != "application/json" {
			problem(w, 415, "Expected JSON")
			return
		}
		var req struct {
			ID   string `json:"id"`
			Lang string `json:"lang"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 2048)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if dec.Decode(&req) != nil || dec.Decode(new(any)) != io.EOF {
			problem(w, 400, "Invalid request")
			return
		}
		lang := language(req.Lang)
		var selected *linuxdesk.Task
		for _, item := range linuxdesk.Tasks(lang) {
			if item.ID == req.ID {
				selected = &item
				break
			}
		}
		if selected == nil {
			problem(w, 404, "Unknown task")
			return
		}
		select {
		case d.busy <- struct{}{}:
			defer func() { <-d.busy }()
		default:
			problem(w, 409, "A diagnostic is already running")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		var res report.Result
		if d.demo {
			res = demoResult(*selected, lang)
		} else {
			res = d.runTask(ctx, selected.ID, lang)
		}
		id, err := token(16)
		if err != nil {
			problem(w, 500, "Unable to prepare report")
			return
		}
		d.mu.Lock()
		d.results[id] = res
		d.order = append(d.order, id)
		if len(d.order) > 8 {
			delete(d.results, d.order[0])
			d.order = d.order[1:]
		}
		d.mu.Unlock()
		rows := make([]rowView, 0, len(res.Rows))
		for _, row := range res.Rows {
			rows = append(rows, rowView{row.Label, row.Value})
		}
		jsonTo(w, 200, resultView{id, res.Title, res.Status, res.Duration.Round(time.Millisecond).String(), rows})
	case endpoint == "api/report" && r.Method == http.MethodGet:
		id := r.URL.Query().Get("id")
		format := r.URL.Query().Get("format")
		if format != "txt" && format != "pdf" {
			problem(w, 400, "Invalid report format")
			return
		}
		d.mu.Lock()
		res, ok := d.results[id]
		d.mu.Unlock()
		if !ok {
			problem(w, 404, "Report unavailable")
			return
		}
		b, e := renderReport(format, res)
		if e != nil {
			problem(w, 500, "Unable to export the report")
			return
		}
		ct := "text/plain; charset=utf-8"
		if format == "pdf" {
			ct = "application/pdf"
		}
		w.Header().Set("Content-Type", ct)
		w.Header().Set("Content-Disposition", `attachment; filename="SpeedDisc-Linux-`+res.Started.Format("20060102-150405")+"."+format+`"`)
		_, _ = w.Write(b)
	case endpoint == "api/close" && r.Method == http.MethodPost:
		select {
		case d.busy <- struct{}{}:
		default:
			problem(w, 409, "A diagnostic is already running")
			return
		}
		jsonTo(w, 200, map[string]bool{"closed": true})
		go d.stop()
	default:
		http.NotFound(w, r)
	}
}
func language(s string) string {
	if s == "en" {
		return "en"
	}
	return "pt"
}

func demoResult(t linuxdesk.Task, lang string) report.Result {
	res := report.Result{Title: t.Title, Status: "Demonstração", Lang: lang, Started: time.Now(), Duration: 125 * time.Millisecond}
	res.Rows = []report.Row{{Label: "Modo", Value: "Valores fictícios; não foi lida informação deste computador."}, {Label: "Armazenamento de exemplo", Value: "256 GiB, 86 GiB livres"}, {Label: "Memória de exemplo", Value: "16 GiB, 10 GiB disponíveis"}}
	if lang == "en" {
		res.Status = "Demo"
		res.Rows = []report.Row{{Label: "Mode", Value: "Fictional values; no information was read from this computer."}, {Label: "Example storage", Value: "256 GiB, 86 GiB free"}, {Label: "Example memory", Value: "16 GiB, 10 GiB available"}}
	}
	return res
}
func main() {
	noOpen := flag.Bool("no-browser", false, "print the local workbench URL without opening a browser")
	demo := flag.Bool("demo", false, "use fictional diagnostics for design tests")
	printVersion := flag.Bool("version", false, "print version")
	flag.Parse()
	if *printVersion {
		fmt.Println("SpeedDisc Linux", version)
		return
	}
	if os.Geteuid() == 0 && !*demo {
		log.Fatal("Run SpeedDisc as a normal user, without sudo. / Execute como utilizador normal, sem sudo.")
	}
	assets, e := fs.Sub(web, "web")
	if e != nil {
		log.Fatal(e)
	}
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		log.Fatal(e)
	}
	defer listener.Close()
	secret, e := token(24)
	if e != nil {
		log.Fatal(e)
	}
	srv := &http.Server{ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	done := make(chan struct{})
	stopped := make(chan struct{})
	var once sync.Once
	shutdown := func() {
		once.Do(func() {
			close(done)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
			close(stopped)
		})
	}
	d := newDesk(listener.Addr().String(), secret, assets, shutdown, *demo)
	d.last.Store(time.Now().UnixNano())
	srv.Handler = d
	url := "http://" + d.host + d.prefix
	fmt.Println("SpeedDisc Linux · Portátil / Portable ·", version)
	fmt.Println(url)
	if !*noOpen && (os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "") {
		go func() { _ = exec.Command("xdg-open", url).Run() }()
	}
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-t.C:
				if time.Since(time.Unix(0, d.last.Load())) > 30*time.Minute && len(d.busy) == 0 {
					shutdown()
					return
				}
			}
		}
	}()
	if err := srv.Serve(listener); errors.Is(err, http.ErrServerClosed) {
		<-stopped
	} else if err != nil {
		log.Fatal(err)
	}
}
