package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDisableStartupFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "Demo.lnk")
	if err := os.WriteFile(src, []byte("lnk"), 0o644); err != nil {
		t.Fatal(err)
	}
	dest, err := disableStartupFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("source still present")
	}
	b, err := os.ReadFile(dest)
	if err != nil || string(b) != "lnk" {
		t.Fatalf("dest %s err %v", dest, err)
	}
	if filepath.Base(filepath.Dir(dest)) != "_SpeedDisc_disabled" {
		t.Fatal(dest)
	}
	other := filepath.Join(t.TempDir(), "notstartup", "x.lnk")
	if err := os.MkdirAll(filepath.Dir(other), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(other, []byte("z"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := disableStartupFile(other); err == nil {
		t.Fatal("expected refusal")
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatal("other moved")
	}
}
