// Copyright 2025 SGNL.ai, Inc.

package testutil

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	v1proxy "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type TestProxyServer struct {
	GrpcErr        error
	ResponseErrStr *string
	Response       *string
	Ci             *connector.ConnectorInfo
	v1proxy.UnimplementedProxyServiceServer
}

func (s *TestProxyServer) ProxyRequest(_ context.Context, req *v1proxy.ProxyRequestMessage) (*v1proxy.Response, error) {
	if s.Ci != nil {
		if req.ClientId != s.Ci.ClientID {
			return nil, fmt.Errorf("Expected %v, got %v client id", req.ClientId, s.Ci.ClientID)
		} else if req.ConnectorId != s.Ci.ID {
			return nil, fmt.Errorf("Expected %v, got %v connector id", req.ConnectorId, s.Ci.ID)
		} else if req.TenantId != s.Ci.TenantID {
			return nil, fmt.Errorf("Expected %v, got %v tenant id", req.TenantId, s.Ci.TenantID)
		}
	}

	// return gRPC error if configured
	if s.GrpcErr != nil {
		return nil, s.GrpcErr
	}

	// return error as part of the response payload
	if s.ResponseErrStr != nil {
		return &v1proxy.Response{
			ResponseType: &v1proxy.Response_LdapSearchResponse{
				LdapSearchResponse: &v1proxy.LDAPSearchResponse{
					Error: *s.ResponseErrStr,
				},
			}}, nil
	}

	// return valid response payload
	return &v1proxy.Response{
		ResponseType: &v1proxy.Response_LdapSearchResponse{
			LdapSearchResponse: &v1proxy.LDAPSearchResponse{
				Response: *s.Response,
			},
		}}, nil
}

func ProxyTestCommonSetup(t *testing.T, proxy *TestProxyServer) (v1proxy.ProxyServiceClient, func()) {
	// Create in-memory listener
	lis := bufconn.Listen(1024 * 1024)

	// Create gRPC test server and serve from the in-memory listener
	srv := grpc.NewServer()
	v1proxy.RegisterProxyServiceServer(srv, proxy)

	go srv.Serve(lis)

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(_ context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	client := v1proxy.NewProxyServiceClient(conn)

	return client, func() {
		conn.Close()
		srv.Stop()
		lis.Close()
	}
}
