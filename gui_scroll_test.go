package main

import "testing"

func TestScrollGeometry(t *testing.T) {
	start, width, last := scrollThumb(100, 25, 0, 19, 5, 10)
	if start != 50 || width != 25 || last != 15 {
		t.Fatalf("unexpected geometry: %d %d %d", start, width, last)
	}
	if got := scrollTarget(start, 100-width, 0, last); got != 10 {
		t.Fatalf("thumb round trip = %d", got)
	}
	if scrollTarget(-8, 75, 3, 18) != 3 || scrollTarget(85, 75, 3, 18) != 18 {
		t.Fatal("drag was not clamped to native range")
	}
	start, width, last = scrollThumb(12, 25, 0, 9, 3, 5)
	if start != 0 || width != 12 || last != 7 {
		t.Fatalf("tiny thumb geometry: %d %d %d", start, width, last)
	}
	start, width, last = scrollThumb(99, 25, 2, 9, 8, 8)
	if start != 0 || width != 99 || last != 2 {
		t.Fatalf("all items fit: %d %d %d", start, width, last)
	}
}
