package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAdaptiveMenu(t *testing.T) {
	for _, lang := range []Lang{LangPT, LangEN} {
		for _, width := range []int{35, 40, 60, 80, 120} {
			var got bytes.Buffer
			ui := NewUI(&got, strings.NewReader(""))
			ui.cols = width
			for _, rows := range []int{8, 14, 17, 25} {
				got.Reset()
				ui.DrawMenu(lang, rows, false)
				if !strings.Contains(got.String(), "SpeedDisc") || !strings.Contains(got.String(), "A ") {
					t.Fatalf("missing brand or access to actions: %d × %d: %q", width, rows, got.String())
				}
				if n := strings.Count(got.String(), "\n") + 2; n > rows {
					t.Fatalf("%d lines do not fit %d × %d", n, width, rows)
				}
				if strings.ContainsRune(got.String(), '\x1b') {
					t.Fatal("unexpected terminal control")
				}
			}
			got.Reset()
			ui.DrawMenu(lang, 8, true)
			if !strings.Contains(got.String(), T(lang, "opt_flush_dns")) && width >= 40 {
				t.Fatal("full action list missing after explicit selection")
			}
		}
	}
}

func TestAdaptiveMenuFitsAndUnknownWidth(t *testing.T) {
	for _, height := range []int{0, 1000} {
		var got bytes.Buffer
		ui := NewUI(&got, strings.NewReader(""))
		ui.cols = 80
		ui.DrawMenu(LangPT, height, false)
		if !strings.Contains(got.String(), T(LangPT, "opt_flush_dns")) {
			t.Fatal("existing complete menu not preserved")
		}
	}
}
