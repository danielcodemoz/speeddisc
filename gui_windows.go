//go:build windows && speeddiscgui

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/danielcodemoz/speeddisc/internal/report"
	"golang.org/x/sys/windows"
)

// The native window owns all controls and results. Workers notify it with a
// message; they never touch controls, and no browser or network is involved.
var (
	dUser   = windows.NewLazySystemDLL("user32.dll")
	dGDI    = windows.NewLazySystemDLL("gdi32.dll")
	dWin    = windows.NewLazySystemDLL("kernel32.dll")
	dShell  = windows.NewLazySystemDLL("shell32.dll")
	dCommon = windows.NewLazySystemDLL("comdlg32.dll")
	dDWM    = windows.NewLazySystemDLL("dwmapi.dll")
	dGrad   = windows.NewLazySystemDLL("msimg32.dll")
	dCreate = dUser.NewProc("CreateWindowExW")
	dSend   = dUser.NewProc("SendMessageW")
	dSet    = dUser.NewProc("SetWindowTextW")
	dGet    = dUser.NewProc("GetWindowTextW")
	dSize   = dUser.NewProc("GetClientRect")
	dMove   = dUser.NewProc("MoveWindow")
	dBrush  = dGDI.NewProc("CreateSolidBrush")
	dSelect = dGDI.NewProc("SelectObject")
	dDelete = dGDI.NewProc("DeleteObject")
	dDraw   = dUser.NewProc("DrawTextW")
	dFill   = dUser.NewProc("FillRect")
	dProc   = syscall.NewCallback(deskProc)
	desk    *deskState
)

const (
	idSearch   = 101
	idList     = 102
	idRun      = 103
	idTxt      = 104
	idPDF      = 105
	idLang     = 106
	idTheme    = 107
	idOptions  = 108
	idOutput   = 109
	idProgress = 110
	idMarquee  = 111
	idNav      = 120
	wmDone     = 0x8001
	wmStartup  = 0x8002
	wmProgress = 0x8003
)

type dRect struct{ L, T, R, B int32 }
type dPoint struct{ X, Y int32 }
type dMsg struct {
	HWND    uintptr
	Msg     uint32
	W, L    uintptr
	Time    uint32
	Pt      dPoint
	Private uint32
}
type dClass struct {
	Size, Style                   uint32
	WndProc                       uintptr
	ClassExtra, WinExtra          int32
	Instance, Icon, Cursor, Brush uintptr
	Menu, Name                    *uint16
	IconSmall                     uintptr
}
type dDrawItem struct {
	CtlType, CtlID, ItemID, Action, State uint32
	HWND, DC                              uintptr
	Rect                                  dRect
	Data                                  uintptr
}
type dVertex struct {
	X, Y       int32
	R, G, B, A uint16
}
type dGradient struct{ First, Second uint32 }
type dMinMax struct{ Reserved, MaxSize, MaxPosition, MinTrack, MaxTrack dPoint }

type deskPalette struct{ bg, panel, raised, line, ink, muted, accent, onAccent, h1, h2, heroInk, heroMuted, ok, warn uintptr }
type deskState struct {
	scroll                              [2]scrollLayout
	scrollPosted                        bool
	win, instance, icon                 uintptr
	h                                   map[int]uintptr
	dpi                                 int32
	fSmall, fBody, fBold, fHead, fHero  uintptr
	brush                               uintptr
	ui                                  deskPalette
	light, admin, busy, preparing, demo bool
	demoRunning                         bool
	progress                            chan guiProgress
	stage                               guiProgress
	lang                                Lang
	mode                                string
	rows                                []guiItem
	item                                guiItem
	runResult                           *report.Result
	work                                chan report.Result
	startup                             chan startupChoice
	startupItems                        []guiStartupItem
	startupNotes                        []string
	resetWinsock, picking               bool
	picked                              int
	pendingID                           string
}
type startupChoice struct {
	items []guiStartupItem
	notes []string
}

var deskCats = []string{"all", "packages", "diagnostics", "cleanup", "network", "advanced"}
var deskCatNames = [][2]string{{"Todas as ferramentas", "All tools"}, {"Pacotes", "Packages"}, {"Diagnóstico", "Diagnostics"}, {"Limpeza", "Cleanup"}, {"Rede", "Network"}, {"Avançado", "Advanced"}}
var deskCatMarks = []string{"01", "02", "03", "04", "05", "06"}

