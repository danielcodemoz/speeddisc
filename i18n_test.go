package main

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

func formatVerbs(s string) []string {
	var out []string
	for i := 0; i < len(s); i++ {
		if s[i] != '%' {
			continue
		}
		if i+1 < len(s) && s[i+1] == '%' {
			i++
			continue
		}
		j := i + 1
		for j < len(s) && !unicode.IsLetter(rune(s[j])) {
			j++
		}
		if j < len(s) {
			out = append(out, s[j:j+1])
			i = j
		}
	}
	return out
}

func TestI18nComplete(t *testing.T) {
	if len(text) < 40 {
		t.Fatalf("too few strings: %d", len(text))
	}
	banned := []string{"100%", "50%", "99%", "% faster", "% mais rápido", "boost"}
	for key, pair := range text {
		if strings.TrimSpace(pair[LangPT]) == "" || strings.TrimSpace(pair[LangEN]) == "" {
			t.Fatalf("empty %s", key)
		}
		vp := formatVerbs(pair[LangPT])
		ve := formatVerbs(pair[LangEN])
		if strings.Join(vp, ",") != strings.Join(ve, ",") {
			t.Fatalf("verb mismatch %s pt=%v en=%v", key, vp, ve)
		}
		args := make([]any, len(vp))
		for i, v := range vp {
			switch v {
			case "d", "i", "x":
				args[i] = 1
			default:
				args[i] = "x"
			}
		}
		for _, lang := range []Lang{LangPT, LangEN} {
			s := pair[lang]
			out := s
			if len(args) > 0 {
				out = fmt.Sprintf(s, args...)
			}
			if strings.Contains(out, "%!") {
				t.Fatalf("bad format %s %d: %s", key, lang, out)
			}
			low := strings.ToLower(out)
			for _, b := range banned {
				if strings.Contains(low, b) {
					t.Fatalf("banned %q in %s", b, key)
				}
			}
		}
	}
}
