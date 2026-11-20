//go:build windows && speeddiscgui

package main

import (
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Keep the real Win32 scrollbar and its keyboard/UI Automation semantics.
// A non-tab-stop sibling masks only its on-screen bar in normal colour modes.
// If registration, subclassing, or native geometry fails, leave Win32 alone.
const (
	idListRail   = 112
	idResultRail = 113
	wmScrollSync = 0x8004
)

var (
	scrollControls = windows.NewLazySystemDLL("comctl32.dll")
	scrollSubclass = scrollControls.NewProc("SetWindowSubclass")
	scrollDef      = scrollControls.NewProc("DefSubclassProc")
	scrollCallback uintptr
	scrollRailCall = syscall.NewCallback(scrollRailProc)
)

type scrollLayout struct {
	handle, target uintptr
	height         int32
	margin         int32
	thumb, length  int32
	low, last      int32
	pos            int32
	dragging       bool
	dragOffset     int32
}

type scrollBarInfo struct {
	size                        uint32
	bounds                      dRect
	button, thumb1, thumb2, pad int32
	state                       [6]uint32
}
type scrollInfo struct {
	size, mask uint32
	low, high  int32
	page       uint32
	pos, track int32
}
type highContrast struct {
	size, flags uint32
	scheme      *uint16
}

func (s *deskState) initScrollRails() {
	if scrollCallback == 0 {
		scrollCallback = syscall.NewCallback(scrollControlProc)
	}
	if scrollSubclass.Find() != nil || scrollDef.Find() != nil {
		return
	}
	name := u16("SpeedDisc.Native.ScrollRail")
	cursor, _, _ := dUser.NewProc("LoadCursorW").Call(0, 32512)
	wc := dClass{Size: uint32(unsafe.Sizeof(dClass{})), WndProc: scrollRailCall, Instance: s.instance, Cursor: cursor, Name: name}
	registered, _, _ := dUser.NewProc("RegisterClassExW").Call(uintptr(unsafe.Pointer(&wc)))
	if registered == 0 {
		return
	}
	for i, id := range []int{idList, idOutput} {
		target := s.h[id]
		if target == 0 {
			continue
		}
		h, _, _ := dCreate.Call(0, ptr(name), 0, 0x44000000, 0, 0, 0, 0, s.win, uintptr(idListRail+i), s.instance, 0) // WS_CHILD | WS_CLIPSIBLINGS
		if h == 0 {
			continue
		}
		ok, _, _ := scrollSubclass.Call(target, scrollCallback, 1, 0)
		if ok == 0 {
			dUser.NewProc("DestroyWindow").Call(h)
			continue
		}
		s.scroll[i] = scrollLayout{handle: h, target: target}
	}
	runtime.KeepAlive(name)
}

func scrollControlProc(hwnd uintptr, message uint32, w, l, subclassID, _ uintptr) uintptr {
	if message == 0x0082 { // WM_NCDESTROY
		scrollControls.NewProc("RemoveWindowSubclass").Call(hwnd, scrollCallback, subclassID)
	}
	result, _, _ := scrollDef.Call(hwnd, uintptr(message), w, l)
	switch message {
	case 0x0005, 0x000C, 0x0018, 0x0083, 0x0100, 0x0102, 0x0115, 0x0201, 0x0202, 0x020A,
		0x0180, 0x0181, 0x0182, 0x0184, 0x0186, 0x0197, 0x01A0, // list population/size, selection and top index
		0x00B1, 0x00B6, 0x00B7: // edit selection, line scroll and scroll caret
		if desk != nil {
			desk.scheduleScrollSync()
		}
	}
	return result
}

func (s *deskState) scheduleScrollSync() {
	if s.scrollPosted || s.win == 0 || (s.scroll[0].handle == 0 && s.scroll[1].handle == 0) {
		return
	}
	ok, _, _ := dUser.NewProc("PostMessageW").Call(s.win, wmScrollSync, 0, 0)
	s.scrollPosted = ok != 0
}

func (s *deskState) syncScrollRails() {
	s.scrollPosted = false
	hc := highContrast{size: uint32(unsafe.Sizeof(highContrast{}))}
	known, _, _ := dUser.NewProc("SystemParametersInfoW").Call(0x42, uintptr(hc.size), uintptr(unsafe.Pointer(&hc)), 0)
	for i := range s.scroll {
		r := &s.scroll[i]
		if r.handle == 0 {
			continue
		}
		bar := scrollBarInfo{size: uint32(unsafe.Sizeof(scrollBarInfo{}))}
		info := scrollInfo{size: uint32(unsafe.Sizeof(scrollInfo{})), mask: 7}
		barOK, _, _ := dUser.NewProc("GetScrollBarInfo").Call(r.target, 0xFFFFFFFB, uintptr(unsafe.Pointer(&bar))) // OBJID_VSCROLL
		rangeOK, _, _ := dUser.NewProc("GetScrollInfo").Call(r.target, 1, uintptr(unsafe.Pointer(&info)))          // SB_VERT
		visible, _, _ := dUser.NewProc("IsWindowVisible").Call(r.target)
		if known == 0 || hc.flags&1 != 0 || visible == 0 || barOK == 0 || rangeOK == 0 ||
			bar.state[0]&0x18001 != 0 || bar.bounds.R <= bar.bounds.L || bar.bounds.B <= bar.bounds.T || int64(info.high)-int64(info.low)+1 <= int64(max(1, info.page)) {
			dUser.NewProc("ShowWindow").Call(r.handle, 0)
			continue
		}
		b := bar.bounds
		width, height := b.R-b.L, b.B-b.T
		point := dPoint{b.L, b.T}
		ok, _, _ := dUser.NewProc("ScreenToClient").Call(s.win, uintptr(unsafe.Pointer(&point)))
		if ok == 0 {
			dUser.NewProc("ShowWindow").Call(r.handle, 0)
			continue
		}
		r.margin = min(s.px(5), height/6)
		r.height = height - 2*r.margin
		r.low, r.pos = info.low, info.pos
		r.thumb, r.length, r.last = scrollThumb(r.height, s.px(25), info.low, info.high, int32(info.page), info.pos)
		// HWND_TOP + SWP_NOACTIVATE | SWP_SHOWWINDOW. Coordinate conversion is
		// physical on the GUI's monitor-aware thread, matching GetScrollBarInfo.
		dUser.NewProc("SetWindowPos").Call(r.handle, 0, uintptr(point.X), uintptr(point.Y), uintptr(width), uintptr(height), 0x50)
		dUser.NewProc("InvalidateRect").Call(r.handle, 0, 0)
	}
}

func scrollRailProc(hwnd uintptr, message uint32, w, l uintptr) uintptr {
	s := desk
	if s == nil {
		return dDefault(hwnd, message, w, l)
	}
	var r *scrollLayout
	for i := range s.scroll {
		if s.scroll[i].handle == hwnd {
			r = &s.scroll[i]
		}
	}
	if r == nil {
		return dDefault(hwnd, message, w, l)
	}
	paint := func(dc uintptr) {
		var bounds dRect
		dSize.Call(hwnd, uintptr(unsafe.Pointer(&bounds)))
		fill(dc, bounds, s.ui.panel)
		if r.height < 1 || bounds.R < 4 {
			return
		}
		pad := max(1, min(s.px(3), bounds.R/4))
		s.round(dc, dRect{pad, r.margin, bounds.R - pad, r.margin + r.height}, s.ui.raised, 9)
		s.round(dc, dRect{pad, r.margin + r.thumb, bounds.R - pad, r.margin + r.thumb + r.length}, s.ui.accent, 9)
	}
	y := int32(int16(l>>16)) - r.margin
	switch message {
	case 0x0014: // WM_ERASEBKGND
		return 1
	case 0x000F: // WM_PAINT
		var p dPaint
		dc, _, _ := dUser.NewProc("BeginPaint").Call(hwnd, uintptr(unsafe.Pointer(&p)))
		if dc != 0 {
			paint(dc)
		}
		dUser.NewProc("EndPaint").Call(hwnd, uintptr(unsafe.Pointer(&p)))
		return 0
	case 0x0318: // WM_PRINTCLIENT
		if w != 0 {
			paint(w)
		}
		return 0
	case 0x0201: // Mouse-down: drag the thumb, or page along the track.
		if y >= r.thumb && y < r.thumb+r.length {
			r.dragging = true
			r.dragOffset = y - r.thumb
			dUser.NewProc("SetCapture").Call(hwnd)
		} else {
			page := uintptr(2) // SB_PAGEUP
			if y > r.thumb {
				page = 3 // SB_PAGEDOWN
			}
			dMessage(r.target, 0x0115, page, 0)
		}
		return 0
	case 0x0200: // WM_MOUSEMOVE
		if r.dragging {
			target := scrollTarget(y-r.dragOffset, r.height-r.length, r.low, r.last)
			if r.target == s.h[idList] {
				dMessage(r.target, 0x0197, uintptr(target), 0)
			} else {
				top := int32(dMessage(r.target, 0x00CE, 0, 0))     // EM_GETFIRSTVISIBLELINE
				dMessage(r.target, 0x00B6, 0, uintptr(target-top)) // EM_LINESCROLL
			}
		}
		return 0
	case 0x0202: // WM_LBUTTONUP
		if r.dragging {
			r.dragging = false
			dUser.NewProc("ReleaseCapture").Call()
		}
		return 0
	case 0x0215: // WM_CAPTURECHANGED
		r.dragging = false
		return 0
	case 0x020A: // WM_MOUSEWHEEL
		dMessage(r.target, message, w, l)
		return 0
	}
	return dDefault(hwnd, message, w, l)
}
