// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package ldap_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	ldap_v3 "github.com/go-ldap/ldap/v3"
	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap/v2.0.0"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
	"github.com/testcontainers/testcontainers-go"
)

type LDAPTestSuite struct {
	testutil.CommonSuite
	ldapContainer testcontainers.Container
	ldapHost      string
	ldapPort      nat.Port
	ctx           context.Context
}

func Test_LDAPTestSuite(t *testing.T) {
	testutil.Run(t, new(LDAPTestSuite))
}

func (s *LDAPTestSuite) SetupSuite() {
	var cancel context.CancelFunc

	s.ctx, cancel = context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	s.ldapContainer, s.ldapPort = s.StartLDAPServer(s.ctx, false)
	s.ldapHost = "localhost:" + s.ldapPort.Port()

	time.Sleep(10 * time.Second)
}

func (s *LDAPTestSuite) TearDownSuite() {
	s.ldapContainer.Terminate(s.ctx)
}

func (s *LDAPTestSuite) Test_AdapterGetPage() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)

	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[ldap_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"invalid_request_with_invalid_creds": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth: &framework.DatasourceAuthCredentials{
					Basic: &framework.BasicAuthCredentials{
						Username: "cn=user,dc=example,dc=org",
						Password: "asdasd",
					},
				},
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Failed to bind credentials: LDAP Result Code 49 \"Invalid Credentials\": .",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
				},
			},
		},
		"invalid_request_with_missing_config": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Active Directory config is invalid: baseDN is not set.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func (s *LDAPTestSuite) Test_AdapterGetUserPage() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)
	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[ldap_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=inetOrgPerson)(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "uid",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
						{
							ExternalId: "mobile",
							Type:       framework.AttributeTypeString,
							List:       true,
						},
						// TODO: add tests for remaining types once we able to load sample data
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"dn":     "cn=marpontes,ou=People,dc=example,dc=org",
							"uid":    int64(1001),
							"mobile": []string{"+1 408 555 1234", "+1 408 555 4564"},
						},
						{
							"dn":     "cn=zach,ou=People,dc=example,dc=org",
							"uid":    int64(1002),
							"mobile": []string{"+1 408 555 8933", "+1 408 555 2722"},
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpqYjJ4c1pXTjBhVzl1SWpwdWRXeHNMQ0p1WlhoMFVHRm5aVU4xY25OdmNpSTZJa1JSUVVGQlFVRkJRVUZCUFNKOSJ9",
				},
			},
		},
		"valid_request_no_result": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=lorem)(objectClass=lorem))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "uid",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
						{
							ExternalId: "mobile",
							Type:       framework.AttributeTypeString,
							List:       true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func (s *LDAPTestSuite) Test_AdapterGetGroupPage() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)
	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[ldap_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=groupofuniquenames)",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Group",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"dn": "cn=Administrator,ou=Groups,dc=example,dc=org",
						},
						{
							"dn": "cn=Developers,ou=Groups,dc=example,dc=org",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpqYjJ4c1pXTjBhVzl1SWpwdWRXeHNMQ0p1WlhoMFVHRm5aVU4xY25OdmNpSTZJa1JSUVVGQlFVRkJRVUZCUFNKOSJ9",
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func (s *LDAPTestSuite) Test_AdapterGetGroupMemberPage() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)
	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[ldap_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"valid_request": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=groupofuniquenames)",
						},
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							CollectionAttribute:       testutil.GenPtr("entryDN"),
							Query:                     "(&(objectClass=groupofuniquenames)({{CollectionAttribute}}={{CollectionId}}))",
							MemberUniqueIDAttribute:   testutil.GenPtr("dn"),
							MemberOfUniqueIDAttribute: testutil.GenPtr("dn"),
							MemberAttribute:           testutil.GenPtr("uniqueMember"),
							MemberOfGroupBatchSize:    10,
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "group_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "member_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "cn=marpontes,ou=People,dc=example,dc=org-cn=Administrator,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Administrator,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=marpontes,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=leonardo,ou=People,dc=example,dc=org-cn=Administrator,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Administrator,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=leonardo,ou=People,dc=example,dc=org",
						},
					},
					NextCursor: "eyJjdXJzb3IiOiJleUpqYjJ4c1pXTjBhVzl1SWpwdWRXeHNMQ0p1WlhoMFVHRm5aVU4xY25OdmNpSTZiblZzYkN3aWJtVjRkRWR5YjNWd1VISnZZMlZ6YzJWa0lqb2lZMjQ5UkdWMlpXeHZjR1Z5Y3l4dmRUMUhjbTkxY0hNc1pHTTlaWGhoYlhCc1pTeGtZejF2Y21jaWZRPT0ifQ==",
				},
			},
		},
		"valid_request_all_pages": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=groupofuniquenames)",
						},
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							CollectionAttribute:       testutil.GenPtr("entryDN"),
							Query:                     "(&(objectClass=groupofuniquenames)({{CollectionAttribute}}={{CollectionId}}))",
							MemberUniqueIDAttribute:   testutil.GenPtr("dn"),
							MemberOfUniqueIDAttribute: testutil.GenPtr("dn"),
							MemberAttribute:           testutil.GenPtr("uniqueMember"),
							MemberOfGroupBatchSize:    10,
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "group_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "member_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 20,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "cn=marpontes,ou=People,dc=example,dc=org-cn=Administrator,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Administrator,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=marpontes,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=leonardo,ou=People,dc=example,dc=org-cn=Administrator,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Administrator,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=leonardo,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=marpontes,ou=People,dc=example,dc=org-cn=Developers,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Developers,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=marpontes,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=zach,ou=People,dc=example,dc=org-cn=Developers,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Developers,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=zach,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=leonardo,ou=People,dc=example,dc=org-cn=Developers,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Developers,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=leonardo,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=lorem,ou=People,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Science,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=lorem,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=bobby,ou=People,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Science,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=bobby,ou=People,dc=example,dc=org",
						},
					},
					NextCursor: "",
				},
			},
		},
		"valid_request_no_parent_objects": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(objectClass=missing)",
						},
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							CollectionAttribute:       testutil.GenPtr("entryDN"),
							Query:                     "(&(objectClass=groupofuniquenames)({{CollectionAttribute}}={{CollectionId}}))",
							MemberUniqueIDAttribute:   testutil.GenPtr("dn"),
							MemberOfUniqueIDAttribute: testutil.GenPtr("dn"),
							MemberAttribute:           testutil.GenPtr("uniqueMember"),
							MemberOfGroupBatchSize:    10,
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "group_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "member_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Error: &framework.Error{
					Message: "Datasource rejected request, returned status code: 404.",
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				},
			},
		},
		"valid_request_check_last_page_group_id": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Group": {
							Query: "(&(objectClass=groupofuniquenames)(cn=Science))",
						},
						"GroupMember": {
							MemberOf:                  testutil.GenPtr("Group"),
							CollectionAttribute:       testutil.GenPtr("entryDN"),
							Query:                     "(&(objectClass=groupofuniquenames)({{CollectionAttribute}}=cn=Science,ou=Groups,dc=example,dc=org))",
							MemberUniqueIDAttribute:   testutil.GenPtr("dn"),
							MemberOfUniqueIDAttribute: testutil.GenPtr("dn"),
							MemberAttribute:           testutil.GenPtr("uniqueMember"),
							MemberOfGroupBatchSize:    10,
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "GroupMember",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "id",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							ExternalId: "group_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							ExternalId: "member_dn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"id":        "cn=lorem,ou=People,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Science,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=lorem,ou=People,dc=example,dc=org",
						},
						{
							"id":        "cn=bobby,ou=People,dc=example,dc=org-cn=Science,ou=Groups,dc=example,dc=org",
							"group_dn":  "cn=Science,ou=Groups,dc=example,dc=org",
							"member_dn": "cn=bobby,ou=People,dc=example,dc=org",
						},
					},
					NextCursor: "",
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					t.Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			gotResponse := adapter.GetPage(tt.ctx, tt.request)
			if tt.wantResponse.Success != nil && gotResponse.Success != nil {
				if diff := cmp.Diff(tt.wantResponse.Success.Objects, gotResponse.Success.Objects); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tt.wantResponse.Success.NextCursor, gotResponse.Success.NextCursor); diff != "" {
					t.Errorf("Response mismatch (-want +got):\n%s", diff)
				}

				if !reflect.DeepEqual(gotResponse.Success.Objects, tt.wantResponse.Success.Objects) {
					t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
				}
			} else if tt.wantResponse.Success != nil || gotResponse.Success != nil {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse, tt.wantResponse)
			}

			if !reflect.DeepEqual(gotResponse.Error, tt.wantResponse.Error) {
				t.Errorf("gotResponse: %v, wantResponse: %v", gotResponse.Error, tt.wantResponse.Error)
			}

			// We already check the b64 encoded cursor in the response, but it's not easy to
			// decipher the cursor just by reading the test case.
			// So in addition, decode the b64 cursor and compare structs.
			if gotResponse.Success != nil && tt.wantCursor != nil {
				var gotCursor pagination.CompositeCursor[string]

				decodedCursor, err := base64.StdEncoding.DecodeString(gotResponse.Success.NextCursor)
				if err != nil {
					t.Errorf("error decoding cursor: %v", err)
				}

				if err := json.Unmarshal(decodedCursor, &gotCursor); err != nil {
					t.Errorf("error unmarshalling cursor: %v", err)
				}

				if !reflect.DeepEqual(&gotCursor, tt.wantCursor) {
					t.Errorf("gotCursor: %v, wantCursor: %v", gotCursor, tt.wantCursor)
				}
			}
		})
	}
}

