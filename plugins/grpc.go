package plugins

import (
	"context"
	"github.com/Debug-Studios/kinsyn/pkg/commons"
	"github.com/Debug-Studios/kinsyn/proto"
	"time"
)

type GRPCClient struct {
	client proto.InputPluginClient
}

func (c *GRPCClient) SyncHighlights() ([]commons.Highlight, error) {
	highlights, err := c.client.SyncHighlights(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}

	createdAtParsed, err := time.Parse("2006-01-02T15:04:05Z", highlights.Highlights[0].CreatedAt)
	if err != nil {
		return nil, err
	}

	commonsHighlights := make([]commons.Highlight, len(highlights.Highlights))
	for i, h := range highlights.Highlights {
		commonsHighlights[i] = commons.Highlight{
			BookTitle:         h.BookTitle,
			BookAuthor:        h.BookAuthor,
			BookLocationStart: int(h.BookLocationStart),
			BookLocationEnd:   int(h.BookLocationEnd),
			CreatedAt:         createdAtParsed,
			Content:           h.Content,
		}
	}

	return commonsHighlights, nil
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
