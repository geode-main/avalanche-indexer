package indexer

import (
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

type PersistorTask struct {
	db     *store.DB
	logger *logrus.Logger
}

func (t PersistorTask) GetName() string {
	return taskPersistor
}

func (t PersistorTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	payload := p.(*Payload)

	skip, err := shouldSkipHeight(t.db, payload)
	if err != nil {
		return err
	}
	if skip {
		t.logger.Info("no height changes detected, archiver skipped")
		return nil
	}

	return runChain(
		payload,
		t.createAddresses,
		t.createValidators,
		t.createDelegations,
		t.createNetworkRecords,
		t.createStats,
	)
}

func (t PersistorTask) createAddresses(payload *Payload) error {
	t.logger.Debug("creating addresses")

	addrMap := map[string]bool{}
	addresses := []model.Address{}

	for _, v := range payload.CurrentValidators {
		for _, addr := range v.RewardOwner.Addresses {
			if addrMap[addr] {
				continue
			}
			addrMap[addr] = true

			addresses = append(addresses, model.Address{
				Value:     addr,
				CreatedAt: payload.SyncTime,
				UpdatedAt: payload.SyncTime,
			})
		}
	}

	for _, d := range payload.CurrentDelegators {
		for _, addr := range d.RewardOwner.Addresses {
			if addrMap[addr] {
				continue
			}
			addrMap[addr] = true

			addresses = append(addresses, model.Address{
				Value:     addr,
				CreatedAt: payload.SyncTime,
				UpdatedAt: payload.SyncTime,
			})
		}
	}

	return t.db.Addresses.Import(addresses)
}

func (t PersistorTask) createValidators(payload *Payload) error {
	t.logger.Debug("creating validators")
	if err := t.db.Validators.Import(payload.Validators); err != nil {
		return err
	}

	t.logger.Debug("creating validator sequences")
	if err := t.db.Validators.ImportSeq(payload.ValidatorSeq); err != nil {
		return err
	}

	return nil
}

func (t PersistorTask) createDelegations(payload *Payload) error {
	t.logger.Debug("creating delegations")
	return t.db.Delegators.Import(payload.Delegations)
}

func (t PersistorTask) createNetworkRecords(payload *Payload) error {
	t.logger.Debug("creating network metrics")
	return t.db.Networks.CreateMetric(payload.NetworkMetric)
}

func (t PersistorTask) createStats(payload *Payload) error {
	for _, bucket := range []string{"h", "d"} {
		t.logger.WithField("bucket", bucket).Debug("creating network stats")
		if err := t.db.Networks.CreateStats(payload.SyncTime, bucket); err != nil {
			return err
		}

		t.logger.WithField("bucket", bucket).Debug("creating validator stats")
		if err := t.db.Validators.CreateStats(payload.SyncTime, bucket); err != nil {
			return err
		}
	}

	t.logger.Debug("resetting table counters")
	if err := t.db.ResetTableSeqCounters(); err != nil {
		return nil
	}

	return nil
}
