package filepath

import (
	"kinsyn/pkg/commons"
	"kinsyn/pkg/config"
	"kinsyn/pkg/parser"
	"kinsyn/plugins"
	"os"
)

type FilePathPlugin struct{}

var _ plugins.InputPlugin = &FilePathPlugin{}

const PluginName = "filepath"

func (f *FilePathPlugin) GetConfig() map[string]interface{} {
	config, err := config.GetPluginConfig(config.InputPluginType, PluginName)
	if err != nil {
		panic(err)
	}

	return config
}

func (f *FilePathPlugin) SetConfig(c map[string]interface{}) {
	err := config.SetPluginConfig(config.InputPluginType, PluginName, c)
	if err != nil {
		panic(err)
	}
}

func (f *FilePathPlugin) SyncHighlights() ([]commons.Highlight, error) {
	config := f.GetConfig()
	path, ok := config["path"].(string)
	if !ok {
		panic("invalid config")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	highlights, err := parser.ParseHighlights(file)
	if err != nil {
		return nil, err
	}

	return highlights, nil
}
