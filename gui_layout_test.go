package main

import "testing"

func TestGUIKeepsResultAndExportVisible(t *testing.T) {
	for _, size := range [][2]int32{{1200, 800}, {962, 683}, {1040, 535}, {1050, 565}, {1050, 570}, {800, 460}} {
		w, h := size[0], size[1]
		g := layoutGUI(w, h)
		if g.resultBottom-g.resultTop < 55 || g.resultBottom >= g.exportTop || g.exportTop+34 > h {
			t.Errorf("%dx%d: result [%d,%d), exports [%d,%d)", w, h, g.resultTop, g.resultBottom, g.exportTop, g.exportTop+34)
		}
		if g.contentTop+g.optionsTop+32 >= g.resultTop || g.navTop+5*g.navPitch+37 >= h-74 {
			t.Errorf("%dx%d: controls overlap the result or session footer", w, h)
		}
	}
}
