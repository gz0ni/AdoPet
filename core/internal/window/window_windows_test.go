package window

import (
	"testing"

	"github.com/gz0ni/claude-pet/core/internal/config"
)

func TestScreenRoundTrip(t *testing.T) {
	mon := monitor{work: rect{Left: -1920, Top: 0, Right: 0, Bottom: 1032}}
	win := size{96, 104}
	for _, dsf := range []float64{1, 1.25, 1.5, 2} {
		p := config.Position{X: 40, Y: -40, Anchor: "bottom-left"}
		x, y := toScreen(p, mon, dsf, win)
		got := fromScreen(p.Anchor, x, y, mon, dsf, win)
		if got.X != p.X || got.Y != p.Y || got.Anchor != p.Anchor {
			t.Errorf("dsf %v: %+v -> %d,%d -> %+v", dsf, p, x, y, got)
		}
	}
}

func TestToScreenOffsetsByWorkArea(t *testing.T) {
	mon := monitor{work: rect{Left: 100, Top: 50, Right: 1100, Bottom: 850}}
	x, y := toScreen(config.Position{X: 10, Y: 20, Anchor: "top-left"}, mon, 2, size{50, 50})
	if x != 120 || y != 90 {
		t.Errorf("toScreen = %d,%d, want 120,90", x, y)
	}
}
