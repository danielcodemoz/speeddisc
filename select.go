package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Action int

const (
	ActInfo Action = iota
	ActRestore
	ActTemp
	ActWU
	ActThumb
	ActDO
	ActWER
	ActPower
	ActSFC
	ActDISM
	ActStartup
)

var actionOrder = []Action{
	ActInfo,
	ActRestore,
	ActTemp,
	ActWU,
	ActThumb,
	ActDO,
	ActWER,
	ActPower,
	ActSFC,
	ActDISM,
	ActStartup,
}

type menuItem struct {
	num     int
	title   string
	desc    string
	acts    []Action
	section string
}

func menuCatalog() []menuItem {
	quick := []Action{ActTemp, ActWU, ActThumb, ActPower}
	deep := []Action{ActTemp, ActWU, ActThumb, ActPower, ActDO, ActWER, ActRestore, ActSFC, ActDISM}
	return []menuItem{
		{1, "opt_quick", "desc_quick", quick, "pkg"},
		{2, "opt_deep", "desc_deep", deep, "pkg"},
		{3, "opt_look", "desc_look", []Action{ActInfo}, "pkg"},
		{4, "opt_temp", "desc_temp", []Action{ActTemp}, "act"},
		{5, "opt_wu", "desc_wu", []Action{ActWU}, "act"},
		{6, "opt_thumb", "desc_thumb", []Action{ActThumb}, "act"},
		{7, "opt_do", "desc_do", []Action{ActDO}, "act"},
		{8, "opt_power", "desc_power", []Action{ActPower}, "act"},
		{9, "opt_wer", "desc_wer", []Action{ActWER}, "act"},
		{10, "opt_sfc", "desc_sfc", []Action{ActSFC}, "act"},
		{11, "opt_dism", "desc_dism", []Action{ActDISM}, "act"},
		{12, "opt_restore", "desc_restore", []Action{ActRestore}, "act"},
		{13, "opt_startup", "desc_startup", []Action{ActStartup}, "act"},
		{14, "opt_info", "desc_info", []Action{ActInfo}, "act"},
	}
}

func tokenize(input string) []string {
	s := strings.TrimSpace(input)
	if s == "" {
		return nil
	}
	if strings.ContainsAny(s, ",;") {
		s = strings.ReplaceAll(s, ";", ",")
		var out []string
		for _, p := range strings.Split(s, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	}
	return strings.Fields(s)
}

func canonical(set map[Action]bool) []Action {
	var out []Action
	for _, a := range actionOrder {
		if set[a] {
			out = append(out, a)
		}
	}
	return out
}

// ParseSelection accepts menu numbers separated by commas, semicolons, or spaces.
// Packages expand to their actions. The result is de-duplicated in execution order.
func ParseSelection(input string) ([]Action, error) {
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty")
	}
	byNum := map[int]menuItem{}
	for _, it := range menuCatalog() {
		byNum[it.num] = it
	}
	set := map[Action]bool{}
	for _, tok := range tokens {
		n, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("bad token")
		}
		it, ok := byNum[n]
		if !ok {
			return nil, fmt.Errorf("unknown option")
		}
		for _, a := range it.acts {
			set[a] = true
		}
	}
	return canonical(set), nil
}

func Contains(acts []Action, want Action) bool {
	for _, a := range acts {
		if a == want {
			return true
		}
	}
	return false
}

func Mutates(acts []Action) bool {
	for _, a := range acts {
		if a != ActInfo {
			return true
		}
	}
	return false
}

func parseIndexes(input string, n int) ([]int, error) {
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty")
	}
	seen := map[int]bool{}
	var out []int
	for _, tok := range tokens {
		v, err := strconv.Atoi(tok)
		if err != nil || v < 1 || v > n {
			return nil, fmt.Errorf("bad index")
		}
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v-1)
	}
	return out, nil
}

func isQuit(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "q", "quit", "sair", "exit":
		return true
	default:
		return false
	}
}

func isLangSwitch(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "l", "lang", "language", "idioma", "língua", "lingua":
		return true
	default:
		return false
	}
}

func isYes(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "s", "sim", "y", "yes":
		return true
	default:
		return false
	}
}
