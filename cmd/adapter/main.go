// Copyright 2025 SGNL.ai, Inc.
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector/client"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	aws "github.com/sgnl-ai/adapters/pkg/aws"
	aws_s3 "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/github"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/jira"
	jiradatacenter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	mysql_0_0_1_alpha "github.com/sgnl-ai/adapters/pkg/my-sql/0.0.1-alpha"
	mysql_0_0_2_alpha "github.com/sgnl-ai/adapters/pkg/my-sql/0.0.2-alpha"
	"github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagerduty"
	"github.com/sgnl-ai/adapters/pkg/rootly"
	"github.com/sgnl-ai/adapters/pkg/salesforce"
	"github.com/sgnl-ai/adapters/pkg/scim"
	"github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/sgnl-ai/adapters/pkg/workday"
	"go.uber.org/zap"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MiB = 1024 * 1024

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ADAPTER")

	// ADAPTER_PORT: The port at which the gRPC server will listen (default: 8080)
	viper.SetDefault("PORT", 8080)
	// ADAPTER_TIMEOUT: The timeout for the HTTP client used to make requests to the datasource, in seconds (default: 30)
	viper.SetDefault("TIMEOUT", 30)
	// ADAPTER_MAX_CONCURRENCY: The number of goroutines run concurrently in AWS adapter (default: 20)
	viper.SetDefault("MAX_CONCURRENCY", 20)
	// ADAPTER_MAX_S3_CSV_ROW_SIZE_BYTES: The maximum size of a CSV row in bytes (default: 1MiB)
	viper.SetDefault("MAX_S3_CSV_ROW_SIZE_BYTES", 1*MiB)
	// ADAPTER_MAX_S3_BYTES_TO_PROCESS_PER_PAGE: The maximum number of bytes to process per page (default: 10MiB)
	viper.SetDefault("MAX_S3_BYTES_TO_PROCESS_PER_PAGE", 10*MiB)
	// Read config from environment variables
	var (
		port                     = viper.GetInt("PORT")                        // ADAPTER_PORT
		timeout                  = viper.GetInt("TIMEOUT")                     // ADAPTER_TIMEOUT
		maxConcurrency           = viper.GetInt("MAX_CONCURRENCY")             // ADAPTER_MAX_CONCURRENCY
		connectorServiceURL      = viper.GetString("CONNECTOR_SERVICE_URL")    // ADAPTER_CONNECTOR_SERVICE_URL
		maxCSVRowSizeBytes       = viper.GetInt64("MAX_S3_CSV_ROW_SIZE_BYTES") // ADAPTER_MAX_S3_CSV_ROW_SIZE_BYTES
		maxBytesToProcessPerPage = viper.GetInt64(
			"MAX_S3_BYTES_TO_PROCESS_PER_PAGE") // ADAPTER_MAX_S3_BYTES_TO_PROCESS_PER_PAGE
	)

	if connectorServiceURL == "" {
		log.Fatal("ADAPTER_CONNECTOR_SERVICE_URL environment variable is required")
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Failed to open server port: %v.", zap.Error(err))
	}

	timeoutDuration := time.Duration(timeout) * time.Second

	s := grpc.NewServer()
	stop := make(chan struct{})
	adapterServer := server.New(stop, server.WithLogger(zaplogger.NewFrameworkLogger(logger)))

	connectorServiceClient, err := grpc.NewClient(
		connectorServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("Failed to create a grpc client to the connector service: %v", zap.Error(err))
	}

	// Initialize the client to fetch data from AWS S3.
	s3Client, err := aws_s3.NewClient(
		client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-S3/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		),
		nil,
		maxCSVRowSizeBytes,
		maxBytesToProcessPerPage,
	)
	if err != nil {
		logger.Fatal("Failed to create a datasource to query AWS S3: %v.", zap.Error(err))
	}

	// Initialize the client to fetch data from AWS.
	awsClient, err := aws.NewClient(
		client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-AWS/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		), nil, maxConcurrency,
	)
	if err != nil {
		logger.Fatal("Failed to create a datasource to query AWS: %v.", zap.Error(err))
	}

	// Register adapters here alphabetically.
	server.RegisterAdapter(adapterServer, "AWS-1.0.0", aws.NewAdapter(awsClient))
	server.RegisterAdapter(
		adapterServer,
		"AzureAD-1.0.1",
		azuread.NewAdapter(azuread.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-AzureAD/1.0.1",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			),
		)),
	)
	server.RegisterAdapter(
		adapterServer,
		"BambooHR-1.0.0",
		bamboohr.NewAdapter(bamboohr.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-BambooHR/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"CrowdStrike-1.0.0",
		crowdstrike.NewAdapter(
			crowdstrike.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-CrowdStrike/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Duo-1.0.0",
		duo.NewAdapter(duo.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Duo/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"GitHub-1.0.0",
		github.NewAdapter(github.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-GitHub/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"GoogleWorkspace-1.0.0",
		googleworkspace.NewAdapter(
			googleworkspace.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-GoogleWorkspace/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"HashiCorpBoundary-1.0.0",
		hashicorp.NewAdapter(
			hashicorp.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-HashiCorpBoundary/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"IdentityNow-1.0.0",
		identitynow.NewAdapter(identitynow.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-IdentityNow/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			), identitynow.DefaultAccountCollectionPageSize,
		)),
	)
	server.RegisterAdapter(
		adapterServer,
		"Jira-1.0.0",
		jira.NewAdapter(jira.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Jira/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"JiraDatacenter-1.0.0",
		jiradatacenter.NewAdapter(jiradatacenter.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-JiraDatacenter/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			),
		)),
	)
	server.RegisterAdapter(
		adapterServer,
		"MySQL-0.0.1-alpha",
		mysql_0_0_1_alpha.NewAdapter(mysql_0_0_1_alpha.NewClient(mysql_0_0_1_alpha.NewDefaultSQLClient(
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"MySQL-0.0.2-alpha",
		mysql_0_0_2_alpha.NewAdapter(mysql_0_0_2_alpha.NewClient(mysql_0_0_2_alpha.NewDefaultSQLClient(
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"Okta-1.0.1",
		okta.NewAdapter(okta.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Okta/1.0.1",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"PagerDuty-1.0.0",
		pagerduty.NewAdapter(pagerduty.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-PagerDuty/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Rootly-1.0.0",
		rootly.NewAdapter(rootly.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Rootly/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Salesforce-1.0.1",
		salesforce.NewAdapter(salesforce.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Salesforce/1.0.1",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			)),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"SCIM2.0-1.0.0",
		scim.NewAdapter(scim.NewClient(client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-SCIM2.0/1.0.0",
			grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
		))),
	)
	server.RegisterAdapter(
		adapterServer,
		"S3-1.0.0",
		aws_s3.NewAdapter(s3Client),
	)
	server.RegisterAdapter(
		adapterServer,
		"ServiceNow-1.0.1",
		servicenow.NewAdapter(servicenow.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-ServiceNow/1.0.1",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			),
		)),
	)
	server.RegisterAdapter(
		adapterServer,
		"Workday-1.0.0",
		workday.NewAdapter(workday.NewClient(
			client.NewSGNLHTTPClientWithProxy(timeoutDuration, "sgnl-Workday/1.0.0",
				grpc_proxy_v1.NewProxyServiceClient(connectorServiceClient),
			),
		)),
	)

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Info("Started adapter gRPC server", zap.Int("port", port))

	if err := s.Serve(listener); err != nil {
		close(stop)

		logger.Fatal("Failed to listen on server port: %v.", zap.Error(err))
	}
}
