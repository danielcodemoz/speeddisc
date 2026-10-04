//go:build windows

package main

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// runToolOutput runs name so that it believes it has a console, without
// attaching it to this process's console.
//
// SFC and DISM draw their percent with console writes. A pipe is fully
// buffered and is not a console, so v0.3.0's parser never saw a percent
// until the process exited (and DISM often wrote none at all). A pseudoconsole
// is a real console from the child's point of view, and it is not the window
// SpeedDisc is using, so the child cannot resize this buffer or this font.
func runToolOutput(name string, args []string, out io.Writer) (int, error) {
	code, started, err := runPseudoConsole(name, args, out)
	if started {
		return code, err
	}
	// Pseudoconsole missing (very old Windows) or it failed before the tool started.
	// Still run the tool off this console. The caller keeps a spinner up.
	return runDetachedPipe(name, args, out)
}

func runDetachedPipe(name string, args []string, out io.Writer) (int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	withoutParentConsole(cmd)
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	return exitCode(err), err
}

func runPseudoConsole(name string, args []string, out io.Writer) (code int, started bool, err error) {
	sa := &windows.SecurityAttributes{
		Length:        uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		InheritHandle: 0,
	}
	var inRead, inWrite windows.Handle
	var outRead, outWrite windows.Handle
	if err = windows.CreatePipe(&inRead, &inWrite, sa, 0); err != nil {
		return 0, false, err
	}
	if err = windows.CreatePipe(&outRead, &outWrite, sa, 64*1024); err != nil {
		windows.CloseHandle(inRead)
		windows.CloseHandle(inWrite)
		return 0, false, err
	}

	cols, rows := int16(80), int16(28)
	if c := consoleColumns(); c > 40 {
		w := c - 1
		if w > 120 {
			w = 120
		}
		cols = int16(w)
	}
	var hPC windows.Handle
	err = windows.CreatePseudoConsole(windows.Coord{X: cols, Y: rows}, inRead, outWrite, 0, &hPC)
	// The pseudoconsole keeps its own copies of these ends.
	windows.CloseHandle(inRead)
	windows.CloseHandle(outWrite)
	if err != nil {
		windows.CloseHandle(inWrite)
		windows.CloseHandle(outRead)
		return 0, false, err
	}
	var pcOnce sync.Once
	closePC := func() { pcOnce.Do(func() { windows.ClosePseudoConsole(hPC) }) }
	defer closePC()

	inFile := os.NewFile(uintptr(inWrite), "pty-in")
	outFile := os.NewFile(uintptr(outRead), "pty-out")
	defer inFile.Close()

	attr, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		outFile.Close()
		return 0, false, err
	}
	defer attr.Delete()
	if err = attr.Update(windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE, unsafe.Pointer(&hPC), unsafe.Sizeof(hPC)); err != nil {
		outFile.Close()
		return 0, false, err
	}

	argv := append([]string{name}, args...)
	cl, err := windows.UTF16FromString(windows.ComposeCommandLine(argv))
	if err != nil {
		outFile.Close()
		return 0, false, err
	}
	app, err := windows.UTF16PtrFromString(name)
	if err != nil {
		outFile.Close()
		return 0, false, err
	}
	si := &windows.StartupInfoEx{
		StartupInfo:             windows.StartupInfo{Cb: uint32(unsafe.Sizeof(windows.StartupInfoEx{}))},
		ProcThreadAttributeList: attr.List(),
	}
	pi := new(windows.ProcessInformation)
	err = windows.CreateProcess(
		app,
		&cl[0],
		nil,
		nil,
		false,
		windows.CREATE_UNICODE_ENVIRONMENT|windows.EXTENDED_STARTUPINFO_PRESENT,
		nil,
		nil,
		&si.StartupInfo,
		pi,
	)
	runtime.KeepAlive(hPC)
	runtime.KeepAlive(cl)
	if err != nil {
		outFile.Close()
		return 0, false, err
	}
	defer windows.CloseHandle(pi.Thread)
	defer windows.CloseHandle(pi.Process)

	var closeOut sync.Once
	closeOutput := func() { closeOut.Do(func() { _ = outFile.Close() }) }
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		_, _ = io.Copy(out, outFile)
	}()

	_, waitErr := windows.WaitForSingleObject(pi.Process, windows.INFINITE)
	// Closing the pseudoconsole makes the output pipe reach EOF.
	closePC()
	select {
	case <-readDone:
	case <-time.After(3 * time.Second):
		closeOutput()
		<-readDone
	}
	closeOutput()

	if waitErr != nil {
		return 0, true, waitErr
	}
	var exit uint32
	if err = windows.GetExitCodeProcess(pi.Process, &exit); err != nil {
		return 0, true, err
	}
	return int(exit), true, nil
}
