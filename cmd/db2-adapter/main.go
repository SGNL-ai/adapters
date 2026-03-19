// Copyright 2026 SGNL.ai, Inc.

// Entrypoint for the DB2 adapter gRPC service.
// Runs as a separate container with IBM DB2 client libraries for direct database connectivity.

package main

import (
	"fmt"
	"log"
	"net"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/adapters/pkg/db2"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const MiB = 1024 * 1024

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("DB2_ADAPTER")

	// DB2_ADAPTER_PORT: The port at which the gRPC server will listen (default: 8080)
	viper.SetDefault("PORT", 8080)
	// DB2_ADAPTER_MAX_CALL_RECV_MSG_SIZE_MB: Maximum gRPC receive message size in MB (default: 8MB)
	viper.SetDefault("MAX_CALL_RECV_MSG_SIZE_MB", 8)
	// DB2_ADAPTER_MAX_CALL_SEND_MSG_SIZE_MB: Maximum gRPC send message size in MB (default: 8MB)
	viper.SetDefault("MAX_CALL_SEND_MSG_SIZE_MB", 8)

	var (
		port                 = viper.GetInt("PORT")                      // DB2_ADAPTER_PORT
		maxCallRecvMsgSizeMB = viper.GetInt("MAX_CALL_RECV_MSG_SIZE_MB") // DB2_ADAPTER_MAX_CALL_RECV_MSG_SIZE_MB
		maxCallSendMsgSizeMB = viper.GetInt("MAX_CALL_SEND_MSG_SIZE_MB") // DB2_ADAPTER_MAX_CALL_SEND_MSG_SIZE_MB
	)

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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to open server port: %d", port), zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxCallRecvMsgSizeMB*MiB),
		grpc.MaxSendMsgSize(maxCallSendMsgSizeMB*MiB),
	)
	stop := make(chan struct{})
	adapterServer := server.New(stop, server.WithLogger(zaplogger.NewFrameworkLogger(logger)))

	// Register DB2 adapter. The DB2 adapter connects directly to the database
	// and does not require a connector service proxy.
	if err := server.RegisterAdapter(
		adapterServer,
		"DB2-1.0.0",
		db2.NewAdapter(db2.NewClient(db2.NewDefaultSQLClient())),
	); err != nil {
		logger.Fatal("Failed to register DB2 adapter", zap.Error(err))
	}

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Info(fmt.Sprintf("Started DB2 adapter gRPC server on port %d", port))

	if err := s.Serve(listener); err != nil {
		close(stop)

		logger.Fatal(fmt.Sprintf("Failed to listen on server port: %d", port), zap.Error(err))
	}

	logger.Info("Cleanup complete, exiting")
}
