# PriceKinetics Code Test

This repository contains a Go service for merging, transforming, and serving event data.

## Prerequisites

- Go >= 1.25
- Docker

Other necessary tools (linters, formatters, helpers) can be installed via:

```
./updategotools.sh
```

## Quick Start

Run the service locally:

```
./run_local.sh
```

## Development

### Lint

```
make lint
```

### Unit Tests

```
make test
```

### Integration Tests

Integration tests require Docker services (MongoDB). The Make target will start Docker Compose for you.

```
make test-integration
```

You can also run integration tests directly with Go:

```
go test -tags=integration ./...
```

### Update Go Tools

```
./updategotools.sh
```

## Layout

### Model

This contains a proto definition of the Event model used by the code test. You can run `core/gen-proto.sh` to regenerate the model.

### Merger

This package merges two partial events together. When adding new fields you will need to make sure you have updated the code in this package to merge the new fields correctly.

### Core

The main service of the code test. This spins up a gRPC server and exposes RPCs to update and query events.

It persists the data in MongoDB and merges a partial update to an event with the existing copy of the event in the database, runs some transformations on the event, and saves back to the database.

## Service Flow (High Level)

1. Update request arrives with an event.
2. Existing event is loaded from the repository (if there exists one with the same event ID).
3. The merger combines existing event (if any) and the new event.
4. Transforms enrich/normalize data.
5. The updated event is persisted and served.

## Generated Code

Proto definitions live in `*.proto` files.
If you need to re-generate proto files, see run `go generate ./...`.

## API (gRPC)

See `API-README.md` for full RPC documentation and example payloads.
You can also use `test.http` to run example requests directly from GoLand.

## Notes

- Unit tests live under `core/service/test/unit`.
- Integration tests live under `core/service/test/integration` and require `-tags=integration`.
