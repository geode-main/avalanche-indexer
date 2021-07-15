package indexer

import (
	"context"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"
)

type indexingPipeline struct {
	db       *store.DB
	rpc      *client.Client
	log      *logrus.Logger
	pipeline pipeline.CustomPipeline
}

func NewPipeline(db *store.DB, rpc *client.Client, logger *logrus.Logger) (*indexingPipeline, error) {
	p := pipeline.NewCustom(NewPayloadFactory())
	p.SetLogger(NewLogger(logger))

	fetcherStage := pipeline.NewStageWithTasks(
		pipeline.StageFetcher,
		NewFetcherTask(rpc, logger), // fetch data from the network
	)

	parserStage := pipeline.NewStageWithTasks(
		pipeline.StageParser,
		NewParserTask(logger), // map all client data to the indexer models
	)

	persistorStage := pipeline.NewStageWithTasks(
		pipeline.StagePersistor,
		NewPersistorTask(db, logger), // save stuff into db
	)

	cleanupStage := pipeline.NewStageWithTasks(
		pipeline.StageCleanup,
		NewCleanupTask(db, logger), // internal cleanup, etc
	)

	p.AddStage(fetcherStage)
	p.AddStage(parserStage)
	p.AddStage(persistorStage)
	p.AddStage(cleanupStage)

	return &indexingPipeline{
		db:       db,
		rpc:      rpc,
		log:      logger,
		pipeline: p,
	}, nil
}

func (p *indexingPipeline) Start() error {
	source := &indexSource{}
	sink := &indexSink{}

	return p.pipeline.Start(
		context.Background(),
		source,
		sink,
		&pipeline.Options{},
	)
}
