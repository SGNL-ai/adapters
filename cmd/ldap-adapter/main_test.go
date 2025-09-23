// Copyright 2025 SGNL.ai, Inc.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	attrUID = "uid"

	// Seed test user using ldapadd.
	ldifPath = "testdata/test.ldif"
)

func TestMainFunction_NoPanic(t *testing.T) {
	t.Setenv("LDAP_ADAPTER_CONNECTOR_SERVICE_URL", "localhost:1234")

	// Use t.TempDir for temporary directory
	tmpDir := t.TempDir()
	authTokensPath := tmpDir + "/fake-auth-tokens"

	// Set required env vars using t.Setenv
	t.Setenv("AUTH_TOKENS_PATH", authTokensPath)

	// Create the dummy file so the watcher does not fail
	f, err := os.Create(authTokensPath)
	if err != nil {
		t.Fatalf("failed to create dummy auth tokens file: %v", err)
	}

	f.Close()

	// Run main in a goroutine and recover from panic
	panicChan := make(chan interface{}, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}

			close(panicChan)
		}()

		main()
	}()

	// Wait briefly to see if panic occurs (service will block on Serve)
	select {
	case p := <-panicChan:
		if p != nil {
			t.Fatalf("main panicked: %v", p)
		}
	case <-time.After(100 * time.Millisecond): // 100ms, Success: no panic in short time
	}
}

