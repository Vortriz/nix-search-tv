package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents configuration options stored in the
// config file
type Config struct {
	UpdateInterval       Duration `json:"update_interval"`
	CacheDir             string   `json:"cache_dir"`
	EnableWaitingMessage bool     `json:"enable_waiting_message"`
	Indexes              []string `json:"indexes"`
}

type config struct {
	UpdateInterval       *Duration `json:"update_interval"`
	CacheDir             *string   `json:"cache_dir"`
	EnableWaitingMessage *bool     `json:"enable_waiting_message"`
	Indexes              *[]string `json:"indexes"`
}

// Keep the constants below in sync with the `Config` json tags
const (
	UpdateIntervalTag       = "update_interval"
	EnableWaitingMessageTag = "enable_waiting_message"
)

func LoadPath(path string) (Config, error) {
	var err error
	if path == "" {
		path, err = defaultConfigDir()
		if err != nil {
			return Config{}, fmt.Errorf("get default config path: %w", err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	loaded := config{}
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}

	conf := defaults()
	if loaded.CacheDir != nil {
		conf.CacheDir = *loaded.CacheDir
	}
	if loaded.UpdateInterval != nil {
		conf.UpdateInterval = *loaded.UpdateInterval
	}
	if loaded.Indexes != nil {
		conf.Indexes = *loaded.Indexes
	}
	if loaded.EnableWaitingMessage != nil {
		conf.EnableWaitingMessage = *loaded.EnableWaitingMessage
	}

	return conf, nil
}

func defaults() Config {
	cacheDir, err := defaultCacheDir()
	if err != nil {
		// The code reach this point only if
		// $HOME is not defined. And my question is,
		// what environment does not have $HOME?
		panic(err)
	}
	return Config{
		UpdateInterval:       Duration(time.Hour * 24 * 7),
		CacheDir:             cacheDir,
		EnableWaitingMessage: true,
		// TODO: use constants from the indexes package
		Indexes: []string{"nixpkgs"},
	}
}

func defaultCacheDir() (string, error) {
	var err error
	// Check xdg first because `os.UserCacheDir`
	// ignores XDG_CACHE_HOME on darwin
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		cacheDir, err = os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user cache dir: %w", err)
		}
	}

	return filepath.Join(cacheDir, "nix-search-tv"), nil
}

func defaultConfigDir() (string, error) {
	var err error
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("cannot get user config dir: %w", err)
		}
	}

	path := filepath.Join(configDir, "nix-search-tv", "config.json")
	return path, nil
}

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	dur, err := time.ParseDuration(string(b))
	if err != nil {
		return fmt.Errorf("parse duration string: %w", err)
	}

	*d = Duration(dur)
	return nil
}
