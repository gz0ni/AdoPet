package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const SchemaVersion = 1

type Position struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Anchor  string `json:"anchor"`
	Display int    `json:"display"`
}

type Config struct {
	Version     int      `json:"version"`
	Pet         string   `json:"pet"`
	Scale       float64  `json:"scale"`
	Position    Position `json:"position"`
	AlwaysOnTop bool     `json:"alwaysOnTop"`

	raw map[string]any
}

func Default() Config {
	return Config{
		Version:     SchemaVersion,
		Pet:         "chibi-ado",
		Scale:       0.5,
		Position:    Position{X: 40, Y: -40, Anchor: "bottom-left"},
		AlwaysOnTop: true,
	}
}

func Load(path string) (Config, error) {
	c := Default()
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return Default(), fmt.Errorf("parse %s: %w", path, err)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&c.raw); err != nil {
		return Default(), fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Scale <= 0 {
		c.Scale = Default().Scale
	}
	return c, nil
}

func Save(path string, c Config) error {
	known, err := toMap(c)
	if err != nil {
		return err
	}
	out := merge(clone(c.raw), known)
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(path, append(data, '\n'))
}

func toMap(c Config) (map[string]any, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	err = dec.Decode(&m)
	return m, err
}

func merge(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = map[string]any{}
	}
	for k, v := range src {
		sm, sok := v.(map[string]any)
		dm, dok := dst[k].(map[string]any)
		if sok && dok {
			dst[k] = merge(dm, sm)
		} else {
			dst[k] = v
		}
	}
	return dst
}

func clone(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			v = clone(vm)
		}
		out[k] = v
	}
	return out
}

func writeAtomic(path string, data []byte) (err error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".config-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(tmp)
		}
	}()
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
