package main

import (
	"fmt"
	"kinsyn/pkg/commons"
	"kinsyn/pkg/config"
	"kinsyn/pkg/parser"
	"kinsyn/plugins"
	"os"

	"github.com/hashicorp/go-plugin"
)

type FilePathPlugin struct{}

var _ plugins.InputPlugin = &FilePathPlugin{}

const PluginName = "filepath"

func (f *FilePathPlugin) SyncHighlights() ([]commons.Highlight, error) {
	pluginConf, err := config.GetPluginConfig(config.InputPluginType, PluginName)
	if err != nil {
		config.SetPluginConfig(config.InputPluginType, PluginName, map[string]interface{}{"path": "/Users/hd/Downloads/highlights.txt"})
		pluginConf, err = config.GetPluginConfig(config.InputPluginType, PluginName)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("pluginConf: %v\n", pluginConf)

	path, ok := pluginConf["path"].(string)
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

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"input": &plugins.InputPluginGRPC{Impl: &FilePathPlugin{}},
		},

		GRPCServer: plugin.DefaultGRPCServer,
	})
}
