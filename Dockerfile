ARG GOLANG_IMAGE=golang:1.25-bookworm
ARG USE_BAZEL_VERSION=6.1.1
ARG BUILD_TAGS=""

# STAGE 1: build...
FROM ${GOLANG_IMAGE} AS build

# Re-declare BUILD_TAGS for this build stage
ARG BUILD_TAGS=""

WORKDIR /app
COPY . ./

# Install DB2 CLI driver if BUILD_TAGS includes db2
RUN echo "BUILD_TAGS is: ${BUILD_TAGS}" && \
    if echo "${BUILD_TAGS}" | grep -q "db2"; then \
        echo "Installing IBM DB2 CLI Driver for db2 build..."; \
        apt-get update && \
        apt-get install -y wget ca-certificates && \
        cd /tmp && \
        wget -v https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/linuxx64_odbc_cli.tar.gz && \
        tar -xzf linuxx64_odbc_cli.tar.gz && \
        mkdir -p /opt/ibm && \
        mv clidriver /opt/ibm/ && \
        rm -f linuxx64_odbc_cli.tar.gz && \
        apt-get clean && \
        rm -rf /var/lib/apt/lists/* && \
        echo "DB2 CLI Driver installation completed"; \
    else \
        echo "Skipping DB2 driver installation for non-db2 build (BUILD_TAGS=${BUILD_TAGS})"; \
        mkdir -p /opt/ibm/clidriver/lib; \
    fi

# Set DB2 environment variables for conditional compilation
ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV CGO_CFLAGS=-I$IBM_DB_HOME/include
ENV CGO_LDFLAGS=-L$IBM_DB_HOME/lib
ENV LD_LIBRARY_PATH=$IBM_DB_HOME/lib:$LD_LIBRARY_PATH

RUN go mod download

ARG GOPS_VERSION=v0.3.27
RUN CGO_ENABLED=0 go install -ldflags "-s -w" github.com/google/gops@${GOPS_VERSION}

# Re-declare BUILD_TAGS ARG for this stage
ARG BUILD_TAGS=""

# Build adapter with appropriate flags based on BUILD_TAGS
RUN echo "BUILD_TAGS in build stage: ${BUILD_TAGS}" && \
    if echo "${BUILD_TAGS}" | grep -q "db2"; then \
        echo "Building with DB2 support (CGO_ENABLED=1, BUILD_TAGS=${BUILD_TAGS})..."; \
        CGO_ENABLED=1 GOOS=linux go build -C /app/cmd/adapter -tags "${BUILD_TAGS}" -o /sgnl/adapter; \
    else \
        echo "Building without DB2 support (CGO_ENABLED=0)..."; \
        CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/adapter -o /sgnl/adapter; \
    fi

# LDAP adapter never needs DB2 support
RUN CGO_ENABLED=0 GOOS=linux go build -C /app/cmd/ldap-adapter -o /sgnl/ldap-adapter

# STAGE 2: run...
# Use debian slim instead of distroless to support DB2 shared libraries
FROM debian:bookworm-slim AS run

# Install runtime dependencies for DB2
RUN apt-get update && \
    apt-get install -y libxml2 libssl3 libc6 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy DB2 runtime libraries from build stage
RUN mkdir -p /opt/ibm/clidriver/lib
COPY --from=build /opt/ibm/clidriver /opt/ibm/clidriver

# Set runtime library path for DB2
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib

# Create non-root user (debian doesn't have nonroot by default)
RUN groupadd --gid 65532 nonroot && \
    useradd --uid 65532 --gid nonroot --shell /bin/false --create-home nonroot

# Fixture files are loaded from `pkg/mock/...`, but we copy the files to `/sgnl/pkg/mock/...`
# so we change the working directory to `/sgnl` to make sure the files are found.
WORKDIR /sgnl

COPY --from=build --chown=nonroot:nonroot /go/bin/gops /sgnl/gops
COPY --from=build --chown=nonroot:nonroot /sgnl/adapter /sgnl/adapter
COPY --from=build --chown=nonroot:nonroot /sgnl/ldap-adapter /sgnl/ldap-adapter
COPY --from=build --chown=nonroot:nonroot /app/pkg/mock/servicenow/fixtures/*.yaml /sgnl/pkg/mock/servicenow/fixtures/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/sgnl/adapter" ]
