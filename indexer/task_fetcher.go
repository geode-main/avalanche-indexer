package indexer

import (
	"context"
	"sync"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
)

type FetcherTask struct {
	rpc    *client.Client
	logger *logrus.Logger
}

func (t FetcherTask) GetName() string {
	return taskFetcher
}

func (t FetcherTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	payload := p.(*Payload)

	networkName, err := t.rpc.Info.NetworkName()
	if err != nil {
		return err
	}

	nodeVersion, err := t.rpc.Info.NodeVersion()
	if err != nil {
		return err
	}

	validatorsResp, err := t.rpc.Platform.GetCurrentValidators()
	if err != nil {
		return err
	}

	pendingValidatorsResp, err := t.rpc.Platform.GetPendingValidators()
	if err != nil {
		return err
	}

	height, err := t.rpc.Platform.GetCurrentHeight()
	if err != nil {
		return err
	}

	stakeResp, err := t.rpc.Platform.GetMinStake()
	if err != nil {
		return err
	}

	blockchainsResp, err := t.rpc.Platform.GetBlockchains()
	if err != nil {
		return err
	}

	peers, err := t.rpc.Info.Peers()
	if err != nil {
		return err
	}

	txFee, err := t.rpc.Info.TxFee()
	if err != nil {
		return err
	}

	payload.NetworkName = networkName
	payload.NodeVersion = nodeVersion
	payload.Height = height
	payload.CurrentValidators = validatorsResp.Validators
	payload.CurrentDelegators = validatorsResp.Delegators
	payload.PendingValidators = pendingValidatorsResp.Validators
	payload.PendingDelegators = pendingValidatorsResp.Delegators
	payload.Blockchains = blockchainsResp.Blockchains
	payload.Peers = peers
	payload.MinStake = stakeResp
	payload.RawTxFee = txFee

	// this will make 100s of http calls...
	// if err := t.fetchBalances(payload); err != nil {
	// 	return err
	// }

	return nil
}

func (t FetcherTask) fetchBalances(payload *Payload) error {
	addrChan := make(chan string)
	balances := map[string]*client.Balance{}
	balancesLock := sync.Mutex{}
	ids := []string{}
	wg := sync.WaitGroup{}
	numWorkers := 10

	for _, v := range payload.CurrentValidators {
		if len(v.RewardOwner.Addresses) == 0 {
			continue
		}
		ids = append(ids, v.RewardOwner.Addresses[0])
	}

	for _, d := range payload.CurrentDelegators {
		if len(d.RewardOwner.Addresses) == 0 {
			continue
		}
		ids = append(ids, d.RewardOwner.Addresses[0])
	}

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				select {
				case addr, ok := <-addrChan:
					if !ok {
						return
					}

					balance, err := t.rpc.Platform.GetBalance(addr)
					if err != nil {
						t.logger.Error(err)
					}

					balancesLock.Lock()
					balances[addr] = balance
					balancesLock.Unlock()

					wg.Done()
				}
			}
		}()
	}

	wg.Add(len(ids))
	for _, id := range ids {
		addrChan <- id
	}

	wg.Wait()
	close(addrChan)

	return nil
}
