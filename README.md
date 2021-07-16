# Avalanche Indexer

Blockchain data indexer and API for Avalanche network.

## Requirements

- PostgreSQL 12.x+
- Go 1.16+
- Access to Avalanchego full node

## Installation

Please see the sections below for all available methods of installation.

### Build from Source

Clone the repository:

```bash
git@github.com:figment-networks/avalanche-indexer.git
cd avalanche-indexer
```

Install dependencies:

```bash
make setup
```

Build the binary:

```bash
make
```

### Binary Releases

See [Github Releases](https://github.com/figment-networks/avalanche-indexer/releases) page for details.

### Docker

Build the docker image with:

```
make docker-build
```

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
  "log_level": "info",
  "rpc_endpoint": "http://localhost:9650",
  "sync_interval": "60s",
  "purge_interval": "5m",
  "server_addr": "localhost:8080",
  "network_id": 1,
  "evm_network_id": 1,
  "evm_chain_id": 43114
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
avalanche-indexer -config path/to/config.json -cmd=server
```

## API Reference

| Method | Path                            | Description
|--------|---------------------------------|------------------------------------
| GET    | /health                         | Healthcheck endpoint
| GET    | /status                         | App version info and sync status
| GET    | /network_stats                  | List of network stats for a time bucket
| GET    | /validators                     | List of active validators
| GET    | /validators/:id                 | Validator details
| GET    | /delegations                    | List of active delegations
| GET    | /address/:id                    | Get address balance (X-chain/P-chain)
| GET    | /assets                         | Get all available assets
| GET    | /assets/:id                     | Get asset details by ID
| GET    | /chains                         | List of existing chains
| GET    | /chain_sync_statuses            | Get primary chain (X/P/C) sync statuses
| GET    | /blocks                         | Get blocks by chain
| GET    | /blocks/:hash                   | Get block by hash (P/C)
| GET    | /transactions                   | Transactions search
| POST   | /transactions                   | Alternative transaction search endpoint
| GET    | /transactions/:hash             | Get transaction details by hash
| GET    | /transaction_outputs/:id        | Get a transaction output details by ID
| GET    | /transaction_types              | Get a summary of all transcation types

## License

Apache License v2.0
