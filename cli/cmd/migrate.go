package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/avalanche-indexer/store/migrations"
)

type MigrateCommand struct {
	connStr string
	command string
	logger  *logrus.Logger
}

func NewMigrateCommand(command string, connStr string, logger *logrus.Logger) MigrateCommand {
	return MigrateCommand{
		command: command,
		connStr: connStr,
		logger:  logger,
	}
}

func (cmd MigrateCommand) Run() error {
	cmd.logger.Info("starting migration")
	defer cmd.logger.Info("finished migration")

	dbConn, err := store.NewRaw(cmd.connStr)
	if err != nil {
		return err
	}
	conn, err := dbConn.DB()
	if err != nil {
		return err
	}
	defer conn.Close()

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for path, f := range migrations.Assets.Files {
		if filepath.Ext(path) != ".sql" {
			continue
		}

		extPath := filepath.Join(tmpDir, filepath.Base(path))
		if err := ioutil.WriteFile(extPath, f.Data, 0755); err != nil {
			return err
		}
	}

	dir := "up"
	if chunks := strings.Split(cmd.command, ":"); len(chunks) > 1 {
		dir = chunks[1]
	}

	switch dir {
	case "migrate", "up":
		err = goose.Up(conn, tmpDir)
	case "down":
		err = goose.Down(conn, tmpDir)
	case "redo":
		if err = goose.Down(conn, tmpDir); err != nil {
			return err
		}
		err = goose.Up(conn, tmpDir)
	}

	return err
}
