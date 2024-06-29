package main

import (
	"kinsyn/pkg/config"

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
		// Do something with the new config.
	})

	log.Info().Msgf("Loaded config: %+v", cfg)

}
