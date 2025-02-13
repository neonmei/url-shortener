# URL Shortener Service

A high-performance URL shortening service built with Go, featuring DynamoDB storage, Ristretto caching, and OpenTelemetry observability.

## Features

- URL shortening with base62 encoding
- DynamoDB-backed persistence with local development support
- In-memory caching using Ristretto
- OpenTelemetry instrumentation for tracing and metrics
- RESTful API with authentication
- Graceful shutdown handling
- Health check endpoints
- Configurable via environment variables

## Prerequisites

- Go 1.22 or later
- Docker (for local DynamoDB)
- AWS credentials (for production DynamoDB)
- [just](https://github.com/casey/just) command runner (optional)

## Quick Start

1. Set up local DynamoDB with dynamodb local or ScyllaDB (alternator):
```bash
just ddb-run         # Start local DynamoDB
just ddb-create      # Create required table
```

```bash
just scylla          # Start local DynamoDB
just ddb-create      # Create required table
```

2. Configure environment variables (see `.env` file for examples):
```bash
export SHORTENER_API_KEY="your-api-key"
export SHORTENER_BASE_URL="http://localhost:8080"
```

3. Run the service:
```bash
just run
```

## API Endpoints

### Public Endpoints

- `GET /:url_id` - Redirect to original URL

### Administrative Endpoints (Requires API Key)

- `POST /v1/urls/short` - Create short URL
- `DELETE /v1/urls/short/:url_id` - Delete short URL
- `GET /v1/urls/short/:url_id` - Fetch URL details

### Platform Endpoints

- `GET /platform/healthz` - Health check

## Development

### Running Tests
```bash
just test          # Run all tests
just coverage      # Generate coverage report
```

### Code Quality
```bash
just fmt           # Format code
just lint          # Run linter
just vet          # Run Go vet
```

### Performance Testing

Performance testing scripts are in a separate K6 repo. But the initial dataset can be generated with:

```bash
just perftest-dataset  # Generate test dataset
```

## Configuration

The service is configured via environment variables. Key configurations include:

- `SHORTENER_PORT` - HTTP server port (default: 8080)
- `SHORTENER_BASE_URL` - Base URL for shortened links
- `SHORTENER_API_KEY` - Authentication token for admin endpoints
- `SHORTENER_CACHE_METRICS_ENABLED` - Enable cache metrics
- `AWS_ENDPOINT_URL_DYNAMODB` - DynamoDB endpoint
- `OTEL_*` - OpenTelemetry configuration

See `config/config.go` for all available options.

## Architecture

The service follows a clean architecture pattern with the following layers:

- `domain` - Core business logic and interfaces
- `application` - Application service layer
- `platform` - Infrastructure implementations
- `cmd/http` - Startup ,API handlers and routing

### Key Components

- **URL Generation**: Uses base62 encoding with configurable collision handling
- **Storage**: DynamoDB with local development support
- **Caching**: Ristretto in-memory cache with optional metrics
- **Observability**: OpenTelemetry integration for tracing and metrics
- **API**: Gin web framework with middleware support
