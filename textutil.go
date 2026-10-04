package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func runeLen(s string) int { return utf8.RuneCountInString(s) }

func padRight(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		if n <= 0 {
			return ""
		}
		return string(r[:n])
	}
	return s + strings.Repeat(" ", n-len(r))
}

func wrap(s string, width int) []string {
	if width < 8 {
		width = 8
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	cur := ""
	for _, w := range words {
		if cur == "" {
			cur = w
			continue
		}
		if runeLen(cur)+1+runeLen(w) <= width {
			cur += " " + w
			continue
		}
		lines = append(lines, cur)
		cur = w
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func clip(s string, n int) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.Join(strings.Fields(s), " ")
	r := []rune(s)
	if n <= 0 || len(r) <= n {
		return s
	}
	if n == 1 {
		return string(r[:1])
	}
	return string(r[:n-1]) + "…"
}

// FormatBytes reports a binary size. The unit is KiB/MiB (1024), not a decimal megabyte.
func FormatBytes(n int64, lang Lang) string {
	if n < 0 {
		n = 0
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	f := float64(n)
	i := 0
	for f >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	var num string
	if i == 0 {
		num = strconv.FormatInt(n, 10)
	} else {
		num = fmt.Sprintf("%.1f", f)
		if lang == LangPT {
			num = strings.ReplaceAll(num, ".", ",")
		}
	}
	return num + " " + units[i]
}
