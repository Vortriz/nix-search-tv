package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	UpdateInterval Duration `json:"update_interval"`
	CacheDir       string   `json:"cache_dir"`
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

func LoadPath(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	return conf, json.Unmarshal(data, &conf)
}

func Default() Config {
	cacheDir, err := defaultCacheDir()
	if err != nil {
		// The code reach this point only if
		// $HOME is not defined. And my question is,
		// what environment does not have $HOME?
		panic(err)
	}
	return Config{
		UpdateInterval: Duration(time.Hour * 24),
		CacheDir:       cacheDir,
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

func ConfigDir() (string, error) {
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
