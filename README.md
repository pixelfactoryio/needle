# Needle

[![Go Reference](https://pkg.go.dev/badge/go.pixelfactory.io/needle.svg)](https://pkg.go.dev/go.pixelfactory.io/needle)
[![Go Report Card](https://goreportcard.com/badge/go.pixelfactory.io/needle)](https://goreportcard.com/report/go.pixelfactory.io/needle)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Small HTTP/1.1, HTTP/2 server with TLS support that blocks ads and trackers by responding to all requests with a transparent 1x1 GIF pixel. Server certificates for requested domains are generated automatically on first request and cached on disk.

## Features

- **Ad & Tracker Blocking**: Responds to all requests with a transparent 1x1 GIF pixel
- **Automatic TLS**: Generates and caches certificates on-demand for requested domains
- **HTTP/2 Support**: Full support for HTTP/1.1 and HTTP/2 protocols
- **Built-in DNS Server**: Optional CoreDNS integration for local DNS resolution
- **Certificate Management**: Automatic certificate generation using a root CA
- **Persistent Storage**: BoltDB-based certificate caching
- **Observability**: Structured logging with support for multiple log levels

## Installation

### Using Go

```bash
go install go.pixelfactory.io/needle@latest
```

### Using Homebrew

```bash
brew install pixelfactoryio/tools/needle
```

### Using Docker

```bash
docker pull pixelfactory/needle:latest
```

### From Source

```bash
git clone https://github.com/pixelfactoryio/needle.git
cd needle
make build
```

## Usage

### Basic Usage

Start the server with default settings:

```bash
needle
```

This starts HTTP on port 80 and HTTPS on port 443 (requires root/sudo).

### Custom Configuration

Run with custom ports and log level:

```bash
needle --http-port 8080 --https-port 8443 --log-level debug
```

### With CoreDNS

Enable the built-in DNS server for local resolution:

```bash
needle --coredns-enabled --coredns-port 5353 --coredns-upstreams 8.8.8.8:53,1.1.1.1:53
```

### Using Configuration File

Create a `config.yaml` file:

```yaml
log-level: info
http-port: "8080"
https-port: "8443"
ca: data/certs/root-ca.crt
ca-key: data/certs/root-ca.key
db-file: data/cache.db
server-timeout: 60s
server-shutdown-timeout: 30s
coredns-enabled: true
coredns-port: 5353
coredns-upstreams:
  - 8.8.8.8:53
  - 1.1.1.1:53
```

Then run:

```bash
needle --config config.yaml
```

### Docker Usage

```bash
docker run -d \
  -p 80:80 \
  -p 443:443 \
  -v $(pwd)/data:/data \
  pixelfactory/needle:latest
```

## Configuration Options

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--log-level` | `LOG_LEVEL` | `info` | Log level (debug, info, warn, error, fatal, panic) |
| `--ca` | `CA` | `data/certs/root-ca.crt` | Root CA certificate path |
| `--ca-key` | `CA_KEY` | `data/certs/root-ca.key` | Root CA private key path |
| `--db-file` | `DB_FILE` | `data/cache.db` | Certificate cache database path |
| `--http-port` | `HTTP_PORT` | `80` | HTTP server port |
| `--https-port` | `HTTPS_PORT` | `443` | HTTPS server port |
| `--server-timeout` | `SERVER_TIMEOUT` | `60s` | Server read/write timeout |
| `--server-shutdown-timeout` | `SERVER_SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |
| `--coredns-enabled` | `COREDNS_ENABLED` | `false` | Enable built-in CoreDNS server |
| `--coredns-port` | `COREDNS_PORT` | `5353` | CoreDNS server port |
| `--coredns-hosts-file` | `COREDNS_HOSTS_FILE` | - | Custom hosts file path |
| `--coredns-upstreams` | `COREDNS_UPSTREAMS` | - | Upstream DNS servers (comma-separated) |
| `--coredns-corefile` | `COREDNS_COREFILE` | - | Custom Corefile path |
| `--config` | - | - | Configuration file path |

## How It Works

1. **Certificate Generation**: On first HTTPS request to a domain, Needle generates a certificate signed by the root CA
2. **Certificate Caching**: Generated certificates are cached in BoltDB for fast subsequent requests
3. **Response**: All HTTP/HTTPS requests receive a transparent 1x1 GIF pixel response
4. **DNS Resolution**: Optional CoreDNS server can redirect ad/tracker domains to localhost

## Development

### Prerequisites

- Go 1.24 or later
- Make
- golangci-lint 2.8.0 or later

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Code Formatting

```bash
make fmt
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Maintainer

Amine Benseddik ([@amine7536](https://github.com/amine7536))
