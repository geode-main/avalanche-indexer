package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Server struct {
	engine *gin.Engine
	logger *logrus.Logger
	db     *store.DB
	rpc    *client.Client
}

func NewServer(db *store.DB, rpc *client.Client, logger *logrus.Logger) *Server {
	srv := &Server{
		engine: gin.New(),
		db:     db,
		logger: logger,
		rpc:    rpc,
	}

	srv.setupMiddleware()
	srv.setupRoutes()

	return srv
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}

func (s *Server) setupRoutes() {
	s.engine.GET("/", s.handleIndex)
	s.engine.GET("/health", s.handleHealth)
	s.engine.GET("/status", s.handleStatus)
	s.engine.GET("/network_stats", s.handleNetworkStats)
	s.engine.GET("/validators", s.handleValidators)
	s.engine.GET("/validators/:id", s.handleValidator)
}

func (s *Server) setupMiddleware() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(requestLogger(s.logger))
}

func (s *Server) handleIndex(c *gin.Context) {
	jsonOk(c, gin.H{
		"endpoints": []string{
			"/health",
			"/status",
			"/validators",
			"/validators/:id",
		},
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
	validators, err := s.db.Validators.FindAll()
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

	delegations, err := s.db.Delegators.FindByNodeID(validator.NodeID)
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
