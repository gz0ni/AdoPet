package sprite

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"math"
	"time"

	"golang.org/x/image/draw"
)

type Cell struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Animation struct {
	Row    int    `json:"row"`
	Rows   []int  `json:"rows,omitempty"`
	Frames int    `json:"frames"`
	FPS    int    `json:"fps"`
	Loop   bool   `json:"loop"`
	Next   string `json:"next,omitempty"`
	Mode   string `json:"mode,omitempty"`
}

type Pet struct {
	Name       string               `json:"name"`
	Sheet      string               `json:"sheet"`
	Cell       Cell                 `json:"cell"`
	Animations map[string]Animation `json:"animations"`
}

func ParsePet(data []byte) (*Pet, error) {
	var p Pet
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse pet.json: %w", err)
	}
	if p.Cell.W <= 0 || p.Cell.H <= 0 {
		return nil, errors.New("pet.json: cell size must be positive")
	}
	if _, ok := p.Animations["idle"]; !ok {
		return nil, errors.New("pet.json: idle animation is required")
	}
	return &p, nil
}

func FrameRect(a Animation, i int, cell Cell) image.Rectangle {
	row, col := a.Row, i
	if len(a.Rows) > 0 {
		per := a.Frames / len(a.Rows)
		if per <= 0 {
			per = 1
		}
		idx := min(i/per, len(a.Rows)-1)
		row, col = a.Rows[idx], i-idx*per
	}
	x, y := col*cell.W, row*cell.H
	return image.Rect(x, y, x+cell.W, y+cell.H)
}

func FrameAt(a Animation, elapsed time.Duration) int {
	if a.Frames <= 0 || elapsed < 0 {
		return 0
	}
	fps := max(a.FPS, 1)
	idx := int(elapsed * time.Duration(fps) / time.Second)
	if a.Loop {
		return idx % a.Frames
	}
	return min(idx, a.Frames-1)
}

func ScaledCell(cell Cell, scale float64) Cell {
	return Cell{
		W: max(1, int(math.Round(float64(cell.W)*scale))),
		H: max(1, int(math.Round(float64(cell.H)*scale))),
	}
}

func ScaleSheet(src image.Image, cell Cell, scale float64) (*image.RGBA, Cell) {
	b := src.Bounds()
	cols, rows := b.Dx()/cell.W, b.Dy()/cell.H
	dc := ScaledCell(cell, scale)
	dst := image.NewRGBA(image.Rect(0, 0, cols*dc.W, rows*dc.H))
	for r := range rows {
		for c := range cols {
			sr := image.Rect(c*cell.W, r*cell.H, (c+1)*cell.W, (r+1)*cell.H).Add(b.Min)
			dr := image.Rect(c*dc.W, r*dc.H, (c+1)*dc.W, (r+1)*dc.H)
			draw.CatmullRom.Scale(dst, dr, src, sr, draw.Src, nil)
		}
	}
	return dst, dc
}
