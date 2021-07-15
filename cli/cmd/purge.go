package cmd

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/sirupsen/logrus"
)

type PurgeCommand struct {
	db     *store.DB
	logger *logrus.Logger
}

func NewPurgeCommand(db *store.DB, logger *logrus.Logger) PurgeCommand {
	return PurgeCommand{
		db:     db,
		logger: logger,
	}
}

func (cmd PurgeCommand) Run() error {
	return runPurge(cmd.db, cmd.logger)
}

func runPurge(db *store.DB, logger *logrus.Logger) error {
	// TODO: make this configurable
	before := time.Unix(time.Now().Unix()-172800, 0)

	logger.WithField("before_time", before).Info("purging validators")
	num, err := db.Validators.PurgeSeq(before)
	if err != nil {
		return err
	}
	logger.WithField("count", num).Info("purged validator records")

	return nil
}
