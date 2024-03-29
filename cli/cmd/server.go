package cmd

import (
	"github.com/figment-networks/avalanche-indexer/api"
	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/sirupsen/logrus"
)

type ServerCommand struct {
	db     *store.DB
	addr   string
	logger *logrus.Logger
	rpc    *client.Client
}

func NewServerCommand(db *store.DB, addr string, logger *logrus.Logger, rpc *client.Client) ServerCommand {
	return ServerCommand{
		db:     db,
		addr:   addr,
		logger: logger,
		rpc:    rpc,
	}
}

func (cmd ServerCommand) Run() error {
	cmd.logger.Info("starting http server on ", cmd.addr)

	server := api.NewServer(cmd.db, cmd.rpc, cmd.logger)
	return server.Run(cmd.addr)
}
