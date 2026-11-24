package window

import "github.com/gz0ni/claude-pet/core/internal/config"

type size struct{ w, h int }

func place(p config.Position, mon, win size) (x, y int) {
	x, y = p.X, p.Y
	if isRight(p.Anchor) {
		x += mon.w - win.w
	}
	if isBottom(p.Anchor) {
		y += mon.h - win.h
	}
	return clamp(x, 0, mon.w-win.w), clamp(y, 0, mon.h-win.h)
}

func relative(anchor string, x, y int, mon, win size) config.Position {
	if anchor == "" {
		anchor = "top-left"
	}
	p := config.Position{X: x, Y: y, Anchor: anchor}
	if isRight(anchor) {
		p.X -= mon.w - win.w
	}
	if isBottom(anchor) {
		p.Y -= mon.h - win.h
	}
	return p
}

func isRight(anchor string) bool  { return anchor == "top-right" || anchor == "bottom-right" }
func isBottom(anchor string) bool { return anchor == "bottom-left" || anchor == "bottom-right" }

func clamp(v, lo, hi int) int {
	if hi < lo {
		return lo
	}
	return max(lo, min(v, hi))
}
