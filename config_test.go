package i2ptui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	p := DefaultConfigPath()
	if p == "" {
		t.Error("expected non-empty path")
	}
	if filepath.Base(p) != "config.json" {
		t.Errorf("expected config.json, got %s", filepath.Base(p))
	}
}

func TestLoadConfigMissing(t *testing.T) {
	cfg := LoadConfig("/nonexistent/path/config.json")
	if cfg.Host != "" {
		t.Errorf("expected empty host for missing file, got %q", cfg.Host)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.json")

	want := Config{
		Host:     "10.0.0.1",
		Port:     "7658",
		Password: "secret",
		Theme:    "light",
	}
	if err := SaveConfig(path, want); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	got := LoadConfig(path)
	if got.Host != want.Host {
		t.Errorf("host: got %q, want %q", got.Host, want.Host)
	}
	if got.Port != want.Port {
		t.Errorf("port: got %q, want %q", got.Port, want.Port)
	}
	if got.Theme != want.Theme {
		t.Errorf("theme: got %q, want %q", got.Theme, want.Theme)
	}

	// Check permissions.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600 permissions, got %o", info.Mode().Perm())
	}
}
