package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Server struct {
	annotations []routeAnnotation
	engine      *gin.Engine
	logger      *logrus.Logger
	db          *store.DB
	rpc         *client.Client
}

type routeAnnotation struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

func NewServer(db *store.DB, rpc *client.Client, logger *logrus.Logger) *Server {
	srv := &Server{
		engine:      gin.New(),
		annotations: []routeAnnotation{},
		db:          db,
		logger:      logger,
		rpc:         rpc,
	}

	srv.setupMiddleware()
	srv.setupRoutes()

	return srv
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *Server) setupRoutes() {
	s.addRoute(http.MethodGet, "/", "Index", s.handleIndex)
	s.addRoute(http.MethodGet, "/health", "Get indexer health", s.handleHealth)
	s.addRoute(http.MethodGet, "/status", "Get indexer status", s.handleStatus)
	s.addRoute(http.MethodGet, "/network_stats", "Get network stats", s.handleNetworkStats)
	s.addRoute(http.MethodGet, "/validators", "Get current validator set", s.handleValidators)
	s.addRoute(http.MethodGet, "/validators/:id", "Get validator details", s.handleValidator)
	s.addRoute(http.MethodGet, "/delegations", "Get active delegations", s.handleDelegations)
	s.addRoute(http.MethodGet, "/address/:id", "Get address details", s.handleAddress)
	s.addRoute(http.MethodGet, "/chains", "Get all blockchains", s.handleBlockchains)
	s.addRoute(http.MethodGet, "/chain_sync_statuses", "Get indexer sync status", s.handleSyncStatus)
	s.addRoute(http.MethodGet, "/assets", "Get all assets", s.handleAssets)
	s.addRoute(http.MethodGet, "/assets/:id", "Get asset details", s.handleAsset)
	s.addRoute(http.MethodGet, "/blocks", "Get blocks", s.handleBlocks)
	s.addRoute(http.MethodGet, "/blocks/:id", "Get block", s.handleBlock)
	s.addRoute(http.MethodGet, "/transactions", "Transactions search", s.handleTransactions)
	s.addRoute(http.MethodPost, "/transactions", "Transactions search", s.handleTransactions)
	s.addRoute(http.MethodGet, "/transactions/:id", "Get transaction details", s.handleTransaction)
	s.addRoute(http.MethodGet, "/transactions/:id/trace", "Get transaction trace", s.handleTransactionTrace)
	s.addRoute(http.MethodGet, "/transaction_outputs/:id", "Get transaction output", s.handleTransactionOutput)
	s.addRoute(http.MethodGet, "/transaction_types", "Get transaction types", s.handleTransactionTypeCounts)
	s.addRoute(http.MethodGet, "/events", "Events search", s.handleEvents)
	s.addRoute(http.MethodGet, "/events/:id", "Event details", s.handleEvent)

}

func (s *Server) addRoute(method, path, description string, handlers ...gin.HandlerFunc) {
	s.engine.Handle(method, path, handlers...)
	s.annotations = append(s.annotations, routeAnnotation{path, description})
}

func (s *Server) setupMiddleware() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(requestLogger(s.logger))
}

// check performs indexer healthcheck
func (s *Server) check() error {
	if err := s.db.Test(); err != nil {
		return err
	}
	if _, err := s.rpc.Info.NodeVersion(); err != nil {
		return err
	}
	return nil
}

// handleIndex returns all available endpoints
func (s *Server) handleIndex(c *gin.Context) {
	jsonOk(c, gin.H{"endpoints": s.annotations})
}

// handleHealth returns the indexer health
func (s *Server) handleHealth(c *gin.Context) {
	if err := s.check(); err != nil {
		jsonError(c, 400, err)
		return
	}

	jsonOk(c, gin.H{"healthy": true})
}

