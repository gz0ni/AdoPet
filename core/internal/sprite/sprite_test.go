package sprite

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/gz0ni/claude-pet/core/assets"
)

func TestParseEmbeddedPet(t *testing.T) {
	p, err := ParsePet(assets.PetJSON)
	if err != nil {
		t.Fatal(err)
	}
	if p.Cell != (Cell{W: 192, H: 208}) {
		t.Fatalf("cell = %+v", p.Cell)
	}
	if a := p.Animations["idle"]; a.Frames != 6 || !a.Loop {
		t.Fatalf("idle = %+v", a)
	}
}

func TestParsePetRejectsInvalid(t *testing.T) {
	for _, in := range []string{`{`, `{"cell":{"w":0,"h":1},"animations":{"idle":{}}}`, `{"cell":{"w":1,"h":1},"animations":{}}`} {
		if _, err := ParsePet([]byte(in)); err == nil {
			t.Errorf("ParsePet(%s) succeeded", in)
		}
	}
}

func TestFrameRect(t *testing.T) {
	cell := Cell{W: 192, H: 208}
	tests := []struct {
		a    Animation
		i    int
		want image.Rectangle
	}{
		{Animation{Row: 0, Frames: 6}, 2, image.Rect(384, 0, 576, 208)},
		{Animation{Row: 4, Frames: 5}, 0, image.Rect(0, 832, 192, 1040)},
		{Animation{Rows: []int{9, 10}, Frames: 16}, 7, image.Rect(1344, 1872, 1536, 2080)},
		{Animation{Rows: []int{9, 10}, Frames: 16}, 9, image.Rect(192, 2080, 384, 2288)},
	}
	for _, tt := range tests {
		if got := FrameRect(tt.a, tt.i, cell); got != tt.want {
			t.Errorf("FrameRect(%+v, %d) = %v, want %v", tt.a, tt.i, got, tt.want)
		}
	}
}

func TestFrameAt(t *testing.T) {
	loop := Animation{Frames: 6, FPS: 6, Loop: true}
	once := Animation{Frames: 5, FPS: 10}
	tests := []struct {
		a    Animation
		d    time.Duration
		want int
	}{
		{loop, 0, 0},
		{loop, 500 * time.Millisecond, 3},
		{loop, 1100 * time.Millisecond, 0},
		{once, 300 * time.Millisecond, 3},
		{once, 5 * time.Second, 4},
		{Animation{}, time.Second, 0},
	}
	for _, tt := range tests {
		if got := FrameAt(tt.a, tt.d); got != tt.want {
			t.Errorf("FrameAt(%+v, %v) = %d, want %d", tt.a, tt.d, got, tt.want)
		}
	}
}

func TestScaleSheet(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 40, 20))
	for y := range 20 {
		for x := range 20 {
			src.Set(x, y, color.NRGBA{255, 0, 0, 255})
		}
	}
	dst, dc := ScaleSheet(src, Cell{W: 20, H: 20}, 0.5)
	if dc != (Cell{W: 10, H: 10}) || dst.Bounds() != image.Rect(0, 0, 20, 10) {
		t.Fatalf("cell %+v bounds %v", dc, dst.Bounds())
	}
	if c := dst.RGBAAt(5, 5); c.R != 255 || c.A != 255 {
		t.Errorf("opaque cell pixel = %v", c)
	}
	if c := dst.RGBAAt(15, 5); c.A != 0 {
		t.Errorf("empty cell pixel = %v", c)
	}
}
