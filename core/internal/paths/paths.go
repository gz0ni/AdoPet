package paths

import (
	"os"
	"path/filepath"
)

const appName = "ClaudePet"

// DataDir returns %APPDATA%\ClaudePet on Windows and
// ~/Library/Application Support/ClaudePet on macOS.
func DataDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

// ConfigFile returns the path to config.json.
func ConfigFile() (string, error) {
	d, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "config.json"), nil
}