// handleStatus returns the indexer status
func (s *Server) handleStatus(c *gin.Context) {
	resp := StatusResponse{
		AppName:    indexer.AppName,
		AppVersion: indexer.AppVersion,
		GitCommit:  indexer.GitCommit,
		GoVersion:  indexer.GoVersion,
		SyncStatus: "stale",
	}

	lastTime, err := s.db.Validators.LastTime()
	if err != nil {
		s.logger.WithError(err).Error("cant fetch last validator time")
	}
	if lastTime != nil {
		resp.SyncTime = lastTime
		if time.Since(*lastTime) < time.Minute*5 {
			resp.SyncStatus = "current"
		}
	} else {
		resp.SyncStatus = "error"
	}

	nodeVersion, err := s.rpc.Info.NodeVersion()
	if err == nil {
		resp.NodeVersion = nodeVersion
	} else {
		resp.NodeVersion = "-"
		s.logger.WithError(err).Error("cant fetch node version")
	}

	networkName, err := s.rpc.Info.NetworkName()
	if err == nil {
		resp.NetworkName = networkName
	} else {
		resp.NetworkName = "-"
		s.logger.WithError(err).Error("cant fetch network name")
	}

	jsonOk(c, resp)
}

// handleSyncStatus returns chain sync status information
func (s *Server) handleSyncStatus(c *gin.Context) {
	result, err := s.db.Platform.GetSyncStatuses()
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// handleNetworkStats returns network stats for a given time bucket
func (s *Server) handleNetworkStats(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		bucket = "h"
	}
	if bucket != "h" && bucket != "d" {
		jsonError(c, 400, "invalid bucket value")
		return
	}

	limit := 0
	fmt.Sscanf(c.Query("limit"), "%d", &limit)
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		switch bucket {
		case "h":
			limit = 24
		case "d":
			limit = 30
		}
	}

	stats, err := s.db.Networks.GetStats(bucket, limit)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, stats)
}

// handleValidators returns validators records
func (s *Server) handleValidators(c *gin.Context) {
	search := store.ValidatorsSearch{}

	if err := c.Bind(&search); err != nil {
		badRequest(c, err)
		return
	}
	if err := search.Validate(); err != nil {
		badRequest(c, err)
		return
	}

	validators, err := s.db.Validators.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, validators)
}

// handleValidator returns validator details
func (s *Server) handleValidator(c *gin.Context) {
	validator, err := s.db.Validators.FindByNodeID(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	delegations, err := s.db.Delegators.Search(store.DelegationsSearch{NodeID: validator.NodeID})
	if shouldReturn(c, err) {
		return
	}

	hourStats, err := s.db.Validators.GetStats(validator.NodeID, "h", 24)
	if shouldReturn(c, err) {
		return
	}

	dayStats, err := s.db.Validators.GetStats(validator.NodeID, "d", 30)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, ValidatorResponse{
		Validator:   validator,
		Delegations: delegations,
		HourlyStats: hourStats,
		DailyStats:  dayStats,
	})
}

// handleDelegations renders all available delegations
func (s *Server) handleDelegations(c *gin.Context) {
	search := store.DelegationsSearch{}

	if err := c.Bind(&search); err != nil {
		badRequest(c, err)
		return
	}

	delegations, err := s.db.Delegators.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, delegations)
}

// handleAddress returns account balance on a given chain
func (s *Server) handleAddress(c *gin.Context) {
	address := c.Param("id")

	switch address[0] {
	case '0': // 0x.... address format
		var height *big.Int

		if heighVal := c.Query("height"); heighVal != "" {
			height = big.NewInt(0)
			_, ok := height.SetString(heighVal, 10)
			if !ok {
				badRequest(c, "invalid height value")
				return
			}
		}

		balance, err := s.rpc.Evm.BalanceAt(context.Background(), common.HexToAddress(address), height)
		if shouldReturn(c, err) {
			return
		}
		jsonOk(c, CBalanceResponse{
			Balance: balance.String(),
			Height:  height,
		})

	case 'P':
		balance, err := s.rpc.Platform.GetBalance(address)
		if shouldReturn(c, err) {
			return
		}
		jsonOk(c, balance)

	case 'X':
		resp, err := s.rpc.Avm.GetAllBalances(address)
		if shouldReturn(c, err) {
			return
		}
		jsonOk(c, resp.Balances)

	default:
		balance, err := s.rpc.Platform.GetBalance("P-" + address)
		if shouldReturn(c, err) {
			return
		}

		stakedResp, err := s.rpc.Platform.GetStake([]string{"P-" + address})
		if shouldReturn(c, err) {
			return
		}
		balance.Staked = stakedResp.Staked

		resp, err := s.rpc.Avm.GetAllBalances("X-" + address)
		if shouldReturn(c, err) {
			return
		}

		jsonOk(c, AddressBalancesResponse{
			Platform: *balance,
			Exchance: resp.Balances,
		})
	}
}

