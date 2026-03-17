# Dockerfile for testing with IBM DB2 client libraries.
# This creates an interactive development environment with DB2 support.

FROM --platform=linux/amd64 golang:1.25-bookworm

# Install dependencies
RUN apt-get update && \
    apt-get install -y \
        wget \
        ca-certificates \
        libxml2 \
        libssl3 \
        libc6 \
        && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download and install IBM DB2 CLI Driver (x86_64)
RUN cd /tmp && \
    wget -q https://public.dhe.ibm.com/ibmdl/export/pub/software/data/db2/drivers/odbc_cli/linuxx64_odbc_cli.tar.gz && \
    tar -xzf linuxx64_odbc_cli.tar.gz && \
    mkdir -p /opt/ibm && \
    mv clidriver /opt/ibm/ && \
    rm -f linuxx64_odbc_cli.tar.gz && \
    echo "DB2 CLI Driver installed successfully"

# Set DB2 environment variables
ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV CGO_CFLAGS=-I/opt/ibm/clidriver/include
ENV CGO_LDFLAGS=-L/opt/ibm/clidriver/lib
ENV LD_LIBRARY_PATH=/opt/ibm/clidriver/lib
ENV CGO_ENABLED=1

WORKDIR /app

# Default command: interactive shell
CMD ["/bin/bash"]
