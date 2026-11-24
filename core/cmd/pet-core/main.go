// Command pet-core is the always-running desktop pet process.
package main

import (
	"bytes"
	"image/png"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/gz0ni/claude-pet/core/assets"
	"github.com/gz0ni/claude-pet/core/internal/config"
	"github.com/gz0ni/claude-pet/core/internal/paths"
	"github.com/gz0ni/claude-pet/core/internal/sprite"
	"github.com/gz0ni/claude-pet/core/internal/window"
)

var version = "0.0.0-dev"

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(log); err != nil {
		log.Error("pet-core", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	debug.SetMemoryLimit(48 << 20)

	cfgPath, err := paths.ConfigFile()
	if err != nil {
		return err
	}
	cfg, err := config.Load(cfgPath)
	persist := err == nil
	if err != nil {
		log.Warn("config not loaded, using defaults; changes will not be saved", "err", err)
	}

	pet, err := sprite.ParsePet(assets.PetJSON)
	if err != nil {
		return err
	}
	sheet, err := png.Decode(bytes.NewReader(assets.Spritesheet))
	if err != nil {
		return err
	}

	log.Info("pet-core starting", "version", version, "config", cfgPath)
	return window.Run(window.Options{
		Pet:         pet,
		Sheet:       sheet,
		Scale:       cfg.Scale,
		Position:    cfg.Position,
		AlwaysOnTop: cfg.AlwaysOnTop,
		OnMoved: func(p config.Position) {
			cfg.Position = p
			if !persist {
				return
			}
			if err := config.Save(cfgPath, cfg); err != nil {
				log.Error("save config", "err", err)
			}
		},
	})
}