// TestLDAPAdapter_GetPage_WithRealLDAP tests the LDAP adapter end-to-end with a real OpenLDAP instance.
//
// Given a running OpenLDAP server with a seeded user,
// When the LDAP adapter GetPage is called for that user entity,
// Then the adapter should return the user in the results.
func TestGivenOpenLDAPWithUser_WhenGetPageIsCalled_ThenUserIsReturned(t *testing.T) {
	t.Setenv("LDAP_ADAPTER_CONNECTOR_SERVICE_URL", "localhost:1234")

	// Arrange
	ctx := context.Background()

	// Start OpenLDAP in Docker using testcontainers-go
	ldapReq := testcontainers.ContainerRequest{
		Image:        "osixia/openldap:1.5.0",
		ExposedPorts: []string{"389/tcp"},
		Env: map[string]string{
			"LDAP_ORGANISATION":   "Example Org",
			"LDAP_DOMAIN":         "example.org",
			"LDAP_ADMIN_PASSWORD": "admin",
		},
		WaitingFor: wait.ForListeningPort("389/tcp").WithStartupTimeout(30 * time.Second),
	}

	ldapC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: ldapReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start LDAP container: %s", err)
	}
	defer ldapC.Terminate(ctx)

	ldapHost, err := ldapC.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get LDAP container host: %s", err)
	}

	ldapPort, err := ldapC.MappedPort(ctx, "389/tcp")
	if err != nil {
		ports, perr := ldapC.Ports(ctx)
		if perr == nil {
			t.Logf("Available ports: %v", ports)
		}

		logs, logErr := ldapC.Logs(ctx)
		if logErr == nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(logs)
			t.Logf("Container logs:\n%s", buf.String())
		}

		ldapPort, err = ldapC.MappedPort(ctx, "389")
		if err != nil {
			t.Fatalf("Could not get LDAP container port: tried '389/tcp' and '389', available ports: %v, error: %v", ports, err)
		}
	}

	// Wait a bit to ensure LDAP server is ready for ldapadd
	t.Log("Sleeping 2s to ensure LDAP server is ready for ldapadd...")
	time.Sleep(2 * time.Second)

	// Print LDIF file contents for debugging
	ldifBytes, readErr := os.ReadFile(ldifPath)
	if readErr != nil {
		t.Fatalf("Failed to read test LDIF file: %v", readErr)
	}

	t.Logf("LDIF file contents before copy:\n%s", string(ldifBytes))

	addCmd := []string{
		"ldapadd",
		"-x",
		"-D", "cn=admin,dc=example,dc=org",
		"-w", "admin",
		"-H", "ldap://localhost:389",
		"-f", "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif",
	}

	err = ldapC.CopyFileToContainer(ctx, ldifPath, "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif", 0644)
	if err != nil {
		t.Fatalf("Failed to copy test.ldif to container: %v", err)
	}

	exitCode, output, err := ldapC.Exec(ctx, addCmd)
	buf := new(bytes.Buffer)

	if output != nil {
		_, _ = buf.ReadFrom(output)
	}

	outputStr := buf.String()
	t.Logf("ldapadd exit code: %d, error: %v, output: %s", exitCode, err, outputStr)

	if err != nil || exitCode != 0 {
		t.Fatalf("Failed to exec ldapadd: %v, exit code: %d, output: %s", err, exitCode, outputStr)
	}

	t.Logf("ldapadd output: %s", outputStr)

	// Wait a moment to ensure LDAP server is ready for search
	time.Sleep(500 * time.Millisecond)

	searchCmd := []string{
		"ldapsearch",
		"-x",
		"-D", "cn=admin,dc=example,dc=org",
		"-w", "admin",
		"-H", "ldap://127.0.0.1:389",
		"-b", "ou=users,dc=example,dc=org",
		"(uid=john)",
	}

	sExit, sOutput, sErr := ldapC.Exec(ctx, searchCmd)
	if sErr != nil || sExit != 0 {
		buf = new(bytes.Buffer)
		_, _ = buf.ReadFrom(sOutput)
		outputBytes := buf.Bytes()
		t.Fatalf("ldapsearch failed: %v, exit code: %d, output: %s", sErr, sExit, string(outputBytes))
	}

	buf = new(bytes.Buffer)
	_, _ = buf.ReadFrom(sOutput)

	outputBytes := buf.Bytes()
	if !bytes.Contains(outputBytes, []byte("uid: john")) {
		t.Fatalf("ldapsearch did not find user 'john'. Output: %s", string(outputBytes))
	}

	// Set up the adapter port
	adapterPort := 54321
	tmpDir := t.TempDir()
	authTokensPath := tmpDir + "/fake-auth-tokens"

	t.Setenv("LDAP_ADAPTER_PORT", fmt.Sprintf("%d", adapterPort))
	// Set up auth tokens for the adapter
	t.Setenv("AUTH_TOKENS_PATH", authTokensPath)
	_ = os.WriteFile(authTokensPath, []byte("[\"test-token\"]"), 0644)

	go func() {
		main()
	}()

	time.Sleep(500 * time.Millisecond) // Wait for server to start

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", adapterPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial grpc server: %v", err)
	}
	defer conn.Close()

	client := api_adapter_v1.NewAdapterClient(conn)

	// Prepare minimal valid LDAP config for the adapter
	ldapConfig := map[string]interface{}{
		"baseDN": "dc=example,dc=org",
		"entityConfig": map[string]interface{}{
			"Person": map[string]interface{}{
				"query": "(objectClass=inetOrgPerson)",
			},
		},
	}

	configBytes, err := json.Marshal(ldapConfig)
	if err != nil {
		t.Fatalf("failed to marshal ldap config: %v", err)
	}

	// Act
	// Prepare and send GetPageRequest
	req := &api_adapter_v1.GetPageRequest{
		Datasource: &api_adapter_v1.DatasourceConfig{
			Id:      "test-ldap",
			Address: fmt.Sprintf("%s:%s", ldapHost, ldapPort.Port()),
			Type:    "LDAP-1.0.0",
			Config:  configBytes,
			Auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_Basic_{
					Basic: &api_adapter_v1.DatasourceAuthCredentials_Basic{
						Username: "cn=admin,dc=example,dc=org",
						Password: "admin",
					},
				},
			},
		},
		Entity: &api_adapter_v1.EntityConfig{
			Id:         "Person",
			ExternalId: "Person",
			Attributes: []*api_adapter_v1.AttributeConfig{
				{
					Id:         attrUID,
					ExternalId: attrUID,
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "cn",
					ExternalId: "cn",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	t.Logf("LDAP Address used in DatasourceConfig: %s", req.Datasource.Address)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"token": "test-token"})
	ctxWithToken := metadata.NewOutgoingContext(ctxTimeout, md)
	resp, err := client.GetPage(ctxWithToken, req)

	// Assert
	if err != nil {
		t.Fatalf("GetPage failed: %v", err)
	}

	page := resp.GetSuccess()
	if page == nil || len(page.Objects) == 0 {
		t.Fatalf("expected at least one user, got none. Full response: %+v", resp)
	}

	found := false

	for _, obj := range page.Objects {
		for _, attr := range obj.Attributes {
			if attr.Id == attrUID && len(attr.Values) > 0 && attr.Values[0].GetStringValue() == "john" {
				found = true

				break
			}
		}
	}

	if !found {
		t.Fatalf("expected to find user 'john' in LDAP results")
	}
}

