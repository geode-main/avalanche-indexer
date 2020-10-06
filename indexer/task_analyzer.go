package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/store"
)

type AnalyzerTask struct {
	db     *store.DB
	logger *logrus.Logger
}

func (t AnalyzerTask) GetName() string {
	return taskAnalyzier
}

func (t AnalyzerTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	// TODO: this does not do anything right now, but the goal is to have the events
	// implemented as part of this task.

	return nil
}
