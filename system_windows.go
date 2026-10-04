//go:build windows

package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type tailBuf struct {
	mu  sync.Mutex
	b   []byte
	max int
}

func (t *tailBuf) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.b = append(t.b, p...)
	if t.max > 0 && len(t.b) > t.max {
		t.b = append([]byte(nil), t.b[len(t.b)-t.max:]...)
	}
	return len(p), nil
}

func (t *tailBuf) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return string(t.b)
}

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

func runCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runLive(name string, args ...string) (int, string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	tb := &tailBuf{max: 6000}
	cmd.Stdout = io.MultiWriter(os.Stdout, tb)
	cmd.Stderr = io.MultiWriter(os.Stderr, tb)
	err := cmd.Run()
	tail := clip(tb.String(), 240)
	if tail == "" && err != nil {
		tail = clip(err.Error(), 240)
	}
	return exitCode(err), tail
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

func runSFC() (int, string) {
	return runLive(sysExe("sfc.exe"), "/scannow")
}

func runDISM() (int, string) {
	return runLive(sysExe("dism.exe"), "/Online", "/Cleanup-Image", "/RestoreHealth")
}
