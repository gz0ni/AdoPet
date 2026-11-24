// Command pet-core is the always-running desktop pet process.
package main

import (
	"log/slog"
	"os"

	"github.com/gz0ni/claude-pet/core/assets"
	"github.com/gz0ni/claude-pet/core/internal/paths"
)

var version = "0.0.0-dev"

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	dir, err := paths.DataDir()
	if err != nil {
		log.Error("data dir", "err", err)
		os.Exit(1)
	}
	log.Info("pet-core skeleton", "version", version, "dataDir", dir, "sheetBytes", len(assets.Spritesheet))
	// M1: window; M2: server + state; M3: tray, hotkeys, autostart.
}
