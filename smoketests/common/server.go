// Copyright 2025 SGNL.ai, Inc.
package common

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server"
	"github.com/sgnl-ai/adapters/pkg/aws"
	s3 "github.com/sgnl-ai/adapters/pkg/aws-s3"
	"github.com/sgnl-ai/adapters/pkg/azuread"
	"github.com/sgnl-ai/adapters/pkg/bamboohr"
	"github.com/sgnl-ai/adapters/pkg/crowdstrike"
	"github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/github"
	googleworkspace "github.com/sgnl-ai/adapters/pkg/google-workspace"
	"github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/jira"
	jiradatacenter "github.com/sgnl-ai/adapters/pkg/jira-datacenter"
	"github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagerduty"
	"github.com/sgnl-ai/adapters/pkg/salesforce"
	"github.com/sgnl-ai/adapters/pkg/scim"
	"github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/sgnl-ai/adapters/pkg/workday"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const testMaxConcurrency = 1

func StartAdapterServer(t *testing.T, client *http.Client, port int) chan struct{} {
	validTokensPath := "./TOKENS_0"

	tokens := []byte(`["dGhpc2lzYXRlc3R0b2tlbg==","dGhpc2lzYWxzb2F0ZXN0dG9rZW4="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}

	t.Setenv("AUTH_TOKENS_PATH", validTokensPath)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("Failed to open server port: %v", err)
	}

	s := grpc.NewServer()

	stop := make(chan struct{})

	adapterServer := server.New(stop)

	// Setup to record-replay AWS sdk
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithHTTPClient(client), // inject the recorder client
	)
	if err != nil {
		t.Fatalf("Failed to create a client for AWS S3 SoR: %v", err)
	}

	s3Client, err := s3.NewClient(client, &cfg)
	if err != nil {
		t.Fatalf("Failed to create a client for AWS S3 SoR: %v", err)
	}

	awsClient, err := aws.NewClient(client, &cfg, testMaxConcurrency)
	if err != nil {
		t.Fatalf("Failed to create a datasource to query AWS: %v", err)
	}

	// Register adapters here alphabetically.
	server.RegisterAdapter(adapterServer, "AWS-1.0.0", aws.NewAdapter(awsClient))
	server.RegisterAdapter(adapterServer, "AzureAD-1.0.1", azuread.NewAdapter(azuread.NewClient(client)))
	server.RegisterAdapter(adapterServer, "BambooHR-1.0.0", bamboohr.NewAdapter(bamboohr.NewClient(client)))
	server.RegisterAdapter(adapterServer, "CrowdStrike-1.0.0", crowdstrike.NewAdapter(crowdstrike.NewClient(client)))
	server.RegisterAdapter(adapterServer, "Duo-1.0.0", duo.NewAdapter(duo.NewClient(client)))
	server.RegisterAdapter(adapterServer, "GitHub-1.0.0", github.NewAdapter(github.NewClient(client)))
	server.RegisterAdapter(adapterServer, "GoogleWorkspace-1.0.0",
		googleworkspace.NewAdapter(googleworkspace.NewClient(client)))
	server.RegisterAdapter(adapterServer, "IdentityNow-1.0.0",
		identitynow.NewAdapter(identitynow.NewClient(client, identitynow.DefaultAccountCollectionPageSize)))
	server.RegisterAdapter(adapterServer, "Jira-1.0.0", jira.NewAdapter(jira.NewClient(client)))
	server.RegisterAdapter(adapterServer, "JiraDatacenter-1.0.0",
		jiradatacenter.NewAdapter(jiradatacenter.NewClient(client)))
	server.RegisterAdapter(adapterServer, "Okta-1.0.1", okta.NewAdapter(okta.NewClient(client)))
	server.RegisterAdapter(adapterServer, "PagerDuty-1.0.0", pagerduty.NewAdapter(pagerduty.NewClient(client)))
	server.RegisterAdapter(adapterServer, "Salesforce-1.0.1", salesforce.NewAdapter(salesforce.NewClient(client)))
	server.RegisterAdapter(adapterServer, "SCIM2.0-1.0.0", scim.NewAdapter(scim.NewClient(client)))
	server.RegisterAdapter(adapterServer, "ServiceNow-1.0.1", servicenow.NewAdapter(servicenow.NewClient(client)))
	server.RegisterAdapter(adapterServer, "S3-1.0.0", s3.NewAdapter(s3Client))
	server.RegisterAdapter(adapterServer, "Workday-1.0.0", workday.NewAdapter(workday.NewClient(client)))

	api_adapter_v1.RegisterAdapterServer(s, adapterServer)

	if err := s.Serve(listener); err != nil {
		t.Fatalf("Failed to listen on server port: %v", err)
	}

	return stop
}

func GetNewAdapterClient(t *testing.T, port int) (api_adapter_v1.AdapterClient, *grpc.ClientConn) {
	var err error

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	dialCtx, dialCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dialCancel()

	conn, err := grpc.DialContext(dialCtx, fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		t.Fatal(err)
	}

	return api_adapter_v1.NewAdapterClient(conn), conn
}
