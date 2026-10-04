package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var thumbPatterns = []string{"thumbcache_*.db", "iconcache_*.db"}

type Stats struct {
	Removed int
	Freed   int64
	Skipped int
	Missing bool
	Err     error
}

func Clean(kind CleanKind, path string) Stats {
	if err := AllowClean(kind, path); err != nil {
		return Stats{Err: err}
	}
	st, err := cleanTree(path)
	if err != nil {
		st.Err = err
	}
	return st
}

func CleanMatching(path string, patterns []string) Stats {
	if err := AllowClean(CleanThumb, path); err != nil {
		return Stats{Err: err}
	}
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Stats{Missing: true}
		}
		return Stats{Err: err}
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return Stats{Err: ErrSymlink}
	}
	if !fi.IsDir() {
		return Stats{Err: ErrNotAllowed}
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return Stats{Err: err}
	}
	var st Stats
	for _, e := range entries {
		if e.IsDir() || e.Type()&os.ModeSymlink != 0 {
			continue
		}
		if !matchAny(e.Name(), patterns) {
			continue
		}
		full := filepath.Join(path, e.Name())
		if skipProtected(full) {
			st.Skipped++
			continue
		}
		var sz int64
		if info, err := e.Info(); err == nil && info.Mode().IsRegular() {
			sz = info.Size()
		}
		if err := os.Remove(full); err != nil {
			st.Skipped++
			continue
		}
		st.Removed++
		st.Freed += sz
	}
	return st
}

func matchAny(name string, patterns []string) bool {
	low := strings.ToLower(name)
	for _, p := range patterns {
		ok, err := pathMatch(strings.ToLower(p), low)
		if err == nil && ok {
			return true
		}
	}
	return false
}

func pathMatch(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

func measureTree(root string) (int64, error) {
	if pathTooShallow(root) {
		return 0, ErrShallow
	}
	fi, err := os.Lstat(root)
	if err != nil {
		return 0, err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return 0, ErrSymlink
	}
	if !fi.IsDir() {
		return 0, ErrNotAllowed
	}
	var total int64
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || path == root {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total, err
}

func cleanTree(root string) (Stats, error) {
	if pathTooShallow(root) {
		return Stats{}, ErrShallow
	}
	root = filepath.Clean(root)
	fi, err := os.Lstat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return Stats{Missing: true}, nil
		}
		return Stats{}, err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return Stats{}, ErrSymlink
	}
	if !fi.IsDir() {
		return Stats{}, ErrNotAllowed
	}
	var st Stats
	var dirs []string
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if path == root {
			return nil
		}
		if walkErr != nil {
			st.Skipped++
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			st.Skipped++
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			dirs = append(dirs, path)
			return nil
		}
		if skipProtected(path) {
			st.Skipped++
			return nil
		}
		var sz int64
		if info, err := d.Info(); err == nil && info.Mode().IsRegular() {
			sz = info.Size()
		}
		if err := os.Remove(path); err != nil {
			st.Skipped++
			return nil
		}
		st.Removed++
		st.Freed += sz
		return nil
	})
	if err != nil {
		return st, err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, d := range dirs {
		_ = os.Remove(d)
	}
	return st, nil
}

func skipProtected(path string) bool {
	b := strings.ToLower(filepath.Base(path))
	if b == "speeddisc.exe" {
		return true
	}
	if strings.HasPrefix(b, "speeddisc-report-") {
		return true
	}
	return false
}
