package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Worker struct {
	db             *store.DB
	rpc            *client.Client
	logger         *logrus.Logger
	syncInterval   time.Duration
	purgeInterval  time.Duration
	archiverConfig string
}

func NewWorkerCommand(
	db *store.DB,
	rpc *client.Client,
	logger *logrus.Logger,
	interval time.Duration,
	purgeInterval time.Duration,
	archiverConfig string) Worker {
	return Worker{
		db:             db,
		rpc:            rpc,
		logger:         logger,
		syncInterval:   interval,
		purgeInterval:  purgeInterval,
		archiverConfig: archiverConfig,
	}
}

func (cmd Worker) Run() error {
	pipeline, err := indexer.NewPipeline(cmd.db, cmd.rpc, cmd.logger, cmd.archiverConfig)
	if err != nil {
		return err
	}

	cmd.logger.WithField("interval", cmd.syncInterval).Info("starting worker")
	defer cmd.logger.Info("stopping worker")

	busy := false

	syncTicker := time.NewTicker(cmd.syncInterval)
	purgeTicker := time.NewTicker(cmd.purgeInterval)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			cmd.logger.WithField("signal", sig).Info("stopping worker")
			return nil
		case <-purgeTicker.C:
			if err := runPurge(cmd.db, cmd.logger); err != nil {
				cmd.logger.WithError(err).Error("purge failed")
			}
		case <-syncTicker.C:
			if busy {
				cmd.logger.Debug("sync is already in progress, skipping run")
				continue
			}

			busy = true
			go func() {
				cmd.logger.Info("starting sync")

				ts := time.Now()
				defer func() {
					busy = false
					cmd.logger.WithField("duration", time.Since(ts).Milliseconds()).Info("sync finished")
				}()

				if err := pipeline.Start(); err != nil {
					cmd.logger.WithError(err).Error("pipeline error")
				}
			}()
		}
	}
}
