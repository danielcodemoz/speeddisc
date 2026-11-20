//go:build windows

package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type largeFile struct {
	Path string
	Size int64
}

func findTopLargeFiles(root string, topN int, timeout time.Duration) []largeFile {
	deadline := time.Now().Add(timeout)
	var files []largeFile

	filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if time.Now().After(deadline) {
			return fs.SkipAll
		}
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() < 100*1024*1024 {
			return nil
		}
		files = append(files, largeFile{Path: path, Size: info.Size()})
		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > topN {
		files = files[:topN]
	}
	return files
}

func cleanBrowserCaches(paths []string, tick func(int)) []Stats {
	var out []Stats
	base := 0
	for _, p := range paths {
		if p == "" {
			out = append(out, Stats{Missing: true})
			continue
		}
		if strings.HasSuffix(normPath(p), "/mozilla/firefox/profiles") {
			st := cleanFirefoxProfiles(p, func(n int) {
				if tick != nil {
					tick(base + n)
				}
			})
			base += st.Removed + st.Skipped
			out = append(out, st)
			continue
		}
		st := cleanBrowserLeaf(p, func(n int) {
			if tick != nil {
				tick(base + n)
			}
		})
		base += st.Removed + st.Skipped
		out = append(out, st)
	}
	return out
}

func cleanFirefoxProfiles(profilesDir string, tick func(int)) Stats {
	if !strings.HasSuffix(normPath(profilesDir), "/mozilla/firefox/profiles") {
		return Stats{Err: ErrNotAllowed}
	}
	fi, err := os.Lstat(profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return Stats{Missing: true}
		}
		return Stats{Err: err}
	}
	if fi.Mode()&os.ModeSymlink != 0 || unsafePathHasReparse(profilesDir, false) {
		return Stats{Err: ErrSymlink}
	}
	if !fi.IsDir() {
		return Stats{Missing: true}
	}

	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return Stats{Err: err}
	}

	var combined Stats
	base := 0
	for _, e := range entries {
		if !e.IsDir() || e.Type()&os.ModeSymlink != 0 {
			continue
		}
		cachePath := filepath.Join(profilesDir, e.Name(), "cache2")
		st := cleanBrowserLeaf(cachePath, func(n int) {
			if tick != nil {
				tick(base + n)
			}
		})
		base += st.Removed + st.Skipped
		combined.Removed += st.Removed
		combined.Freed += st.Freed
		combined.Skipped += st.Skipped
		if st.Err != nil && combined.Err == nil && !st.Missing {
			combined.Err = st.Err
		}
	}
	return combined
}

func cleanBrowserLeaf(path string, tick func(int)) Stats {
	if err := AllowClean(CleanBrowserCache, path); err != nil {
		return Stats{Err: err}
	}
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return Stats{Missing: true}
		}
		return Stats{Err: err}
	}
	if unsafePathHasReparse(path, false) {
		return Stats{Err: ErrSymlink}
	}
	return CleanNotify(CleanBrowserCache, path, tick)
}
