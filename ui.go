package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const boxInner = 64

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

func (u *UI) info(s string) { fmt.Fprintln(u.out, "  "+s) }
func (u *UI) ok(s string)   { fmt.Fprintln(u.out, u.paint(ansiGreen, "  "+s)) }
func (u *UI) warn(s string) { fmt.Fprintln(u.out, u.paint(ansiYellow, "  "+s)) }
func (u *UI) err(s string)  { fmt.Fprintln(u.out, u.paint(ansiRed, "  "+s)) }
func (u *UI) step(s string) { fmt.Fprintln(u.out, u.paint(ansiBoldCyan, "  "+s)) }

func (u *UI) paragraph(s string) {
	for _, line := range wrap(s, 76) {
		fmt.Fprintln(u.out, "  "+line)
	}
}

func (u *UI) Banner(lang Lang) {
	fmt.Fprintln(u.out)
	fmt.Fprintln(u.out, u.paint(ansiBoldCyan, "  SpeedDisc"))
	mark := "·"
	if !u.utf8 {
		mark = "-"
	}
	fmt.Fprintf(u.out, "  v%s %s %s\n", version, mark, T(lang, "subtitle"))
	fmt.Fprintln(u.out, "  "+T(lang, "lang_now"))
}

func (u *UI) Menu(lang Lang) {
	fmt.Fprintln(u.out)
	u.box(T(lang, "section_packages"), u.itemLines(lang, "pkg"))
	fmt.Fprintln(u.out)
	u.box(T(lang, "section_actions"), u.itemLines(lang, "act"))
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "   L   %s\n", T(lang, "opt_lang"))
	fmt.Fprintf(u.out, "   Q   %s\n", T(lang, "opt_quit"))
	fmt.Fprintln(u.out)
	u.info(T(lang, "menu_hint"))
}

func (u *UI) itemLines(lang Lang, section string) []string {
	var lines []string
	for _, it := range menuCatalog() {
		if it.section != section {
			continue
		}
		title := fmt.Sprintf(" %2d  %s", it.num, T(lang, it.title))
		lines = append(lines, title)
		for _, dl := range wrap(T(lang, it.desc), boxInner-6) {
			lines = append(lines, "      "+dl)
		}
		lines = append(lines, "")
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func (u *UI) box(title string, lines []string) {
	horiz, vert, tl, tr, bl, br := "-", "|", "+", "+", "+", "+"
	if u.utf8 {
		horiz, vert, tl, tr, bl, br = "─", "│", "┌", "┐", "└", "┘"
	}
	t := title
	topFill := boxInner - runeLen(t) - 2
	if topFill < 1 {
		t = padRight(t, boxInner-3)
		topFill = 1
	}
	top := tl + horiz + t + " " + strings.Repeat(horiz, topFill) + tr
	fmt.Fprintln(u.out, u.paint(ansiCyan, top))
	for _, line := range lines {
		fmt.Fprintf(u.out, "%s%s%s\n", u.paint(ansiCyan, vert), padRight(line, boxInner), u.paint(ansiCyan, vert))
	}
	bottom := bl + strings.Repeat(horiz, boxInner) + br
	fmt.Fprintln(u.out, u.paint(ansiCyan, bottom))
}

func (u *UI) section(lang Lang, titleKey string, lines []msg) {
	fmt.Fprintln(u.out)
	u.info(T(lang, titleKey))
	if len(lines) == 0 {
		u.info(T(lang, "rep_none"))
		return
	}
	for _, m := range lines {
		u.info(m.Text(lang))
	}
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
