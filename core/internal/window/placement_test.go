package window

import (
	"testing"

	"github.com/gz0ni/claude-pet/core/internal/config"
)

func TestPlace(t *testing.T) {
	mon, win := size{1920, 1080}, size{96, 104}
	tests := []struct {
		p            config.Position
		wantX, wantY int
	}{
		{config.Position{X: 40, Y: -40, Anchor: "bottom-left"}, 40, 936},
		{config.Position{X: 10, Y: 20, Anchor: "top-left"}, 10, 20},
		{config.Position{X: -10, Y: 20, Anchor: "top-right"}, 1814, 20},
		{config.Position{X: -10, Y: -10, Anchor: "bottom-right"}, 1814, 966},
		{config.Position{X: 5000, Y: -5000, Anchor: "top-left"}, 1824, 0},
		{config.Position{X: 1, Y: 2}, 1, 2},
	}
	for _, tt := range tests {
		x, y := place(tt.p, mon, win)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("place(%+v) = %d,%d, want %d,%d", tt.p, x, y, tt.wantX, tt.wantY)
		}
	}
}

func TestRelativeRoundTrip(t *testing.T) {
	mon, win := size{2560, 1440}, size{144, 156}
	for _, anchor := range []string{"top-left", "top-right", "bottom-left", "bottom-right"} {
		p := relative(anchor, 300, 700, mon, win)
		if p.Anchor != anchor {
			t.Fatalf("anchor = %q", p.Anchor)
		}
		if x, y := place(p, mon, win); x != 300 || y != 700 {
			t.Errorf("%s: round trip = %d,%d", anchor, x, y)
		}
	}
}

func TestRelativeDefaultsAnchor(t *testing.T) {
	if p := relative("", 1, 2, size{100, 100}, size{10, 10}); p.Anchor != "top-left" || p.X != 1 || p.Y != 2 {
		t.Errorf("relative = %+v", p)
	}
}
