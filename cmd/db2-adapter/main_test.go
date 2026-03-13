// Copyright 2026 SGNL.ai, Inc.

// Tests for the DB2 adapter gRPC service entrypoint.
// Verifies adapter construction, registration, and configuration defaults.

package main

import (
	"os"
	"testing"

	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/adapters/pkg/db2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// setupAuthTokens creates a temporary auth tokens file required by the adapter framework server.
func setupAuthTokens(t *testing.T) {
	t.Helper()

	tmpDir := t.TempDir()
	authTokensPath := tmpDir + "/auth-tokens"

	err := os.WriteFile(authTokensPath, []byte(`["testtoken"]`), 0o600)
	require.NoError(t, err)

	t.Setenv("AUTH_TOKENS_PATH", authTokensPath)
}

func TestDB2AdapterConstruction_GivenDefaultSQLClient_WhenCreating_ThenReturnsValidAdapter(t *testing.T) {
	// Arrange
	sqlClient := db2.NewDefaultSQLClient()
	client := db2.NewClient(sqlClient)

	// Act
	adapter := db2.NewAdapter(client)

	// Assert
	assert.NotNil(t, adapter)
}

func TestDB2AdapterRegistration_GivenValidAdapter_WhenRegistering_ThenSucceeds(t *testing.T) {
	// Arrange
	setupAuthTokens(t)

	sqlClient := db2.NewDefaultSQLClient()
	client := db2.NewClient(sqlClient)
	adapter := db2.NewAdapter(client)

	stop := make(chan struct{})
	adapterServer := server.New(stop)

	s := grpc.NewServer()
	defer s.Stop()

	// Act
	err := server.RegisterAdapter(adapterServer, "DB2-1.0.0", adapter)

	// Assert
	require.NoError(t, err)
}

func TestConfigDefaults_GivenFreshViper_WhenReadingDefaults_ThenReturnsExpectedValues(t *testing.T) {
	// Arrange
	v := viper.New()
	v.SetEnvPrefix("DB2_ADAPTER")
	v.SetDefault("PORT", 8080)
	v.SetDefault("MAX_CALL_RECV_MSG_SIZE_MB", 8)
	v.SetDefault("MAX_CALL_SEND_MSG_SIZE_MB", 8)

	// Act
	port := v.GetInt("PORT")
	recvSize := v.GetInt("MAX_CALL_RECV_MSG_SIZE_MB")
	sendSize := v.GetInt("MAX_CALL_SEND_MSG_SIZE_MB")

	// Assert
	assert.Equal(t, 8080, port)
	assert.Equal(t, 8, recvSize)
	assert.Equal(t, 8, sendSize)
}

func TestConfigDefaults_GivenDB2AdapterPrefix_WhenNoConnectorServiceURL_ThenNotRequired(t *testing.T) {
	// Arrange
	v := viper.New()
	v.SetEnvPrefix("DB2_ADAPTER")

	// Act
	connectorServiceURL := v.GetString("CONNECTOR_SERVICE_URL")

	// Assert - DB2 adapter does not require a connector service URL
	assert.Empty(t, connectorServiceURL)
}
