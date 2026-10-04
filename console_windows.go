//go:build windows

package main

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var procSetConsoleTitleW = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetConsoleTitleW")

func setupConsole() (color bool, utf8 bool) {
	setConsoleTitle("SpeedDisc")
	utf8 = true
	if err := windows.SetConsoleOutputCP(65001); err != nil {
		utf8 = false
	}
	_ = windows.SetConsoleCP(65001)
	return enableVT(os.Stdout), utf8
}

func setConsoleTitle(title string) {
	p, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return
	}
	_, _, _ = procSetConsoleTitleW.Call(uintptr(unsafe.Pointer(p)))
}

func enableVT(f *os.File) bool {
	h := windows.Handle(f.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return false
	}
	if err := windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
		return false
	}
	return true
}

func elevated() bool {
	tok, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer tok.Close()
	return tok.IsElevated()
}
