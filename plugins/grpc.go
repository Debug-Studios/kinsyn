package plugins

import (
	"context"
	"kinsyn/proto"
)

type GRPCClient struct {
	client proto.InputPluginClient
}

func (s *GRPCClient) SyncHighlights(ctx context.Context, req *proto.Empty) (*proto.HighlightList, error) {
	return s.client.SyncHighlights(ctx, req)
}

type GRPCServer struct {
	Impl InputPlugin
}

func (s *GRPCServer) SyncHighlights(ctx context.Context, req *proto.Empty) (*proto.HighlightList, error) {
	highlights, err := s.Impl.SyncHighlights()
	if err != nil {
		return nil, err
	}

	protoHighlights := make([]*proto.Highlight, len(highlights))
	for i, h := range highlights {
		protoHighlights[i] = &proto.Highlight{
			BookTitle:         h.BookTitle,
			BookAuthor:        h.BookAuthor,
			BookLocationStart: int32(h.BookLocationStart),
			BookLocationEnd:   int32(h.BookLocationEnd),
			CreatedAt:         h.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Content:           h.Content,
		}
	}

	return &proto.HighlightList{Highlights: protoHighlights}, nil
}
