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
	
	if len(files) != 2 {
		t.Errorf("Expected 2 files >= 100MB, got %d", len(files))
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
