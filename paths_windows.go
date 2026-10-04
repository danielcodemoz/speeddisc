//go:build windows

package main

import (
	"os"
	"path/filepath"
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

func sysExe(rel string) string {
	p := filepath.Join(windowsRoot(), "System32", filepath.FromSlash(rel))
	if st, err := os.Stat(p); err == nil && !st.IsDir() {
		return p
	}
	return filepath.Base(filepath.FromSlash(rel))
}
