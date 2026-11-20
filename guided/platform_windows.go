//go:build windows

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var kernel = syscall.NewLazyDLL("kernel32.dll")
var driveType = kernel.NewProc("GetDriveTypeW")

func unsafeEntry(p string) bool {
	u, e := syscall.UTF16PtrFromString(p)
	if e != nil {
		return true
	}
	a, e := syscall.GetFileAttributes(u)
	return e != nil || a&0x400 != 0
}
func localRoot(p string) error {
	if strings.HasPrefix(p, `\\`) {
		return fmt.Errorf("UNC/device paths excluded")
	}
	vol := filepath.VolumeName(p)
	if len(vol) != 2 || vol[1] != ':' {
		return fmt.Errorf("local drive required")
	}
	u, _ := syscall.UTF16PtrFromString(vol + `\`)
	t, _, _ := driveType.Call(uintptr(unsafe.Pointer(u)))
	if t != 3 {
		return fmt.Errorf("fixed local drive required")
	}
	for q := p; ; q = filepath.Dir(q) {
		if unsafeEntry(q) {
			return fmt.Errorf("reparse point or inaccessible path excluded")
		}
		if filepath.Dir(q) == q {
			break
		}
	}
	return nil
}
func startup(ctx context.Context) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var system [32768]uint16
	n, _, _ := kernel.NewProc("GetSystemDirectoryW").Call(uintptr(unsafe.Pointer(&system[0])), uintptr(len(system)))
	if n == 0 || n >= uintptr(len(system)) {
		return nil, fmt.Errorf("system directory unavailable")
	}
	shell := filepath.Join(syscall.UTF16ToString(system[:n]), "WindowsPowerShell", "v1.0", "powershell.exe")
	// Fixed read-only query. Names only: never return Command (which may contain secrets).
	cmd := exec.CommandContext(ctx, shell, "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", `[Console]::OutputEncoding=[System.Text.Encoding]::UTF8; $ErrorActionPreference='Stop'; @(Get-CimInstance Win32_StartupCommand | Select-Object -ExpandProperty Name) | ConvertTo-Json -Compress`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	b, e := cmd.Output()
	if e != nil {
		return nil, fmt.Errorf("startup inventory unavailable (timeout/policy/provider)")
	}
	var names any
	if len(strings.TrimSpace(string(b))) == 0 {
		return []string{}, nil
	}
	if err := json.Unmarshal(b, &names); err != nil {
		return nil, fmt.Errorf("startup response unreadable")
	}
	return names, nil
}
