//go:build windows

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindTopLargeFiles(t *testing.T) {
	dir := t.TempDir()

	sizes := []int64{
		150 * 1024 * 1024,
		200 * 1024 * 1024,
		100 * 1024 * 1024,
	}

	for i, sz := range sizes {
		name := filepath.Join(dir, "file"+string(rune('A'+i))+".dat")
		f, err := os.Create(name)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := f.Truncate(sz); err != nil {
			f.Close()
			t.Fatalf("Failed to truncate test file: %v", err)
		}
		f.Close()
	}

	files := findTopLargeFiles(dir, 10, 5*time.Second)

	if len(files) != 3 {
		t.Errorf("Expected 3 files >= 100 MiB, got %d", len(files))
	}

	if len(files) > 0 && files[0].Size < files[len(files)-1].Size {
		t.Errorf("Files not sorted by size descending")
	}
}

func TestCleanBrowserCaches(t *testing.T) {
	stats := cleanBrowserCaches([]string{}, nil)
	if len(stats) != 0 {
		t.Errorf("Expected 0 stats for empty paths, got %d", len(stats))
	}

	stats = cleanBrowserCaches([]string{""}, nil)
	if len(stats) != 1 || !stats[0].Missing {
		t.Errorf("Expected 1 missing stat for empty path")
	}
}

func TestBrowserCacheExactLeaves(t *testing.T) {
	dir := t.TempDir()
	chrome := filepath.Join(dir, "Google", "Chrome", "User Data", "Default", "Cache")
	firefox := filepath.Join(dir, "Mozilla", "Firefox", "Profiles", "example", "cache2")
	fixtures := []string{filepath.Join(chrome, "tmp.bin"), filepath.Join(firefox, "tmp.bin")}
	keep := filepath.Join(filepath.Dir(chrome), "Cookies")
	for _, p := range append(fixtures, keep) {
		if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("fixture"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := cleanBrowserCaches([]string{chrome, filepath.Dir(filepath.Dir(firefox))}, nil)
	if len(got) != 2 || got[0].Err != nil || got[1].Err != nil || got[0].Removed != 1 || got[1].Removed != 1 {
		t.Fatalf("unexpected results: %#v", got)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("profile file changed: %v", err)
	}
}
