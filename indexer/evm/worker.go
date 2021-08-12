package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	corethTypes "github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Worker struct {
	log           *logrus.Logger
	rpc           *client.Client
	db            *store.DB
	chain         string
	status        *model.SyncStatus
	syncStatusKey string

	errWaitTime time.Duration
	syncTime    time.Duration
	cycleTime   time.Duration
}

type fetchData struct {
	txID    string
	err     error
	receipt *corethTypes.Receipt
	trace   *client.Call
}

func NewWorker(db *store.DB, rpc *client.Client, log *logrus.Logger, chain string) Worker {
	return Worker{
		db:            db,
		rpc:           rpc,
		log:           log,
		chain:         chain,
		syncStatusKey: fmt.Sprintf("%s_evm", chain),

		errWaitTime: time.Second,
		syncTime:    time.Second * 3,
		cycleTime:   time.Millisecond * 10,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.WithField("chain", w.syncStatusKey).Info("starting worker")

	timer := time.NewTimer(time.Second)
	defer func() {
		timer.Stop()
		w.log.WithField("chain", w.syncStatusKey).Info("worker stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			w.log.WithField("chain", w.syncStatusKey).Info("stopping worker")
			return
		case startTime := <-timer.C:
			if err := w.Run(); err != nil {
				w.log.WithField("chain", w.syncStatusKey).WithError(err).Info("worker run failed")
				timer.Reset(w.errWaitTime)
				break
			}

			w.log.
				WithField("chain", w.syncStatusKey).
				WithField("index", w.status.IndexID).
				WithField("lag", w.status.Lag()).
				WithField("duration", time.Since(startTime)).
				Info("finished run")

			if w.status.AtTip() {
				timer.Reset(w.syncTime)
			} else {
				timer.Reset(w.cycleTime)
			}
		}
	}
}

func (w *Worker) Run() error {
	status, err := w.getSyncStatus()
	if err != nil {
		return err
	}
	w.status = status

	if w.status.AtTip() {
		return nil
	}

	startHeight := w.status.NextID()
	page := 1
	heightRange := 100
	maxConcurrency := 100

	for {
		search := store.TxSearchInput{
			Chain:       w.chain,
			Type:        model.TxTypeEvm,
			StartHeight: int(startHeight),
			EndHeight:   int(startHeight) + heightRange,
			Order:       "height_asc",
			Limit:       heightRange,
			Page:        page,
		}

		txSearch, err := w.db.Transactions.Search(&search)
		if err != nil {
			return err
		}
		if len(txSearch.Transactions) == 0 {
			break
		}

		results := []fetchData{}
		resultsLock := sync.Mutex{}

		txIDS := make([]string, len(txSearch.Transactions))
		for idx, tx := range txSearch.Transactions {
			txIDS[idx] = tx.ID
		}

		doConcurrently(txIDS, maxConcurrency, func(txID string) {
			data := fetchData{}
			w.fetchTxData(txID, &data)

			resultsLock.Lock()
			results = append(results, data)
			resultsLock.Unlock()
		})

		for _, result := range results {
			if result.err != nil {
				w.log.WithField("tx_id", result.txID).WithError(result.err).Error("data fetch failed")
				return err
			}

			traceData, err := json.Marshal(result.trace)
			if err != nil {
				return err
			}

			trace := &model.EvmTrace{
				ID:        result.txID,
				Data:      string(traceData),
				Timestamp: time.Now(),
			}

			if err := w.db.Platform.CreateEvmTrace(trace); err != nil {
				return err
			}

			if err := w.createReceiptAndLogs(&result); err != nil {
				return err
			}

			w.status.IndexID = result.receipt.BlockNumber.Int64()
			w.status.IndexTime = time.Now()
		}

		// Continue onto next page only if we have max results on the current page
		if len(txSearch.Transactions) == search.Limit {
			page++
		} else {
			break
		}
	}

	return w.db.Platform.UpdateSyncStatus(status)
}

func (w *Worker) fetchTxData(txID string, data *fetchData) {
	data.txID = txID

	receipt, err := w.rpc.Evm.TransactionReceipt(context.Background(), common.HexToHash(txID))
	if err != nil {
		data.err = err
		return
	}
	data.receipt = receipt

	trace, err := w.rpc.Evm.TraceTransaction(context.Background(), txID)
	if err != nil {
		data.err = err
		return
	}
	data.trace = trace
}

var (
	cachedBlock     *corethTypes.Header
	cachedBlockTime time.Time
)

func (w *Worker) getSyncStatus() (*model.SyncStatus, error) {
	var lastBlock *corethTypes.Header
	var err error

	if time.Since(cachedBlockTime) < time.Second*10 {
		lastBlock = cachedBlock
	} else {
		lastBlock, err = w.rpc.Evm.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		cachedBlock = lastBlock
		cachedBlockTime = time.Now()
	}

	status, err := w.db.Platform.GetSyncStatus(w.syncStatusKey)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}

		status = &model.SyncStatus{
			ID:      w.syncStatusKey,
			IndexID: 0,
			TipID:   lastBlock.Number.Int64(),
			TipTime: time.Unix(int64(lastBlock.Time), 0),
		}

		if err := w.db.Platform.UpdateSyncStatus(status); err != nil {
			return nil, err
		}
	}

	status.TipID = lastBlock.Number.Int64()
	status.TipTime = time.Unix(int64(lastBlock.Time), 0)

	return status, nil
}

func (w *Worker) createReceiptAndLogs(data *fetchData) error {
	logsBatch := make([]model.EvmLog, len(data.receipt.Logs))

	for idx, logEntry := range data.receipt.Logs {
		topics := pq.StringArray{}
		for _, topic := range logEntry.Topics {
			topics = append(topics, topic.String())
		}

		logsBatch[idx] = model.EvmLog{
			Idx:     int(logEntry.Index),
			TxIdx:   int(logEntry.TxIndex),
			Address: logEntry.Address.String(),
			Removed: logEntry.Removed,
			Topics:  topics,
			Data:    common.Bytes2Hex(logEntry.Data),
		}
	}

	logsData, err := json.Marshal(logsBatch)
	if err != nil {
		return err
	}

	receipt := &model.EvmReceipt{
		ID:              data.receipt.TxHash.String(),
		ContractAddress: data.receipt.ContractAddress.String(),
		Type:            int(data.receipt.Type),
		Status:          int(data.receipt.Status),
		Logs:            string(logsData),
	}

	return w.db.Platform.CreateEvmReceipt(receipt)
}

func doConcurrently(items []string, maxConcurrency int, workFn func(string)) {
	queue := make(chan string)

	wg := &sync.WaitGroup{}
	wg.Add(maxConcurrency)

	for i := 0; i < maxConcurrency; i++ {
		go func() {
			defer wg.Done()

			for item := range queue {
				workFn(item)
			}
		}()
	}

	for _, item := range items {
		queue <- item
	}

	close(queue)
	wg.Wait()
}
