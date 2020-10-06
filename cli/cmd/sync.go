package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Sync struct {
	logger         *logrus.Logger
	db             *store.DB
	rpc            *client.Client
	archiverConfig string
}

func NewSyncCommand(logger *logrus.Logger, db *store.DB, rpc *client.Client, archiverConfig string) Sync {
	return Sync{
		logger:         logger,
		db:             db,
		rpc:            rpc,
		archiverConfig: archiverConfig,
	}
}

func (cmd Sync) Run() error {
	// show all the debug details when running a one-off sync command
	if cmd.logger.Level != logrus.DebugLevel {
		cmd.logger.SetLevel(logrus.DebugLevel)
	}

	cmd.logger.Info("starting sync")
	defer cmd.logger.Info("finished sync")

	pipeline, err := indexer.NewPipeline(cmd.db, cmd.rpc, cmd.logger, cmd.archiverConfig)
	if err != nil {
		return err
	}

	return pipeline.Start()
}
