// Package report exports local SpeedDisc operation results; it does not run operations.
package report

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

// Row is one name/value pair from a completed operation.
type Row struct{ Label, Value string }

// Result carries the displayed, localised operation data. Lang defaults to PT.
type Result struct {
	Title, Status, Lang string
	Started             time.Time
	Duration            time.Duration
	Rows                []Row
}

type labels struct {
	report, operation, status, started, duration, field, value, page, none, omitted, truncated string
}

var pt = labels{"Relatório de operação", "Operação", "Estado", "Início", "Duração", "Campo", "Resultado", "Página", "Não há detalhes adicionais.", "Linhas omitidas", " [texto truncado]"}
var en = labels{"Operation report", "Operation", "Status", "Started", "Elapsed", "Field", "Result", "Page", "No additional details.", "Omitted rows", " [text truncated]"}

const maxRows = 500

func clean(s string, max int, suffix string) string {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	var b strings.Builder
	n := 0
	for _, r := range s {
		if n == max {
			b.WriteString(suffix)
			break
		}
		n++
		switch {
		case r == '\n':
			b.WriteRune(r)
		case r == '\t':
			b.WriteByte(' ')
		case unicode.IsControl(r) || r == 0x2028 || r == 0x2029 || r == 0xfeff || r >= 0x202a && r <= 0x202e || r >= 0x2066 && r <= 0x2069:
			fmt.Fprintf(&b, "<U+%04X>", r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalise(r Result) (Result, labels) {
	l := pt
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.Lang)), "en") {
		l = en
		r.Lang = "en"
	} else {
		r.Lang = "pt-PT"
	}
	oneLine := func(v string, n int) string {
		return strings.TrimSpace(strings.ReplaceAll(clean(v, n, l.truncated), "\n", " "))
	}
	r.Title = oneLine(r.Title, 112)
	if r.Title == "" {
		r.Title = l.operation
	}
	r.Status = oneLine(r.Status, 128)
	if r.Status == "" {
		r.Status = "—"
	}
	n := min(len(r.Rows), maxRows)
	rows := make([]Row, 0, n+1)
	for _, row := range r.Rows[:n] {
		rows = append(rows, Row{oneLine(row.Label, 124), clean(row.Value, 2048, l.truncated)})
	}
	if len(r.Rows) > n {
		rows = append(rows, Row{l.omitted, fmt.Sprintf("%d", len(r.Rows)-n)})
	}
	if len(rows) == 0 {
		rows = append(rows, Row{"—", l.none})
	}
	r.Rows = rows
	return r, l
}

func metadata(r Result, l labels) []Row {
	at := "—"
	if !r.Started.IsZero() {
		at = r.Started.Format("2006-01-02 15:04:05 -07:00")
		if r.Lang != "en" {
			at = r.Started.Format("02/01/2006 15:04:05 -07:00")
		}
	}
	elapsed := "—"
	if r.Duration >= 0 {
		elapsed = r.Duration.Round(time.Millisecond).String()
		if r.Duration > 0 && r.Duration < time.Millisecond {
			elapsed = "<1 ms"
		}
	}
	return []Row{{l.status, r.Status}, {l.started, at}, {l.duration, elapsed}}
}

// Write renders a UTF-8 TXT file or an A4 PDF. It never opens a file nor executes
// operations. PDF embeds the site font with WinAnsi coverage: unsupported code points are
// displayed losslessly as <U+XXXX>; TXT preserves those original characters.
func Write(w io.Writer, format string, r Result) error {
	data, err := render(format, r)
	if err != nil {
		return err
	}
	n, err := w.Write(data)
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func render(format string, r Result) ([]byte, error) {
	r, l := normalise(r)
	switch strings.ToLower(strings.TrimPrefix(strings.TrimSpace(format), ".")) {
	case "txt":
		var b bytes.Buffer
		fmt.Fprintf(&b, "SpeedDisc | %s\r\n%s\r\n\r\n", l.report, r.Title)
		for _, row := range metadata(r, l) {
			fmt.Fprintf(&b, "%s: %s\r\n", row.Label, row.Value)
		}
		fmt.Fprintf(&b, "\r\n%s / %s\r\n%s\r\n", l.field, l.value, strings.Repeat("─", 44))
		for _, row := range r.Rows {
			fmt.Fprintf(&b, "%s: %s\r\n", row.Label, strings.ReplaceAll(row.Value, "\n", "\r\n  "))
		}
		return b.Bytes(), nil
	case "pdf":
		return pdf(r, l)
	default:
		return nil, fmt.Errorf("report: unsupported format %q", format)
	}
}

// Save creates path exclusively (fs.ErrExist if already present). It does not
// create directories, append extensions or replace an existing user's report.
func Save(path, format string, r Result) error {
	data, err := render(format, r)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	identity, statErr := f.Stat()
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if err = errors.Join(writeErr, closeErr); err != nil && statErr == nil {
		existing, lookupErr := os.Lstat(path)
		if lookupErr == nil && os.SameFile(identity, existing) {
			_ = os.Remove(path)
		}
	}
	return err
}
