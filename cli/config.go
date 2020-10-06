package cli

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Environment string `json:"env"`
	DatabaseURL string `json:"database_url"`
	RPCEndpoint string `json:"rpc_endpoint"`
	ServerAddr  string `json:"server_addr"`
	LogLevel    string `json:"log_level"`
	LogSQL      bool   `json:"log_sql"`

	SyncEnabled  bool   `json:"sync_enabled"`
	SyncInterval string `json:"sync_interval"`

	PurgeEnabled  bool   `json:"purge_enabled"`
	PurgeInterval string `json:"purge_interval"`
	PurgePeriod   string `json:"purge_period"`

	Archiver        string `json:"archiver"`
	ArchiverEnabled bool   `json:"archiver_enabled"`

	syncInterval  time.Duration
	purgeInterval time.Duration
}

func readConfig(path string) (*Config, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	config := &Config{}
	return config, json.NewDecoder(reader).Decode(config)
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return errors.New("database connection string is required")
	}

	if c.RPCEndpoint == "" {
		return errors.New("rpc endpoint is required")
	}

	if c.SyncInterval == "" {
		return errors.New("sync interval is required")
	}

	if c.PurgeInterval == "" {
		return errors.New("purge interval is required")
	}

	dur, err := time.ParseDuration(c.SyncInterval)
	if err != nil {
		return err
	}
	c.syncInterval = dur

	purgeDur, err := time.ParseDuration(c.PurgeInterval)
	if err != nil {
		return err
	}
	c.purgeInterval = purgeDur

	if c.Archiver != "" {
		uri, err := url.Parse(c.Archiver)
		if err != nil {
			return errors.New("invalida archiver configuration")
		}

		if !(uri.Scheme == "dir" || uri.Scheme == "s3") {
			return errors.New("unsupported archiver type")
		}
	}

	return nil
}

func (c *Config) GetSyncInterval() time.Duration {
	return c.syncInterval
}

func (c *Config) GetPurgeInterval() time.Duration {
	return c.purgeInterval
}
