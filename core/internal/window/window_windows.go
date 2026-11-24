package window

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"runtime/debug"
	"slices"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/gz0ni/claude-pet/core/internal/config"
	"github.com/gz0ni/claude-pet/core/internal/sprite"
)

type monitor struct {
	handle  uintptr
	work    rect
	primary bool
}

type petWindow struct {
	hwnd    uintptr
	memDC   uintptr
	bitmap  uintptr
	bits    []byte
	sheet   *image.RGBA
	anim    sprite.Animation
	render  sprite.Cell
	win     size
	start   time.Time
	frame   int
	anchor  string
	last    rect
	onMoved func(config.Position)
}

var (
	current    *petWindow
	wndProcPtr = syscall.NewCallback(wndProc)
	enumResult []monitor
	enumProc   = syscall.NewCallback(func(h, _, _, _ uintptr) uintptr {
		mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
		if r, _, _ := procGetMonitorInfoW.Call(h, ptr(&mi)); r != 0 {
			enumResult = append(enumResult, monitor{handle: h, work: mi.Work, primary: mi.Flags&monitorPrimary != 0})
		}
		return 1
	})
)

func Run(o Options) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_, _ = call(procSetProcessDpiAwarenessContext, dpiAwarenessPerMonitorV2)

	mons, err := monitors()
	if err != nil {
		return err
	}
	mon := mons[0]
	if d := o.Position.Display; d >= 0 && d < len(mons) {
		mon = mons[d]
	}
	dsf := monitorScale(mon.handle)

	winCell := sprite.ScaledCell(o.Pet.Cell, o.Scale)
	sheet, render := sprite.ScaleSheet(o.Sheet, o.Pet.Cell, o.Scale*dsf)
	o.Sheet = nil
	debug.FreeOSMemory()

	w := &petWindow{
		sheet:   sheet,
		anim:    o.Pet.Animations["idle"],
		render:  render,
		win:     size{winCell.W, winCell.H},
		start:   time.Now(),
		anchor:  o.Position.Anchor,
		onMoved: o.OnMoved,
	}
	x, y := toScreen(o.Position, mon, dsf, w.win)
	if err := w.create(x, y, o.AlwaysOnTop); err != nil {
		w.release()
		return err
	}
	defer w.release()
	current = w
	defer func() { current = nil }()

	if err := w.draw(&point{int32(x), int32(y)}); err != nil {
		return err
	}
	w.last = w.rect()
	_, _, _ = procShowWindow.Call(w.hwnd, swShowNoActivate)
	if _, err := call(procSetTimer, w.hwnd, 1, uintptr(1000/max(w.anim.FPS, 1)), 0); err != nil {
		return fmt.Errorf("SetTimer: %w", err)
	}
	return loop()
}

func loop() error {
	var m msg
	for {
		r, _, err := procGetMessageW.Call(ptr(&m), 0, 0, 0)
		switch int32(r) {
		case -1:
			return fmt.Errorf("GetMessage: %w", err)
		case 0:
			return nil
		}
		_, _, _ = procTranslateMessage.Call(ptr(&m))
		_, _, _ = procDispatchMessageW.Call(ptr(&m))
	}
}

func (w *petWindow) create(x, y int, topmost bool) error {
	className, _ := windows.UTF16PtrFromString("ClaudePetWindow")
	title, _ := windows.UTF16PtrFromString("Claude Pet")
	inst, _, _ := procGetModuleHandleW.Call(0)
	cursor, _, _ := procLoadCursorW.Call(0, idcArrow)

	wc := wndClassEx{WndProc: wndProcPtr, Instance: inst, Cursor: cursor, ClassName: className}
	wc.Size = uint32(unsafe.Sizeof(wc))
	if _, err := call(procRegisterClassExW, ptr(&wc)); err != nil {
		return fmt.Errorf("RegisterClassEx: %w", err)
	}

	ex := uintptr(wsExLayered | wsExToolWindow | wsExNoActivate)
	if topmost {
		ex |= wsExTopmost
	}
	hwnd, err := call(procCreateWindowExW, ex, ptr(className), ptr(title), wsPopup,
		uintptr(x), uintptr(y), uintptr(w.render.W), uintptr(w.render.H), 0, 0, inst, 0)
	if err != nil {
		return fmt.Errorf("CreateWindowEx: %w", err)
	}
	w.hwnd = hwnd

	screen, _, _ := procGetDC.Call(0)
	defer func() { _, _, _ = procReleaseDC.Call(0, screen) }()
	if w.memDC, err = call(procCreateCompatibleDC, screen); err != nil {
		return fmt.Errorf("CreateCompatibleDC: %w", err)
	}
	bmi := bitmapInfo{Header: bitmapInfoHeader{
		Width:       int32(w.render.W),
		Height:      -int32(w.render.H),
		Planes:      1,
		BitCount:    32,
		Compression: biRGB,
	}}
	bmi.Header.Size = uint32(unsafe.Sizeof(bmi.Header))
	var bits unsafe.Pointer
	if w.bitmap, err = call(procCreateDIBSection, w.memDC, ptr(&bmi), dibRGB, uintptr(unsafe.Pointer(&bits)), 0, 0); err != nil {
		return fmt.Errorf("CreateDIBSection: %w", err)
	}
	w.bits = unsafe.Slice((*byte)(bits), w.render.W*w.render.H*4)
	_, _, _ = procSelectObject.Call(w.memDC, w.bitmap)
	return nil
}

