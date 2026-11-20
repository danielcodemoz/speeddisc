//go:build windows && speeddiscgui

package main

import (
	"fmt"
	"slices"
	"testing"
)

func TestGUIActionCatalogueRetainsLegacyIDs(t *testing.T) {
	items := guiItems(LangEN)
	if len(items) != len(menuCatalog()) {
		t.Fatalf("catalogue size %d != %d", len(items), len(menuCatalog()))
	}
	for i, x := range items {
		if x.ID != fmt.Sprint(i+1) || x.Title == "" || x.Detail == "" || x.Category == "" {
			t.Fatalf("invalid catalogue item %#v", x)
		}
		for _, lang := range []Lang{LangPT, LangEN} {
			if out := guiItems(lang)[i]; out.ID != x.ID || out.Category != x.Category || out.Mutates != x.Mutates || out.NeedsAdmin != x.NeedsAdmin {
				t.Fatalf("language changed identity of %#v", x)
			}
		}
	}
	if x := items[13]; !x.Guided || guiRequiresAdmin(x, guiOptions{}) ||
		!guiRequiresAdmin(x, guiOptions{StartupIDs: []string{"fixture"}}) {
		t.Fatalf("startup privileges incorrect for %#v", x)
	}
}

func TestGUIMutationsRequireExplicitConfirmation(t *testing.T) {
	for _, tc := range []struct {
		id      string
		options guiOptions
	}{
		{"1", guiOptions{}},
		{"5", guiOptions{}},
		{"14", guiOptions{StartupIDs: []string{"never-run"}}},
		{"19", guiOptions{Winsock: true}},
		{"15", guiOptions{Winsock: true, Confirmed: true}},
		{"bogus", guiOptions{Confirmed: true}},
	} {
		tc.options.Progress = func(guiProgress) {
			t.Fatal("rejected operation emitted running progress")
		}
		res := guiRun(tc.id, LangEN, tc.options)
		if res.Status != "Not run" || len(res.Rows) != 1 {
			t.Fatalf("unexpected result for %q: %#v", tc.id, res)
		}
	}
}

func TestGUIStartupFingerprintsByExactRegistryValue(t *testing.T) {
	e := StartupEntry{
		Name: "fixture", Where: "HKCU fake", Command: "fixture-command", Key: "fake", Hive: "HKCU", CanDisable: true,
	}
	id, ok := guiStartupID(e)
	if !ok || len(id) != 64 {
		t.Fatalf("invalid id %q for fixture", id)
	}
	e.Command += " changed"
	other, ok := guiStartupID(e)
	if !ok || other == id {
		t.Fatal("changed registry command retained same approval identity")
	}
	e.CanDisable = false
	if id, ok := guiStartupID(e); ok || id != "" {
		t.Fatal("non-text values must not be actionable")
	}
}

func TestGUIProgressOnlyUsesParserPercentages(t *testing.T) {
	var console, gui []int
	p := newPercentFilter(nil)
	p.onPercent = scanPercentSink(func(n int) { console = append(console, n) }, func(n int) { gui = append(gui, n) })
	for _, line := range []string{"waiting\r", "12.5%\r", "12.5%\r", "105%\r", "73 percent\r", "100%\r"} {
		if _, err := p.Write([]byte(line)); err != nil {
			t.Fatal(err)
		}
	}
	p.flush()
	want := []int{12, 100}
	if !slices.Equal(console, want) || !slices.Equal(gui, want) {
		t.Fatalf("percent forwarding: console %v; GUI %v", console, gui)
	}
}

func TestGUIProgressUsesLocalizedCatalogLabels(t *testing.T) {
	for _, a := range actionOrder {
		if a == ActSFC || a == ActDISM {
			continue
		}
		pt, en := actionProgressLabel(LangPT, a), actionProgressLabel(LangEN, a)
		if pt == "" || en == "" || pt == T(LangPT, "working") || en == T(LangEN, "working") {
			t.Fatalf("action %v missing progress label", a)
		}
	}
	if actionProgressLabel(LangPT, ActStartup) == actionProgressLabel(LangEN, ActStartup) {
		t.Fatal("startup progress label not localized")
	}
}
