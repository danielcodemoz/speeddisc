package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	ansiReset    = "\x1b[0m"
	ansiCyan     = "\x1b[36m"
	ansiGreen    = "\x1b[32m"
	ansiYellow   = "\x1b[33m"
	ansiRed      = "\x1b[31m"
	ansiBoldCyan = "\x1b[1;36m"
)

type UI struct {
	out   io.Writer
	in    *bufio.Reader
	color bool
	utf8  bool
	// cols is an optional visible width in columns. Zero means the real console width.
	cols int
}

func NewUI(out io.Writer, in io.Reader) *UI {
	return &UI{out: out, in: bufio.NewReader(in), utf8: true}
}

func (u *UI) paint(code, s string) string {
	if !u.color || code == "" {
		return s
	}
	return code + s + ansiReset
}

func (u *UI) Ask(prompt string) string {
	fmt.Fprint(u.out, prompt)
	line, err := u.in.ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		return "q"
	}
	return strings.TrimSpace(line)
}

func (u *UI) info(s string) { u.emit("", s) }
func (u *UI) ok(s string)   { u.emit(ansiGreen, s) }
func (u *UI) warn(s string) { u.emit(ansiYellow, s) }
func (u *UI) err(s string)  { u.emit(ansiRed, s) }
func (u *UI) step(s string) { u.emit(ansiBoldCyan, s) }

func (u *UI) emit(code, s string) {
	for _, line := range wrap(s, u.textWidth()) {
		fmt.Fprintln(u.out, u.paint(code, "  "+line))
	}
}

func (u *UI) paragraph(s string) {
	u.emit("", s)
}

// rawWrap prints s without going past the visible window.
func (u *UI) rawWrap(s string) {
	limit := u.termWidth() - 1
	if limit < 1 {
		limit = 1
	}
	for _, line := range wrap(s, limit) {
		fmt.Fprintln(u.out, line)
	}
}

// termWidth is the visible console width in columns.
// It is never used to resize the buffer, the window, or the font.
func (u *UI) termWidth() int {
	if u.cols > 0 {
		return u.cols
	}
	if n := consoleColumns(); n > 0 {
		return n
	}
	return 80
}

// textWidth is how many runes of text fit after the two-space indent,
// leaving the last column unused so a line cannot widen the buffer.
func (u *UI) textWidth() int {
	w := u.termWidth() - 3
	if w < 1 {
		return 1
	}
	return w
}

func (u *UI) boxInner() int {
	total := u.termWidth() - 1
	if total < 4 {
		total = u.termWidth()
	}
	inner := total - 2
	if inner < 1 {
		return 1
	}
	return inner
}

func (u *UI) Banner(lang Lang) {
	fmt.Fprintln(u.out)
	fmt.Fprintln(u.out, u.paint(ansiBoldCyan, "  SpeedDisc"))
	mark := "·"
	if !u.utf8 {
		mark = "-"
	}
	head := fmt.Sprintf("v%s %s %s", version, mark, T(lang, "subtitle"))
	for _, line := range wrap(head, u.textWidth()) {
		fmt.Fprintln(u.out, "  "+line)
	}
	for _, line := range wrap(T(lang, "lang_now"), u.textWidth()) {
		fmt.Fprintln(u.out, "  "+line)
	}
}

func (u *UI) Menu(lang Lang) {
	fmt.Fprintln(u.out)
	inner := u.boxInner()
	u.box(T(lang, "section_packages"), u.sectionColumns(lang, "pkg", columnsThatFit(inner, 2), inner))
	fmt.Fprintln(u.out)
	u.box(T(lang, "section_actions"), u.sectionColumns(lang, "act", columnsThatFit(inner, 3), inner))
	fmt.Fprintln(u.out)
	u.rawWrap(fmt.Sprintf("   L   %s", T(lang, "opt_lang")))
	u.rawWrap(fmt.Sprintf("   Q   %s", T(lang, "opt_quit")))
	fmt.Fprintln(u.out)
	u.info(T(lang, "menu_hint"))
}

// columnsThatFit picks one column on a narrow console, and two or three
// only when each column still has room. It does not assume 114 or 120.
func columnsThatFit(inner, want int) int {
	if want < 1 {
		want = 1
	}
	const minCol = 28
	const gap = 2
	best := 1
	for c := 1; c <= want; c++ {
		if c == 1 {
			best = 1
			continue
		}
		if inner >= c*minCol+(c-1)*gap {
			best = c
		}
	}
	return best
}

