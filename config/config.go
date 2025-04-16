package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/3timeslazy/nix-search-tv/indexes/indices"
)

// Config represents configuration options stored in the
// config file
type Config struct {
	UpdateInterval       Duration     `json:"update_interval"`
	CacheDir             string       `json:"cache_dir"`
	EnableWaitingMessage bool         `json:"enable_waiting_message"`
	Indexes              []string     `json:"indexes"`
	Experimental         Experimental `json:"experimental"`
}

type config struct {
	UpdateInterval       *Duration    `json:"update_interval"`
	CacheDir             *string      `json:"cache_dir"`
	EnableWaitingMessage *bool        `json:"enable_waiting_message"`
	Indexes              *[]string    `json:"indexes"`
	Experimental         Experimental `json:"experimental"`
}

type Experimental struct {
	RenderDocsIndexes map[string]string `json:"render_docs_indexes"`
}

// Keep the constants below in sync with the `Config` json tags
const (
	UpdateIntervalTag       = "update_interval"
	EnableWaitingMessageTag = "enable_waiting_message"
)

func LoadDefault() (Config, error) {
	path, err := defaultConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("get default config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		data = []byte("{}")
		err = nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	loaded := config{}
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}

	return mergeDefaults(loaded), nil
}

func LoadPath(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	loaded := config{}
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}

	return mergeDefaults(loaded), nil
}

func mergeDefaults(loaded config) Config {
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

	conf.Experimental = Experimental{
		RenderDocsIndexes: loaded.Experimental.RenderDocsIndexes,
	}

	return conf
}

func defaults() Config {
	cacheDir, err := defaultCacheDir()
	if err != nil {
		// The code reach this point only if
		// $HOME is not defined. And my question is,
		// what environment does not have $HOME?
		panic(err)
	}

	indexes := []string{indices.Nixpkgs, indices.HomeManager, indices.Nur}
	if runtime.GOOS == "linux" {
		indexes = append(indexes, indices.NixOS)
	}
	if runtime.GOOS == "darwin" {
		indexes = append(indexes, indices.Darwin)
	}

	return Config{
		UpdateInterval:       Duration(time.Hour * 24 * 7),
		CacheDir:             cacheDir,
		EnableWaitingMessage: true,
		Indexes:              indexes,
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
