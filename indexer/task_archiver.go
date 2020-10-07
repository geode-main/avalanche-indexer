package indexer

import (
	"context"
	"fmt"

	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/indexer/archiver"
)

type ArchiverTask struct {
	db     *store.DB
	arc    archiver.Archiver
	logger *logrus.Logger
}

func (t ArchiverTask) GetName() string {
	return taskArchiver
}

func (t ArchiverTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	if t.arc == nil {
		t.logger.Debug("archiver is not provided, skipping")
		return nil
	}

	payload := p.(*Payload)

	skip, err := shouldSkipHeight(t.db, payload)
	if err != nil {
		return err
	}
	if skip {
		t.logger.Info("no height changes detected, archiver skipped")
		return nil
	}

	id := fmt.Sprintf("%v", payload.SyncTime.Unix())
	snapshot, err := archiver.NewSnapshot(id)
	if err != nil {
		return err
	}

	snapshot.Meta.AppName = AppName
	snapshot.Meta.AppVersion = AppVersion
	snapshot.Meta.ChainName = "avalanche"
	snapshot.Meta.ChainNetwork = payload.NetworkName
	snapshot.Meta.ChainVersion = payload.NodeVersion
	snapshot.Meta.Time = payload.SyncTime
	snapshot.Meta.Height = &payload.Height

	snapshot.Add("peers_count", len(payload.Peers))
	snapshot.Add("blockchains", payload.Blockchains)
	snapshot.Add("current_validators", payload.CurrentValidators)
	snapshot.Add("current_delegators", payload.CurrentDelegators)
	snapshot.Add("min_stake", payload.MinStake)
	snapshot.Add("tx_fee", payload.RawTxFee)

	t.logger.
		WithField("id", snapshot.ID).
		Debug("saving archiver snapshot")

	if err := t.arc.Commit(snapshot); err != nil {
		t.logger.WithError(err).Error("unable to commit archiver snapshot")
		// dont return err here since we dont want archiver to fail the further tasks
	}

	return nil
}
