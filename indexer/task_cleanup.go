package indexer

import (
	"context"

	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"
)

type CleanupTask struct {
	db     *store.DB
	logger *logrus.Logger
}

func (t CleanupTask) GetName() string {
	return taskCleanup
}

func (t CleanupTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	return nil
}
