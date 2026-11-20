package main

import (
	"fmt"
	"io"
	"strings"
)

// lineCounter records the physical row budget of the already wrapped menu.
// The existing UI limits output to visible columns; two more rows leave
// room for the separation and input prompt.
type lineCounter int

func (c *lineCounter) Write(p []byte) (int, error) {
	*c += lineCounter(strings.Count(string(p), "\n"))
	return len(p), nil
}

func (u *UI) fullMenuRows(lang Lang) int {
	var lines lineCounter
	dry := *u
	dry.out = &lines
	dry.cols = u.termWidth()
	dry.Banner(lang)
	dry.Menu(lang)
	return int(lines) + 2
}

func (u *UI) DrawMenu(lang Lang, rows int, forceFull bool) {
	if forceFull || rows <= 0 || u.fullMenuRows(lang) <= rows {
		u.Banner(lang)
		u.Menu(lang)
		return
	}
	u.compactMenu(lang, rows)
}

func (u *UI) compactMenu(lang Lang, rows int) {
	if rows >= 14 && u.termWidth() >= 40 {
		u.Banner(lang)
		for _, item := range menuCatalog()[:4] {
			u.rawWrap(fmt.Sprintf("  %d  %s", item.num, T(lang, item.title)))
		}
		if lang == LangEN {
			u.rawWrap("  A  All actions  |  L  Language  |  Q  Quit")
			if rows >= 17 {
				u.info("Enter an option number. System changes require confirmation.")
			}
		} else {
			u.rawWrap("  A  Todas as acções  |  L  Idioma  |  Q  Sair")
			if rows >= 17 {
				u.info("Indique o número da opção. Alterações ao sistema requerem confirmação.")
			}
		}
		return
	}
	// Very short/narrow view keeps the identity and path to the full
	// menu visible rather than attempting to draw every description.
	fmt.Fprintln(u.out, u.paint(ansiBoldCyan, "  SpeedDisc"))
	if lang == LangEN {
		u.rawWrap("  1 Quick / 2 Deep / 3 Check / 4 Full")
		u.rawWrap("  A All / L Lang / Q Quit")
	} else {
		u.rawWrap("  1 Rápido / 2 Fundo / 3 Ver / 4 Total")
		u.rawWrap("  A Tudo / L Idioma / Q Sair")
	}
}

func isActionMenu(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "a", "all", "acções", "ações", "actions":
		return true
	default:
		return false
	}
}

var _ io.Writer = (*lineCounter)(nil)
