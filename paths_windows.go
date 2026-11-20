//go:build windows

package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func windowsRoot() string {
	if s := os.Getenv("SystemRoot"); s != "" {
		return s
	}
	return `C:\Windows`
}

func windowsTemp() string {
	return filepath.Join(windowsRoot(), "Temp")
}

func wuDownload() string {
	return filepath.Join(windowsRoot(), "SoftwareDistribution", "Download")
}

func doCache() string {
	return filepath.Join(windowsRoot(), "ServiceProfiles", "NetworkService", "AppData", "Local", "Microsoft", "Windows", "DeliveryOptimization", "Cache")
}

func thumbDir() string {
	la := os.Getenv("LOCALAPPDATA")
	if la == "" {
		return ""
	}
	return filepath.Join(la, "Microsoft", "Windows", "Explorer")
}

func rawUserTemps() []string {
	var in []string
	for _, k := range []string{"TEMP", "TMP"} {
		if v := os.Getenv(k); v != "" {
			in = append(in, v)
		}
	}
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		in = append(in, filepath.Join(la, "Temp"))
	}
	return in
}

func userTempTargets() []string {
	win := normPath(windowsTemp())
	seen := map[string]bool{}
	var out []string
	for _, p := range rawUserTemps() {
		n := normPath(p)
		if n == "" || n == win || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, p)
	}
	return out
}

func werTargets() []string {
	pd := os.Getenv("ProgramData")
	if pd == "" {
		pd = `C:\ProgramData`
	}
	var bases []string
	bases = append(bases, filepath.Join(pd, `Microsoft\Windows\WER`))
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		bases = append(bases, filepath.Join(la, `Microsoft\Windows\WER`))
	}
	var out []string
	for _, b := range bases {
		for _, sub := range []string{"ReportArchive", "ReportQueue", "Temp"} {
			out = append(out, filepath.Join(b, sub))
		}
	}
	return out
}

func programDataOrDefault() string {
	if pd := os.Getenv("ProgramData"); pd != "" {
		return pd
	}
	return `C:\ProgramData`
}

// Helper programs are always resolved relative to the OS-reported system
// directory. A missing helper must not fall back to PATH or the working directory.
func sysExe(rel string) string {
	dir, err := windows.GetSystemDirectory()
	if err != nil || !filepath.IsAbs(dir) {
		return ""
	}
	return filepath.Join(dir, filepath.FromSlash(rel))
}

func browserCachePaths() []string {
	var paths []string
	la := os.Getenv("LOCALAPPDATA")
	if la == "" {
		return paths
	}
	paths = append(paths,
		filepath.Join(la, "Google", "Chrome", "User Data", "Default", "Cache"),
		filepath.Join(la, "Google", "Chrome", "User Data", "Default", "Code Cache"),
		filepath.Join(la, "Microsoft", "Edge", "User Data", "Default", "Cache"),
		filepath.Join(la, "Microsoft", "Edge", "User Data", "Default", "Code Cache"),
		filepath.Join(la, "Mozilla", "Firefox", "Profiles"),
	)
	return paths
}

// Refuse reparse points along an existing path. In a startup destination,
// the final directory may be missing, but its parents must exist.
func unsafePathHasReparse(path string, allowMissingLeaf bool) bool {
	current := filepath.Clean(path)
	volume := filepath.VolumeName(current)
	if len(volume) != 2 || volume[1] != ':' || !filepath.IsAbs(current) {
		return true
	}
	leaf := true
	for {
		p, err := windows.UTF16PtrFromString(current)
		if err != nil {
			return true
		}
		attrs, err := windows.GetFileAttributes(p)
		if err != nil {
			if !(leaf && allowMissingLeaf && err == windows.ERROR_FILE_NOT_FOUND) {
				return true
			}
		} else if attrs&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 || (allowMissingLeaf && leaf && attrs&windows.FILE_ATTRIBUTE_DIRECTORY == 0) {
			return true
		}
		up := filepath.Dir(current)
		if up == current {
			return false
		}
		current, leaf = up, false
	}
}
