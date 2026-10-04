//go:build windows

package main

import (
	"errors"
	"io"
	"os/exec"

	"regexp"
	"strings"

	"golang.org/x/sys/windows"
)

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

// withoutParentConsole keeps helper tools off this console.
// Inheriting it lets PowerShell, SFC, and DISM resize the buffer, which makes Windows shrink the font.
func withoutParentConsole(cmd *exec.Cmd) {
	cmd.SysProcAttr = &windows.SysProcAttr{CreationFlags: windows.DETACHED_PROCESS}
}

func runCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	withoutParentConsole(cmd)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runScan(out io.Writer, label string, max int, paint func(string) string, working string, utf8 bool, name string, args ...string) (int, string, *scanLine) {
	// The live row is drawn by setPercent. The filter must not also call show(),
	// or the spinner and the plain percent fight over the same row.
	live := startLive(out, label, working, "", max, paint, utf8)
	filt := newPercentFilter(nil)
	filt.onPercent = live.setPercent
	code, err := runToolOutput(name, args, filt)
	filt.flush()
	live.stop()
	tail := clip(filt.rawText(), 240)
	if tail == "" && err != nil {
		tail = clip(err.Error(), 240)
	}
	if err != nil && code == 0 {
		code = -1
	}
	return code, tail, live.line
}

var guidRe = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func firstGUID(s string) string {
	return guidRe.FindString(s)
}

const highPerformanceGUID = "8c5e7fda-e8bf-4a96-9a85-a6e23a8c635c"

func powercfg(args ...string) (string, error) {
	return runCapture(sysExe("powercfg.exe"), args...)
}

func activeSchemeHas(guid string) bool {
	out, err := powercfg("/getactivescheme")
	if err != nil && strings.TrimSpace(out) == "" {
		return false
	}
	return strings.Contains(strings.ToLower(out), strings.ToLower(guid))
}

func setHighPerformance() (string, error) {
	out, err := powercfg("/setactive", highPerformanceGUID)
	if err == nil && activeSchemeHas(highPerformanceGUID) {
		return highPerformanceGUID, nil
	}
	out2, err2 := powercfg("/duplicatescheme", highPerformanceGUID)
	guid := firstGUID(out2)
	detail := clip(strings.TrimSpace(out+" "+out2), 220)
	if guid == "" {
		if err2 != nil && detail == "" {
			return "", err2
		}
		if detail == "" && err != nil {
			return "", err
		}
		if detail == "" {
			return "", ErrPowercfg
		}
		return "", errors.New(detail)
	}
	if _, err := powercfg("/setactive", guid); err != nil {
		return "", err
	}
	if !activeSchemeHas(guid) {
		return "", ErrScheme
	}
	return guid, nil
}

func createRestorePoint() error {
	ps := sysExe("WindowsPowerShell/v1.0/powershell.exe")
	out, err := runCapture(ps,
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command",
		"Checkpoint-Computer -Description 'SpeedDisc' -RestorePointType 'MODIFY_SETTINGS'",
	)
	if err != nil {
		detail := clip(strings.TrimSpace(out), 220)
		if detail == "" {
			return err
		}
		return errors.New(detail)
	}
	return nil
}

func runSFC(out io.Writer, max int, paint func(string) string, working string, utf8 bool) (int, string, *scanLine) {
	return runScan(out, "SFC", max, paint, working, utf8, sysExe("sfc.exe"), "/scannow")
}

func runDISM(out io.Writer, max int, paint func(string) string, working string, utf8 bool) (int, string, *scanLine) {
	return runScan(out, "DISM", max, paint, working, utf8, sysExe("dism.exe"), "/Online", "/Cleanup-Image", "/RestoreHealth")
}
