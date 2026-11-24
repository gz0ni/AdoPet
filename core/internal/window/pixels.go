package window

import "image"

func blitBGRA(dst []byte, src *image.RGBA, r image.Rectangle) {
	w := r.Dx()
	for y := range r.Dy() {
		s := src.Pix[src.PixOffset(r.Min.X, r.Min.Y+y):]
		d := dst[y*w*4:]
		for x := range w {
			i := x * 4
			d[i], d[i+1], d[i+2], d[i+3] = s[i+2], s[i+1], s[i], s[i+3]
		}
	}
}
