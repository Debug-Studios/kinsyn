package main

import (
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}

/*
pluginConf, err := config.GetPluginConfig(config.InputPluginType, PluginName)
	if err != nil || pluginConf == nil {
		config.SetPluginConfig(config.InputPluginType, PluginName, map[string]interface{}{"path": "/Users/hd/Downloads/highlights.txt"})
		pluginConf, err = config.GetPluginConfig(config.InputPluginType, PluginName)
		if err != nil {
			panic(err)
		}
	}

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
*/