func rgb(r, g, b byte) uintptr { return uintptr(r) | uintptr(g)<<8 | uintptr(b)<<16 }
func rgbHex(x uint32) uintptr  { return rgb(byte(x>>16), byte(x>>8), byte(x)) }
func deskColors(light bool) deskPalette {
	if light {
		return deskPalette{
			bg: rgbHex(0xF2F5F6), panel: rgbHex(0xFFFFFF), raised: rgbHex(0xE6EEEE), line: rgbHex(0xA8BCBB),
			ink: rgbHex(0x16272B), muted: rgbHex(0x51676A), accent: rgbHex(0x176653), onAccent: rgbHex(0xFFFFFF),
			h1: rgbHex(0x15343C), h2: rgbHex(0x215A47), heroInk: rgbHex(0xFFFFFF), heroMuted: rgbHex(0xD7F2EA),
			ok: rgbHex(0x126552), warn: rgbHex(0x8B3F1B),
		}
	}
	return deskPalette{
		bg: rgbHex(0x10181D), panel: rgbHex(0x1B252A), raised: rgbHex(0x26343A), line: rgbHex(0x4D6265),
		ink: rgbHex(0xECF4F3), muted: rgbHex(0xAFC5C6), accent: rgbHex(0x85E7BF), onAccent: rgbHex(0x112820),
		h1: rgbHex(0x173640), h2: rgbHex(0x215744), heroInk: rgbHex(0xFFFFFF), heroMuted: rgbHex(0xD7F2EA),
		ok: rgbHex(0xA7E6CA), warn: rgbHex(0xFFCD99),
	}
}
func u16(x string) *uint16 {
	v, _ := windows.UTF16PtrFromString(strings.ReplaceAll(x, "\x00", ""))
	return v
}
func ptr(p *uint16) uintptr     { return uintptr(unsafe.Pointer(p)) }
func dText(h uintptr, x string) { p := u16(x); dSet.Call(h, ptr(p)); runtime.KeepAlive(p) }
func dRead(h uintptr) string {
	var b [512]uint16
	dGet.Call(h, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	return windows.UTF16ToString(b[:])
}
func dMessage(h uintptr, code uint32, w, l uintptr) uintptr {
	v, _, _ := dSend.Call(h, uintptr(code), w, l)
	return v
}
func (s *deskState) l(pt, en string) string {
	if s.lang == LangPT {
		return pt
	}
	return en
}
func (s *deskState) px(v int32) int32 { return (v*s.dpi + 48) / 96 }
func (s *deskState) dp(v int32) int32 { return (v*96 + s.dpi/2) / s.dpi }
func (s *deskState) box(x, y, w, h int32) dRect {
	return dRect{s.px(x), s.px(y), s.px(x + w), s.px(y + h)}
}
func (s *deskState) lay(id int, x, y, w, h int32) {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	dMove.Call(s.h[id], uintptr(s.px(x)), uintptr(s.px(y)), uintptr(s.px(w)), uintptr(s.px(h)), 1)
}
func (s *deskState) invalid() { dUser.NewProc("InvalidateRect").Call(s.win, 0, 0) }
func (s *deskState) hShow(id int, show bool) {
	n := uintptr(0)
	if show {
		n = 5
	}
	dUser.NewProc("ShowWindow").Call(s.h[id], n)
}
func (s *deskState) enable(id int, on bool) {
	n := uintptr(0)
	if on {
		n = 1
	}
	dUser.NewProc("EnableWindow").Call(s.h[id], n)
}
func dDefault(hwnd uintptr, msg uint32, w, l uintptr) uintptr {
	v, _, _ := dUser.NewProc("DefWindowProcW").Call(hwnd, uintptr(msg), w, l)
	return v
}
func dFont(h int32, weight uintptr, face string) uintptr {
	f := u16(face)
	v, _, _ := dGDI.NewProc("CreateFontW").Call(uintptr(-h), 0, 0, 0, weight, 0, 0, 0, 1, 0, 0, 5, 0, ptr(f))
	runtime.KeepAlive(f)
	return v
}
func (s *deskState) fonts() {
	old := []uintptr{s.fSmall, s.fBody, s.fBold, s.fHead, s.fHero}
	s.fSmall = dFont(s.px(13), 500, "Segoe UI")
	s.fBody = dFont(s.px(16), 450, "Segoe UI")
	s.fBold = dFont(s.px(15), 650, "Segoe UI")
	s.fHead = dFont(s.px(25), 650, "Segoe UI")
	s.fHero = dFont(s.px(31), 700, "Segoe UI")
	for id, h := range s.h {
		f := s.fBody
		if id == idOutput {
			f = s.fSmall
		}
		dMessage(h, 0x30, f, 1)
	}
	dMessage(s.h[idList], 0x01A0, 0, uintptr(s.px(34)))
	for _, f := range old {
		if f != 0 {
			dDelete.Call(f)
		}
	}
}

func (s *deskState) add(id int, cls string, style uint32) {
	c, t := u16(cls), u16("")
	base := uint32(0x50010000)
	if id == idProgress || id == idMarquee {
		base = 0x40000000
	}
	h, _, _ := dCreate.Call(0, ptr(c), ptr(t), uintptr(style|base), 0, 0, 0, 0, s.win, uintptr(id), s.instance, 0)
	runtime.KeepAlive(c)
	runtime.KeepAlive(t)
	s.h[id] = h
}
func (s *deskState) controls() {
	for i := range deskCats {
		s.add(idNav+i, "BUTTON", 0x0000000B)
	}
	for _, id := range []int{idRun, idTxt, idPDF, idTheme, idLang, idOptions} {
		s.add(id, "BUTTON", 0x0000000B)
	}
	s.add(idSearch, "EDIT", 0x0080)
	s.add(idList, "LISTBOX", 0x04200000|0x0001|0x0100|0x0010|0x0040)
	s.add(idOutput, "EDIT", 0x04200000|0x0004|0x0040|0x0800)
	ic := [2]uint32{8, 0x20}
	windows.NewLazySystemDLL("comctl32.dll").NewProc("InitCommonControlsEx").Call(uintptr(unsafe.Pointer(&ic)))
	s.add(idProgress, "msctls_progress32", 0x01)
	s.add(idMarquee, "msctls_progress32", 0x08)
	s.initScrollRails()
	dMessage(s.h[idProgress], 0x406, 0, 100)
	s.fonts()
	s.updateTheme()
	s.captions()
	s.filter("3")
	dText(s.h[idOutput], s.l("O resultado aparecerá aqui após a tarefa.", "The result will appear here after the task."))
}
func (s *deskState) captions() {
	for i, n := range deskCatNames {
		dText(s.h[idNav+i], s.l(n[0], n[1]))
	}
	dText(s.h[idLang], s.l("EN", "PT"))
	dText(s.h[idTheme], s.l("Tema", "Theme"))
	dText(s.h[idTxt], s.l("Exportar TXT", "Export TXT"))
	dText(s.h[idPDF], s.l("Exportar PDF", "Export PDF"))

}
func squashSearch(x string) string {
	return strings.NewReplacer("á", "a", "à", "a", "ã", "a", "â", "a", "é", "e", "ê", "e", "í", "i", "ó", "o", "ô", "o", "õ", "o", "ú", "u", "ç", "c").Replace(strings.ToLower(x))
}
func (s *deskState) filter(prefer string) {
	s.rows = nil
	q := squashSearch(strings.TrimSpace(dRead(s.h[idSearch])))
	h := s.h[idList]
	dMessage(h, 0x0184, 0, 0)
	sel := 0
	for _, a := range guiItems(s.lang) {
		if s.mode != "all" && s.mode != a.Category {
			continue
		}
		if q != "" && !strings.Contains(squashSearch(a.Title+" "+a.Detail+" "+a.ID), q) {
			continue
		}
		if a.ID == prefer {
			sel = len(s.rows)
		}
		s.rows = append(s.rows, a)
		p := u16(fmt.Sprintf("%02s  ·  %s", a.ID, a.Title))
		dMessage(h, 0x0180, 0, ptr(p))
		runtime.KeepAlive(p)
	}
	if len(s.rows) > 0 {
		dMessage(h, 0x0186, uintptr(sel), 0)
		if sel <= 3 {
			dMessage(h, 0x0197, 0, 0)
		}
		s.selectItem(sel)
	} else {
		s.selectItem(-1)
	}
}
func (s *deskState) selectItem(i int) {
	prev := s.item.ID
	s.item = guiItem{}
	if i >= 0 && i < len(s.rows) {
		s.item = s.rows[i]
	}
	if prev != s.item.ID {
		s.resetWinsock = false
		s.startupItems = nil
		s.startupNotes = nil
	}
	s.sync()
}
func (s *deskState) updateTheme() {
	old := s.brush
	s.ui = deskColors(s.light)
	s.brush, _, _ = dBrush.Call(s.ui.panel)
	if old != 0 {
		dDelete.Call(old)
	}
	v := int32(1)
	if s.light {
		v = 0
	}
	dDWM.NewProc("DwmSetWindowAttribute").Call(s.win, 20, uintptr(unsafe.Pointer(&v)), 4)
	s.invalid()
	s.scheduleScrollSync()
}
func (s *deskState) sync() {
	block := s.busy || s.preparing
	dText(s.h[idRun], s.l("Executar ferramenta", "Run tool"))
	switch {
	case s.busy:
		dText(s.h[idRun], s.l("Em execução…", "Running…"))
	case s.preparing:
		dText(s.h[idRun], s.l("A ler o arranque…", "Reading startup…"))
	case !s.admin && (s.item.NeedsAdmin || s.resetWinsock):
		dText(s.h[idRun], s.l("Abrir como administrador", "Open as administrator"))
	case s.item.ID == "14":
		dText(s.h[idRun], s.l("Analisar arranque", "Analyse startup"))
	case !s.item.Mutates:
		dText(s.h[idRun], s.l("Iniciar análise", "Start analysis"))
	}
	for _, id := range []int{idLang, idSearch, idList, idOptions} {
		s.enable(id, !block)
	}
	for i := range deskCats {
		s.enable(idNav+i, !block)
	}
	s.enable(idSearch, !block && !s.picking)
	s.hShow(idSearch, !s.picking)
	canRun := s.item.ID != "" && !block
	if s.picking {
		canRun = !block && s.picked >= 0 && s.picked < len(s.startupItems) && s.startupItems[s.picked].CanDisable
		if s.admin {
			dText(s.h[idRun], s.l("Desactivar entrada", "Disable entry"))
		} else {
			dText(s.h[idRun], s.l("Abrir como administrador", "Open as administrator"))
		}
	}
	s.enable(idRun, canRun && !s.demo)
	if s.demo {
		dText(s.h[idRun], s.l("Modo de demonstração", "Demo mode"))
	}
	s.enable(idTxt, s.runResult != nil && !block)
	s.enable(idPDF, s.runResult != nil && !block)
	showOptions := s.item.ID == "14" || s.item.ID == "19"
	if s.item.ID == "14" {
		dText(s.h[idOptions], s.l("Gerir entradas…", "Manage entries…"))
	}
	if s.item.ID == "19" {
		x := "  ○"
		if s.resetWinsock {
			x = "  ●"
		}
		dText(s.h[idOptions], s.l("Incluir Winsock", "Include Winsock")+x)
	}
	if s.picking {
		dText(s.h[idOptions], s.l("Voltar às ferramentas", "Back to tools"))
	}
	s.hShow(idOptions, showOptions && !s.demo && !s.busy)
	s.hShow(idRun, !s.busy)
	s.refreshProgress()
	s.invalid()
}
func (s *deskState) layout() {
	var r dRect
	dSize.Call(s.win, uintptr(unsafe.Pointer(&r)))
	w, h := s.dp(r.R), s.dp(r.B)
	g := layoutGUI(w, h)
	l := g.side
	x := l + 22
	y := g.contentTop
	body := w - x - 22
	left := (body * 42) / 100
	right := x + left + 12
	rw := body - left - 12
	s.lay(idLang, w-143, 19, 50, 34)
	s.lay(idTheme, w-84, 19, 62, 34)
	for i := range deskCats {
		s.lay(idNav+i, 12, g.navTop+int32(i)*g.navPitch, l-24, 37)
	}
	s.lay(idSearch, x+13, y+63, left-26, 24)
	s.lay(idList, x+9, y+114, left-18, h-y-163)
	by := y + g.runTop
	s.lay(idRun, right+13, by, rw-26, 38)
	s.lay(idProgress, right+13, by+33, rw-26, 16)
	s.lay(idMarquee, right+13, by+33, rw-26, 16)
	s.lay(idOptions, right+13, y+g.optionsTop, rw-26, 32)
	s.lay(idOutput, right+13, g.resultTop, rw-26, g.resultBottom-g.resultTop)
	s.lay(idTxt, right+13, g.exportTop, (rw-34)/2, 34)
	s.lay(idPDF, right+21+(rw-34)/2, g.exportTop, (rw-34)/2, 34)
	s.scheduleScrollSync()
	s.invalid()
}
func fill(dc uintptr, r dRect, c uintptr) {
	b, _, _ := dBrush.Call(c)
	dFill.Call(dc, uintptr(unsafe.Pointer(&r)), b)
	dDelete.Call(b)
}
func (s *deskState) round(dc uintptr, r dRect, c uintptr, rad int32) {
	b, _, _ := dBrush.Call(c)
	old, _, _ := dSelect.Call(dc, b)
	nullP, _, _ := dGDI.NewProc("GetStockObject").Call(8)
	oldP, _, _ := dSelect.Call(dc, nullP)
	dGDI.NewProc("RoundRect").Call(dc, uintptr(r.L), uintptr(r.T), uintptr(r.R), uintptr(r.B), uintptr(s.px(rad)), uintptr(s.px(rad)))
	dSelect.Call(dc, oldP)
	dSelect.Call(dc, old)
	dDelete.Call(b)
}
func plot(dc, font, color uintptr, r dRect, text string, flags uintptr) {
	old, _, _ := dSelect.Call(dc, font)
	dGDI.NewProc("SetBkMode").Call(dc, 1)
	dGDI.NewProc("SetTextColor").Call(dc, color)
	p := u16(text)
	dDraw.Call(dc, ptr(p), ^uintptr(0), uintptr(unsafe.Pointer(&r)), flags|0x0800)
	runtime.KeepAlive(p)
	dSelect.Call(dc, old)
}
func gradient(dc uintptr, r dRect, a, b uintptr) {
	v := func(x, y int32, c uintptr) dVertex {
		return dVertex{x, y, uint16(c&255) << 8, uint16((c>>8)&255) << 8, uint16((c>>16)&255) << 8, 0}
	}
	verts := [2]dVertex{v(r.L, r.T, a), v(r.R, r.B, b)}
	mesh := dGradient{0, 1}
	dGrad.NewProc("GradientFill").Call(dc, uintptr(unsafe.Pointer(&verts[0])), 2, uintptr(unsafe.Pointer(&mesh)), 1, 0)
}
func (s *deskState) panel(dc uintptr, r dRect) {
	s.round(dc, r, s.ui.line, 18)
	r.L++
	r.R--
	r.T++
	r.B--
	s.round(dc, r, s.ui.panel, 18)
}
func (s *deskState) paint(dc uintptr) {
	var r dRect
	dSize.Call(s.win, uintptr(unsafe.Pointer(&r)))
	w, h := s.dp(r.R), s.dp(r.B)
	fill(dc, r, s.ui.bg)
	g := layoutGUI(w, h)
	l := g.side
	x := l + 22
	y := g.contentTop
	body := w - x - 22
	left := body * 42 / 100
	right := x + left + 12
	rw := body - left - 12
	fill(dc, s.box(0, 0, l, h), s.ui.panel)
	fill(dc, s.box(l-1, 0, 1, h), s.ui.line)
	if s.icon != 0 {
		dUser.NewProc("DrawIconEx").Call(dc, uintptr(s.px(18)), uintptr(s.px(20)), s.icon, uintptr(s.px(31)), uintptr(s.px(31)), 0, 0, 3)
	} else {
		s.round(dc, s.box(18, 20, 31, 31), s.ui.accent, 12)
	}
	plot(dc, s.fHead, s.ui.ink, s.box(57, 18, l-62, 38), "SpeedDisc", 0x20|4)
	plot(dc, s.fSmall, s.ui.muted, s.box(18, 67, l-24, 36), s.l("OFICINA PORTÁTIL", "PORTABLE WORKBENCH"), 0x10)
	plot(dc, s.fSmall, s.ui.muted, s.box(20, g.navTop-28, l-30, 24), s.l("ESPAÇO DE TRABALHO", "WORKSPACE"), 0x20)
	line := s.l("Sessão normal", "Standard session")
	if s.admin && !s.demo {
		line = s.l("Sessão de administrador", "Administrator session")
	}
	if s.demo {
		line = s.l("Demonstração · dados fictícios", "Demo · fictional data")
	}
	plot(dc, s.fSmall, s.ui.ok, s.box(18, h-74, l-26, 32), line, 0x10)
	plot(dc, s.fSmall, s.ui.muted, s.box(18, h-32, l-26, 20), "danielpro.dev", 0x20)
	plot(dc, s.fSmall, s.ui.muted, s.box(x, 24, body-155, 27), s.l("Windows  /  Portátil · Sem instalação", "Windows  /  Portable · No installation"), 0x20|4)
	if g.heroHeight > 0 {
		hero := s.box(x, 70, body, g.heroHeight)
		save, _, _ := dGDI.NewProc("SaveDC").Call(dc)
		region, _, _ := dGDI.NewProc("CreateRoundRectRgn").Call(uintptr(hero.L), uintptr(hero.T), uintptr(hero.R+1), uintptr(hero.B+1), uintptr(s.px(20)), uintptr(s.px(20)))
		if region != 0 && save != 0 {
			dGDI.NewProc("ExtSelectClipRgn").Call(dc, region, 1)
		}
		gradient(dc, hero, s.ui.h1, s.ui.h2)
		if save != 0 {
			dGDI.NewProc("RestoreDC").Call(dc, save)
		}
		if region != 0 {
			dDelete.Call(region)
		}
		plot(dc, s.fSmall, s.ui.heroMuted, s.box(x+19, 85, body-39, 24), s.l("DIAGNÓSTICO E MANUTENÇÃO  /  OFFLINE", "DIAGNOSTICS & MAINTENANCE  /  OFFLINE"), 0x20)
		plot(dc, s.fHero, s.ui.heroInk, s.box(x+17, 111, body-35, 43), s.l("O teu PC, sob controlo.", "Your PC, in your control."), 0x20)
		if g.heroHeight >= 128 {
			plot(dc, s.fSmall, s.ui.heroMuted, s.box(x+19, 162, body-38, 34), s.l("Analisa primeiro. Decide o que mudar. Guarda o resultado.", "Analyse first. Choose what changes. Keep the result."), 0x10)
		}
	}
	s.panel(dc, s.box(x, y, left, h-y-26))
	s.panel(dc, s.box(right, y, rw, h-y-26))
	heading := s.l("Ferramentas", "Tools")
	if s.picking {
		heading = s.l("Entradas de arranque", "Startup entries")
	}
	plot(dc, s.fHead, s.ui.ink, s.box(x+13, y+11, left-26, 35), heading, 0x20)
	if !s.picking {
		plot(dc, s.fSmall, s.ui.muted, s.box(x+12, y+42, left-24, 19), s.l("Pesquisar ferramentas", "Search tools"), 0x20)
		s.round(dc, s.box(x+8, y+60, left-16, 28), s.ui.raised, 10)
	}
	count := fmt.Sprintf(s.l("%d ferramentas encontradas", "%d tools found"), len(s.rows))
	if s.picking {
		count = fmt.Sprintf(s.l("%d entradas encontradas", "%d entries found"), len(s.startupItems))
	}
	plot(dc, s.fSmall, s.ui.muted, s.box(x+12, y+94, left-24, 16), count, 0x20)
	title := s.item.Title
	if title == "" {
		title = s.l("Sem correspondências", "No matches")
	}
	kind := s.l("FERRAMENTA  /  SÓ LEITURA", "TOOL  /  READ ONLY")
	if s.item.Mutates {
		kind = s.l("FERRAMENTA  /  ALTERA O SISTEMA", "TOOL  /  CHANGES THE SYSTEM")
	}
	if s.item.ID == "14" {
		kind = s.l("ARRANQUE  /  ANÁLISE OU ALTERAÇÃO", "STARTUP  /  REVIEW OR CHANGE")
	}
	plot(dc, s.fSmall, s.ui.accent, s.box(right+13, y+13, rw-26, 21), kind, 0x20)
	detail := s.item.Detail
	if s.picking && s.picked >= 0 && s.picked < len(s.startupItems) {
		entry := s.startupItems[s.picked]
		title = entry.Name
		detail = entry.Detail
	}
	plot(dc, s.fHead, s.ui.ink, s.box(right+12, y+36, rw-24, 58), title, 0x10)
	if detail == "" {
		detail = s.l("Tente outro termo ou categoria.", "Try another term or category.")
	}
	plot(dc, s.fBody, s.ui.muted, s.box(right+13, y+96, rw-26, 65), detail, 0x10)
	fill(dc, s.box(right+13, y+g.divTop, rw-26, 1), s.ui.line)
	status := s.l("Último resultado", "Last result")
	if s.busy {
		status = s.l("Operação em curso", "Operation in progress")
	}
	if !s.busy && s.runResult != nil {
		status = s.runResult.Status
	}
	statusInk := s.ui.ink
	if s.runResult != nil && !s.busy {
		switch s.runResult.Status {
		case "Concluído", "Completed", "Exemplo fictício", "Fictional example":
			statusInk = s.ui.ok
		default:
			statusInk = s.ui.warn
		}
	}
	plot(dc, s.fBold, statusInk, s.box(right+13, y+g.resultHeading, rw-26, 23), status, 0x20|0x8000)
	if s.busy {
		plot(dc, s.fBold, s.ui.accent, s.box(right+13, y+g.runTop, rw-26, 28), s.progressCaption(), 0x20|0x8000)
	}
	plot(dc, s.fSmall, s.ui.muted, s.box(x, h-24, body, 20), s.l("RELATÓRIOS LOCAIS  ·  NADA É ENVIADO", "LOCAL REPORTS  ·  NOTHING IS UPLOADED"), 0x20)
}
func (s *deskState) drawControl(v *dDrawItem) {
	if v.CtlType == 2 && v.CtlID == idList {
		i := int(v.ItemID)
		r := v.Rect
		fill(v.DC, r, s.ui.panel)
		if v.State&1 != 0 {
			s.round(v.DC, r, s.ui.raised, 10)
		}
		if v.State&16 != 0 {
			fill(v.DC, dRect{r.L, r.T + 2, r.L + s.px(2), r.B - 2}, s.ui.accent)
		}
		name, mark := "", ""
		if s.picking {
			if i >= 0 && i < len(s.startupItems) {
				entry := s.startupItems[i]
				name = entry.Name
				mark = fmt.Sprintf("%02d", i+1)
			}
		} else if i >= 0 && i < len(s.rows) {
			name = s.rows[i].Title
			mark = fmt.Sprintf("%02s", s.rows[i].ID)
		}
		m := r
		m.L += s.px(9)
		m.R = m.L + s.px(27)
		plot(v.DC, s.fSmall, s.ui.accent, m, mark, 0x20|4)
		r.L += s.px(37)
		r.R -= s.px(6)
		plot(v.DC, s.fBold, s.ui.ink, r, name, 0x20|4|0x8000)
		return
	}
	if v.CtlType != 4 {
		return
	}
	r := v.Rect
	fill(v.DC, r, s.ui.panel)
	c, t := s.ui.raised, s.ui.ink
	if v.CtlID == idRun {
		c, t = s.ui.accent, s.ui.onAccent
	}
	if v.CtlID >= idNav && int(v.CtlID) < idNav+len(deskCats) {
		if deskCats[int(v.CtlID)-idNav] != s.mode {
			c = s.ui.panel
		} else {
			t = s.ui.accent
		}
	}
	if v.State&1 != 0 {
		c, t = s.ui.line, s.ui.ink
	}
	if v.State&4 != 0 {
		c, t = s.ui.panel, s.ui.muted
	}
	if v.State&16 != 0 {
		s.round(v.DC, r, s.ui.accent, 13)
		r.L += 2
		r.R -= 2
		r.T += 2
		r.B -= 2
	}
	s.round(v.DC, r, c, 12)
	if v.CtlID >= idNav && int(v.CtlID) < idNav+len(deskCats) {
		m := r
		m.L += s.px(9)
		m.R = m.L + s.px(27)
		plot(v.DC, s.fSmall, s.ui.accent, m, deskCatMarks[int(v.CtlID)-idNav], 0x20|4)
		r.L += s.px(37)
		plot(v.DC, s.fSmall, t, r, dRead(v.HWND), 0x20|4|0x8000)
	} else {
		plot(v.DC, s.fBold, t, r, dRead(v.HWND), 0x20|4|1|0x8000)
	}
}

func (s *deskState) message(body, title string, flags uintptr) uintptr {
	a, b := u16(body), u16(title)
	rc, _, _ := dUser.NewProc("MessageBoxW").Call(s.win, ptr(a), ptr(b), flags)
	runtime.KeepAlive(a)
	runtime.KeepAlive(b)
	return rc
}
func (s *deskState) command(id, notify int) {
	if id == idTheme && notify == 0 {
		s.light = !s.light
		s.updateTheme()
		s.sync()
		return
	}
	if s.busy || s.preparing {
		return
	}
	switch {
	case id == idSearch && notify == 0x300:
		s.leavePicker()
		s.filter(s.item.ID)
	case id == idList && notify == 1:
		i := int(dMessage(s.h[idList], 0x0188, 0, 0))
		if s.picking {
			s.pick(i)
		} else {
			s.selectItem(i)
		}
	case id == idLang && notify == 0:
		s.leavePicker()
		s.lang = otherLang(s.lang)
		s.captions()
		s.filter(s.item.ID)
	case id >= idNav && id < idNav+len(deskCats) && notify == 0:
		s.leavePicker()
		s.mode = deskCats[id-idNav]
		s.filter(s.item.ID)
	case id == idRun && notify == 0 && !s.demo:
		s.run()
	case id == idOptions && notify == 0 && !s.demo:
		if s.item.ID == "14" {
			if s.picking {
				s.leavePicker()
				s.filter("14")
			} else {
				s.chooseStartup()
			}
		}
		if s.item.ID == "19" {
			s.resetWinsock = !s.resetWinsock
			s.sync()
		}
	case id == idTxt && notify == 0:
		s.export("txt")
	case id == idPDF && notify == 0:
		s.export("pdf")
	}
}
func (s *deskState) leavePicker() {
	if s.picking {
		s.picking = false
		s.picked = -1
	}
}
func (s *deskState) chooseStartup() {
	s.preparing = true
	s.sync()
	lang := s.lang
	ch := s.startup
	win := s.win
	go func() {
		items, notes := guiStartup(lang)
		ch <- startupChoice{items, notes}
		dUser.NewProc("PostMessageW").Call(win, wmStartup, 0, 0)
	}()
}
func (s *deskState) showStartup(c startupChoice) {
	s.preparing = false
	s.startupItems = c.items
	s.startupNotes = c.notes
	if len(c.items) == 0 {
		m := s.l("Não há entradas disponíveis para rever.", "No entries are available to review.")
		if len(c.notes) > 0 {
			m += "\n\n" + strings.Join(c.notes, "\n")
		}
		s.sync()
		s.message(m, "SpeedDisc", 0x40)
		return
	}
	s.picking = true
	h := s.h[idList]
	dMessage(h, 0x0184, 0, 0)
	for _, a := range c.items {
		t := a.Name
		if !a.CanDisable {
			t += "  ·  " + s.l("Só leitura", "Read only")
		}
		p := u16(t)
		dMessage(h, 0x0180, 0, ptr(p))
		runtime.KeepAlive(p)
	}
	dMessage(h, 0x0186, 0, 0)
	s.pick(0)
	dUser.NewProc("SetFocus").Call(h)
}
func (s *deskState) pick(i int) {
	s.picked = i
	s.sync()
}
func (s *deskState) openElevated() {
	if s.message(s.l("O Windows vai pedir permissão para abrir outra janela do SpeedDisc. Nenhuma operação será executada automaticamente. A nova janela usa o perfil da conta escolhida no Windows. Confirme a ferramenta nessa janela.", "Windows will ask to open another SpeedDisc window. No operation will run automatically. The new window uses the Windows account you choose. Confirm the tool in that window."), "SpeedDisc", 0x24|0x100) != 6 {
		return
	}
	path, err := os.Executable()
	if err == nil {
		verb, exe, args := u16("runas"), u16(path), u16("--focus="+s.item.ID+" --lang="+s.l("pt", "en"))
		n, _, _ := dShell.NewProc("ShellExecuteW").Call(s.win, ptr(verb), ptr(exe), ptr(args), 0, 1)
		runtime.KeepAlive(verb)
		runtime.KeepAlive(exe)
		runtime.KeepAlive(args)
		if n > 32 {
			return
		}
	}
	s.message(s.l("A janela de administrador não foi aberta. Nenhuma operação foi executada.", "The administrator window did not open. No operation was run."), "SpeedDisc", 0x30)
}
func (s *deskState) run() {
	if s.demo {
		return
	}
	a := s.item
	if a.ID == "" {
		return
	}
	opts := guiOptions{Winsock: s.resetWinsock}
	var chosen string
	if s.picking {
		if s.picked < 0 || s.picked >= len(s.startupItems) {
			return
		}
		choice := s.startupItems[s.picked]
		if !choice.CanDisable {
			return
		}
		opts.StartupIDs = []string{choice.ID}
		chosen = choice.Name + "\n" + choice.Detail
	}
	if !s.admin && guiRequiresAdmin(a, opts) {
		s.openElevated()
		return
	}
	if (a.Mutates && a.ID != "14") || len(opts.StartupIDs) > 0 {
		detail := a.Detail
		if chosen != "" {
			detail = s.l("Será desactivada apenas a entrada indicada:", "Only the following entry will be disabled:") + "\n\n" + chosen
		}
		if a.ID == "19" && !opts.Winsock {
			detail = s.l("Limpa apenas a cache DNS. O Winsock não será reposto.", "Flushes the DNS cache only. Winsock will not be reset.")
		}
		if a.ID == "19" && opts.Winsock {
			detail = s.l("Limpa a cache DNS e repõe o Winsock. Pode ser necessário reiniciar o Windows.", "Flushes the DNS cache and resets Winsock. A Windows restart may be required.")
		}
		if s.message(a.Title+"\n\n"+detail+"\n\n"+s.l("Confirma esta alteração?", "Confirm this change?"), s.l("Rever antes de executar", "Review before running"), 0x34|0x100) != 6 {
			return
		}
		opts.Confirmed = true
	}
	s.busy = true
	s.runResult = nil
	s.progress = make(chan guiProgress, 1)
	progress := s.progress
	s.showProgress(guiProgress{Label: a.Title, Percent: -1})
	s.sync()
	lang := s.lang
	ch := s.work
	win := s.win
	opts.Progress = func(stage guiProgress) {
		select {
		case <-progress:
		default:
		}
		select {
		case progress <- stage:
		default:
		}
		dUser.NewProc("PostMessageW").Call(win, wmProgress, 0, 0)
	}
	go func() {
		result := guiRun(a.ID, lang, opts)
		ch <- result
		dUser.NewProc("PostMessageW").Call(win, wmDone, 0, 0)
	}()
}
func (s *deskState) progressCaption() string {
	p := s.stage
	label := s.l("A preparar…", "Preparing…")
	if p.Step > 0 && p.Total > 0 {
		label = fmt.Sprintf(s.l("Etapa %d de %d", "Step %d of %d"), p.Step, p.Total)
	}
	if p.Percent >= 0 {
		label += fmt.Sprintf(" · %d%%", min(p.Percent, 100))
	}
	return label
}

func (s *deskState) refreshProgress() {
	known := s.busy && s.stage.Percent >= 0
	unknown := s.busy && !known
	s.hShow(idProgress, known)
	s.hShow(idMarquee, unknown)
	if known {
		dMessage(s.h[idProgress], 0x402, uintptr(min(100, s.stage.Percent)), 0)
	}
	m := uintptr(0)
	if unknown {
		m = 1
	}
	dMessage(s.h[idMarquee], 0x40a, m, 45)
}

func (s *deskState) showProgress(progress guiProgress) {
	s.stage = progress
	note := s.l("O Windows não forneceu uma percentagem para esta etapa.", "Windows has not reported a percentage for this step.")
	if s.stage.Percent >= 0 {
		note = s.l("A percentagem refere-se à etapa actual, não ao pacote inteiro.", "The percentage applies to the current step, not the whole package.")
	}
	if s.demoRunning {
		note = s.l("Demonstração com dados fictícios. Nenhuma operação foi executada.", "Demo with fictional data. No operation was executed.")
	}
	dText(s.h[idOutput], s.progressCaption()+"\r\n"+s.stage.Label+"\r\n\r\n"+note)
	s.refreshProgress()
	s.invalid()
}

func (s *deskState) finish(r report.Result) {
	s.busy = false
	s.runResult = &r
	var b strings.Builder
	b.WriteString(r.Title + "\r\n")
	b.WriteString(fmt.Sprintf(s.l("%d registos · %s", "%d entries · %s"), len(r.Rows), r.Duration.Round(time.Millisecond)))
	b.WriteString("\r\n\r\n")
	for _, row := range r.Rows {
		if b.Len() > 48000 {
			b.WriteString("\r\n" + s.l("Consulte o ficheiro exportado para o relatório completo.", "See the exported file for the full report."))
			break
		}
		b.WriteString(row.Label + ": " + row.Value + "\r\n")
	}
	dText(s.h[idOutput], b.String())
	dMessage(s.h[idOutput], 0x00B1, 0, 0)
	dMessage(s.h[idOutput], 0x00B7, 0, 0)
	if s.picking {
		s.leavePicker()
		s.filter("14")
	}
	s.sync()
}

// OPENFILENAMEW, for a native save picker. Existing reports are never overwritten.
type dOpenFile struct {
	Size                         uint32
	Owner, Instance              uintptr
	Filter, CustomFilter         *uint16
	MaxCustomFilter, FilterIndex uint32
	File                         *uint16
	MaxFile                      uint32
	FileTitle                    *uint16
	MaxFileTitle                 uint32
	InitialDir, Title            *uint16
	Flags                        uint32
	FileOffset, FileExtension    uint16
	DefaultExt                   *uint16
	CustomData, Hook             uintptr
	Template                     *uint16
	Reserved                     uintptr
	ReservedFlags, FlagsEx       uint32
}

func (s *deskState) export(format string) {
	if s.runResult == nil {
		return
	}
	base := "SpeedDisc-" + s.runResult.Started.Format("20060102-150405") + "." + format
	buf := make([]uint16, 1024)
	copy(buf, syscall.StringToUTF16(base))
	extension, title := u16(format), u16(s.l("Guardar relatório SpeedDisc", "Save SpeedDisc report"))
	filters := syscall.StringToUTF16(strings.ToUpper(format) + " (*." + format + ")")
	filters = append(filters, syscall.StringToUTF16("*."+format)...)
	filters = append(filters, 0)
	d := dOpenFile{Owner: s.win, Filter: &filters[0], FilterIndex: 1, File: &buf[0], MaxFile: uint32(len(buf)), Title: title, DefaultExt: extension, Flags: 0x00080000 | 0x00000800 | 0x00000008 | 0x02000000}
	d.Size = uint32(unsafe.Sizeof(d))
	ok, _, _ := dCommon.NewProc("GetSaveFileNameW").Call(uintptr(unsafe.Pointer(&d)))
	if ok == 0 {
		n, _, _ := dCommon.NewProc("CommDlgExtendedError").Call()
		if n != 0 {
			s.message(fmt.Sprintf(s.l("A janela de gravação falhou (código %d).", "The save dialog failed (code %d)."), n), "SpeedDisc", 0x30)
		}
		return
	}
	path := windows.UTF16ToString(buf)
	err := report.Save(path, format, *s.runResult)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			s.message(s.l("O ficheiro já existe. Escolha outro nome para preservar o original.", "The file already exists. Choose a new name to preserve it."), "SpeedDisc", 0x30)
		} else {
			s.message(s.l("Não foi possível guardar o relatório: ", "Could not save the report: ")+err.Error(), "SpeedDisc", 0x30)
		}
		return
	}
	s.message(s.l("Relatório guardado em:", "Report saved to:")+"\n\n"+path, "SpeedDisc", 0x40)
}

