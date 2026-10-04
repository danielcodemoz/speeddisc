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

func TestMenuColumnsFitConsole(t *testing.T) {
	var b bytes.Buffer
	ui := NewUI(&b, strings.NewReader(""))
	ui.color = false
	ui.utf8 = true
	ui.Banner(LangPT)
	ui.Menu(LangPT)
	s := b.String()
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			continue
		}
		n++
		w := len([]rune(line))
		if w > 118 {
			t.Fatalf("line wider than 118 (%d): %q", w, line)
		}
	}
	if n > 46 {
		t.Fatalf("menu still too tall: %d lines", n)
	}
	foundPkg, foundAct := false, false
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, "Pacote Rápido") && strings.Contains(line, "Pacote Profundo") {
			foundPkg = true
		}
		if strings.Contains(line, "Ficheiros temporários") && strings.Contains(line, "Cache do Windows Update") {
			foundAct = true
		}
		if strings.Contains(line, "SoftwareDistribution") && !strings.Contains(line, "SoftwareDistribution\\Download") && !strings.Contains(line, "SoftwareDistribution\\Download.") {
			// wrapped pieces are ok; a clipped token would drop the end
		}
	}
	if !strings.Contains(s, "SoftwareDistribution\\Download") {
		t.Fatal("description was clipped")
	}
	if !foundPkg || !foundAct {
		t.Fatalf("expected side-by-side columns pkg=%v act=%v", foundPkg, foundAct)
	}
	for _, leak := range []string{"Quick package", "Temporary files", "Full package", "Look only"} {
		if strings.Contains(s, leak) {
			t.Fatalf("english %q in pt menu", leak)
		}
	}
	b.Reset()
	ui.out = &b
	ui.Banner(LangEN)
	ui.Menu(LangEN)
	en := b.String()
	for _, leak := range []string{"Pacote", "Ficheiros", "Sair", "Temporários", "completo"} {
		if strings.Contains(en, leak) {
			t.Fatalf("portuguese %q in en menu", leak)
		}
	}
}
