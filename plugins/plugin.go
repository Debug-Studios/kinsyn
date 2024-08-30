package plugins

import "kinsyn/pkg/commons"

type Plugin interface {
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
}

type InputPlugin interface {
	Plugin
	SyncHighlights() ([]commons.Highlight, error)
}

type OutputPlugin interface {
	Plugin
	SendNotification([]commons.Highlight) error
}
