package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// percentRe matches a percentage such as 27%, 27.1% or 27,1%.
// The integer part is what the console shows. Decimals do not get their own line.
var percentRe = regexp.MustCompile(`(\d{1,3})(?:[.,]\d+)?\s*%`)

// scanLine is one console row that updates with a carriage return.
// It does not resize the console buffer, window, or font.
type scanLine struct {
	mu    sync.Mutex
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
	s.mu.Lock()
	defer s.mu.Unlock()
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
// An empty text only clears the row.
func (s *scanLine) finish(text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(text) == "" {
		if s.held > 0 {
			fmt.Fprint(s.out, "\r"+strings.Repeat(" ", s.held)+"\r")
			s.held = 0
		}
		return
	}
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
	fmt.Fprintf(s.out, "\r%s%s\n", shown, strings.Repeat(" ", pad))
	s.held = 0
}

// percentFilter keeps SFC/DISM output off the screen except for one
// whole-percent progress row. Raw text is kept only for the failure tail.
type percentFilter struct {
	mu        sync.Mutex
	buf       []byte
	raw       []byte
	line      *scanLine
	last      int
	have      bool
	onPercent func(int)
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
			break
		}
		chunk := f.buf[:i]
		rest := f.buf[i+1:]
		if f.buf[i] == '\r' && len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}
		f.buf = append([]byte(nil), rest...)
		f.onChunk(string(chunk))
	}
	// ConPTY often rewrites the row with cursor sequences and no newline yet.
	if bytes.Contains(f.buf, []byte("%")) {
		f.onChunk(string(f.buf))
	}
}

func (f *percentFilter) onChunk(chunk string) {
	chunk = stripANSI(chunk)
	chunk = strings.TrimSpace(chunk)
	if chunk == "" {
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
	if f.onPercent != nil {
		f.onPercent(n)
		return
	}
	if f.line != nil {
		f.line.show(n)
	}
}

// liveLine keeps one row moving while a step runs.
// A real whole percent replaces the spinner text. The spinner stays so a
// stalled percent (DISM often sits on one number) still looks alive.
// It never invents a percent.
type liveLine struct {
	mu        sync.Mutex
	line      *scanLine
	label     string
	working   string
	countFmt  string
	frames    []string
	frame     int
	pct       int
	have      bool
	count     int
	haveCount bool
	lastDraw  time.Time
	stopCh    chan struct{}
	done      chan struct{}
	stopped   bool
}

func startLive(out io.Writer, label, working, countFmt string, max int, paint func(string) string, utf8 bool) *liveLine {
	frames := []string{"|", "/", "-", "\\"}
	if utf8 {
		frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	}
	if countFmt == "" {
		countFmt = "%d"
	}
	l := &liveLine{
		line:     newScanLine(out, label, max, paint),
		label:    label,
		working:  working,
		countFmt: countFmt,
		frames:   frames,
		stopCh:   make(chan struct{}),
		done:     make(chan struct{}),
	}
	l.redraw(true)
	go l.loop()
	return l
}

func (l *liveLine) loop() {
	defer close(l.done)
	t := time.NewTicker(180 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-l.stopCh:
			return
		case <-t.C:
			l.mu.Lock()
			if l.stopped {
				l.mu.Unlock()
				return
			}
			l.frame++
			text := l.renderLocked()
			l.mu.Unlock()
			l.line.draw(text)
		}
	}
}

func (l *liveLine) renderLocked() string {
	mark := l.frames[l.frame%len(l.frames)]
	switch {
	case l.have:
		return fmt.Sprintf("  %s %d%% %s", l.label, l.pct, mark)
	case l.haveCount:
		return fmt.Sprintf("  %s %s %s", l.label, fmt.Sprintf(l.countFmt, l.count), mark)
	default:
		return fmt.Sprintf("  %s %s %s", l.label, l.working, mark)
	}
}

func (l *liveLine) redraw(force bool) {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	if !force && !l.lastDraw.IsZero() && time.Since(l.lastDraw) < 80*time.Millisecond {
		l.mu.Unlock()
		return
	}
	l.lastDraw = time.Now()
	text := l.renderLocked()
	l.mu.Unlock()
	l.line.draw(text)
}

func (l *liveLine) setPercent(n int) {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	if l.have && l.pct == n {
		l.mu.Unlock()
		return
	}
	l.have = true
	l.pct = n
	l.lastDraw = time.Now()
	text := l.renderLocked()
	l.mu.Unlock()
	l.line.draw(text)
}

func (l *liveLine) setCount(n int) {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	l.count = n
	l.haveCount = true
	if l.have {
		l.mu.Unlock()
		return
	}
	if n > 0 && n%12 != 0 && !l.lastDraw.IsZero() && time.Since(l.lastDraw) < 120*time.Millisecond {
		l.mu.Unlock()
		return
	}
	l.lastDraw = time.Now()
	text := l.renderLocked()
	l.mu.Unlock()
	l.line.draw(text)
}

func (l *liveLine) stop() {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return
	}
	l.stopped = true
	close(l.stopCh)
	l.mu.Unlock()
	<-l.done
}

// stopClear ends the ticker and erases the row so the final result can use it.
func (l *liveLine) stopClear() {
	l.stop()
	l.line.finish("")
}
