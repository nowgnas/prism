package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all user-configurable prism settings.
type Config struct {
	Shell         string   `toml:"shell"`
	StartupTabs   int      `toml:"startup_tabs"`
	SidebarWidth  int      `toml:"sidebar_width"`
	PollInterval  int      `toml:"poll_interval_ms"`
	AIProcesses   []string `toml:"ai_processes"`
	Theme         string   `toml:"theme"`
	MouseEnabled  bool     `toml:"mouse_enabled"`
}

// Load reads the config file at the given path.
// If the file doesn't exist, it returns defaults.
func Load(path string) (Config, error) {
	cfg := Defaults()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// First run — write defaults
		if err2 := Write(path, cfg); err2 != nil {
			// Non-fatal: just return defaults
			return cfg, nil
		}
		return cfg, nil
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Write serializes the config to the given path, creating directories as needed.
func Write(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := toml.NewEncoder(f)
	return enc.Encode(cfg)
}

// DefaultPath returns the platform-appropriate config file path.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "prism", "config.toml")
}
