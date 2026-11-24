package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingReturnsDefaults(t *testing.T) {
	c, err := Load(filepath.Join(t.TempDir(), "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Scale != 0.5 || c.Position.Anchor != "bottom-left" || !c.AlwaysOnTop {
		t.Fatalf("defaults = %+v", c)
	}
}

func TestLoadFillsMissingKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	write(t, path, `{"scale": 0.75, "position": {"x": 5}}`)
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.Scale != 0.75 || c.Position.X != 5 || c.Position.Y != -40 || c.Position.Anchor != "bottom-left" || c.Pet != "chibi-ado" {
		t.Fatalf("loaded = %+v", c)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	write(t, path, `{"scale":`)
	c, err := Load(path)
	if err == nil {
		t.Fatal("expected error")
	}
	if c.Scale != Default().Scale {
		t.Fatalf("expected defaults, got %+v", c)
	}
}

func TestSavePreservesUnknownKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	write(t, path, `{"scale": 0.5, "future": {"a": 1}, "position": {"x": 1, "extra": true}, "big": 12345678901234567890}`)
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	c.Position.X = 99
	if err := Save(path, c); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if string(m["big"]) != "12345678901234567890" {
		t.Errorf("big = %s", m["big"])
	}
	var pos map[string]any
	if err := json.Unmarshal(m["position"], &pos); err != nil {
		t.Fatal(err)
	}
	if pos["x"] != float64(99) || pos["extra"] != true || pos["anchor"] != "bottom-left" {
		t.Errorf("position = %v", pos)
	}
	if _, ok := m["future"]; !ok {
		t.Error("unknown top-level key dropped")
	}
	entries, _ := os.ReadDir(filepath.Dir(path))
	if len(entries) != 1 {
		t.Errorf("temp files left: %d entries", len(entries))
	}
}

func TestSaveCreatesDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.json")
	if err := Save(path, Default()); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.Position != Default().Position {
		t.Fatalf("round trip = %+v", c)
	}
}

func write(t *testing.T, path, s string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}
