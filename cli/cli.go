package cli

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/cli/cmd"
	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer"
	"github.com/figment-networks/avalanche-indexer/store"
)

var cliOpts struct {
	command    string
	configPath string
	version    bool
}

func init() {
	flag.StringVar(&cliOpts.command, "cmd", "", "Command to execute")
	flag.StringVar(&cliOpts.configPath, "config", "", "Path to configuration file")
	flag.BoolVar(&cliOpts.version, "v", false, "Show version")
	flag.Parse()
}

func Run() {
	if cliOpts.version {
		fmt.Println(indexer.VersionString())
		return
	}

	log := logrus.StandardLogger()

	if cliOpts.command == "" {
		log.Fatal("command is required")
	}

	if cliOpts.configPath == "" {
		log.Fatal("config path option is required")
	}

	config, err := readConfig(cliOpts.configPath)
	if err != nil {
		log.Fatal("config read error:", err)
	}
	if err := config.Validate(); err != nil {
		log.Fatal("config validation error:", err)
	}

	switch config.LogLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	}

	db, err := store.New(config.DatabaseURL)
	if err != nil {
		log.Fatal("db init error:", err)
	}

	rpc := client.New(config.RPCEndpoint)
	var command cliCommand

	switch cliOpts.command {
	case "status":
		command = cmd.NewStatusCommand(rpc, log)
	case "sync":
		command = cmd.NewSyncCommand(log, db, rpc, config.Archiver)
	case "worker":
		command = cmd.NewWorkerCommand(db, rpc, log, config.GetSyncInterval(), config.GetPurgeInterval(), config.Archiver)
	case "listener":
		command = cmd.NewListenerCommand(db, rpc, log)
	case "server":
		command = cmd.NewServerCommand(db, config.ServerAddr, log, rpc)
	case "migrate", "migrate:up", "migrate:down", "migrate:redo":
		command = cmd.NewMigrateCommand(cliOpts.command, config.DatabaseURL, log)
	case "purge":
		command = cmd.NewPurgeCommand(db, log)
	default:
		log.Fatal("invalid command")
	}

	err = command.Run()
	if err != nil {
		log.WithError(err).Fatal("command failed")
	}
}

type cliCommand interface {
	Run() error
}
