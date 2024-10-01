package plugins

import (
	"context"
	"kinsyn/pkg/commons"
	"kinsyn/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type InputPlugin interface {
	SyncHighlights() ([]commons.Highlight, error)
}

type OutputPlugin interface {
	SendNotification([]commons.Highlight) error
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "kinsyn",
	MagicCookieValue: "kinsyn",
}

var PluginMap = map[string]plugin.Plugin{
	"input": &InputPluginGRPC{},
}

type InputPluginGRPC struct {
	plugin.Plugin
	Impl InputPlugin
}

func (p *InputPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterInputPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *InputPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewInputPluginClient(c)}, nil
}
