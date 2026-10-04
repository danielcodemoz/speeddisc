package main

import (
	"reflect"
	"testing"
)

func TestParseSelection(t *testing.T) {
	quick := []Action{ActTemp, ActWU, ActThumb, ActPower}
	deep := []Action{ActRestore, ActTemp, ActWU, ActThumb, ActDO, ActWER, ActPower, ActSFC, ActDISM}
	cases := []struct {
		in   string
		want []Action
		bad  bool
	}{
		{"1", quick, false},
		{"2", deep, false},
		{"3", []Action{ActInfo}, false},
		{"4, 8", []Action{ActTemp, ActPower}, false},
		{"1;4", quick, false},
		{"14", []Action{ActInfo}, false},
		{"2,13", append(append([]Action{}, deep...), ActStartup), false},
		{"1,1,4", quick, false},
		{"3,4", []Action{ActInfo, ActTemp}, false},
		{"1 4 8", []Action{ActTemp, ActWU, ActThumb, ActPower}, false},
		{"", nil, true},
		{"0", nil, true},
		{"99", nil, true},
		{"foo", nil, true},
		{"1,foo", nil, true},
	}
	for _, tc := range cases {
		got, err := ParseSelection(tc.in)
		if tc.bad {
			if err == nil {
				t.Fatalf("ParseSelection(%q) expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseSelection(%q): %v", tc.in, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("ParseSelection(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestMenuNumbersUnique(t *testing.T) {
	seen := map[int]bool{}
	for _, it := range menuCatalog() {
		if seen[it.num] {
			t.Fatalf("duplicate menu number %d", it.num)
		}
		seen[it.num] = true
		if it.title == "" || it.desc == "" || len(it.acts) == 0 {
			t.Fatalf("incomplete item %d", it.num)
		}
	}
	if len(seen) != 14 {
		t.Fatalf("got %d items", len(seen))
	}
}

func TestMutates(t *testing.T) {
	if Mutates([]Action{ActInfo}) {
		t.Fatal("info should not mutate")
	}
	if !Mutates([]Action{ActInfo, ActTemp}) {
		t.Fatal("temp mutates")
	}
}

func TestYesQuitLang(t *testing.T) {
	for _, s := range []string{"s", "sim", "y", "YES"} {
		if !isYes(s) {
			t.Fatal(s)
		}
	}
	if isYes("n") || isYes("não") {
		t.Fatal("no")
	}
	if !isQuit("sair") || !isQuit("Q") {
		t.Fatal("quit")
	}
	if !isLangSwitch("idioma") || !isLangSwitch("L") {
		t.Fatal("lang")
	}
}

func TestParseIndexes(t *testing.T) {
	got, err := parseIndexes("1, 3, 3", 5)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int{0, 2}) {
		t.Fatal(got)
	}
	if _, err := parseIndexes("6", 5); err == nil {
		t.Fatal("expected error")
	}
}
