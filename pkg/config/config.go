package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	yamlv3 "gopkg.in/yaml.v3"
)

var k = koanf.New(".")

type PluginType int

const (
	InputPluginType PluginType = iota
	OutputPluginType
)

func (p PluginType) String() string {
	return [...]string{"input", "output"}[p]
}

type Plugin struct {
	Path   string                 `koanf:"path" yaml:"path"`
	Config map[string]interface{} `koanf:"config" yaml:"config"`
}

type PluginConfig struct {
	InputPlugins  map[string]Plugin `koanf:"input_plugins" yaml:"input_plugins"`
	OutputPlugins map[string]Plugin `koanf:"output_plugins" yaml:"output_plugins"`
}

type Config struct {
	Plugins PluginConfig `koanf:"plugins" yaml:"plugins"`
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

// SaveConfig saves the configuration to a file.
func SaveConfig(filename string, cfg *Config) error {
	data, err := yamlv3.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func createDefaultConfig(filename string) error {
	defaultConfig := &Config{
		Plugins: PluginConfig{
			InputPlugins:  map[string]Plugin{},
			OutputPlugins: map[string]Plugin{},
		},
	}

	return SaveConfig(filename, defaultConfig)
}

// GetPluginConfig gets the configuration for a specific plugin.
func GetPluginConfig(pluginType PluginType, pluginName string) (map[string]interface{}, error) {
	cfg, err := LoadConfig(GetConfigFilePath())
	if err != nil {
		return nil, err
	}

	var pluginConfig map[string]interface{}
	if pluginType == InputPluginType {
		pluginConfig = cfg.Plugins.InputPlugins[pluginName].Config
	} else if pluginType == OutputPluginType {
		pluginConfig = cfg.Plugins.OutputPlugins[pluginName].Config
	}

	return pluginConfig, nil
}

// SetPluginConfig sets the configuration for a specific plugin.
func SetPluginConfig(pluginType PluginType, pluginName string, pluginCfg map[string]interface{}) error {
	cfg, err := LoadConfig(GetConfigFilePath())
	if err != nil {
		return err
	}

	if pluginType == InputPluginType {
		plugin := cfg.Plugins.InputPlugins[pluginName]
		plugin.Config = pluginCfg
		cfg.Plugins.InputPlugins[pluginName] = plugin
	} else if pluginType == OutputPluginType {
		plugin := cfg.Plugins.OutputPlugins[pluginName]
		plugin.Config = pluginCfg
		cfg.Plugins.OutputPlugins[pluginName] = plugin
	}

	return SaveConfig(GetConfigFilePath(), cfg)
}
