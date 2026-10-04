package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPercentProgressWholeStepsOnly(t *testing.T) {
	var b bytes.Buffer
	line := newScanLine(&b, "SFC", 80, nil)
	f := newPercentFilter(line)
	in := strings.Join([]string{
		"Beginning system scan. This process will take some time.",
		"Verification 26.7% complete.",
		"Verification 26.9% complete.",
		"Verification 27.0% complete.",
		"Verification 27.1% complete.",
		"Verification 27.9% complete.",
		"Verification 28% complete.",
		"Verification 28.2% complete.",
		"Verificação 28,6% concluída.",
	}, "\n") + "\n"
	if _, err := f.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	f.flush()
	got := b.String()
	if strings.Contains(got, "\n") {
		t.Fatalf("progress introduced a newline: %q", got)
	}
	if strings.Contains(got, "Beginning") || strings.Contains(got, "26.7") || strings.Contains(got, "27.1") || strings.Contains(got, "27.9") {
		t.Fatalf("decimal or status leaked: %q", got)
	}
	for _, want := range []string{"SFC 26%", "SFC 27%", "SFC 28%"} {
		if strings.Count(got, want) != 1 {
			t.Fatalf("count %s in %q", want, got)
		}
	}
	if !strings.Contains(got, "\r") {
		t.Fatalf("expected carriage return: %q", got)
	}
	line.finish("  SFC concluído")
	final := b.String()
	if strings.Count(final, "\n") != 1 {
		t.Fatalf("want one final line: %q", final)
	}
	if !strings.HasSuffix(final, "SFC concluído\n") {
		t.Fatalf("final line: %q", final)
	}
	if strings.Contains(final, "26.9") {
		t.Fatal(final)
	}
}

func TestDISMProgressSameRow(t *testing.T) {
	var b bytes.Buffer
	line := newScanLine(&b, "DISM", 80, nil)
	f := newPercentFilter(line)
	chunks := []string{
		"[==                         ] 42.3%\r",
		"[==                         ] 42.8%\r",
		"[===                        ] 43,0%\r",
		"[===                        ] 43.4%\n",
	}
	for _, c := range chunks {
		if _, err := f.Write([]byte(c)); err != nil {
			t.Fatal(err)
		}
	}
	f.flush()
	got := b.String()
	if strings.Contains(got, "\n") {
		t.Fatalf("newline in dism progress: %q", got)
	}
	if strings.Count(got, "DISM 42%") != 1 || strings.Count(got, "DISM 43%") != 1 {
		t.Fatalf("%q", got)
	}
	if strings.Contains(got, "42.3") || strings.Contains(got, "[==") {
		t.Fatalf("raw dism leaked: %q", got)
	}
	line.finish("  DISM finished")
	if !strings.HasSuffix(b.String(), "DISM finished\n") {
		t.Fatal(b.String())
	}
}

func TestPercentCarriageReturnDoesNotGrow(t *testing.T) {
	var b bytes.Buffer
	line := newScanLine(&b, "SFC", 40, nil)
	f := newPercentFilter(line)
	_, _ = f.Write([]byte("9.9%\r100%\r"))
	f.flush()
	line.finish("  SFC finished")
	if strings.Count(b.String(), "\n") != 1 {
		t.Fatal(b.String())
	}
	// The shorter "SFC 9%" row must be padded so "100" does not leave leftovers,
	// and the final line is still a single row.
	if !strings.Contains(b.String(), "SFC 9%") || !strings.Contains(b.String(), "SFC 100%") {
		t.Fatal(b.String())
	}
}

func TestPercentANSIAndNoNewlineYet(t *testing.T) {
	var b bytes.Buffer
	line := newScanLine(&b, "DISM", 80, nil)
	f := newPercentFilter(line)
	if _, err := f.Write([]byte("\x1b[2K\x1b[1G[==] 41.2%")); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("\x1b[1G[==] 42.0%\r")); err != nil {
		t.Fatal(err)
	}
	f.flush()
	got := b.String()
	if strings.Contains(got, "\n") || strings.Contains(got, "\x1b") || strings.Contains(got, "[==") {
		t.Fatalf("%q", got)
	}
	if strings.Count(got, "DISM 41%") != 1 || strings.Count(got, "DISM 42%") != 1 {
		t.Fatalf("%q", got)
	}
}

func TestLiveRowStaysOneLine(t *testing.T) {
	var b bytes.Buffer
	live := startLive(&b, "SFC", "a trabalhar", "%d ficheiros", 80, nil, true)
	live.setPercent(27)
	live.setPercent(27)
	live.setPercent(28)
	live.stop()
	live.line.finish("  SFC concluído")
	got := b.String()
	if strings.Count(got, "\n") != 1 {
		t.Fatalf("want one final line: %q", got)
	}
	if !strings.Contains(got, "27%") || !strings.Contains(got, "28%") {
		t.Fatal(got)
	}
	if strings.Contains(got, "27.1") || strings.Contains(got, "100%") {
		t.Fatal(got)
	}
	last := strings.TrimRight(got, " \t\r\n")
	if !strings.HasSuffix(last, "SFC concluído") {
		t.Fatal(got)
	}
}

func TestLiveCountNotAPercent(t *testing.T) {
	var b bytes.Buffer
	live := startLive(&b, "Temp", "a trabalhar", "%d ficheiros", 80, nil, false)
	for i := 1; i <= 20; i++ {
		live.setCount(i)
	}
	live.stopClear()
	got := b.String()
	if strings.Contains(got, "%") {
		t.Fatalf("count invented a percent: %q", got)
	}
	if strings.Contains(got, "\n") {
		t.Fatalf("count opened new lines: %q", got)
	}
	if !strings.Contains(got, "ficheiros") {
		t.Fatal(got)
	}
}
