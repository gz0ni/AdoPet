package window

import (
	"math"

	"github.com/gz0ni/claude-pet/core/internal/config"
)

type screen struct {
	x, y, w, h float64
	scale      float64
}

func (s screen) size() size {
	return size{int(math.Round(s.w)), int(math.Round(s.h))}
}

func toCocoa(p config.Position, s screen, win size) (float64, float64) {
	x, y := place(p, s.size(), win)
	return s.x + float64(x), s.y + s.h - float64(y+win.h)
}

func fromCocoa(anchor string, x, y float64, s screen, win size) config.Position {
	dx := int(math.Round(x - s.x))
	dy := int(math.Round(s.y + s.h - y - float64(win.h)))
	return relative(anchor, dx, dy, s.size(), win)
}
