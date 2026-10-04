package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type tr string

type size int64

type msg struct {
	key  string
	args []any
}

func renderArgs(lang Lang, args []any) []any {
	if len(args) == 0 {
		return nil
	}
	out := make([]any, len(args))
	for i, a := range args {
		switch v := a.(type) {
		case tr:
			out[i] = T(lang, string(v))
		case size:
			out[i] = FormatBytes(int64(v), lang)
		default:
			out[i] = a
		}
	}
	return out
}

func (m msg) Text(lang Lang) string {
	format := T(lang, m.key)
	args := renderArgs(lang, m.args)
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

type Report struct {
	mu       sync.Mutex
	Started  time.Time
	Computer string
	Version  string
	lang     Lang
	Found    []msg
	Done     []msg
	Failed   []msg
	Path     string
	HasSnap  bool
}

func newReport() *Report {
	return &Report{
		Started:  time.Now(),
		Computer: computerName(),
		Version:  version,
		lang:     LangPT,
	}
}

func computerName() string {
	if s := os.Getenv("COMPUTERNAME"); s != "" {
		return s
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return "?"
}

func (r *Report) setLang(l Lang) {
	r.mu.Lock()
	r.lang = l
	r.mu.Unlock()
}

func (r *Report) markSnap() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.HasSnap {
		return false
	}
	r.HasSnap = true
	return true
}

func (r *Report) addFound(m msg) { r.add(&r.Found, m) }
func (r *Report) addDone(m msg)  { r.add(&r.Done, m) }
func (r *Report) addFail(m msg)  { r.add(&r.Failed, m) }

func (r *Report) add(dst *[]msg, m msg) {
	r.mu.Lock()
	*dst = append(*dst, m)
	r.mu.Unlock()
}

func (r *Report) FoundCopy() []msg { return r.copyOf(func() []msg { return r.Found }) }
func (r *Report) DoneCopy() []msg  { return r.copyOf(func() []msg { return r.Done }) }

func (r *Report) copyOf(fn func() []msg) []msg {
	r.mu.Lock()
	defer r.mu.Unlock()
	src := fn()
	out := make([]msg, len(src))
	copy(out, src)
	return out
}

func (r *Report) touched() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Found)+len(r.Done)+len(r.Failed) > 0
}

func (r *Report) render(lang Lang) string {
	var b strings.Builder
	line := func(s string) {
		b.WriteString(s)
		b.WriteString("\r\n")
	}
	section := func(title string, lines []msg) {
		line(title)
		if len(lines) == 0 {
			line("- " + T(lang, "rep_none"))
			return
		}
		for _, m := range lines {
			line("- " + m.Text(lang))
		}
	}
	line(T(lang, "rep_title"))
	line(T(lang, "rep_version", r.Version))
	line(T(lang, "rep_date", r.Started.Format("2006-01-02 15:04")))
	line(T(lang, "rep_lang", T(lang, "rep_lang_name")))
	line(T(lang, "rep_computer", r.Computer))
	line("")
	section(T(lang, "rep_found"), r.Found)
	line("")
	section(T(lang, "rep_done"), r.Done)
	line("")
	section(T(lang, "rep_failed"), r.Failed)
	line("")
	line(T(lang, "rep_notes"))
	line("- " + T(lang, "note_recycle"))
	line("- " + T(lang, "note_services"))
	line("- " + T(lang, "note_prefetch"))
	line("- " + T(lang, "note_sysmain"))
	line("")
	return b.String()
}

func (r *Report) Write() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Path == "" {
		r.Path = chooseReportPath(r.Started)
	}
	body := r.render(r.lang)
	data := append([]byte{0xEF, 0xBB, 0xBF}, []byte(body)...)
	if err := os.WriteFile(r.Path, data, 0o644); err != nil {
		return "", err
	}
	return r.Path, nil
}

func chooseReportPath(when time.Time) string {
	name := "SpeedDisc-report-" + when.Format("20060102-150405") + ".txt"
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		if writableDir(dir) {
			return filepath.Join(dir, name)
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return filepath.Join(cwd, name)
}

func writableDir(dir string) bool {
	f, err := os.CreateTemp(dir, ".speeddisc-write-*")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}
