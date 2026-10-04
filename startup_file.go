package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func skipStartupName(name string) bool {
	switch strings.ToLower(name) {
	case "desktop.ini", "thumbs.db", "_speeddisc_disabled":
		return true
	default:
		return false
	}
}

// disableStartupFile moves a file into a subfolder of its Startup directory.
// Windows does not run items in subfolders of Startup. The target of a shortcut is not touched.
func disableStartupFile(file string) (string, error) {
	dir := filepath.Dir(file)
	if !IsStartupDir(dir) {
		return "", ErrNotAllowed
	}
	base := filepath.Base(file)
	if base == "." || base == ".." || skipStartupName(base) || strings.ContainsAny(base, `/\`) {
		return "", ErrNotAllowed
	}
	destDir := filepath.Join(dir, "_SpeedDisc_disabled")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(destDir, base)
	if _, err := os.Lstat(dest); err == nil {
		dest = filepath.Join(destDir, fmt.Sprintf("%d-%s", time.Now().Unix(), base))
	}
	if err := os.Rename(file, dest); err != nil {
		return "", err
	}
	return dest, nil
}
