// Copyright 2026 SGNL.ai, Inc.
package testutil

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	LDAPPort  = "389"
	LDAPSPort = "636"
)

func setupLDAPContainer(ctx context.Context, isLDAPS bool) (testcontainers.Container, error) {
	absPath, err := filepath.Abs(filepath.Join("..", "..", "..", "dev", "active_directory", "directory.ldif"))
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}

	var waitForPort wait.Strategy
	if isLDAPS {
		waitForPort = wait.ForListeningPort(LDAPSPort + "/tcp")
	} else {
		waitForPort = wait.ForListeningPort(LDAPPort + "/tcp")
	}

	request := testcontainers.ContainerRequest{
		Image:        "osixia/openldap:1.5.0",
		ExposedPorts: []string{"389/tcp", "636/tcp"},
		AutoRemove:   true,
		WaitingFor:   waitForPort,
		Files: []testcontainers.ContainerFile{
			{
				Reader:            reader,
				HostFilePath:      absPath,
				ContainerFilePath: "/container/service/slapd/assets/config/bootstrap/ldif/custom/directory.ldif",
				FileMode:          0o700,
			},
		},
	}

	return testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Logger:           log.Default(),
			Started:          true,
		},
	)
}

// StartLDAPServer runs an instance of active directory over LDAP protocol in a local container for testing.
// It returns the container and open port.
// May fail the test internally if setup fails.
func (s *CommonSuite) StartLDAPServer(ctx context.Context, isLDAPS bool) (testcontainers.Container, nat.Port) {
	container, err := setupLDAPContainer(ctx, isLDAPS)
	if err != nil {
		s.T().Fatalf("Failed to setup LDAP container: %v", err)
	}

	var port nat.Port

	if isLDAPS {
		port, err = container.MappedPort(ctx, LDAPSPort)
		if err != nil {
			s.T().Fatalf("Failed to get mapped port of LDAPS container: %v", err)
		}
	} else {
		port, err = container.MappedPort(ctx, LDAPPort)
		if err != nil {
			s.T().Fatalf("Failed to get mapped port of LDAP container: %v", err)
		}
	}

	return container, port
}
