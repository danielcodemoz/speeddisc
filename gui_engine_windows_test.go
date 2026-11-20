//go:build windows && speeddiscgui

package main

import (
	"fmt"
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
