// Copyright 2026 SGNL.ai, Inc.

package ldap_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	ldap_v3 "github.com/go-ldap/ldap/v3"
	framework "github.com/sgnl-ai/adapter-framework"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap/v2.0.0"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper to wait for LDAP server readiness.
func waitForLDAPReady(addr, bindDN, bindPassword string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := ldap_v3.DialURL("ldap://" + addr)
		if err == nil {
			err = conn.Bind(bindDN, bindPassword)
			conn.Close()

			if err == nil {
				return nil // Ready!
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("LDAP server at %s not ready after %s", addr, timeout)
}

// Helper to start a real OpenLDAP container for integration tests.
func startOpenLDAPContainer(t *testing.T) (testcontainers.Container, string, func()) {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "osixia/openldap:1.5.0",
		ExposedPorts: []string{"389/tcp"},
		WaitingFor:   wait.ForListeningPort("389/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start LDAP container: %v", err)
	}

	port, err := container.MappedPort(ctx, "389/tcp")
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("failed to get mapped port: %v", err)
	}

	addr := "localhost:" + port.Port()

	// Wait for LDAP server to be ready
	if err := waitForLDAPReady(addr, "cn=admin,dc=example,dc=org", "admin", 10*time.Second); err != nil {
		container.Terminate(ctx)
		t.Fatalf("LDAP server not ready: %v", err)
	}

	cleanup := func() { container.Terminate(ctx) }

	return container, addr, cleanup
}

func TestIntegration_SessionExpiryTTL(t *testing.T) {
	// Arrange
	_, addr, cleanup := startOpenLDAPContainer(t)
	defer cleanup()

	client := ldap_adapter.NewLDAPRequester(100*time.Millisecond, 10*time.Millisecond)
	bindDN := "cn=admin,dc=example,dc=org"
	bindPassword := "admin"
	baseDN := "dc=example,dc=org"
	request := &ldap_adapter.Request{
		ConnectionParams: ldap_adapter.ConnectionParams{
			BindDN:       bindDN,
			BindPassword: bindPassword,
			BaseDN:       baseDN,
		},
		BaseURL:          "ldap://" + addr,
		PageSize:         1,
		EntityExternalID: "Person",
		EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
			"Person": {Query: "(objectClass=inetOrgPerson)"},
		},
		Attributes: []*framework.AttributeConfig{},
	}

	_, err := client.Request(context.Background(), request)
	if err != nil {
		t.Fatalf("initial request failed: %v", err)
	}

	// Act: Simulate TTL expiry
	time.Sleep(1 * time.Second) // Wait for cleanup goroutine

	_, err2 := client.Request(context.Background(), request)
	// Assert: Should succeed, and a new session should be created
	if err2 != nil {
		t.Fatalf("request after TTL expiry failed: %v", err2)
	}
}

func TestIntegration_ConnectionDropRecovery(t *testing.T) {
	// Arrange
	container, addr, cleanup := startOpenLDAPContainer(t)
	defer cleanup()

	client := ldap_adapter.NewLDAPRequester(100*time.Millisecond, 10*time.Millisecond)
	bindDN := "cn=admin,dc=example,dc=org"
	bindPassword := "admin"
	baseDN := "dc=example,dc=org"
	request := &ldap_adapter.Request{
		ConnectionParams: ldap_adapter.ConnectionParams{
			BindDN:       bindDN,
			BindPassword: bindPassword,
			BaseDN:       baseDN,
		},
		BaseURL:          "ldap://" + addr,
		PageSize:         1,
		EntityExternalID: "Person",
		EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
			"Person": {Query: "(objectClass=inetOrgPerson)"},
		},
		Attributes: []*framework.AttributeConfig{},
	}

	_, err := client.Request(context.Background(), request)
	if err != nil {
		t.Fatalf("initial request failed: %v", err)
	}

	// Act: Simulate connection drop by stopping and recreating the container
	if err := container.Stop(context.Background(), nil); err != nil {
		t.Fatalf("failed to stop container: %v", err)
	}

	if err := container.Terminate(context.Background()); err != nil {
		t.Fatalf("failed to terminate container: %v", err)
	}

	time.Sleep(1 * time.Second)
	// Start a new container
	newContainer, newAddr, _ := startOpenLDAPContainer(t)
	defer newContainer.Terminate(context.Background())

	// Wait for the new container to be ready
	time.Sleep(2 * time.Second)

	// Update the request's BaseURL
	request.BaseURL = "ldap://" + newAddr
	_, err2 := client.Request(context.Background(), request)
	// Assert: Should recover and succeed
	if err2 != nil {
		t.Fatalf("request after connection drop failed: %v", err2)
	}
}