// handleTransactions performs transactions search
func (s *Server) handleTransactions(c *gin.Context) {
	input := &store.TxSearchInput{}
	if err := c.Bind(input); err != nil {
		badRequest(c, err)
		return
	}

	output, err := s.db.Transactions.Search(input)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, output.Transactions)
}

// handleTransaction loads and renders a single transaction details
func (s *Server) handleTransaction(c *gin.Context) {
	tx, err := s.db.Transactions.GetByID(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, tx)
}

// handleTransaction returns EVM trace details for a transaction
func (s *Server) handleTransactionTrace(c *gin.Context) {
	resp := TxTraceResponse{}

	trace, err := s.db.Platform.GetEvmTrace(c.Param("id"))
	if err != nil && err != store.ErrNotFound {
		serverError(c, err)
		return
	}
	if trace != nil {
		traceCall := &client.Call{}
		err = json.Unmarshal([]byte(trace.Data), traceCall)
		if shouldReturn(c, err) {
			return
		}
		resp.Trace = traceCall
	}

	receipt, err := s.db.Platform.GetEvmReceipt(c.Param("id"))
	if err != nil && err != store.ErrNotFound {
		serverError(c, err)
		return
	}
	if receipt != nil {
		logs := []model.EvmLog{}
		err = json.Unmarshal([]byte(receipt.Logs), &logs)
		if shouldReturn(c, err) {
			return
		}

		resp.Receipt = receipt
		resp.Logs = logs
	}

	jsonOk(c, resp)
}

// handleTransactionTypeCounts returns all available transaction typed and associated counts
func (s *Server) handleTransactionTypeCounts(c *gin.Context) {
	result, err := s.db.Transactions.GetTypeCounts(c.Query("chain"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, result)
}

// handleTransactionOutput returns a single transaction output record
func (s *Server) handleTransactionOutput(c *gin.Context) {
	result, err := s.db.Platform.GetTransactionOutput(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, result)
}

// handleBlockchains renders all available blockchains
func (s *Server) handleBlockchains(c *gin.Context) {
	chains, err := s.db.Platform.Chains()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, chains)
}

// handleBlockchains renders all available blockchains
func (s *Server) handleAssets(c *gin.Context) {
	var (
		assets []model.Asset
		err    error
	)

	if assetType := c.Query("type"); assetType != "" {
		assets, err = s.db.Assets.GetByType(assetType)
	} else {
		assets, err = s.db.Assets.GetAll()
	}

	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, assets)
}

// handleAsset renders asset details
func (s *Server) handleAsset(c *gin.Context) {
	asset, err := s.db.Assets.Get(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	count, err := s.db.Assets.GetTransactionsCount(asset.AssetID)
	if err != nil {
		s.logger.WithError(err).Error("cant fetch transactions count for asset")
	}
	if count != nil {
		asset.TransactionsCount = count
	}

	jsonOk(c, asset)
}

// handleBlocks renders blocks matching the search parameters
func (s Server) handleBlocks(c *gin.Context) {
	input := &store.BlocksSearch{}
	if err := c.Bind(input); err != nil {
		badRequest(c, err)
		return
	}

	blocks, err := s.db.Platform.GetBlocks(input)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, blocks)
}

// handleBlock renders a single block details
func (s Server) handleBlock(c *gin.Context) {
	block, err := s.db.Platform.GetBlock(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, block)
}

// handleEvents renders events matching the search parameters
func (s Server) handleEvents(c *gin.Context) {
	input := eventsSearchInput(c)
	if input == nil {
		return
	}

	events, err := s.db.Events.Search(input)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, events)
}

// handleEvent renders a single event details
func (s Server) handleEvent(c *gin.Context) {
	event, err := s.db.Events.FindByID(c.Param("id"))
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, event)
}
