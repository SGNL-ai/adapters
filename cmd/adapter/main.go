// Copyright 2025 SGNL.ai, Inc.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	aws "github.com/sgnl-ai/adapters/pkg/aws"
	aws_s3 "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/client"
	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/github"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/hashicorp"
	"github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/jira"
	jiradatacenter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/ldap"
	mysql "github.com/sgnl-ai/adapters/pkg/my-sql"
	"github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagerduty"
	"github.com/sgnl-ai/adapters/pkg/salesforce"
	"github.com/sgnl-ai/adapters/pkg/scim"
	"github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/sgnl-ai/adapters/pkg/workday"

	"google.golang.org/grpc"
)

var (
	// Port is the port at which the gRPC server will listen.
	Port = flag.Int("port", 8080, "The server port")

	// Timeout is the timeout for the HTTP client used to make requests to the datasource (seconds).
	Timeout = flag.Int("timeout", 30, "The timeout for the HTTP client used to make requests to the datasource (seconds)")

	// MaxConcurrency is the number of goroutines run concurrently in AWS adapter.
	MaxConcurrency = flag.Int("max_concurrency", 20, "The number of goroutines run concurrently in AWS adapter")
)

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *Port))
	if err != nil {
		logger.Fatalf("Failed to open server port: %v", err)
	}

	timeout := time.Duration(*Timeout) * time.Second

	s := grpc.NewServer()
	stop := make(chan struct{})
	adapterServer := server.New(stop)

	// Initialize the client to fetch data from AWS S3.
	s3Client, err := aws_s3.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-S3/1.0.0"), nil)
	if err != nil {
		logger.Fatalf("Failed to create a datasource to query AWS S3: %v", err)
	}

	// Initialize the client to fetch data from AWS.
	awsClient, err := aws.NewClient(
		client.NewSGNLHttpClient(timeout, "sgnl-AWS/1.0.0"), nil, *MaxConcurrency,
	)
	if err != nil {
		logger.Fatalf("Failed to create a datasource to query AWS: %v", err)
	}

	// Register adapters here alphabetically.
	server.RegisterAdapter(adapterServer, "AWS-1.0.0", aws.NewAdapter(awsClient))
	server.RegisterAdapter(
		adapterServer,
		"AzureAD-1.0.1",
		azuread.NewAdapter(azuread.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-AzureAD/1.0.1"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"BambooHR-1.0.0",
		bamboohr.NewAdapter(bamboohr.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-BambooHR/1.0.0"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"CrowdStrike-1.0.0",
		crowdstrike.NewAdapter(
			crowdstrike.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-CrowdStrike/1.0.0")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Duo-1.0.0",
		duo.NewAdapter(duo.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-Duo/1.0.0"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"GitHub-1.0.0",
		github.NewAdapter(github.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-GitHub/1.0.0"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"GoogleWorkspace-1.0.0",
		googleworkspace.NewAdapter(
			googleworkspace.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-GoogleWorkspace/1.0.0")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"HashiCorpBoundary-1.0.0",
		hashicorp.NewAdapter(
			hashicorp.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-HashiCorpBoundary/1.0.0")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"IdentityNow-1.0.0",
		identitynow.NewAdapter(identitynow.NewClient(
			client.NewSGNLHttpClient(timeout, "sgnl-IdentityNow/1.0.0"),
			identitynow.DefaultAccountCollectionPageSize,
		)),
	)
	server.RegisterAdapter(
		adapterServer,
		"Jira-1.0.0",
		jira.NewAdapter(jira.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-Jira/1.0.0"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"JiraDatacenter-1.0.0",
		jiradatacenter.NewAdapter(jiradatacenter.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-JiraDatacenter/1.0.0"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"LDAP-1.0.0",
		ldap.NewAdapter(),
	)
	server.RegisterAdapter(
		adapterServer,
		"MySQL-0.0.1-alpha",
		mysql.NewAdapter(mysql.NewClient(mysql.NewDefaultSQLClient())),
	)
	server.RegisterAdapter(
		adapterServer,
		"Okta-1.0.1",
		okta.NewAdapter(okta.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-Okta/1.0.1"))),
	)
	server.RegisterAdapter(
		adapterServer,
		"PagerDuty-1.0.0",
		pagerduty.NewAdapter(pagerduty.NewClient(
			client.NewSGNLHttpClient(timeout, "sgnl-PagerDuty/1.0.0")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Salesforce-1.0.1",
		salesforce.NewAdapter(salesforce.NewClient(
			client.NewSGNLHttpClient(timeout, "sgnl-Salesforce/1.0.1")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"SCIM2.0-1.0.0",
		scim.NewAdapter(scim.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-SCIM2.0/1.0.0"))),
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
			client.NewSGNLHttpClient(timeout, "sgnl-ServiceNow/1.0.1")),
		),
	)
	server.RegisterAdapter(
		adapterServer,
		"Workday-1.0.0",
		workday.NewAdapter(workday.NewClient(client.NewSGNLHttpClient(timeout, "sgnl-Workday/1.0.0"))),
	)

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	logger.Printf("Started adapter gRPC server on port %d", *Port)

	if err := s.Serve(listener); err != nil {
		close(stop)

		logger.Fatalf("Failed to listen on server port: %v", err)
	}
}
