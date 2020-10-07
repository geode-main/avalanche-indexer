package indexer

import (
	"fmt"
	"net/url"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer/archiver"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"
)

const (
	taskFetcher   = "fetcher"
	taskParser    = "parser"
	taskPersistor = "persistor"
	taskAnalyzier = "analyzer"
	taskArchiver  = "archiver"
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

func NewArchiverTask(configStr string, db *store.DB, logger *logrus.Logger) pipeline.Task {
	var arc archiver.Archiver

	if configStr != "" {
		uri, err := url.Parse(configStr)
		if err != nil {
			// TODO: dont panic here
			panic(err)
		}

		switch uri.Scheme {
		case "dir":
			dir := fmt.Sprintf("%s%s", uri.Host, uri.Path)

			logger.WithField("dir", dir).Debug("configuring file archiver")

			arc = archiver.NewFileArchiver(dir)
		case "s3":
			region := uri.Host
			bucket := uri.Path[1:]

			logger.
				WithField("region", region).
				WithField("bucket", bucket).
				Debug("configuring s3 archiver")

			arc = archiver.NewS3Archiver(region, bucket)
		}
	}

	// TODO: dont panic here
	if arc != nil {
		if err := arc.Test(); err != nil {
			panic(err)
		}
	}

	return &ArchiverTask{
		logger: logger,
		arc:    arc,
		db:     db,
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
