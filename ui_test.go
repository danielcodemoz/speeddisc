package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestMenuReadableWithoutANSI(t *testing.T) {
	var b bytes.Buffer
	ui := NewUI(&b, strings.NewReader(""))
	ui.color = false
	ui.utf8 = true
	ui.Banner(LangPT)
	ui.Menu(LangPT)
	s := b.String()
	if strings.Contains(s, "\x1b") {
		t.Fatal("ansi leaked")
	}
	for _, needle := range []string{"SpeedDisc", "Pacote Rápido", "Pacote Profundo", "Pacote completo", "Pacote Só ver", "Sair", "English"} {
		if !strings.Contains(s, needle) {
			t.Fatalf("missing %q", needle)
		}
	}
	assertBoxAligned(t, s)
	b.Reset()
	ui.out = &b
	ui.Banner(LangEN)
	ui.Menu(LangEN)
	en := b.String()
	if !strings.Contains(en, "Quick package") || !strings.Contains(en, "Full package") || !strings.Contains(en, "Português") {
		t.Fatal(en)
	}
}

func TestMenuColorFallbackStillAligns(t *testing.T) {
	var b bytes.Buffer
	ui := NewUI(&b, strings.NewReader(""))
	ui.color = true
	ui.utf8 = false
	ui.Menu(LangEN)
	s := b.String()
	if !strings.Contains(s, "\x1b[") {
		t.Fatal("expected color")
	}
	plain := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(s, "")
	if strings.Contains(plain, "┌") {
		t.Fatal("utf box in ascii mode")
	}
	assertBoxAligned(t, plain)
}

func assertBoxAligned(t *testing.T, s string) {
	t.Helper()
	var width int
	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			continue
		}
		r := []rune(line)
		if len(r) == 0 {
			continue
		}
		c := r[0]
		if c != '┌' && c != '│' && c != '└' && c != '+' && c != '|' {
			continue
		}
		if width == 0 {
			width = len(r)
		}
		if len(r) != width {
			t.Fatalf("width %d != %d for %q", len(r), width, line)
		}
	}
	if width == 0 {
		t.Fatal("no box")
	}
}

func menuAt(t *testing.T, cols int) string {
	t.Helper()
	var b bytes.Buffer
	ui := NewUI(&b, strings.NewReader(""))
	ui.color = false
	ui.utf8 = true
	ui.cols = cols
	ui.Banner(LangPT)
	ui.Menu(LangPT)
	s := b.String()
	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			continue
		}
		w := len([]rune(line))
		if w > cols {
			t.Fatalf("cols %d: line wider than window (%d): %q", cols, w, line)
		}
	}
	return s
}

func TestMenuColumnsFollowWidth(t *testing.T) {
	wide := menuAt(t, 140)
	if !strings.Contains(wide, "SoftwareDistribution\\Download") {
		t.Fatal("description was clipped")
	}
	pkg, act2, act3 := false, false, false
	for _, line := range strings.Split(wide, "\n") {
		if strings.Contains(line, "Pacote Rápido") && strings.Contains(line, "Pacote Profundo") {
			pkg = true
		}
		if strings.Contains(line, "Ficheiros temporários") && strings.Contains(line, "Cache do Windows Update") && strings.Contains(line, "Cache de miniaturas") {
			act3 = true
		}
	}
	if !pkg || !act3 {
		t.Fatalf("wide layout pkg=%v act3=%v", pkg, act3)
	}
	for _, leak := range []string{"Quick package", "Temporary files", "Full package", "Look only"} {
		if strings.Contains(wide, leak) {
			t.Fatalf("english %q in pt menu", leak)
		}
	}

	mid := menuAt(t, 90)
	for _, line := range strings.Split(mid, "\n") {
		if strings.Contains(line, "Ficheiros temporários") && strings.Contains(line, "Cache do Windows Update") {
			act2 = true
			if strings.Contains(line, "Cache de miniaturas") {
				t.Fatalf("three actions on a medium line: %q", line)
			}
		}
	}
	if !act2 {
		t.Fatal("expected two action columns at 90")
	}

	narrow := menuAt(t, 48)
	for _, line := range strings.Split(narrow, "\n") {
		if strings.Contains(line, "Ficheiros temporários") && strings.Contains(line, "Cache do Windows Update") {
			t.Fatalf("two actions on a narrow line: %q", line)
		}
		if strings.Contains(line, "Pacote Rápido") && strings.Contains(line, "Pacote Profundo") {
			t.Fatalf("two packages on a narrow line: %q", line)
		}
	}
	if !strings.Contains(narrow, "não") || !strings.Contains(narrow, "Informação") {
		t.Fatal("portuguese accents missing from the menu")
	}

	var b bytes.Buffer
	ui := NewUI(&b, strings.NewReader(""))
	ui.cols = 100
	ui.utf8 = true
	ui.Banner(LangEN)
	ui.Menu(LangEN)
	en := b.String()
	for _, leak := range []string{"Pacote", "Ficheiros", "Sair", "Temporários", "completo"} {
		if strings.Contains(en, leak) {
			t.Fatalf("portuguese %q in en menu", leak)
		}
	}
}
