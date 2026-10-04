package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// percentRe matches a percentage such as 27%, 27.1% or 27,1%.
// The integer part is what the console shows. Decimals do not get their own line.
var percentRe = regexp.MustCompile(`(\d{1,3})(?:[.,]\d+)?\s*%`)

// scanLine is one console row that updates with a carriage return.
// It does not resize the console buffer, window, or font.
type scanLine struct {
	out   io.Writer
	label string
	max   int
	held  int
	paint func(string) string
}

func newScanLine(out io.Writer, label string, max int, paint func(string) string) *scanLine {
	if max < 1 {
		max = 1
	}
	return &scanLine{out: out, label: label, max: max, paint: paint}
}

func (s *scanLine) show(whole int) {
	s.draw(fmt.Sprintf("  %s %d%%", s.label, whole))
}

func (s *scanLine) draw(text string) {
	if s.max > 0 && runeLen(text) > s.max {
		text = clip(text, s.max)
	}
	n := runeLen(text)
	pad := 0
	if n < s.held {
		pad = s.held - n
	}
	shown := text
	if s.paint != nil {
		shown = s.paint(text)
	}
	fmt.Fprintf(s.out, "\r%s%s", shown, strings.Repeat(" ", pad))
	if n > s.held {
		s.held = n
	}
}

// finish overwrites the progress row and leaves a single final line.
func (s *scanLine) finish(text string) {
	if strings.TrimSpace(text) == "" {
		if s.held > 0 {
			fmt.Fprint(s.out, "\r"+strings.Repeat(" ", s.held)+"\r")
		}
		return
	}
	s.draw(text)
	fmt.Fprintln(s.out)
	s.held = 0
}

// percentFilter keeps SFC/DISM output off the screen except for one
// whole-percent progress row. Raw text is kept only for the failure tail.
type percentFilter struct {
	mu   sync.Mutex
	buf  []byte
	raw  []byte
	line *scanLine
	last int
	have bool
}

func newPercentFilter(line *scanLine) *percentFilter {
	return &percentFilter{line: line}
}

func (f *percentFilter) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.raw = append(f.raw, p...)
	if len(f.raw) > 8000 {
		f.raw = append([]byte(nil), f.raw[len(f.raw)-8000:]...)
	}
	f.buf = append(f.buf, p...)
	f.consume()
	return len(p), nil
}

func (f *percentFilter) flush() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.buf) == 0 {
		return
	}
	f.onChunk(string(f.buf))
	f.buf = nil
}

func (f *percentFilter) rawText() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return string(f.raw)
}

func (f *percentFilter) consume() {
	for {
		i := bytes.IndexAny(f.buf, "\r\n")
		if i < 0 {
			if len(f.buf) > 8192 {
				f.onChunk(string(f.buf))
				f.buf = nil
			}
			return
		}
		chunk := f.buf[:i]
		rest := f.buf[i+1:]
		if f.buf[i] == '\r' && len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}
		f.buf = append([]byte(nil), rest...)
		f.onChunk(string(chunk))
	}
}

func (f *percentFilter) onChunk(chunk string) {
	chunk = strings.TrimSpace(chunk)
	if chunk == "" || f.line == nil {
		return
	}
	ms := percentRe.FindAllStringSubmatch(chunk, -1)
	if len(ms) == 0 {
		return
	}
	n, err := strconv.Atoi(ms[len(ms)-1][1])
	if err != nil || n > 100 {
		return
	}
	if f.have && f.last == n {
		return
	}
	f.have = true
	f.last = n
	f.line.show(n)
}
