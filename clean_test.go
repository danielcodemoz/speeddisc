package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanRefusesUnknownTree(t *testing.T) {
	dir := t.TempDir()
	keep := filepath.Join(dir, "keep.txt")
	if err := os.WriteFile(keep, []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	st := Clean(CleanTemp, dir)
	if st.Err == nil {
		t.Fatal("expected refusal")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("file was deleted")
	}
}

func TestCleanTreeDeletesContentsOnly(t *testing.T) {
	root := filepath.Join(t.TempDir(), "inner")
	if err := os.MkdirAll(filepath.Join(root, "a", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "a", "b", "f.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "SpeedDisc.exe"), []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(outside, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link.txt")); err != nil {
		t.Fatal(err)
	}
	linkedDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(linkedDir, "inside.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(linkedDir, filepath.Join(root, "dirlink")); err != nil {
		t.Fatal(err)
	}
	st, err := cleanTree(root)
	if err != nil {
		t.Fatal(err)
	}
	if st.Removed != 1 || st.Freed != 5 {
		t.Fatalf("stats %+v", st)
	}
	if _, err := os.Stat(filepath.Join(root, "a", "b", "f.txt")); !os.IsNotExist(err) {
		t.Fatal("file still there")
	}
	if _, err := os.Stat(filepath.Join(root, "SpeedDisc.exe")); err != nil {
		t.Fatal("exe should remain", err)
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatal("symlink target removed")
	}
	if _, err := os.Stat(filepath.Join(linkedDir, "inside.txt")); err != nil {
		t.Fatal("linked dir contents removed")
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatal("root removed")
	}
}

func TestCleanTreeSymlinkRoot(t *testing.T) {
	real := t.TempDir()
	if err := os.WriteFile(filepath.Join(real, "f.txt"), []byte("z"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	_, err := cleanTree(link)
	if !errors.Is(err, ErrSymlink) {
		t.Fatalf("got %v", err)
	}
	if _, err := os.Stat(filepath.Join(real, "f.txt")); err != nil {
		t.Fatal(err)
	}
}

func TestMeasureDoesNotDelete(t *testing.T) {
	root := filepath.Join(t.TempDir(), "tree")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("abcd"), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := measureTree(root)
	if err != nil || n != 4 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if _, err := os.Stat(filepath.Join(root, "a.txt")); err != nil {
		t.Fatal(err)
	}
}

func TestMatchThumb(t *testing.T) {
	if !matchAny("thumbcache_32.db", thumbPatterns) {
		t.Fatal("thumb")
	}
	if !matchAny("iconcache_256.db", thumbPatterns) {
		t.Fatal("icon")
	}
	if matchAny("thumbcache_32.db.bak", thumbPatterns) || matchAny("settings.dat", thumbPatterns) {
		t.Fatal("overmatch")
	}
}
