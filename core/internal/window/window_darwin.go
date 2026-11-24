package window

/*
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore
#include "cocoa_darwin.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gz0ni/claude-pet/core/internal/config"
	"github.com/gz0ni/claude-pet/core/internal/sprite"
)

func init() {
	runtime.LockOSThread()
}

type petWindow struct {
	anim    sprite.Animation
	win     size
	start   time.Time
	frame   int
	anchor  string
	onMoved func(config.Position)
}

var current *petWindow

func Run(o Options) error {
	C.petInit()
	scr := screens()
	if len(scr) == 0 {
		return errors.New("no screens")
	}
	mon := scr[0]
	if d := o.Position.Display; d >= 0 && d < len(scr) {
		mon = scr[d]
	}

	w := &petWindow{
		anim:    o.Pet.Animations["idle"],
		start:   time.Now(),
		anchor:  o.Position.Anchor,
		onMoved: o.OnMoved,
	}
	winCell := sprite.ScaledCell(o.Pet.Cell, o.Scale)
	w.win = size{winCell.W, winCell.H}

	sheet, render := sprite.ScaleSheet(o.Sheet, o.Pet.Cell, o.Scale*mon.scale)
	o.Sheet = nil
	n := max(w.anim.Frames, 1)
	frameLen := render.W * render.H * 4
	pix := make([]byte, n*frameLen)
	for i := range n {
		blitRGBA(pix[i*frameLen:], sheet, sprite.FrameRect(w.anim, i, render))
	}
	sheet = nil
	debug.FreeOSMemory()

	x, y := toCocoa(o.Position, mon, w.win)
	ok := C.petCreate(C.double(x), C.double(y), C.double(w.win.w), C.double(w.win.h), C.double(mon.scale),
		C.bool(o.AlwaysOnTop), (*C.uint8_t)(&pix[0]), C.int(n), C.int(render.W), C.int(render.H))
	if !ok {
		return errors.New("create pet window")
	}
	runtime.KeepAlive(pix)
	pix = nil
	debug.FreeOSMemory()

	current = w
	C.petRun(C.double(1 / float64(max(w.anim.FPS, 1))))
	return nil
}

func screens() []screen {
	var out []screen
	for i := 0; ; i++ {
		var s screen
		n := int(C.petScreen(C.int(i), (*C.double)(&s.x), (*C.double)(&s.y), (*C.double)(&s.w), (*C.double)(&s.h), (*C.double)(&s.scale)))
		if i >= n {
			return out
		}
		out = append(out, s)
	}
}

//export goPetTick
func goPetTick() {
	w := current
	if w == nil {
		return
	}
	if i := sprite.FrameAt(w.anim, time.Since(w.start)); i != w.frame {
		w.frame = i
		C.petShowFrame(C.int(i))
	}
}

//export goPetMoved
func goPetMoved(x, y C.double, idx C.int) {
	w := current
	if w == nil || w.onMoved == nil {
		return
	}
	scr := screens()
	if len(scr) == 0 {
		return
	}
	d := int(idx)
	if d < 0 || d >= len(scr) {
		d = 0
	}
	p := fromCocoa(w.anchor, float64(x), float64(y), scr[d], w.win)
	p.Display = d
	w.onMoved(p)
}
