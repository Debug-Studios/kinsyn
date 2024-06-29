package plugins

type Plugin interface{}

type InputPlugin interface {
	SyncHighlights() ([]string, error)
}

type OutputPlugin interface {
	SendNotification(string) error
}