func (w *petWindow) release() {
	if w.hwnd != 0 {
		_, _, _ = procKillTimer.Call(w.hwnd, 1)
		_, _, _ = procDestroyWindow.Call(w.hwnd)
		w.hwnd = 0
	}
	if w.memDC != 0 {
		_, _, _ = procDeleteDC.Call(w.memDC)
		w.memDC = 0
	}
	if w.bitmap != 0 {
		_, _, _ = procDeleteObject.Call(w.bitmap)
		w.bitmap = 0
	}
}

func (w *petWindow) draw(dst *point) error {
	blitBGRA(w.bits, w.sheet, sprite.FrameRect(w.anim, w.frame, w.render))
	sz := sizeT{int32(w.render.W), int32(w.render.H)}
	src := point{}
	blend := blendFunction{BlendOp: acSrcOver, SourceConstantAlpha: 255, AlphaFormat: acSrcAlpha}
	screen, _, _ := procGetDC.Call(0)
	defer func() { _, _, _ = procReleaseDC.Call(0, screen) }()
	if _, err := call(procUpdateLayeredWindow, w.hwnd, screen, ptr(dst), ptr(&sz), w.memDC, ptr(&src), 0, ptr(&blend), ulwAlpha); err != nil {
		return fmt.Errorf("UpdateLayeredWindow: %w", err)
	}
	return nil
}

func (w *petWindow) tick() {
	if i := sprite.FrameAt(w.anim, time.Since(w.start)); i != w.frame {
		w.frame = i
		_ = w.draw(nil)
	}
}

func (w *petWindow) rect() rect {
	var r rect
	_, _, _ = procGetWindowRect.Call(w.hwnd, ptr(&r))
	return r
}

func (w *petWindow) savePosition() {
	r := w.rect()
	if r == w.last || w.onMoved == nil {
		return
	}
	w.last = r
	mons, err := monitors()
	if err != nil {
		return
	}
	h, _, _ := procMonitorFromWindow.Call(w.hwnd, monitorNear)
	idx := max(0, slices.IndexFunc(mons, func(m monitor) bool { return m.handle == h }))
	mon := mons[idx]
	p := fromScreen(w.anchor, int(r.Left), int(r.Top), mon, monitorScale(mon.handle), w.win)
	p.Display = idx
	w.onMoved(p)
}

func wndProc(hwnd, m, wp, lp uintptr) uintptr {
	switch m {
	case wmNCHitTest:
		return htCaption
	case wmMouseActivate:
		return maNoActivate
	case wmNCLButtonDblClk, wmNCRButtonDown, wmNCRButtonUp, wmContextMenu:
		return 0
	case wmTimer:
		if current != nil {
			current.tick()
		}
		return 0
	case wmExitSizeMove:
		if current != nil {
			current.savePosition()
		}
		return 0
	case wmDestroy:
		if current != nil && current.hwnd == hwnd {
			current.hwnd = 0
		}
		_, _, _ = procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, m, wp, lp)
	return r
}

func monitors() ([]monitor, error) {
	enumResult = nil
	if _, err := call(procEnumDisplayMonitors, 0, 0, enumProc, 0); err != nil {
		return nil, fmt.Errorf("EnumDisplayMonitors: %w", err)
	}
	if len(enumResult) == 0 {
		return nil, fmt.Errorf("no monitors")
	}
	mons := slices.Clone(enumResult)
	slices.SortStableFunc(mons, func(a, b monitor) int {
		switch {
		case a.primary && !b.primary:
			return -1
		case b.primary && !a.primary:
			return 1
		}
		return 0
	})
	return mons, nil
}

func monitorScale(h uintptr) float64 {
	var dx, dy uint32
	if r, _, _ := procGetDpiForMonitor.Call(h, 0, ptr(&dx), ptr(&dy)); r != 0 || dx == 0 {
		return 1
	}
	return float64(dx) / 96
}

func workSize(m monitor, dsf float64) size {
	return size{
		int(math.Round(float64(m.work.Right-m.work.Left) / dsf)),
		int(math.Round(float64(m.work.Bottom-m.work.Top) / dsf)),
	}
}

func toScreen(p config.Position, m monitor, dsf float64, win size) (int, int) {
	x, y := place(p, workSize(m, dsf), win)
	return int(m.work.Left) + int(math.Round(float64(x)*dsf)), int(m.work.Top) + int(math.Round(float64(y)*dsf))
}

func fromScreen(anchor string, x, y int, m monitor, dsf float64, win size) config.Position {
	dx := int(math.Round(float64(x-int(m.work.Left)) / dsf))
	dy := int(math.Round(float64(y-int(m.work.Top)) / dsf))
	return relative(anchor, dx, dy, workSize(m, dsf), win)
}
