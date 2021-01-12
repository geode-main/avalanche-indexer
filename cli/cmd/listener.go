package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Listener struct {
	db     *store.DB
	rpc    *client.Client
	logger *logrus.Logger
}

func NewListenerCommand(
	db *store.DB,
	rpc *client.Client,
	logger *logrus.Logger,
) Listener {
	return Listener{
		db:     db,
		rpc:    rpc,
		logger: logger,
	}
}

func (cmd Listener) Run() error {
	cmd.logger.Info("fetching availanle blockchains...")
	blockchains, err := cmd.rpc.Platform.GetBlockchains()
	if err != nil {
		return err
	}

	for _, c := range blockchains.Blockchains {
		cmd.logger.Info(c)
	}

	return nil
}
