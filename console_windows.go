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
	// UTF-8 code page only. Do not touch the font, the window, or the screen buffer.
	// A buffer wider than the window is what makes Windows shrink the font.
	if err := windows.SetConsoleOutputCP(65001); err != nil {
		utf8 = false
	}
	_ = windows.SetConsoleCP(65001)
	return enableVT(os.Stdout), utf8
}

// consoleColumns is the visible window width, not the screen-buffer width.
// Nothing here calls SetConsoleScreenBufferSize, SetConsoleWindowInfo, or SetCurrentConsoleFontEx.
func consoleColumns() int {
	h := windows.Handle(os.Stdout.Fd())
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(h, &info); err != nil {
		return 0
	}
	w := visibleWidth(int(info.Window.Right-info.Window.Left+1), int(info.Size.X), 0)
	if w < 1 {
		return 0
	}
	return w
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
	// Wrap at the real edge. Do not change the font, the window, or the buffer.
	want := mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_WRAP_AT_EOL_OUTPUT
	if err := windows.SetConsoleMode(h, want); err != nil {
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
