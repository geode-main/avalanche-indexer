package shared

import (
	"strconv"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/avalanche-indexer/util"
)

func GetSyncStatus(indexClient *client.IndexClient, db *store.DB, containerType string, chain string) (*model.SyncStatus, error) {
	lastAccepted, err := indexClient.GetLastAccepted(containerType)
	if err != nil {
		return nil, err
	}

	lastIndex, err := util.ParseInt64(lastAccepted.Index)
	if err != nil {
		return nil, err
	}

	status, err := db.Platform.GetSyncStatus(chain)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}

		status = &model.SyncStatus{
			ID:      chain,
			IndexID: 0,
			TipID:   lastIndex,
			TipTime: lastAccepted.Timestamp,
		}

		if err := db.Platform.UpdateSyncStatus(status); err != nil {
			return nil, err
		}
	}

	status.TipID = lastIndex
	status.TipTime = lastAccepted.Timestamp

	return status, nil
}

func ProcessContainerRange(
	status *model.SyncStatus,
	indexClient *client.IndexClient,
	db *store.DB,
	containerType string,
	batchSize int,
	handlerFn func(*model.RawMessage) error,
) error {
	resp, err := indexClient.GetContainerRange(containerType, int(status.NextID()), batchSize)
	if err != nil {
		return err
	}
	containerItems := resp.Containers

	for _, c := range containerItems {
		idx, err := strconv.Atoi(c.Index)
		if err != nil {
			return err
		}

		raw := &model.RawMessage{
			IndexID:   idx,
			CreatedAt: c.Timestamp,
			Data:      c.Bytes,
		}

		if err := handlerFn(raw); err != nil {
			return err
		}

		status.IndexID = int64(idx)
		status.IndexTime = c.Timestamp
	}

	return db.Platform.UpdateSyncStatus(status)
}