func (s *LDAPTestSuite) Test_HostnameValidation() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)

	// Wait for LDAP server to be ready
	time.Sleep(10 * time.Second)

	tests := map[string]struct {
		ctx          context.Context
		request      *framework.Request[ldap_adapter.Config]
		wantErrCode  api_adapter_v1.ErrorCode
		wantResponse framework.Response
	}{
		"invalid_hostname_with_port": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: "ldaps://localhost:636",
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN:           "dc=example,dc=org",
					CertificateChain: validCertificateChain,
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantErrCode: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
		"valid_hostname_without_port": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: "localhost:" + s.ldapPort.Port(),
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							Query: "(&(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				PageSize: 2,
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{},
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(t *testing.T) {
			response := adapter.GetPage(tt.ctx, tt.request)

			if tt.wantErrCode != 0 {
				if response.Error == nil {
					t.Fatal("expected error, got nil")
				}

				if response.Error.Code != tt.wantErrCode {
					t.Errorf("expected error code %v, got %v", tt.wantErrCode, response.Error.Code)
				}

				return
			}

			if response.Error != nil {
				t.Fatalf("expected success, got error: %v", response.Error)
			}

			if response.Success == nil {
				t.Fatal("expected success response, got nil")
			}
		})
	}
}

func (s *LDAPTestSuite) Test_AdapterGetPage_MissingAttributes_PanicRegression() {
	adapter := ldap_adapter.NewAdapter(nil, time.Minute, time.Minute)
	tests := map[string]struct {
		ctx                context.Context
		request            *framework.Request[ldap_adapter.Config]
		inputRequestCursor *pagination.CompositeCursor[string]
		wantResponse       framework.Response
		wantCursor         *pagination.CompositeCursor[string]
	}{
		"request_missing_attributes_should_not_panic": {
			ctx: context.Background(),
			request: &framework.Request[ldap_adapter.Config]{
				Address: s.ldapHost,
				Auth:    validAuthCredentials,
				Config: &ldap_adapter.Config{
					BaseDN: "dc=example,dc=org",
					EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
						"Person": {
							// Query all person objects - some may not have all requested attributes
							Query: "(&(objectClass=person))",
						},
					},
				},
				Entity: framework.EntityConfig{
					ExternalId: "Person",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "dn",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
						{
							// This attribute may not exist on all entries - testing empty Values scenario
							ExternalId: "nonExistentAttr",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							// This binary attribute doesn't exist in our test LDIF - testing empty ByteValues scenario
							ExternalId: "objectGUID",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							// Another binary attribute that doesn't exist - testing empty ByteValues scenario
							ExternalId: "objectSid",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
						{
							// Testing with different attribute types that may be missing
							ExternalId: "missingBoolAttr",
							Type:       framework.AttributeTypeBool,
							List:       false,
						},
						{
							ExternalId: "missingIntAttr",
							Type:       framework.AttributeTypeInt64,
							List:       false,
						},
						{
							ExternalId: "missingDoubleAttr",
							Type:       framework.AttributeTypeDouble,
							List:       false,
						},
						{
							ExternalId: "missingDateAttr",
							Type:       framework.AttributeTypeDateTime,
							List:       false,
						},
						{
							// Test existing attribute to ensure normal case still works
							ExternalId: "cn",
							Type:       framework.AttributeTypeString,
							List:       false,
						},
					},
				},
				PageSize: 5, // Get multiple entries to test various scenarios
			},
			wantResponse: framework.Response{
				Success: &framework.Page{
					// Should return entries successfully even with missing attributes
					// The key test is that it doesn't panic - exact objects may vary
					Objects: []framework.Object{
						// We expect at least some entries with dn and cn, missing attrs should be nil or absent
					},
				},
			},
		},
	}

	for name, tt := range tests {
		s.T().Run(name, func(_ *testing.T) {
			if tt.inputRequestCursor != nil {
				encodedCursor, err := pagination.MarshalCursor(tt.inputRequestCursor)
				if err != nil {
					s.T().Error(err)
				}

				tt.request.Cursor = encodedCursor
			}

			// This is the key test - it should not panic
			gotResponse := adapter.GetPage(tt.ctx, tt.request)

			// If we got an error, it should be a proper framework error, not a panic
			if gotResponse.Error != nil {
				s.T().Logf("Got expected error (not panic): %v", gotResponse.Error.Message)
				// Error is acceptable as long as we didn't panic
				return
			}

			// If successful, verify we got some objects
			if gotResponse.Success == nil {
				s.T().Errorf("got nil success response when error was also nil")

				return
			}

			if len(gotResponse.Success.Objects) == 0 {
				s.T().Errorf("expected at least some objects in successful response")
			}

			// Log the objects to see what we got
			s.T().Logf("Successfully retrieved %d objects without panic", len(gotResponse.Success.Objects))

			for i, obj := range gotResponse.Success.Objects {
				s.T().Logf("Object %d: %+v", i, obj)
			}
		})
	}
}

