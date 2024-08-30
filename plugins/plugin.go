package plugins

type Plugin interface {
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
}

type InputPlugin interface {
	Plugin
	SyncHighlights() ([]string, error)
}

type OutputPlugin interface {
	Plugin
	SendNotification(string) error
}
