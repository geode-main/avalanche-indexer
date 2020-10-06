package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
)

// sink does not do anything at this time
type indexSink struct {
}

func (s indexSink) Consume(ctx context.Context, p pipeline.Payload) error {
	return nil
}
