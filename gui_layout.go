package main

type guiGeometry struct {
	side, heroHeight, contentTop int32
	navTop, navPitch             int32
	runTop, optionsTop           int32
	divTop, resultHeading        int32
	resultTop, resultBottom      int32
	exportTop                    int32
}

// All values use client-area logical pixels, not the monitor's physical
// resolution. Painting and child controls must share the same breakpoints.
func layoutGUI(w, h int32) guiGeometry {
	g := guiGeometry{side: 218, heroHeight: 128, contentTop: 214,
		navTop: 170, navPitch: 44, runTop: 165, optionsTop: 210,
		divTop: 246, resultHeading: 252, resultTop: 281}
	if w < 900 {
		g.side = 180
	}
	if h < 680 {
		g.heroHeight, g.contentTop = 76, 158
		g.navTop, g.navPitch = 127, 39
		g.divTop, g.resultHeading, g.resultTop = 246, 252, 275
	}
	if h < 570 {
		g.heroHeight, g.contentTop = 0, 65
		g.optionsTop, g.divTop, g.resultHeading, g.resultTop = 204, 240, 244, 260
		g.navTop, g.navPitch = 127, 39
	}
	g.resultTop += g.contentTop
	g.resultBottom, g.exportTop = h-77, h-70
	return g
}
