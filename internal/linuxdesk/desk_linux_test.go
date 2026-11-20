//go:build linux

package linuxdesk

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixture(t *testing.T, base, p, value string) {
	t.Helper()
	path := filepath.Join(base, p)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(value), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestCatalogueAndFixtureDiagnostics(t *testing.T) {
	dir := t.TempDir()
	fixture(t, dir, "meminfo", "MemTotal: 4096 kB\nMemAvailable: 1024 kB\nSwapTotal: 2048 kB\nSwapFree: 1024 kB\n")
	fixture(t, dir, "net/dev", "iface headings\n     lo: 1024 12 1 0 0 0 0 0 2048 13 2 0 0 0 0 0\n")
	fixture(t, dir, "777/statm", "200 100 20 0 0 0 0\n")
	fixture(t, dir, "777/comm", "fixture-process\n")
	fixture(t, dir, "888/statm", "200 20 2 0 0 0\n")
	fixture(t, dir, "888/comm", "fixture-other\n")
	fs := environment{dir, dir, dir, dir, os.Getuid(), 4096}
	if len(Tasks("en")) != 6 || Tasks("en")[0].Title == Tasks("pt")[0].Title {
		t.Fatal("catalogue not complete or translated")
	}
	mem := run(context.Background(), "memory", "en", fs)
	if mem.Status != "Completed" || len(mem.Rows) != 5 || mem.Rows[0].Value != "4.0 MiB" || mem.Rows[2].Value != "3.0 MiB" {
		t.Fatalf("memory metrics %#v", mem)
	}
	net := run(context.Background(), "network", "en", fs)
	if len(net.Rows) != 2 || net.Rows[0].Label != "lo" || !strings.Contains(net.Rows[0].Value, "1/2") {
		t.Fatalf("network metrics %#v", net)
	}
	ps := run(context.Background(), "processes", "en", fs)
	if len(ps.Rows) != 3 || ps.Rows[0].Label != "PID 777" || !strings.Contains(ps.Rows[0].Value, "400.0 KiB") {
		t.Fatalf("process metrics %#v", ps)
	}
	disk := run(context.Background(), "disk", "pt", fs)
	if disk.Status != "Concluído" || len(disk.Rows) < 3 {
		t.Fatalf("disk metrics %#v", disk)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := run(ctx, "disk", "en", fs); got.Status != "Unavailable" {
		t.Fatalf("cancelled result %#v", got)
	}
	if got := run(context.Background(), "who knows", "pt", fs); got.Status != "Indisponível" {
		t.Fatalf("unknown operation %#v", got)
	}
}

func TestScanIsReadOnlyBoundedAndNeverFollowsSymlinks(t *testing.T) {
	root := t.TempDir()
	out := t.TempDir()
	fixture(t, root, "nested/large", "fixture")
	fixture(t, out, "secret", "secret")
	for path, size := range map[string]int64{filepath.Join(root, "nested", "large"): 75 * 1024 * 1024, filepath.Join(out, "secret"): 80 * 1024 * 1024} {
		if err := os.Truncate(path, size); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(out, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(root, filepath.Join(out, "alias")); err != nil {
		t.Fatal(err)
	}
	limits := scanLimits{8, 100, time.Second}
	got, err := scan(context.Background(), root, os.Getuid(), false, limits)
	if err != nil || len(got.top) != 1 || got.top[0].name != "nested/large" || got.skipped != 1 {
		t.Fatalf("scan result %#v, %v", got, err)
	}
	if _, err := scan(context.Background(), filepath.Join(out, "alias"), os.Getuid(), false, limits); err == nil {
		t.Fatal("followed symlink at the scan root")
	}
	got, err = scan(context.Background(), root, os.Getuid(), true, limits)
	if err != nil || got.count != 1 || got.bytes != 75*1024*1024 {
		t.Fatalf("temp result %#v, %v", got, err)
	}
	got, err = scan(context.Background(), root, os.Getuid(), false, scanLimits{8, 1, time.Second})
	if err != nil || !got.limited {
		t.Fatalf("scan should return partial when budget is spent: %#v, %v", got, err)
	}
	if stat, err := os.Stat(filepath.Join(out, "secret")); err != nil || stat.Size() != 80*1024*1024 {
		t.Fatal("outside fixture changed")
	}
	if stat, err := os.Stat(filepath.Join(root, "nested", "large")); err != nil || stat.Size() != 75*1024*1024 {
		t.Fatal("inside fixture changed")
	}
}
