package main

import (
	"strings"
	"testing"
	"time"
)

func TestReportLanguageAndNotes(t *testing.T) {
	r := &Report{
		Started:  time.Date(2026, 10, 4, 15, 4, 0, 0, time.UTC),
		Computer: "BANCADA",
		Version:  version,
		Found:    []msg{{"found_os", []any{"Windows 11 Pro 24H2"}}},
		Done:     []msg{{"done_clean", []any{tr("name_user_temp"), 2, size(1536), 1}}},
	}
	for _, lang := range []Lang{LangPT, LangEN} {
		body := r.render(lang)
		for _, key := range []string{"rep_title", "rep_found", "rep_done", "rep_failed", "rep_notes", "note_recycle", "note_services", "note_prefetch", "note_sysmain"} {
			if !strings.Contains(body, T(lang, key)) {
				t.Fatalf("lang %d missing %s\n%s", lang, key, body)
			}
		}
		if strings.Contains(body, "%!") {
			t.Fatal(body)
		}
		for _, banned := range []string{"100%", "50%", "% faster", "% mais rápido", "boost"} {
			if strings.Contains(strings.ToLower(body), banned) {
				t.Fatalf("banned %q in report", banned)
			}
		}
	}
	if !strings.Contains(r.render(LangPT), "1,5 KiB") {
		t.Fatal(r.render(LangPT))
	}
	if !strings.Contains(r.render(LangEN), "1.5 KiB") {
		t.Fatal(r.render(LangEN))
	}
	if !strings.Contains(r.render(LangPT), "português") {
		t.Fatal("pt lang name")
	}
	if !strings.Contains(r.render(LangEN), "English") {
		t.Fatal("en lang name")
	}
}

func TestFormatBytes(t *testing.T) {
	if FormatBytes(0, LangEN) != "0 B" {
		t.Fatal(FormatBytes(0, LangEN))
	}
	if FormatBytes(1024, LangPT) != "1,0 KiB" {
		t.Fatal(FormatBytes(1024, LangPT))
	}
	if FormatBytes(1024, LangEN) != "1.0 KiB" {
		t.Fatal(FormatBytes(1024, LangEN))
	}
}
