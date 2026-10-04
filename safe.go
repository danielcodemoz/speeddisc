package main

import (
	"errors"
	"strings"
)

var (
	ErrShallow      = errors.New("refusing a drive root or other shallow path")
	ErrDanger       = errors.New("refusing a protected Windows path")
	ErrNotAllowed   = errors.New("path is not an allowed cleanup target")
	ErrSymlink      = errors.New("refusing to clean a symlink")
	ErrValueChanged = errors.New("value changed since it was listed")
	ErrNotText      = errors.New("not a text value")
)

type CleanKind int

const (
	CleanTemp CleanKind = iota
	CleanWU
	CleanDO
	CleanWER
	CleanThumb
)

// normPath lowercases and resolves "." and ".." without using the host
// filepath package, so Windows paths can be checked on any OS.
func normPath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.ReplaceAll(p, `\`, `/`)
	drive := ""
	if len(p) >= 2 && p[1] == ':' {
		drive = strings.ToLower(p[:2])
		p = p[2:]
	}
	abs := strings.HasPrefix(p, "/")
	var out []string
	for _, s := range strings.Split(p, "/") {
		if s == "" || s == "." {
			continue
		}
		if s == ".." {
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
			continue
		}
		out = append(out, s)
	}
	body := strings.ToLower(strings.Join(out, "/"))
	if drive != "" {
		if body == "" {
			return drive
		}
		return drive + "/" + body
	}
	if abs {
		if body == "" {
			return "/"
		}
		return "/" + body
	}
	return body
}

func pathBody(n string) string {
	if len(n) >= 2 && n[1] == ':' {
		if len(n) == 2 {
			return "/"
		}
		return "/" + n[3:]
	}
	if n == "" {
		return "/"
	}
	if !strings.HasPrefix(n, "/") {
		return "/" + n
	}
	return n
}

func pathTooShallow(p string) bool {
	n := normPath(p)
	if n == "" || n == "/" {
		return true
	}
	body := n
	if len(n) >= 2 && n[1] == ':' {
		if len(n) == 2 {
			return true
		}
		body = n[3:]
	} else {
		body = strings.TrimPrefix(body, "/")
	}
	if body == "" {
		return true
	}
	return len(strings.Split(body, "/")) < 2
}

func refuseDangerous(path string) error {
	body := pathBody(normPath(path))
	badPrefixes := []string{
		"/windows/system32",
		"/windows/syswow64",
		"/windows/winsxs",
		"/windows/servicing",
		"/windows/fonts",
		"/windows/systemapps",
		"/program files",
		"/program files (x86)",
	}
	for _, b := range badPrefixes {
		if body == b || strings.HasPrefix(body, b+"/") {
			return ErrDanger
		}
	}
	exact := []string{
		"/windows",
		"/windows/softwaredistribution",
		"/windows/prefetch",
		"/windows/system32",
		"/programdata",
		"/users",
	}
	for _, b := range exact {
		if body == b {
			return ErrDanger
		}
	}
	return nil
}

// AllowClean is the allowlist for destructive deletes. Callers never pass a
// free-typed path: every kind has a required suffix, and protected trees are refused.
func AllowClean(kind CleanKind, path string) error {
	if strings.TrimSpace(path) == "" || pathTooShallow(path) {
		return ErrShallow
	}
	if err := refuseDangerous(path); err != nil {
		return err
	}
	n := normPath(path)
	switch kind {
	case CleanTemp:
		base := n[strings.LastIndex(n, "/")+1:]
		if base != "temp" && base != "tmp" {
			return ErrNotAllowed
		}
	case CleanWU:
		if !strings.HasSuffix(n, "/windows/softwaredistribution/download") {
			return ErrNotAllowed
		}
	case CleanDO:
		if !strings.HasSuffix(n, "/microsoft/windows/deliveryoptimization/cache") ||
			!strings.Contains(n, "/serviceprofiles/networkservice/") {
			return ErrNotAllowed
		}
	case CleanWER:
		ok := strings.HasSuffix(n, "/microsoft/windows/wer/reportarchive") ||
			strings.HasSuffix(n, "/microsoft/windows/wer/reportqueue") ||
			strings.HasSuffix(n, "/microsoft/windows/wer/temp")
		if !ok {
			return ErrNotAllowed
		}
	case CleanThumb:
		if !strings.HasSuffix(n, "/microsoft/windows/explorer") {
			return ErrNotAllowed
		}
	default:
		return ErrNotAllowed
	}
	return nil
}

func IsStartupDir(path string) bool {
	n := normPath(path)
	return strings.HasSuffix(n, "/microsoft/windows/start menu/programs/startup")
}
