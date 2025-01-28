# go-redis

A Redis server implemented in Go.

## Features

*   Implements the RESP protocol
*   Supports basic Redis commands (SET, GET, DEL, etc.)
*   Supports multiple databases
*   AOF persistence
*   Cluster mode

## Getting Started

### Prerequisites

*   Go 1.18 or higher

### Building

```bash
go build -o go-redis.exe main.go
```

### Running

```bash
./go-redis
```

### Configuration

The server can be configured using the `redis.conf` file. If the file does not exist, the server will use the default configuration.

The following options can be configured:

*   `bind`: The address to bind to (default: `0.0.0.0`)
*   `port`: The port to listen on (default: `6379`)
*   `appendOnly`: Whether to enable AOF persistence (default: `false`)
*   `appendFilename`: The name of the AOF file (default: `appendonly.aof`)
*   `databases`: The number of databases (default: 16)

### Cluster mode

To enable cluster mode, you need to configure the `self` and `peers` options in the `redis.conf` file.

*   `self`: The address of the current node
*   `peers`: A comma-separated list of addresses of other nodes in the cluster

## License

This project is licensed under the MIT License.
