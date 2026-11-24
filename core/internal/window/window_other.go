//go:build !windows

package window

import (
	"runtime/debug"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/gz0ni/claude-pet/core/internal/config"
	"github.com/gz0ni/claude-pet/core/internal/sprite"
)

type game struct {
	anim    sprite.Animation
	frames  []*ebiten.Image
	start   time.Time
	frame   int
	dirty   bool
	win     size
	render  sprite.Cell
	anchor  string
	onMoved func(config.Position)

	dragging     bool
	moved        bool
	grabX, grabY int
}

func Run(o Options) error {
	mons := ebiten.AppendMonitors(nil)
	if o.Position.Display >= 0 && o.Position.Display < len(mons) {
		ebiten.SetMonitor(mons[o.Position.Display])
	}
	mon := ebiten.Monitor()
	dsf := 1.0
	if mon != nil {
		dsf = mon.DeviceScaleFactor()
	}

	winCell := sprite.ScaledCell(o.Pet.Cell, o.Scale)
	sheet, render := sprite.ScaleSheet(o.Sheet, o.Pet.Cell, o.Scale*dsf)
	o.Sheet = nil
	img := ebiten.NewImageFromImage(sheet)

	g := &game{
		anim:    o.Pet.Animations["idle"],
		start:   time.Now(),
		dirty:   true,
		win:     size{winCell.W, winCell.H},
		render:  render,
		anchor:  o.Position.Anchor,
		onMoved: o.OnMoved,
	}
	for i := range g.anim.Frames {
		r := sprite.FrameRect(g.anim, i, render)
		g.frames = append(g.frames, img.SubImage(r).(*ebiten.Image))
	}
	debug.FreeOSMemory()

	ebiten.SetWindowTitle("Claude Pet")
	ebiten.SetWindowSize(g.win.w, g.win.h)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(o.AlwaysOnTop)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum) //nolint:staticcheck
	go scheduleFrames(max(g.anim.FPS, 1))
	if mon != nil {
		mw, mh := mon.Size()
		ebiten.SetWindowPosition(place(o.Position, size{mw, mh}, g.win))
	}

	return ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{
		ScreenTransparent: true,
		SkipTaskbar:       true,
		InitUnfocused:     true,
	})
}

func scheduleFrames(fps int) {
	for range time.Tick(time.Second / time.Duration(fps)) {
		ebiten.ScheduleFrame() //nolint:staticcheck
	}
}

func (g *game) Update() error {
	g.updateDrag()
	if i := sprite.FrameAt(g.anim, time.Since(g.start)); i != g.frame {
		g.frame = i
		g.dirty = true
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if !g.dirty || len(g.frames) == 0 {
		return
	}
	screen.Clear()
	screen.DrawImage(g.frames[g.frame], nil)
	g.dirty = false
}

func (g *game) Layout(int, int) (int, int) {
	return g.render.W, g.render.H
}

func (g *game) updateDrag() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.dragging, g.moved = true, false
		g.grabX, g.grabY = g.cursor()
		return
	}
	if !g.dragging {
		return
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.dragging = false
		if g.moved {
			g.savePosition()
		}
		return
	}
	cx, cy := g.cursor()
	if cx == g.grabX && cy == g.grabY {
		return
	}
	wx, wy := ebiten.WindowPosition()
	ebiten.SetWindowPosition(wx+cx-g.grabX, wy+cy-g.grabY)
	g.moved = true
}

func (g *game) cursor() (int, int) {
	x, y := ebiten.CursorPosition()
	return x * g.win.w / g.render.W, y * g.win.h / g.render.H
}

func (g *game) savePosition() {
	if g.onMoved == nil {
		return
	}
	mon := ebiten.Monitor()
	if mon == nil {
		return
	}
	mw, mh := mon.Size()
	x, y := ebiten.WindowPosition()
	p := relative(g.anchor, x, y, size{mw, mh}, g.win)
	p.Display = max(0, slices.Index(ebiten.AppendMonitors(nil), mon))
	g.onMoved(p)
}
