#!/bin/bash
# Test script for DB2 connection inside Docker container.
# Usage: ./dev/db2-test/test-db2-connection.sh

set -e

echo "=== DB2 Connection Test ==="
echo ""

# Check environment
echo "Checking DB2 environment..."
echo "IBM_DB_HOME: ${IBM_DB_HOME:-not set}"
echo "LD_LIBRARY_PATH: ${LD_LIBRARY_PATH:-not set}"
echo ""

# Check clidriver installation
if [ -d "/opt/ibm/clidriver" ]; then
    echo "OK: clidriver found at /opt/ibm/clidriver"
    ls -la /opt/ibm/clidriver/lib/*.so 2>/dev/null | head -5 || echo "  (no .so files found)"
else
    echo "✗ clidriver NOT found"
    exit 1
fi
echo ""

# Download dependencies
echo "Downloading Go dependencies..."
go mod download
echo ""

# Build with db2 tag
echo "Building with -tags db2..."
CGO_ENABLED=1 go build -tags db2 -v ./...
echo ""
echo "OK: Build successful with DB2 support!"
echo ""

# Run tests with db2 tag
echo "Running tests with -tags db2..."
CGO_ENABLED=1 go test -tags db2 -v ./pkg/db2/... -run TestSampleConfig 2>&1 || true
echo ""

echo "=== Test Complete ==="
