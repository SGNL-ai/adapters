# Multi-stage Dockerfile for the DB2 adapter service.
# Builds with IBM DB2 CLI driver (CGO) and runs on the SGNL Debian FIPS image
# (ghcr.io/sgnl-ai/debian) to eliminate perl/shell-related CVEs present in
# upstream debian:trixie-slim.

ARG GOLANG_IMAGE=golang:1.26-trixie
ARG RUNTIME_IMAGE=ghcr.io/sgnl-ai/debian:trixie-debian13-fips-r0
ARG DB2_CLI_VERSION=v12.1.2
# IBM DB2 CLI driver is x86_64 only; default platform to amd64.
ARG TARGETPLATFORM=linux/amd64

# ---------------------------------------------------------------------------
# Stage 1: Build the Go binary with CGO (needs DB2 CLI headers + libs)
# ---------------------------------------------------------------------------
FROM --platform=${TARGETPLATFORM} ${GOLANG_IMAGE} AS build

ARG DB2_CLI_VERSION

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

# Collect shared libraries needed at runtime that are NOT in the SGNL debian
# base image. Scans the Go binary, DB2 CLI driver libs, and libxml2.
RUN set -eu; \
    mkdir -p /out/usr/lib; \
    { ldd /sgnl/db2-adapter; \
      find /opt/ibm/clidriver/lib -name '*.so*' -exec ldd {} + 2>/dev/null; \
      ldd /usr/lib/x86_64-linux-gnu/libxml2.so.2; \
    } 2>/dev/null | awk '/=>/ {print $3}' | sort -u | while read -r lib; do \
        [ -f "$lib" ] || continue; \
        case "$lib" in \
            */ld-linux*) continue ;; \
            */libc.so*|*/libssl.so*|*/libcrypto.so*|*/libgcc_s.so*) continue ;; \
            */libm.so*|*/libdl.so*|*/libpthread.so*|*/librt.so*)   continue ;; \
            */libz.so*|*/libtinfo.so*|*/libselinux.so*) continue ;; \
            */libzstd.so*|*/libpcre2*.so*) continue ;; \
            */libsystemd.so*|*/libcap.so*) continue ;; \
            /opt/ibm/*) continue ;; \
        esac; \
        cp -L "$lib" /out/usr/lib/; \
    done

# ---------------------------------------------------------------------------
# Stage 2: Runtime
# ---------------------------------------------------------------------------
# Uses the SGNL Debian FIPS image built from
# glab/continuous-identity/infra/builders/debian/Dockerfile.
# This image has no perl-base (eliminating the CRITICAL CVEs) but includes
# shell, curl, CA certs, and core shared libs.
FROM --platform=${TARGETPLATFORM} ${RUNTIME_IMAGE} AS run

LABEL org.opencontainers.image.source="https://github.com/SGNL-ai/adapters" \
      org.opencontainers.image.title="SGNL DB2 Adapter" \
      org.opencontainers.image.description="DB2 adapter on SGNL Debian FIPS base (no perl, no CVE exposure)"

# Shared libs not in the base image (libxml2, libcrypt, libicuuc, etc.)
COPY --from=build /out/usr/lib/ /usr/lib/x86_64-linux-gnu/

# IBM DB2 CLI Driver: runtime libs, message catalogs, and config.
# msg/ is required for human-readable error messages; without it the driver
# returns unhelpful SQL10007N errors on connection failures.
COPY --from=build /opt/ibm/clidriver/lib /opt/ibm/clidriver/lib
COPY --from=build /opt/ibm/clidriver/msg /opt/ibm/clidriver/msg
COPY --from=build /opt/ibm/clidriver/cfg /opt/ibm/clidriver/cfg

ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib

WORKDIR /sgnl

COPY --from=build /go/bin/gops /sgnl/gops
COPY --from=build /sgnl/db2-adapter /sgnl/db2-adapter

EXPOSE 8080

# SGNL debian image ships sgnl:808:808
USER 808:808

ENTRYPOINT [ "/sgnl/db2-adapter" ]
