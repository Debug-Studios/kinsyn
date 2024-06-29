package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	yamlv2 "gopkg.in/yaml.v2"
)

var k = koanf.New(".")

type Config struct {
	InputPlugins  []string `koanf:"input_plugins"`
	OutputPlugins []string `koanf:"output_plugins"`
}

// LoadConfig loads the configuration from a file and environment variables.
func LoadConfig(filename string) (*Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warn().Msg("Configuration file not found, creating default configuration file")
		if err := createDefaultConfig(filename); err != nil {
			return nil, err
		}
	}

	// Load the config from the YAML file.
	f := file.Provider(filename)
	if err := k.Load(f, yaml.Parser()); err != nil {
		return nil, err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// WatchConfig watches the configuration file for changes and reloads it.
func WatchConfig(filename string, onReload func(cfg *Config)) {
	f := file.Provider(filename)
	if err := f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Error().Msgf("Error watching config file: %v", err)
			return
		}
		log.Print("Configuration file changed, reloading...")
		cfg, err := LoadConfig(filename)
		if err != nil {
			log.Error().Msgf("Failed to reload configuration: %v", err)
		} else {
			onReload(cfg)
		}
	}); err != nil {
		log.Error().Msgf("Failed to add file watcher: %v", err)
	}
}

func createDefaultConfig(filename string) error {
	defaultConfig := &Config{
		InputPlugins:  []string{"plugins/input/usb-sync.so"},
		OutputPlugins: []string{"plugins/output/email.so"},
	}

	data, err := yamlv2.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
