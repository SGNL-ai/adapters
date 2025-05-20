// Copyright 2025 SGNL.ai, Inc.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/adapters/pkg/ldap"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LDAP_ADAPTER")

	// LDAP_ADAPTER_PORT: The port at which the gRPC server will listen (default: 8080)
	viper.SetDefault("PORT", 8080)
	// LDAP_ADAPTER_TIMEOUT: The session pool TTL in minutes (default: 30)
	viper.SetDefault("TIMEOUT", 30)
	// LDAP_ADAPTER_CONNECTOR_SERVICE_URL: The URL of the connector service (default: empty string)
	viper.SetDefault("CONNECTOR_SERVICE_URL", "")

	// Read config from environment variables
	port := viper.GetInt("PORT")                                    // LDAP_ADAPTER_PORT
	adapterTTL := viper.GetInt("TIMEOUT")                           // LDAP_ADAPTER_TIMEOUT
	connectorServiceURL := viper.GetString("CONNECTOR_SERVICE_URL") // LDAP_ADAPTER_CONNECTOR_SERVICE_URL

	logger := log.New(
		os.Stdout, "ldap-adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile,
	)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatalf("Failed to open server port: %v.", err)
	}

	s := grpc.NewServer()
	stop := make(chan struct{})
	adapterServer := server.New(stop)

	connectorServiceClient, err := grpc.Dial(
		connectorServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("Failed to create a grpc client to the connector service: %v.", err)
	}

	// Register only the LDAP adapter
	server.RegisterAdapter(
		adapterServer,
		"LDAP-1.0.0",
		ldap.NewAdapter(grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient), time.Duration(adapterTTL)*time.Minute),
	)

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Printf("LDAP Adapter gRPC server listening on %d.", port)

	if err := s.Serve(list); err != nil {
		close(stop)
		logger.Fatalf("Failed to serve: %v.", err)
	}

	logger.Println("Cleanup complete, exiting.")
}
