package main

import (
	"kinsyn/plugins"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/go-plugin"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  plugins.HandshakeConfig,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC, plugin.ProtocolNetRPC},
		Cmd:              exec.Command("./filepath"),
		Plugins:          plugins.PluginMap,
	})

	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Error().Msgf("Error creating RPC client: %v", err)
		return
	}

	raw, err := rpcClient.Dispense("input")
	if err != nil {
		log.Error().Msgf("Error dispensing plugin: %v", err)
		return
	}

	inputPlugin := raw.(plugins.InputPlugin)

	highlights, err := inputPlugin.SyncHighlights()
	if err != nil {
		log.Error().Msgf("Error syncing highlights: %v", err)
		return
	}

	log.Info().Msgf("Synced %d highlights", len(highlights))
}
