package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
)

// dummy implementation since there are no heights available in ava
type indexSource struct {
}

func (s *indexSource) Next(context.Context, pipeline.Payload) bool {
	return false
}

func (s *indexSource) Current() int64 {
	return 0
}

func (s *indexSource) Err() error {
	return nil
}

func (s *indexSource) Len() int64 {
	return 0
}

func (s *indexSource) Skip(stageName pipeline.StageName) bool {
	return false
}
