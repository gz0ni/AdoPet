package window

import (
	"testing"

	"github.com/gz0ni/claude-pet/core/internal/config"
)

func TestCocoaRoundTrip(t *testing.T) {
	s := screen{x: -1470, y: 0, w: 1470, h: 923, scale: 2}
	win := size{96, 104}
	for _, p := range []config.Position{
		{X: 40, Y: -40, Anchor: "bottom-left"},
		{X: 10, Y: 20, Anchor: "top-left"},
		{X: -30, Y: 50, Anchor: "top-right"},
	} {
		x, y := toCocoa(p, s, win)
		if got := fromCocoa(p.Anchor, x, y, s, win); got != p {
			t.Errorf("%+v -> %v,%v -> %+v", p, x, y, got)
		}
	}
}

func TestToCocoaFlipsY(t *testing.T) {
	s := screen{x: 0, y: 70, w: 1470, h: 853, scale: 2}
	x, y := toCocoa(config.Position{X: 10, Y: 20, Anchor: "top-left"}, s, size{50, 60})
	if x != 10 || y != 843 {
		t.Errorf("toCocoa = %v,%v, want 10,843", x, y)
	}
}
