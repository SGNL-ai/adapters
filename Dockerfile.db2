# Multi-stage Dockerfile for the DB2 adapter service.
# Builds with IBM DB2 CLI driver (CGO) and runs on debian:bookworm-slim.

ARG GOLANG_IMAGE=golang:1.26-bookworm
ARG DB2_CLI_VERSION=v12.1.2

# STAGE 1: build
# Note: IBM DB2 CLI driver is x86_64 only. Use --platform linux/amd64 when building on ARM.
FROM ${GOLANG_IMAGE} AS build

ARG DB2_CLI_VERSION

# Install IBM DB2 CLI Driver (x86_64 only)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        wget \
        ca-certificates \
        libxml2 \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN cd /tmp && \
    wget -q https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/${DB2_CLI_VERSION}/linuxx64_odbc_cli.tar.gz && \
    tar -xzf linuxx64_odbc_cli.tar.gz && \
    mkdir -p /opt/ibm && \
    mv clidriver /opt/ibm/ && \
    rm -f linuxx64_odbc_cli.tar.gz

# Set DB2 build environment
ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV CGO_CFLAGS=-I/opt/ibm/clidriver/include
ENV CGO_LDFLAGS=-L/opt/ibm/clidriver/lib
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib
ENV CGO_ENABLED=1

WORKDIR /app
COPY . ./

RUN go mod download

ARG GOPS_VERSION=v0.3.27
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/google/gops@${GOPS_VERSION}
RUN GOOS=linux go build -C /app/cmd/db2-adapter -tags db2 -o /sgnl/db2-adapter

# STAGE 2: run
FROM debian:bookworm-slim AS run

# Install runtime dependencies for DB2 client libraries
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libxml2 \
        libssl3 \
        libc6 \
        ca-certificates \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy IBM DB2 CLI Driver libraries, message catalogs, and config from build stage.
# The msg/ directory is required for human-readable error messages; without it the
# driver returns unhelpful SQL10007N errors on connection failures.
COPY --from=build /opt/ibm/clidriver/lib /opt/ibm/clidriver/lib
COPY --from=build /opt/ibm/clidriver/msg /opt/ibm/clidriver/msg
COPY --from=build /opt/ibm/clidriver/cfg /opt/ibm/clidriver/cfg

ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib

WORKDIR /sgnl

COPY --from=build --chown=nonroot:nonroot /go/bin/gops /sgnl/gops
COPY --from=build --chown=nonroot:nonroot /sgnl/db2-adapter /sgnl/db2-adapter

EXPOSE 8080

RUN groupadd --gid 65532 nonroot && \
    useradd --uid 65532 --gid nonroot --shell /bin/false --create-home nonroot
USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/db2-adapter" ]
