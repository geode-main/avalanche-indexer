# Avalanche Indexer

Data indexer and API service for Avalanche Network

*Project is under active development*

## Requirements

- PostgreSQL 10.x+
- Go 1.14+

## Installation

Please see the sections below for all available methods of installation.

### Binary Releases

See [Github Releases](https://github.com/figment-networks/avalanche-indexer/releases) page for details.

## Usage

```bash
$ avalanche-indexer --help

Usage of ./avalanche-indexer:
  -cmd string
    	Command to execute
  -config string
    	Path to configuration file
  -v	Show version
```

Executing commands:

```bash
avalanche-indexer -config=path/to/config.json -cmd=COMMAND
```

Available commands:

| Name      | Description
|-----------|-----------------------------------------------------
| `status`  | Print out current indexer and node status
| `migrate` | Perform database migration
| `sync`    | Run a one-time indexer sync (for testing purposes)
| `worker`  | Start the indexer sync worker
| `server`  | Start the indexer API server

## Configuration

You can configure the service using a config file or environment variables.

### Config File

Example:

```json
{
  "database_url": "postgres://localhost/avalanche-indexer",
  "log_level": "debug",
  "rpc_endpoint": "http://localhost:9650",
  "sync_interval": "30s",
  "purge_interval": "60s",
  "server_addr": "localhost:8080",
  "archiver": "s3://us-east-1/bucketname",
  "archiver_enabled": true
}
```

## Running Application

Once you have created a database and specified all configuration options, you
need to migrate the database. You can do that by running the command below:

```bash
avalanche-indexer -config path/to/config.json -cmd=migrate
```

Perform the indexer check:

```bash
avalanche-indexer -config path/to/config.josn -cmd=status
```

Perform the initial sync:

```bash
avalanche-indexer -config path/to/config.josn -cmd=sync
```

If previous steps did not produce any errors you can start the indexer worker:

```bash
avalanche-indexer -config path/to/config.json -cmd=worker
```

Start the API server:

```bash
near-indexer -config path/to/config.json -cmd=server
```

## API Reference

| Method | Path                            | Description
|--------|---------------------------------|------------------------------------
| GET    | /health                         | Healthcheck endpoint
| GET    | /status                         | App version info and sync status
| GET    | /network_stats                  | List of network stats for a time bucket
| GET    | /validators                     | List of all active validators
| GET    | /validators/:id                 | Validator details

## License

Apache License v2.0
