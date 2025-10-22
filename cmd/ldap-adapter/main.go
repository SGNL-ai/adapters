// Copyright 2025 SGNL.ai, Inc.

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	adapter_v1 "github.com/sgnl-ai/adapters/pkg/ldap/v1.0.0"
	adapter_v2 "github.com/sgnl-ai/adapters/pkg/ldap/v2.0.0"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LDAP_ADAPTER")

	// LDAP_ADAPTER_PORT: The port at which the gRPC server will listen (default: 8080)
	viper.SetDefault("PORT", 8080)
	// LDAP_ADAPTER_SESSION_TTL: The session pool TTL in minutes (default: 30)
	viper.SetDefault("SESSION_TTL", 30)
	// LDAP_ADAPTER_SESSION_CLEANUP_INTERVAL: The session pool cleanup interval in minutes (default: 1)
	viper.SetDefault("SESSION_CLEANUP_INTERVAL", 1)
	// Read config from environment variables
	port := viper.GetInt("PORT")                                       // LDAP_ADAPTER_PORT
	adapterTTL := viper.GetInt("SESSION_TTL")                          // LDAP_ADAPTER_SESSION_TTL
	adapterCleanupInterval := viper.GetInt("SESSION_CLEANUP_INTERVAL") // LDAP_ADAPTER_SESSION_CLEANUP_INTERVAL
	connectorServiceURL := viper.GetString("CONNECTOR_SERVICE_URL")    // LDAP_ADAPTER_CONNECTOR_SERVICE_URL

	if connectorServiceURL == "" {
		log.Fatal("LDAP_ADAPTER_CONNECTOR_SERVICE_URL environment variable is required")
	}

	loggerCfg, err := zaplogger.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load logger configuration: %v", err)
	}

	logger := zaplogger.New(*loggerCfg, zap.WithCaller(true))

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Failed to open server port", zap.Error(err))
	}

	s := grpc.NewServer()
	stop := make(chan struct{})
	adapterServer := server.New(stop, server.WithLogger(zaplogger.NewFrameworkLogger(logger)))

	connectorServiceClient, err := grpc.Dial(
		connectorServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("Failed to create a grpc client to the connector service", zap.Error(err))
	}

	// Register LDAP-v1.0.0 adapter.
	server.RegisterAdapter(
		adapterServer,
		"LDAP-1.0.0",
		adapter_v1.NewAdapter(
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			time.Duration(adapterTTL)*time.Minute,
			time.Duration(adapterCleanupInterval)*time.Minute),
	)

	// Register LDAP-v2.0.0 adapter.
	server.RegisterAdapter(
		adapterServer,
		"LDAP-2.0.0",
		adapter_v2.NewAdapter(
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			time.Duration(adapterTTL)*time.Minute,
			time.Duration(adapterCleanupInterval)*time.Minute),
	)

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Info("Started LDAP adapter gRPC server", zap.Int("port", port))

	if err := s.Serve(list); err != nil {
		close(stop)
		logger.Fatal("Failed to serve", zap.Error(err))
	}

	logger.Info("Cleanup complete, exiting")
}
