package cmd

import (
	"context"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ava-labs/avalanchego/genesis"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/indexer/avm"
	"github.com/figment-networks/avalanche-indexer/indexer/blocks"
	"github.com/figment-networks/avalanche-indexer/indexer/codec"
	"github.com/figment-networks/avalanche-indexer/indexer/cvm"
	"github.com/figment-networks/avalanche-indexer/indexer/pvm"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

type WorkerCommand struct {
	db             *store.DB
	rpc            *client.Client
	logger         *logrus.Logger
	syncInterval   time.Duration
	purgeInterval  time.Duration
	archiverConfig string

	networkID  uint32
	evmChainID uint32
}

func NewWorkerCommand(
	db *store.DB,
	rpc *client.Client,
	logger *logrus.Logger,
	interval time.Duration,
	purgeInterval time.Duration,
	networkID uint32,
	evmChainID uint32,
) WorkerCommand {
	return WorkerCommand{
		db:            db,
		rpc:           rpc,
		logger:        logger,
		syncInterval:  interval,
		purgeInterval: purgeInterval,
		networkID:     networkID,
		evmChainID:    evmChainID,
	}
}

func (cmd WorkerCommand) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := cmd.startChainWorkers(ctx); err != nil {
			cmd.logger.WithError(err).Error("chain workers failed")
		}
	}()

	go func() {
		defer wg.Done()
		if err := cmd.startPipelineWorker(ctx); err != nil {
			cmd.logger.WithError(err).Error("pipeline worker failed")
		}
	}()

	go func() {
		s := <-initSignals()
		cmd.logger.Info("received signal: ", s)
		cancel()
	}()

	wg.Wait()
	return nil
}

func (cmd WorkerCommand) startChainWorkers(ctx context.Context) error {
	_, assetID, err := genesis.Genesis(cmd.networkID, "")
	if err != nil {
		return err
	}

	pID, err := cmd.rpc.Info.BlockchainID("P")
	if err != nil {
		return err
	}

	xID, err := cmd.rpc.Info.BlockchainID("X")
	if err != nil {
		return err
	}

	cID, err := cmd.rpc.Info.BlockchainID("C")
	if err != nil {
		return err
	}

	cmd.db.Assets.Create(&model.Asset{
		AssetID:      assetID.String(),
		Type:         model.AssetTypeFixed,
		Name:         "Avalanche",
		Symbol:       "AVAX",
		Denomination: 9,
	})

	cmd.db.Platform.CreateChain(&model.Chain{ChainID: pID, Name: "P"})
	cmd.db.Platform.CreateChain(&model.Chain{ChainID: xID, Name: "X"})
	cmd.db.Platform.CreateChain(&model.Chain{ChainID: cID, Name: "C"})

	avmWorker := avm.NewWorker(&cmd.rpc.Index, cmd.db, codec.AVM, xID, assetID.String())
	pvmWorker := pvm.NewWorker(&cmd.rpc.Index, cmd.db, codec.PVM, pID, assetID.String())
	cvmWorker := cvm.NewWorker(cmd.db, codec.EVM, &cmd.rpc.Index, &cmd.rpc.Evm, cID, assetID.String(), big.NewInt(int64(cmd.evmChainID)))
	pblocksWorker := blocks.NewWorker(cmd.db, cmd.rpc, cmd.logger, pID)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		avmWorker.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		pvmWorker.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		cvmWorker.Start(ctx)
	}()

	go func() {
		defer wg.Done()
		pblocksWorker.Start(ctx)
	}()

	wg.Wait()

	return nil
}

func (cmd WorkerCommand) startPipelineWorker(ctx context.Context) error {
	pipeline, err := indexer.NewPipeline(cmd.db, cmd.rpc, cmd.logger)
	if err != nil {
		return err
	}

	cmd.logger.WithField("interval", cmd.syncInterval).Info("starting worker")
	defer cmd.logger.Info("stopping worker")

	busy := false

	syncTicker := time.NewTicker(cmd.syncInterval)
	purgeTicker := time.NewTicker(cmd.purgeInterval)

	for {
		select {
		case <-ctx.Done():
			cmd.logger.Info("stopping worker")
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

func initSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	return c
}
