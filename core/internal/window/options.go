package window

import (
	"image"

	"github.com/gz0ni/claude-pet/core/internal/config"
	"github.com/gz0ni/claude-pet/core/internal/sprite"
)

type Options struct {
	Pet         *sprite.Pet
	Sheet       image.Image
	Scale       float64
	Position    config.Position
	AlwaysOnTop bool
	OnMoved     func(config.Position)
}
