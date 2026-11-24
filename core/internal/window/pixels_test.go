package window

import (
	"image"
	"image/color"
	"slices"
	"testing"
)

func TestBlitBGRA(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 2))
	src.SetRGBA(2, 1, color.RGBA{R: 10, G: 20, B: 30, A: 40})
	src.SetRGBA(3, 1, color.RGBA{R: 1, G: 2, B: 3, A: 4})
	dst := make([]byte, 2*1*4)
	blitBGRA(dst, src, image.Rect(2, 1, 4, 2))
	if want := []byte{30, 20, 10, 40, 3, 2, 1, 4}; !slices.Equal(dst, want) {
		t.Errorf("dst = %v, want %v", dst, want)
	}
}

func TestBlitRGBA(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 2))
	src.SetRGBA(2, 1, color.RGBA{R: 10, G: 20, B: 30, A: 40})
	src.SetRGBA(3, 1, color.RGBA{R: 1, G: 2, B: 3, A: 4})
	dst := make([]byte, 2*1*4)
	blitRGBA(dst, src, image.Rect(2, 1, 4, 2))
	if want := []byte{10, 20, 30, 40, 1, 2, 3, 4}; !slices.Equal(dst, want) {
		t.Errorf("dst = %v, want %v", dst, want)
	}
}