// Test_EmptyAttributeValues_DirectPanicRegression tests for the specific case where
// empty attribute values cause a panic.
func (s *LDAPTestSuite) Test_EmptyAttributeValues_DirectPanicRegression() {
	tests := []struct {
		name        string
		attr        *ldap_v3.EntryAttribute
		attrType    framework.AttributeType
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty_string_attr_values",
			attr: &ldap_v3.EntryAttribute{
				Name:   "description",
				Values: []string{}, // This is the condition that caused the original panic
			},
			attrType:    framework.AttributeTypeString,
			expectError: false,
		},
		{
			name: "empty_values_with_requested_bool",
			attr: &ldap_v3.EntryAttribute{
				Name:   "accountDisabled",
				Values: []string{}, // Empty values that would cause panic on attr.Values[0]
			},
			attrType:    framework.AttributeTypeBool,
			expectError: false,
		},
		{
			name: "empty_values_with_requested_int",
			attr: &ldap_v3.EntryAttribute{
				Name:   "employeeNumber",
				Values: []string{}, // Empty values that would cause panic
			},
			attrType:    framework.AttributeTypeInt64,
			expectError: false,
		},
		{
			name: "empty_values_with_requested_double",
			attr: &ldap_v3.EntryAttribute{
				Name:   "salary",
				Values: []string{}, // Empty values that would cause panic
			},
			attrType:    framework.AttributeTypeDouble,
			expectError: false,
		},
		{
			name: "empty_values_with_requested_datetime",
			attr: &ldap_v3.EntryAttribute{
				Name:   "whenCreated",
				Values: []string{}, // Empty values that would cause panic
			},
			attrType:    framework.AttributeTypeDateTime,
			expectError: false,
		},
		{
			name: "empty_values_with_requested_duration",
			attr: &ldap_v3.EntryAttribute{
				Name:   "lockoutDuration",
				Values: []string{}, // Empty values that would cause panic
			},
			attrType:    framework.AttributeTypeDuration,
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// This directly calls the function that was panicking
			result, err := ldap_adapter.StringAttrValuesToRequestedType(tt.attr, false, tt.attrType)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")

					return
				}

				if err.Message != tt.errorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMsg, err.Message)
				}

				if err.Code != api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE {
					t.Errorf("expected error code ERROR_CODE_INVALID_ATTRIBUTE_TYPE, got %v", err.Code)
				}

				if result != nil {
					t.Errorf("expected nil result when error occurs, got %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// For non-list types with empty values, should return empty string
				if result != "" {
					t.Errorf("expected empty string, got %v", result)
				}
			}
		})
	}
}
