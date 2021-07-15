package indexer

import (
	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"
)

const (
	taskFetcher   = "fetcher"
	taskParser    = "parser"
	taskPersistor = "persistor"
	taskAnalyzier = "analyzer"
	taskCleanup   = "cleanup"
)

func NewFetcherTask(rpc *client.Client, logger *logrus.Logger) pipeline.Task {
	return &FetcherTask{
		rpc:    rpc,
		logger: logger,
	}
}

func NewParserTask(logger *logrus.Logger) pipeline.Task {
	return &ParserTask{
		logger: logger,
	}
}

func NewPersistorTask(db *store.DB, logger *logrus.Logger) pipeline.Task {
	return &PersistorTask{
		db:     db,
		logger: logger,
	}
}

func NewAnalyzerTask(db *store.DB, logger *logrus.Logger) pipeline.Task {
	return &AnalyzerTask{
		db:     db,
		logger: logger,
	}
}

func NewCleanupTask(db *store.DB, logger *logrus.Logger) pipeline.Task {
	return &CleanupTask{
		db:     db,
		logger: logger,
	}
}

func runChain(p *Payload, fns ...func(p *Payload) error) error {
	for _, fn := range fns {
		if err := fn(p); err != nil {
			return err
		}
	}
	return nil
}

func shouldSkipHeight(db *store.DB, p *Payload) (bool, error) {
	// height, err := db.Validators.LastHeight()
	// if err != nil {
	// 	return false, err
	// }
	// return p.HeightChanged(height), nil

	// dont skip heights for now
	return false, nil
}

func logStart(t pipeline.Task, logger *logrus.Logger) {
	logger.WithField("name", t.GetName()).Debug("task started")
}

func logDone(t pipeline.Task, logger *logrus.Logger) {
	logger.WithField("name", t.GetName()).Debug("task finished")
}