// TestGivenOpenLDAPWithMultipleUsers_WhenPagedGetPageIsCalled_ThenAllUsersAreReturnedAcrossPages
//
// Given a running OpenLDAP server with three users,
// When the LDAP adapter GetPage is called with a page size of 2,
// Then the adapter should return two users on the first page and the remaining user on the next page.
func TestGivenOpenLDAPWithMultipleUsers_WhenPagedGetPageIsCalled_ThenAllUsersAreReturnedAcrossPages(t *testing.T) {
	t.Setenv("LDAP_ADAPTER_CONNECTOR_SERVICE_URL", "localhost:1234")

	// Arrange
	ctx := context.Background()

	ldapReq := testcontainers.ContainerRequest{
		Image:        "osixia/openldap:1.5.0",
		ExposedPorts: []string{"389/tcp"},
		Env: map[string]string{
			"LDAP_ORGANISATION":   "Example Org",
			"LDAP_DOMAIN":         "example.org",
			"LDAP_ADMIN_PASSWORD": "admin",
		},
		WaitingFor: wait.ForListeningPort("389/tcp").WithStartupTimeout(30 * time.Second),
	}

	ldapC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: ldapReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start LDAP container: %s", err)
	}
	defer ldapC.Terminate(ctx)

	ldapHost, err := ldapC.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get LDAP container host: %s", err)
	}

	ldapPort, err := ldapC.MappedPort(ctx, "389/tcp")
	if err != nil {
		ports, perr := ldapC.Ports(ctx)
		if perr == nil {
			t.Logf("Available ports: %v", ports)
		}

		logs, logErr := ldapC.Logs(ctx)
		if logErr == nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(logs)
			t.Logf("Container logs:\n%s", buf.String())
		}

		ldapPort, err = ldapC.MappedPort(ctx, "389")
		if err != nil {
			t.Fatalf("Could not get LDAP container port: tried '389/tcp' and '389', available ports: %v, error: %v", ports, err)
		}
	}

	t.Log("Sleeping 2s to ensure LDAP server is ready for ldapadd...")
	time.Sleep(2 * time.Second)

	// Print LDIF file contents for debugging
	ldifBytes, readErr := os.ReadFile(ldifPath)
	if readErr != nil {
		t.Fatalf("Failed to read test LDIF file: %v", readErr)
	}

	t.Logf("LDIF file contents before copy:\n%s", string(ldifBytes))

	addCmd := []string{
		"ldapadd",
		"-x",
		"-D", "cn=admin,dc=example,dc=org",
		"-w", "admin",
		"-H", "ldap://localhost:389",
		"-f", "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif",
	}

	err = ldapC.CopyFileToContainer(ctx, ldifPath, "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif", 0644)
	if err != nil {
		t.Fatalf("Failed to copy test.ldif to container: %v", err)
	}

	exitCode, output, err := ldapC.Exec(ctx, addCmd)
	buf := new(bytes.Buffer)

	if output != nil {
		_, _ = buf.ReadFrom(output)
	}

	outputStr := buf.String()
	t.Logf("ldapadd exit code: %d, error: %v, output: %s", exitCode, err, outputStr)

	if err != nil || exitCode != 0 {
		t.Fatalf("Failed to exec ldapadd: %v, exit code: %d, output: %s", err, exitCode, outputStr)
	}

	t.Logf("ldapadd output: %s", outputStr)

	time.Sleep(500 * time.Millisecond)

	// Set up the adapter port
	adapterPort := 54322
	tmpDir := t.TempDir()
	authTokensPath := tmpDir + "/fake-auth-tokens"

	t.Setenv("LDAP_ADAPTER_PORT", fmt.Sprintf("%d", adapterPort))
	// Set up auth tokens for the adapter
	t.Setenv("AUTH_TOKENS_PATH", authTokensPath)
	_ = os.WriteFile(authTokensPath, []byte("[\"test-token\"]"), 0644)

	go func() {
		main()
	}()
	time.Sleep(500 * time.Millisecond)

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", adapterPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial grpc server: %v", err)
	}

	defer conn.Close()
	client := api_adapter_v1.NewAdapterClient(conn)

	ldapConfig := map[string]interface{}{
		"baseDN": "dc=example,dc=org",
		"entityConfig": map[string]interface{}{
			"Person": map[string]interface{}{
				"query": "(objectClass=inetOrgPerson)",
			},
		},
	}

	configBytes, err := json.Marshal(ldapConfig)
	if err != nil {
		t.Fatalf("failed to marshal ldap config: %v", err)
	}

	// Act
	// First page request
	req := &api_adapter_v1.GetPageRequest{
		Datasource: &api_adapter_v1.DatasourceConfig{
			Id:      "test-ldap",
			Address: fmt.Sprintf("%s:%s", ldapHost, ldapPort.Port()),
			Type:    "LDAP-1.0.0",
			Config:  configBytes,
			Auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_Basic_{
					Basic: &api_adapter_v1.DatasourceAuthCredentials_Basic{
						Username: "cn=admin,dc=example,dc=org",
						Password: "admin",
					},
				},
			},
		},
		Entity: &api_adapter_v1.EntityConfig{
			Id:         "Person",
			ExternalId: "Person",
			Attributes: []*api_adapter_v1.AttributeConfig{
				{
					Id:         attrUID,
					ExternalId: attrUID,
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "cn",
					ExternalId: "cn",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	t.Logf("LDAP Address used in DatasourceConfig: %s", req.Datasource.Address)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	md := metadata.New(map[string]string{"token": "test-token"})
	ctxWithToken := metadata.NewOutgoingContext(ctxTimeout, md)

	resp, err := client.GetPage(ctxWithToken, req)

	// Assert (First Page)
	if err != nil {
		t.Fatalf("GetPage (first page) failed: %v", err)
	}

	page := resp.GetSuccess()

	if page == nil || len(page.Objects) != 2 {
		t.Fatalf("expected two users on first page, got: %+v", resp)
	}

	userSet := make(map[string]bool)

	for _, obj := range page.Objects {
		for _, attr := range obj.Attributes {
			if attr.Id == attrUID && len(attr.Values) > 0 {
				userSet[attr.Values[0].GetStringValue()] = true
			}
		}
	}

	if page.NextCursor == "" {
		t.Fatalf("expected next_cursor to be set for paging, got empty string")
	}

	// Act (Second Page)
	// Use the next_cursor to fetch the next page
	req.Cursor = page.NextCursor
	resp2, err := client.GetPage(ctxWithToken, req)

	// Assert (Second Page)
	if err != nil {
		t.Fatalf("GetPage (second page) failed: %v", err)
	}

	page2 := resp2.GetSuccess()
	if page2 == nil || len(page2.Objects) != 1 {
		t.Fatalf("expected one user on second page, got: %+v", resp2)
	}

	for _, obj := range page2.Objects {
		for _, attr := range obj.Attributes {
			if attr.Id == attrUID && len(attr.Values) > 0 {
				userSet[attr.Values[0].GetStringValue()] = true
			}
		}
	}

	// Assert all users are returned across both pages
	expectedUsers := map[string]bool{"john": true, "alice": true, "bob": true}
	for user := range expectedUsers {
		if !userSet[user] {
			t.Fatalf("expected user '%s' to be returned across pages, but was missing", user)
		}
	}
}

func TestGivenOpenLDAPWithGroupMembers_WhenGetGroupMemberPageIsCalled_ThenGroupDNsAreReturned(t *testing.T) {
	t.Setenv("LDAP_ADAPTER_CONNECTOR_SERVICE_URL", "localhost:1234")

	// Arrange
	ctx := context.Background()

	ldapReq := testcontainers.ContainerRequest{
		Image:        "osixia/openldap:1.5.0",
		ExposedPorts: []string{"389/tcp"},
		Env: map[string]string{
			"LDAP_ORGANISATION":   "Example Org",
			"LDAP_DOMAIN":         "example.org",
			"LDAP_ADMIN_PASSWORD": "admin",
		},
		WaitingFor: wait.ForListeningPort("389/tcp").WithStartupTimeout(30 * time.Second),
	}

	ldapC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: ldapReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start LDAP container: %s", err)
	}
	defer ldapC.Terminate(ctx)

	ldapHost, err := ldapC.Host(ctx)
	if err != nil {
		t.Fatalf("Could not get LDAP container host: %s", err)
	}

	ldapPort, err := ldapC.MappedPort(ctx, "389/tcp")
	if err != nil {
		ports, perr := ldapC.Ports(ctx)
		if perr == nil {
			t.Logf("Available ports: %v", ports)
		}

		logs, logErr := ldapC.Logs(ctx)
		if logErr == nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(logs)
			t.Logf("Container logs:\n%s", buf.String())
		}

		ldapPort, err = ldapC.MappedPort(ctx, "389")
		if err != nil {
			t.Fatalf("Could not get LDAP container port: tried '389/tcp' and '389', available ports: %v, error: %v", ports, err)
		}
	}

	t.Log("Sleeping 2s to ensure LDAP server is ready for ldapadd...")
	time.Sleep(2 * time.Second)

	// Print LDIF file contents for debugging
	ldifBytes, readErr := os.ReadFile(ldifPath)
	if readErr != nil {
		t.Fatalf("Failed to read test LDIF file: %v", readErr)
	}

	t.Logf("LDIF file contents before copy:\n%s", string(ldifBytes))

	addCmd := []string{
		"ldapadd",
		"-x",
		"-D", "cn=admin,dc=example,dc=org",
		"-w", "admin",
		"-H", "ldap://localhost:389",
		"-f", "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif",
	}

	err = ldapC.CopyFileToContainer(ctx, ldifPath, "/container/service/slapd/assets/config/bootstrap/ldif/test.ldif", 0644)
	if err != nil {
		t.Fatalf("Failed to copy directory.ldif to container: %v", err)
	}

	exitCode, output, err := ldapC.Exec(ctx, addCmd)
	buf := new(bytes.Buffer)

	if output != nil {
		_, _ = buf.ReadFrom(output)
	}

	outputStr := buf.String()
	t.Logf("ldapadd exit code: %d, error: %v, output: %s", exitCode, err, outputStr)

	if err != nil || exitCode != 0 {
		t.Fatalf("Failed to exec ldapadd: %v, exit code: %d, output: %s", err, exitCode, outputStr)
	}

	t.Logf("ldapadd output: %s", outputStr)

	time.Sleep(500 * time.Millisecond)

	// Set up the adapter port
	adapterPort := 54323
	tmpDir := t.TempDir()
	authTokensPath := tmpDir + "/fake-auth-tokens"

	t.Setenv("LDAP_ADAPTER_PORT", fmt.Sprintf("%d", adapterPort))
	// Set up auth tokens for the adapter
	t.Setenv("AUTH_TOKENS_PATH", authTokensPath)
	_ = os.WriteFile(authTokensPath, []byte("[\"test-token\"]"), 0644)

	go func() {
		main()
	}()
	time.Sleep(500 * time.Millisecond)

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", adapterPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial grpc server: %v", err)
	}

	defer conn.Close()
	client := api_adapter_v1.NewAdapterClient(conn)

	ldapConfig := map[string]any{
		"baseDN": "dc=example,dc=org",
		"entityConfig": map[string]any{
			"Group": map[string]any{
				"query": "(&(objectClass=groupofuniquenames)(cn=Science))",
			},
			"GroupMember": map[string]any{
				"memberOf":                  "Group",
				"collectionAttribute":       "cn",
				"query":                     "(&(objectClass=groupofuniquenames)({{CollectionAttribute}}=Science))",
				"memberUniqueIDAttribute":   "dn",
				"memberOfUniqueIDAttribute": "dn",
				"memberAttribute":           "uniqueMember",
				"memberOfGroupBatchSize":    10,
			},
		},
	}

	configBytes, err := json.Marshal(ldapConfig)
	if err != nil {
		t.Fatalf("failed to marshal ldap config: %v", err)
	}

	// Act
	// First page request
	req := &api_adapter_v1.GetPageRequest{
		Datasource: &api_adapter_v1.DatasourceConfig{
			Id:      "test-ldap",
			Address: fmt.Sprintf("%s:%s", ldapHost, ldapPort.Port()),
			Type:    "LDAP-1.0.0",
			Config:  configBytes,
			Auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_Basic_{
					Basic: &api_adapter_v1.DatasourceAuthCredentials_Basic{
						Username: "cn=admin,dc=example,dc=org",
						Password: "admin",
					},
				},
			},
		},
		Entity: &api_adapter_v1.EntityConfig{
			Id:         "GroupMember",
			ExternalId: "GroupMember",
			Attributes: []*api_adapter_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
					UniqueId:   true,
				},
				{
					Id:         "group_dn",
					ExternalId: "group_dn",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
				{
					Id:         "member_dn",
					ExternalId: "member_dn",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					List:       false,
				},
			},
		},
		PageSize: 2,
	}

	t.Logf("LDAP Address used in DatasourceConfig: %s", req.Datasource.Address)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	md := metadata.New(map[string]string{"token": "test-token"})
	ctxWithToken := metadata.NewOutgoingContext(ctxTimeout, md)

	resp, err := client.GetPage(ctxWithToken, req)

	// Assert (First Page)
	if err != nil {
		t.Fatalf("GetPage (first page) failed: %v", err)
	}

	page := resp.GetSuccess()

	if page == nil || len(page.Objects) != 2 {
		t.Fatalf("expected two group members on first page, got: %+v", resp)
	}

	wantGroupMemberIDSlice := []string{"uid=john,ou=users,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org",
		"uid=alice,ou=users,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org"}

	var gotGroupMemberIDSlice []string

	for _, obj := range page.Objects {
		for _, attr := range obj.Attributes {
			for _, value := range attr.Values {
				if attr.Id == "id" {
					if strings.Contains(value.GetStringValue(), "<nil>") {
						t.Fatalf("id contains nil indicating no group id for the member, got: %s", value.GetStringValue())
					}
					// Collect the group member IDs
					gotGroupMemberIDSlice = append(gotGroupMemberIDSlice, value.GetStringValue())
				}
			}
		}
	}

	if !reflect.DeepEqual(gotGroupMemberIDSlice, wantGroupMemberIDSlice) {
		t.Fatalf("gotResponse: %v, wantResponse: %v", gotGroupMemberIDSlice, wantGroupMemberIDSlice)
	}

	if page.NextCursor != "" {
		t.Fatalf("expected empty next_cursor but got: %s", page.NextCursor)
	}
}
