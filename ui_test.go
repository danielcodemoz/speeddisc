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
	for _, needle := range []string{"SpeedDisc", "Pacote Rápido", "Pacote Profundo", "Pacote Só ver", "Sair", "English"} {
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
	if !strings.Contains(en, "Quick package") || !strings.Contains(en, "Português") {
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