type dPaint struct {
	DC              uintptr
	Erase           int32
	Rect            dRect
	Restore, Update int32
	Reserved        [32]byte
}

// LPARAM is not always a pointer. Reinterpret it only for the three
// documented messages whose payload is a temporary native structure.
func messageData(l uintptr) unsafe.Pointer { return *(*unsafe.Pointer)(unsafe.Pointer(&l)) }

func deskProc(hwnd uintptr, msg uint32, w, l uintptr) uintptr {
	s := desk
	if s == nil {
		return dDefault(hwnd, msg, w, l)
	}
	switch msg {
	case 1: // WM_CREATE
		s.win = hwnd
		if p := dUser.NewProc("GetDpiForWindow"); p.Find() == nil {
			v, _, _ := p.Call(hwnd)
			if v > 0 {
				s.dpi = int32(v)
			}
		}
		s.controls()
		s.layout()
		if s.pendingID != "" {
			s.filter(s.pendingID)
		}
		if s.demoRunning {
			s.busy = true
			s.showProgress(guiProgress{Step: 2, Total: 3, Percent: 42, Label: s.l("A verificar componentes do Windows", "Checking Windows components")})
			s.sync()
		} else if s.demo {
			s.finish(report.Result{
				Title:  s.l("Análise do sistema · Demonstração", "System analysis · Demo"),
				Status: s.l("Exemplo fictício", "Fictional example"), Lang: s.l("pt", "en"),
				Started: time.Now(), Duration: 175 * time.Millisecond,
				Rows: []report.Row{
					{Label: s.l("Disco principal", "System disk"), Value: "256 GB · " + s.l("84 GB livres", "84 GB available")},
					{Label: s.l("Memória", "Memory"), Value: "16 GB"},
					{Label: s.l("Arranque", "Startup"), Value: s.l("8 entradas · só leitura", "8 entries · read only")},
					{Label: s.l("Ligação de rede", "Network connection"), Value: s.l("Activa · demonstração", "Active · demo")},
					{Label: s.l("Energia", "Power"), Value: s.l("Equilibrado", "Balanced")},
					{Label: s.l("Verificação de armazenamento", "Storage check"), Value: s.l("Sem avisos nesta amostra", "No warnings in this sample")},
					{Label: s.l("Verificação de componentes", "Component check"), Value: s.l("Amostra sem alterações", "Sample without changes")},
					{Label: s.l("Observação", "Note"), Value: s.l("Nenhuma operação foi executada.", "No operation was executed.")},
				},
			})
		}
		return 0
	case 5: // WM_SIZE
		if len(s.h) > 0 {
			s.layout()
		}
		return 0
	case 0x24: // WM_GETMINMAXINFO
		m := (*dMinMax)(messageData(l))
		area := monitorWorkArea(hwnd)
		m.MinTrack.X = min(s.px(800), area.R-area.L)
		m.MinTrack.Y = min(s.px(505), area.B-area.T)
		return 0
	case 0x02E0: // WM_DPICHANGED
		s.dpi = int32(w & 0xffff)
		q := (*dRect)(messageData(l))
		dUser.NewProc("SetWindowPos").Call(hwnd, 0, uintptr(q.L), uintptr(q.T), uintptr(q.R-q.L), uintptr(q.B-q.T), 0x14)
		s.fonts()
		s.layout()
		return 0
	case 0x14: // WM_ERASEBKGND (paint covers the whole client area)
		return 1
	case 0xf: // WM_PAINT
		var p dPaint
		dc, _, _ := dUser.NewProc("BeginPaint").Call(hwnd, uintptr(unsafe.Pointer(&p)))
		if dc != 0 && s.fBody != 0 {
			s.paint(dc)
		}
		dUser.NewProc("EndPaint").Call(hwnd, uintptr(unsafe.Pointer(&p)))
		return 0
	case 0x2b: // WM_DRAWITEM
		if l != 0 {
			s.drawControl((*dDrawItem)(messageData(l)))
			return 1
		}
	case 0x111: // WM_COMMAND
		s.command(int(w&0xffff), int((w>>16)&0xffff))
		return 0
	case 0x133, 0x134, 0x138: // WM_CTLCOLOREDIT/LISTBOX/STATIC
		dGDI.NewProc("SetTextColor").Call(w, s.ui.ink)
		dGDI.NewProc("SetBkColor").Call(w, s.ui.panel)
		return s.brush
	case 0x001A, 0x031A, wmScrollSync: // WM_SETTINGCHANGE, WM_THEMECHANGED, or native control movement
		s.syncScrollRails()
		if msg == wmScrollSync {
			return 0
		}
	case wmProgress:
		select {
		case stage := <-s.progress:
			if s.busy && !s.demoRunning {
				s.showProgress(stage)
			}
		default:
		}
		return 0
	case wmDone:
		select {
		case r := <-s.work:
			s.finish(r)
		default:
		}
		return 0
	case wmStartup:
		select {
		case c := <-s.startup:
			s.showStartup(c)
		default:
		}
		return 0
	case 0x10: // WM_CLOSE
		if (s.busy && !s.demoRunning) || s.preparing {
			s.message(s.l("Aguarde até a operação terminar.", "Wait for the operation to finish."), "SpeedDisc", 0x40)
			return 0
		}
	case 2: // WM_DESTROY
		for _, f := range []uintptr{s.fSmall, s.fBody, s.fBold, s.fHead, s.fHero, s.brush} {
			if f != 0 {
				dDelete.Call(f)
			}
		}
		dUser.NewProc("PostQuitMessage").Call(0)
		return 0
	}
	return dDefault(hwnd, msg, w, l)
}

