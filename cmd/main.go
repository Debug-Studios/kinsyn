package main

import (
	"kinsyn/pkg/config"
	"kinsyn/plugins/input/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configFilePath := config.GetConfigFilePath()

	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	// Watch the config file for changes.
	go config.WatchConfig(configFilePath, func(cfg *config.Config) {
		log.Info().Msgf("Config reloaded: %+v", cfg)
	})

	log.Info().Msgf("Loaded config: %+v", cfg)

	filepathPlugin := &filepath.FilePathPlugin{}

	highlights, err := filepathPlugin.SyncHighlights()
	if err != nil {
		log.Error().Msgf("Failed to sync highlights: %v", err)
	}

	log.Info().Msgf("Highlights: %v", highlights[0])

	// Keep the main goroutine alive.
	select {}

	// TODO: Load some sort of frontend for the app.

}
