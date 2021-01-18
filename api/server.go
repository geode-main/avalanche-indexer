package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
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
}

func (s *Server) addRoute(method, path, description string, handlers ...gin.HandlerFunc) {
	s.engine.Handle(method, path, handlers...)
	s.annotations = append(s.annotations, routeAnnotation{path, description})
}

func (s *Server) setupMiddleware() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(requestLogger(s.logger))
}

func (s *Server) handleIndex(c *gin.Context) {
	jsonOk(c, gin.H{
		"endpoints": s.annotations,
	})
}

func (s *Server) check() error {
	if err := s.db.Test(); err != nil {
		return err
	}
	if _, err := s.rpc.Info.NodeVersion(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleHealth(c *gin.Context) {
	if err := s.check(); err != nil {
		jsonError(c, 400, err)
		return
	}

	jsonOk(c, gin.H{"healthy": true})
}

func (s *Server) handleStatus(c *gin.Context) {
	data := gin.H{
		"app_name":    indexer.AppName,
		"app_version": indexer.AppVersion,
		"git_commit":  indexer.GitCommit,
		"go_version":  indexer.GoVersion,
		"sync_status": "stale",
	}

	lastTime, err := s.db.Validators.LastTime()
	if err != nil {
		s.logger.WithError(err).Error("cant fetch last validator time")
	}
	if lastTime != nil {
		data["sync_time"] = lastTime
		if time.Since(*lastTime) < time.Minute*5 {
			data["sync_status"] = "current"
		}
	} else {
		data["sync_status"] = "error"
	}

	nodeVersion, err := s.rpc.Info.NodeVersion()
	if err == nil {
		data["node_version"] = nodeVersion
	} else {
		data["node_version"] = "-"
		s.logger.WithError(err).Error("cant fetch node version")
	}

	networkName, err := s.rpc.Info.NetworkName()
	if err == nil {
		data["network_name"] = networkName
	} else {
		data["network_name"] = "-"
		s.logger.WithError(err).Error("cant fetch network name")
	}

	jsonOk(c, data)
}

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

	jsonOk(c, gin.H{
		"validator":   validator,
		"delegations": delegations,
		"stats_24h":   hourStats,
		"stats_30d":   dayStats,
	})
}

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

func (s *Server) handleAddress(c *gin.Context) {
	address := c.Param("id")

	switch address[0] {
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
		jsonError(c, 400, "Invalid address")
	}
}