// Use the work area, not the entire panel behind the taskbar.
type desktopMonitorInfo struct {
	Size         uint32
	Bounds, Work dRect
	Flags        uint32
}

func monitorWorkArea(hwnd uintptr) dRect {
	monitor, _, _ := dUser.NewProc("MonitorFromWindow").Call(hwnd, 1)
	m := desktopMonitorInfo{Size: uint32(unsafe.Sizeof(desktopMonitorInfo{}))}
	if monitor != 0 {
		ok, _, _ := dUser.NewProc("GetMonitorInfoW").Call(monitor, uintptr(unsafe.Pointer(&m)))
		if ok != 0 && m.Work.R > m.Work.L && m.Work.B > m.Work.T {
			return m.Work
		}
	}
	w, _, _ := dUser.NewProc("GetSystemMetrics").Call(0)
	h, _, _ := dUser.NewProc("GetSystemMetrics").Call(1)
	return dRect{R: int32(w), B: int32(h)}
}
func main() {
	runtime.LockOSThread()
	dpiInit := dUser.NewProc("SetProcessDpiAwarenessContext")
	if dpiInit.Find() == nil {
		dpiInit.Call(^uintptr(3))
	}
	inst, _, _ := dWin.NewProc("GetModuleHandleW").Call(0)
	icoName := u16("APP")
	icon, _, _ := dUser.NewProc("LoadIconW").Call(inst, ptr(icoName))
	runtime.KeepAlive(icoName)
	cursor, _, _ := dUser.NewProc("LoadCursorW").Call(0, 32512)
	s := &deskState{instance: inst, icon: icon, h: make(map[int]uintptr), dpi: 96, admin: elevated(), lang: LangPT, mode: "all", work: make(chan report.Result, 1), startup: make(chan startupChoice, 1), picked: -1}
	for _, a := range os.Args[1:] {
		if a == "--demo" {
			s.demo = true
		}
		if a == "--demo-running" {
			s.demo = true
			s.demoRunning = true
		}
		if a == "--lang=en" {
			s.lang = otherLang(LangPT)
		}
		if raw, ok := strings.CutPrefix(a, "--focus="); ok {
			if n, e := strconv.Atoi(raw); e == nil && n >= 1 && n <= 19 {
				s.pendingID = strconv.Itoa(n)
			}
		}
	}
	desk = s
	cls := u16("SpeedDisc.Native.Workbench")
	wc := dClass{Size: uint32(unsafe.Sizeof(dClass{})), Style: 3, WndProc: dProc, Instance: inst, Icon: icon, IconSmall: icon, Cursor: cursor, Name: cls}
	rc, _, _ := dUser.NewProc("RegisterClassExW").Call(uintptr(unsafe.Pointer(&wc)))
	if rc == 0 {
		s.message(s.l("Não foi possível iniciar a interface.", "The interface could not be started."), "SpeedDisc", 0x10)
		return
	}
	title := u16("SpeedDisc · " + s.l("Oficina portátil", "Portable workbench"))
	area := monitorWorkArea(0)
	screenW, screenH := uintptr(area.R-area.L), uintptr(area.B-area.T)
	scale := uintptr(96)
	if p := dUser.NewProc("GetDpiForSystem"); p.Find() == nil {
		if n, _, _ := p.Call(); n > 0 {
			scale = n
		}
	}
	w, h := min(uintptr(1200)*scale/96, screenW), min(uintptr(800)*scale/96, screenH)
	s.dpi = int32(scale)
	x, y := area.L+int32((screenW-w)/2), area.T+int32((screenH-h)/2)
	hwnd, _, _ := dCreate.Call(0, ptr(cls), ptr(title), 0x02CF0000, uintptr(x), uintptr(y), w, h, 0, 0, inst, 0)
	runtime.KeepAlive(cls)
	runtime.KeepAlive(title)
	if hwnd == 0 {
		s.message(s.l("A janela não abriu.", "The window did not open."), "SpeedDisc", 0x10)
		return
	}
	dUser.NewProc("ShowWindow").Call(hwnd, 5)
	if s.pendingID == "" {
		dMessage(s.h[idList], 0x0197, 0, 0)
	}
	dUser.NewProc("UpdateWindow").Call(hwnd)
	var m dMsg
	for {
		n, _, _ := dUser.NewProc("GetMessageW").Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(n) <= 0 {
			break
		}
		handled, _, _ := dUser.NewProc("IsDialogMessageW").Call(hwnd, uintptr(unsafe.Pointer(&m)))
		if handled == 0 {
			dUser.NewProc("TranslateMessage").Call(uintptr(unsafe.Pointer(&m)))
			dUser.NewProc("DispatchMessageW").Call(uintptr(unsafe.Pointer(&m)))
		}
	}
}
