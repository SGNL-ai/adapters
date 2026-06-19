# Multi-stage Dockerfile for the DB2 adapter service.
# Builds with IBM DB2 CLI driver (CGO) and runs on CrowdStrike DHI (Distroless
# Hardened Image) to eliminate perl/shell-related CVEs from the runtime.

ARG GOLANG_IMAGE=golang:1.26-trixie
ARG ARTIFACTORY_DOCKER_REGISTRY=docker.artifactory.cicd.dc
ARG DB2_CLI_VERSION=v12.1.2

# ---------------------------------------------------------------------------
# Stage 1: Build the Go binary with CGO (needs DB2 CLI headers + libs)
# ---------------------------------------------------------------------------
# IBM DB2 CLI driver is x86_64 only; pin platform for cross-build correctness.
FROM --platform=linux/amd64 ${GOLANG_IMAGE} AS build

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

# ---------------------------------------------------------------------------
# Stage 2: Resolve runtime shared library dependencies
# ---------------------------------------------------------------------------
# The DHI -fips runtime ships libc, libssl, and libcrypto but not libxml2 or
# the DB2 CLI driver's transitive deps. This stage installs them via apt and
# uses ldd to collect exactly the .so files needed at runtime.
FROM --platform=linux/amd64 ${ARTIFACTORY_DOCKER_REGISTRY}/crowdstrike/dhi-debian-base:trixie-debian13-dev AS deps

ARG APT_MIRROR=

RUN set -eu; \
    if [ -n "$APT_MIRROR" ]; then \
        printf 'Types: deb\nURIs: %s/ext-debian-remote/debian\nSuites: trixie trixie-updates\nComponents: main\nSigned-By: /usr/share/keyrings/debian-archive-keyring.gpg\n\nTypes: deb\nURIs: %s/ext-debian-security-remote/debian-security\nSuites: trixie-security\nComponents: main\nSigned-By: /usr/share/keyrings/debian-archive-keyring.gpg\n' \
            "$APT_MIRROR" "$APT_MIRROR" > /etc/apt/sources.list.d/debian.sources; \
        rm -f /etc/apt/sources.list; \
    fi; \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        libxml2 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /opt/ibm/clidriver/lib /opt/ibm/clidriver/lib
COPY --from=build /opt/ibm/clidriver/msg /opt/ibm/clidriver/msg
COPY --from=build /opt/ibm/clidriver/cfg /opt/ibm/clidriver/cfg
COPY --from=build /sgnl/db2-adapter /sgnl/db2-adapter
COPY --from=build /go/bin/gops /sgnl/gops

# Collect shared libraries required by the Go binary and DB2 CLI driver that
# are NOT already present in the DHI -fips base image.
RUN set -eu; \
    mkdir -p /out/usr/lib /out/lib /out/opt/ibm/clidriver /out/sgnl; \
    \
    LDD_OUT=$( \
        ldd /sgnl/db2-adapter 2>/dev/null; \
        find /opt/ibm/clidriver/lib -name '*.so*' -exec ldd {} \; 2>/dev/null \
    ); \
    \
    echo "$LDD_OUT" | awk '/=>/ {print $3}' | sort -u | while read -r lib; do \
        [ -f "$lib" ] || continue; \
        case "$lib" in \
            */libc.so*|*/libssl.so*|*/libcrypto.so*|*/libgcc_s.so*) continue ;; \
            */libm.so*|*/libdl.so*|*/libpthread.so*|*/librt.so*)   continue ;; \
        esac; \
        cp -L "$lib" /out/usr/lib/; \
    done; \
    \
    # Dynamic linker (ld-linux-x86-64.so.2) - copy if not in base
    echo "$LDD_OUT" | awk '!/=>/ && /^\s*\//' | awk '{print $1}' | sort -u | while read -r lib; do \
        [ -f "$lib" ] && cp -L "$lib" /out/lib/ || true; \
    done; \
    \
    cp -r /opt/ibm/clidriver/lib /out/opt/ibm/clidriver/; \
    cp -r /opt/ibm/clidriver/msg /out/opt/ibm/clidriver/; \
    cp -r /opt/ibm/clidriver/cfg /out/opt/ibm/clidriver/; \
    cp /sgnl/db2-adapter /out/sgnl/; \
    cp /sgnl/gops /out/sgnl/

# ---------------------------------------------------------------------------
# Stage 3: Minimal DHI -fips runtime
# ---------------------------------------------------------------------------
FROM --platform=linux/amd64 ${ARTIFACTORY_DOCKER_REGISTRY}/crowdstrike/dhi-debian-base:trixie-debian13-fips AS run

LABEL org.opencontainers.image.source="https://github.com/SGNL-ai/adapters" \
      org.opencontainers.image.title="SGNL DB2 Adapter" \
      org.opencontainers.image.description="DB2 adapter on CrowdStrike DHI (FIPS-enabled, no perl/shell)"

# Shared libraries not included in DHI -fips (libxml2, libicuuc, etc.)
COPY --from=deps /out/usr/lib/ /usr/lib/x86_64-linux-gnu/
COPY --from=deps /out/lib/ /lib/

# IBM DB2 CLI Driver: runtime libs, message catalogs, and config.
# msg/ is required for human-readable error messages; without it the driver
# returns unhelpful SQL10007N errors on connection failures.
COPY --from=deps /out/opt/ibm/clidriver/ /opt/ibm/clidriver/

ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib

WORKDIR /sgnl

COPY --from=deps /out/sgnl/db2-adapter /sgnl/db2-adapter
COPY --from=deps /out/sgnl/gops /sgnl/gops

EXPOSE 8080

# DHI -fips already ships nonroot:65532, no groupadd/useradd needed.
USER 65532:65532

ENTRYPOINT [ "/sgnl/db2-adapter" ]
