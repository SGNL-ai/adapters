# DB2 Test Environment

Docker-based development environment for building and testing the DB2 adapter with IBM DB2 CLI driver libraries. Required because the DB2 adapter uses CGO and needs the IBM client libraries available at compile time.

## Prerequisites

- Docker with `linux/amd64` platform support (including Apple Silicon via emulation)

## Quick Start

All commands run from the **project root** (`adapters/`):

```bash
# Build the test container
docker compose -f dev/db2-test/docker-compose.yml build

# Start an interactive shell inside the container
docker compose -f dev/db2-test/docker-compose.yml run --rm db2-test bash

# Clean up (remove container and built image)
docker compose -f dev/db2-test/docker-compose.yml down --rmi local
```

## Inside the Container

Once inside the container, the project is mounted at `/app` with DB2 environment variables pre-configured.

```bash
# Build with DB2 support
CGO_ENABLED=1 go build -tags db2 ./...

# Run DB2 tests
CGO_ENABLED=1 go test -tags db2 -v ./pkg/db2/...

# Run the connection test script
./dev/db2-test/test-db2-connection.sh
```

## Recording Test Fixtures

The `db2_record_fixtures.go` script captures real DB2 responses as JSON fixtures for contract testing.

Environment variables (constants and defaults defined in `db2_record_fixtures.go`):

| Variable | Constant | Description | Default |
|---|---|---|---|
| `DB2_PASSWORD` | `EnvDB2Password` | DB2 user password | *(required)* |
| `DB2_CERT_BASE64` | `EnvDB2CertB64` | Base64-encoded TLS certificate | *(optional)* |
| `DB2_DATABASE` | `EnvDB2Database` | Database name | `LMTESTDB` |
| `DB2_USER` | `EnvDB2User` | Database user | `db2inst1` |
| `DB2_HOST` | `EnvDB2Host` | Database host | `localhost` |
| `DB2_PORT` | `EnvDB2Port` | Database port | `50001` |

```bash
DB2_PASSWORD=<password> CGO_ENABLED=1 go run -tags db2 dev/db2-test/db2_record_fixtures.go
```

Fixtures are written to `pkg/db2/testdata/fixtures/`.

## Files

| File | Description |
|---|---|
| `Dockerfile` | Build image with Go toolchain and IBM DB2 CLI Driver (v12.1.2) |
| `docker-compose.yml` | Compose config with volume mounts and DB2 environment |
| `test-db2-connection.sh` | Validates DB2 environment, builds, and runs a smoke test |
| `db2_record_fixtures.go` | Records live DB2 responses as JSON test fixtures |
