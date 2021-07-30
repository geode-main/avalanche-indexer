package cmd

import (
	"math/big"

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

type SyncCommand struct {
	networkID  uint32
	evmChainID uint32

	logger *logrus.Logger
	db     *store.DB
	rpc    *client.Client
}

type syncWorker interface {
	ProcessMessage(*model.RawMessage) error
}

func NewSyncCommand(logger *logrus.Logger, db *store.DB, rpc *client.Client, networkID uint32, evmChainID uint32) SyncCommand {
	return SyncCommand{
		networkID:  networkID,
		evmChainID: evmChainID,
		logger:     logger,
		db:         db,
		rpc:        rpc,
	}
}

func (cmd SyncCommand) Run() error {
	// show all the debug details when running a one-off sync command
	if cmd.logger.Level != logrus.DebugLevel {
		cmd.logger.SetLevel(logrus.DebugLevel)
	}

	cmd.logger.Info("starting sync")
	defer cmd.logger.Info("finished sync")

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

	pipeline, err := indexer.NewPipeline(cmd.db, cmd.rpc, cmd.logger)
	if err != nil {
		return err
	}

	if err := pipeline.Start(); err != nil {
		return err
	}

	avmWorker := avm.NewWorker(&cmd.rpc.Index, cmd.db, codec.AVM, xID, assetID.String())
	pvmWorker := pvm.NewWorker(&cmd.rpc.Index, cmd.db, codec.PVM, pID, assetID.String())
	cvmWorker := cvm.NewWorker(cmd.db, codec.EVM, &cmd.rpc.Index, &cmd.rpc.Evm, cID, assetID.String(), big.NewInt(int64(cmd.evmChainID)))
	pblocksWorker := blocks.NewWorker(cmd.db, cmd.rpc, cmd.logger, pID)

	if err = avmWorker.Run(); err != nil {
		return err
	}

	if err = pvmWorker.Run(); err != nil {
		return err
	}

	if err = cvmWorker.Run(); err != nil {
		return err
	}

	if pblocksWorker.Run(); err != nil {
		return err
	}

	return nil
}
