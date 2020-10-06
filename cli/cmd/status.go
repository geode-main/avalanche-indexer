package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
)

type Status struct {
	rpc    *client.Client
	logger *logrus.Logger
}

func NewStatusCommand(rpc *client.Client, logger *logrus.Logger) Status {
	return Status{
		rpc:    rpc,
		logger: logger,
	}
}

func (cmd Status) Run() error {
	version, err := cmd.rpc.Info.NodeVersion()
	if err != nil {
		cmd.logger.WithError(err).Error("cant fetch node version")
	}

	cmd.logger.Info("node version: ", version)

	return nil
}
