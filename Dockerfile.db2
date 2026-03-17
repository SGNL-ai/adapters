# Multi-stage Dockerfile for the DB2 adapter service.
# Builds with IBM DB2 CLI driver (CGO) and runs on debian:bookworm-slim.

ARG GOLANG_IMAGE=golang:1.25-bookworm

# STAGE 1: build
FROM --platform=linux/amd64 ${GOLANG_IMAGE} AS build

# Install IBM DB2 CLI Driver (x86_64 only)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        wget \
        ca-certificates \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN cd /tmp && \
    wget -q https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/v12.1.2/linuxx64_odbc_cli.tar.gz && \
    tar -xzf linuxx64_odbc_cli.tar.gz && \
    mkdir -p /opt/ibm && \
    mv clidriver /opt/ibm/ && \
    rm -f linuxx64_odbc_cli.tar.gz

# Set DB2 build environment
ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV CGO_CFLAGS=-I/opt/ibm/clidriver/include
ENV CGO_LDFLAGS=-L/opt/ibm/clidriver/lib
ENV CGO_ENABLED=1

WORKDIR /app
COPY . ./

RUN go mod download

ARG GOPS_VERSION=v0.3.27
RUN go install -ldflags "-s -w" github.com/google/gops@${GOPS_VERSION}
RUN GOOS=linux go build -tags db2 -C /app/cmd/db2-adapter -o /sgnl/db2-adapter

# STAGE 2: run
FROM --platform=linux/amd64 debian:bookworm-slim AS run

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

# Copy IBM DB2 CLI Driver libraries from build stage
COPY --from=build /opt/ibm/clidriver/lib /opt/ibm/clidriver/lib

ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib

WORKDIR /sgnl

COPY --from=build /go/bin/gops /sgnl/gops
COPY --from=build /sgnl/db2-adapter /sgnl/db2-adapter

EXPOSE 8080

RUN groupadd --system nonroot && useradd --system --gid nonroot nonroot
USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/db2-adapter" ]
