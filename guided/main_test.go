package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiskMetadataOnly(t *testing.T) {
	p := t.TempDir()
	os.WriteFile(filepath.Join(p, "a.txt"), []byte("abc"), 0600)
	os.Mkdir(filepath.Join(p, "CrossDevice"), 0700)
	os.WriteFile(filepath.Join(p, "CrossDevice", "private"), []byte("never"), 0600)
	linkErr := os.Symlink("a.txt", filepath.Join(p, "link"))
	expectedSkipped := 1
	if linkErr == nil {
		expectedSkipped++
	}
	d, e := scan(context.Background(), p, 100)
	if e != nil || d.Files != 1 || d.Bytes != 3 || d.Skipped != expectedSkipped {
		t.Fatalf("%+v %v", d, e)
	}
	b, _ := os.ReadFile(filepath.Join(p, "a.txt"))
	if string(b) != "abc" {
		t.Fatal("changed source")
	}
}
func TestBoundedAndCancelled(t *testing.T) {
	p := t.TempDir()
	for _, s := range []string{"a", "b", "c"} {
		os.WriteFile(filepath.Join(p, s), []byte("x"), 0600)
	}
	d, _ := scan(context.Background(), p, 1)
	if !d.Stopped || d.Files != 1 {
		t.Fatal(d)
	}
	ctx, c := context.WithCancel(context.Background())
	c()
	d, _ = scan(ctx, p, 100)
	if !d.Stopped || d.Files != 0 {
		t.Fatal(d)
	}
}
func TestCrossDeviceRoot(t *testing.T) {
	if _, e := scan(context.Background(), filepath.Join(t.TempDir(), "cRoSsDeViCe"), 100); e == nil {
		t.Fatal("accepted")
	}
}
func TestNoOverwrite(t *testing.T) {
	p := filepath.Join(t.TempDir(), "report.json")
	if e := saveReport(p, Report{Flow: "disk"}); e != nil {
		t.Fatal(e)
	}
	if e := saveReport(p, Report{Flow: "other"}); e == nil {
		t.Fatal("overwrote")
	}
	b, _ := os.ReadFile(p)
	if !bytes.Contains(b, []byte("disk")) {
		t.Fatal(string(b))
	}
}
func TestMenuRepeatedBackLanguageEOF(t *testing.T) {
	var b bytes.Buffer
	run(strings.NewReader("L\ninvalid\n1\nb\nL\nq\n"), &b)
	if !strings.Contains(b.String(), "Disk full") || !strings.Contains(b.String(), "Disco cheio") {
		t.Fatal(b.String())
	}
	run(strings.NewReader(""), &b)
}
func TestMissingRoot(t *testing.T) {
	_, e := scan(context.Background(), filepath.Join(t.TempDir(), "missing"), 100)
	if e == nil {
		t.Fatal("missing accepted")
	}
}

func TestDeniedDirectoryPartial(t *testing.T) {
	p := t.TempDir()
	denied := filepath.Join(p, "denied")
	if e := os.Mkdir(denied, 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.Chmod(denied, 0000); e != nil {
		t.Skip(e)
	}
	defer os.Chmod(denied, 0700)
	if f, e := os.Open(denied); e == nil {
		f.Close()
		t.Skip("account can read denied fixture")
	}
	d, e := scan(context.Background(), p, 100)
	if e != nil {
		t.Fatal(e)
	}
	if d.Denied != 1 {
		t.Fatalf("expected one denied directory: %+v", d)
	}
}

func TestHumanReport(t *testing.T) {
	var b bytes.Buffer
	showReport(&b, "en", Report{Flow: "disk", Findings: Disk{Files: 1, Bytes: 2048, Largest: []Large{{"photo.jpg", 2048}}}, NextSteps: "Review backups", Partial: true})
	for _, want := range []string{"Largest visited files", "photo.jpg", "Next step: Review backups", "Partial observation"} {
		if !strings.Contains(b.String(), want) {
			t.Fatal(want, b.String())
		}
	}
	if strings.Contains(b.String(), `"Findings"`) {
		t.Fatal("raw JSON screen")
	}
	b.Reset()
	showReport(&b, "pt", Report{Flow: "internet", Findings: []Adapter{{Name: "Local", Up: true, MTU: 1500}}})
	if !strings.Contains(b.String(), "activo") {
		t.Fatal(b.String())
	}
}
func TestReportExclusions(t *testing.T) {
	p := filepath.Join(t.TempDir(), "CrossDevice")
	os.Mkdir(p, 0700)
	if e := saveReport(filepath.Join(p, "report.json"), Report{}); e == nil {
		t.Fatal("CrossDevice report accepted")
	}
}
