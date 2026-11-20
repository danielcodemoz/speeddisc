//go:build windows

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// Query the token, never adjust it. An elevated CI pass is not standard-user proof.
func TestWindowsIdentity(t *testing.T) {
	token, err := syscall.OpenCurrentProcessToken()
	if err != nil {
		t.Fatal(err)
	}
	defer token.Close()
	var elevated, returned uint32
	err = syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&elevated)), 4, &returned)
	if err != nil {
		t.Fatal(err)
	}
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Windows test identity=%s SID=%s elevated=%t", u.Username, u.Uid, elevated != 0)
	if elevated != 0 {
		t.Log("STANDARD_USER_ACCEPTANCE=NOT_VERIFIED (elevated token)")
	} else {
		t.Log("Non-elevated token observed; this alone does not prove a non-admin account or UAC UI acceptance")
	}
}

func TestWindowsRejectNonLocalPaths(t *testing.T) {
	for _, p := range []string{`\\localhost\share`, `\\?\C:\Windows`, `\\.\C:`, `\\example.invalid\share`} {
		if err := localRoot(p); err == nil {
			t.Fatalf("accepted excluded path %q", p)
		}
		if err := saveReport(filepath.Join(p, "report.json"), Report{}); err == nil {
			t.Fatalf("accepted excluded report %q", p)
		}
	}
}

func TestWindowsJunctionExclusion(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "not-visited.txt"), []byte("synthetic"), 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "junction")
	// Fixed fixture-only local junction. No registry, ACL, account or policy changes.
	cmd := exec.Command("cmd.exe", "/c", "mklink", "/J", link, outside)
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("junction fixture could not be created: %v", err)
	}
	d, err := scan(context.Background(), root, 100)
	if err != nil {
		t.Fatal(err)
	}
	if d.Files != 0 || d.Bytes != 0 || d.Skipped != 1 {
		t.Fatalf("junction traversed: %+v", d)
	}
	if err := localRoot(link); err == nil {
		t.Fatal("junction root accepted")
	}
	if err := saveReport(filepath.Join(link, "unexpected.json"), Report{}); err == nil {
		t.Fatal("report through junction accepted")
	}
	if _, err := os.Stat(filepath.Join(outside, "unexpected.json")); !os.IsNotExist(err) {
		t.Fatal("outside report exists")
	}
}

func TestWindowsPortableProcess(t *testing.T) {
	executable := os.Getenv("SPEEDDISC_SMOKE_EXE")
	if executable == "" {
		t.Skip("set SPEEDDISC_SMOKE_EXE to the binary built from this checkout")
	}
	root := t.TempDir()
	fixture := filepath.Join(root, "fixture")
	if err := os.Mkdir(fixture, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "synthetic.txt"), []byte("12345"), 0600); err != nil {
		t.Fatal(err)
	}
	disk := filepath.Join(root, "disk.json")
	internet := filepath.Join(root, "internet.json")
	boot := filepath.Join(root, "startup.json")
	input := "L\n1\n" + fixture + "\n" + disk + "\n2\n" + internet + "\n3\n" + boot + "\nL\nq\n"
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable)
	cmd.Stdin = strings.NewReader(input)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("portable process failed (%v); output intentionally not logged", err)
	}
	for _, pair := range []struct{ path, flow string }{{disk, "disk"}, {internet, "internet"}, {boot, "startup"}} {
		b, err := os.ReadFile(pair.path)
		if err != nil {
			t.Fatalf("%s report missing: %v", pair.flow, err)
		}
		var r Report
		if err = json.Unmarshal(b, &r); err != nil {
			t.Fatal(err)
		}
		if r.Flow != pair.flow || r.Error != "" {
			t.Fatalf("%s collector failed: %s", pair.flow, r.Error)
		}
	}
	for _, want := range []string{"Largest visited files", "Next step:", "SPEEDDISC", "Disco cheio"} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("missing UI marker %q", want)
		}
	}
	b, err := os.ReadFile(filepath.Join(fixture, "synthetic.txt"))
	if err != nil || string(b) != "12345" {
		t.Fatal("fixture changed")
	}
	t.Log("Portable process: PT/EN, all three flows, report creation and source preservation passed. No raw inventory logged.")
}