// sectionColumns lays items left to right, then down, inside the menu box.
func (u *UI) sectionColumns(lang Lang, section string, cols, inner int) []string {
	var items []menuItem
	for _, it := range menuCatalog() {
		if it.section == section {
			items = append(items, it)
		}
	}
	if cols < 1 {
		cols = 1
	}
	gap := 2
	base := (inner - gap*(cols-1)) / cols
	extra := inner - (base*cols + gap*(cols-1))
	widths := make([]int, cols)
	for i := range widths {
		widths[i] = base
		if i >= cols-extra {
			widths[i]++
		}
	}
	blocks := make([][]string, len(items))
	for i, it := range items {
		blocks[i] = itemBlock(lang, it, widths[i%cols])
	}
	rows := (len(blocks) + cols - 1) / cols
	var lines []string
	for r := 0; r < rows; r++ {
		if r > 0 {
			lines = append(lines, "")
		}
		maxH := 0
		for c := 0; c < cols; c++ {
			i := r*cols + c
			if i < len(blocks) && len(blocks[i]) > maxH {
				maxH = len(blocks[i])
			}
		}
		for h := 0; h < maxH; h++ {
			var b strings.Builder
			for c := 0; c < cols; c++ {
				if c > 0 {
					b.WriteString(strings.Repeat(" ", gap))
				}
				i := r*cols + c
				line := ""
				if i < len(blocks) && h < len(blocks[i]) {
					line = blocks[i][h]
				}
				b.WriteString(padRight(line, widths[c]))
			}
			lines = append(lines, b.String())
		}
	}
	return lines
}

func itemBlock(lang Lang, it menuItem, width int) []string {
	if width < 8 {
		width = 8
	}
	var lines []string
	title := fmt.Sprintf("%2d  %s", it.num, T(lang, it.title))
	if runeLen(title) <= width {
		lines = append(lines, title)
	} else {
		for i, tl := range wrap(title, width) {
			if i == 0 {
				lines = append(lines, tl)
				continue
			}
			lines = append(lines, "    "+tl)
		}
	}
	descW := width - 4
	if descW < 4 {
		descW = 4
	}
	for _, dl := range wrap(T(lang, it.desc), descW) {
		lines = append(lines, "    "+dl)
	}
	return lines
}

func (u *UI) box(title string, lines []string) {
	inner := u.boxInner()
	horiz, vert, tl, tr, bl, br := "-", "|", "+", "+", "+", "+"
	if u.utf8 {
		horiz, vert, tl, tr, bl, br = "─", "│", "┌", "┐", "└", "┘"
	}
	label := " " + title + " "
	topFill := inner - runeLen(label) - 1
	if topFill < 1 {
		label = padRight(label, inner-2)
		topFill = 1
	}
	top := tl + horiz + label + strings.Repeat(horiz, topFill) + tr
	fmt.Fprintln(u.out, u.paint(ansiBoldCyan, top))
	for _, line := range lines {
		fmt.Fprintf(u.out, "%s%s%s\n", u.paint(ansiCyan, vert), padRight(line, inner), u.paint(ansiCyan, vert))
	}
	bottom := bl + strings.Repeat(horiz, inner) + br
	fmt.Fprintln(u.out, u.paint(ansiCyan, bottom))
}

func (u *UI) section(lang Lang, titleKey string, lines []msg) {
	fmt.Fprintln(u.out)
	u.step(T(lang, titleKey))
	if len(lines) == 0 {
		u.info(T(lang, "rep_none"))
		return
	}
	for _, m := range lines {
		u.info(m.Text(lang))
	}
}

// track starts one in-place row for a step that has no real percentage.
func (u *UI) track(label, working, countFmt string) *liveLine {
	return startLive(u.out, label, working, countFmt, u.termWidth()-1, func(s string) string {
		return u.paint(ansiBoldCyan, s)
	}, u.utf8)
}

func confirmBody(lang Lang, acts []Action) string {
	parts := []string{T(lang, "confirm_intro")}
	if Contains(acts, ActSFC) || Contains(acts, ActDISM) {
		parts = append(parts, T(lang, "confirm_long"))
	}
	if Contains(acts, ActRestore) {
		parts = append(parts, T(lang, "confirm_restore"))
	}
	if Contains(acts, ActStartup) {
		parts = append(parts, T(lang, "confirm_startup"))
	}
	parts = append(parts, T(lang, "confirm_safe"))
	return strings.Join(parts, " ")
}
