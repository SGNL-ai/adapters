// Copyright 2025 SGNL.ai, Inc.
package smoketests

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/smoketests/common"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestSCIMAdapter_User_Sailpoint(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/scim/user_sailpoint")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	req := adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_Basic_{
					Basic: &adapter_api_v1.DatasourceAuthCredentials_Basic{
						Username: "{{OMITTED}}",
						Password: "{{OMITTED}}",
					},
				},
			},
			Address: "https://example-tenant.example-domain.com:8080/identityiq/scim/v2",
			Id:      "SCIM",
			Type:    "SCIM2.0-1.0.0",
			Config: []byte(`{
				"scimProtocolVersion": 2
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SCIMUser",
			ExternalId: "Users",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userName",
					ExternalId: "userName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"id": "id",
							"values": [
								{
									"string_value": "0a0000fa80ce10c68180ce11410000ff"
								}
							]
						},
						{
							"id": "userName",
							"values": [
								{
									"string_value": "spadmin"
								}
							]
						}
					]
				}
			],
			"nextCursor": "2"
		}
	}`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

func TestSCIMAdapter_User_Curity(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/scim/user_curity")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	req := adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "{{OMITTED}}",
				},
			},
			Address: "https://curity-aws:8443/user-management",
			Id:      "SCIM",
			Type:    "SCIM2.0-1.0.0",
			Config: []byte(`{
				"scimProtocolVersion": 2
			}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SCIMUser",
			ExternalId: "Users",
			Ordered:    false,
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userName",
					ExternalId: "userName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			ChildEntities: nil,
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"id": "id",
							"values": [
								{
									"string_value": "VVNFUjo1OTMyMTg5Ny1iMDNiLTRjNjYtOTY2MS02MDM5ZmU3OThmNmE"
								}
							]
						},
						{
							"id": "userName",
							"values": [
								{
									"string_value": "Rakshith"
								}
							]
						}
					]
				}
			],
			"nextCursor": "2"
		}
	}`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(gotResp, wantResp, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}

// Setup a fictitious server fixture that returns a Full Enterprise User Extension Representation.
// The response from this server is mocked referencing the following RFCs.
// https://datatracker.ietf.org/doc/html/rfc7643#section-8.3
// https://datatracker.ietf.org/doc/html/rfc7643#section-8.7.1
func TestSCIMAdapter_User_Fictitious_Server(t *testing.T) {
	httpClient, recorder := common.StartRecorder(t, "fixtures/scim/user_full_representation_mock_server")
	defer recorder.Stop()

	port := common.AvailableTestPort(t)

	stop := make(chan struct{})

	// Start Adapter Server
	go func() {
		stop = common.StartAdapterServer(t, httpClient, port)
	}()

	time.Sleep(10 * time.Millisecond)

	adapterClient, conn := common.GetNewAdapterClient(t, port)
	defer conn.Close()

	ctx, cancelCtx := common.GetAdapterCtx()
	defer cancelCtx()

	req := adapter_api_v1.GetPageRequest{
		Datasource: &adapter_api_v1.DatasourceConfig{
			Auth: &adapter_api_v1.DatasourceAuthCredentials{
				AuthMechanism: &adapter_api_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "{{OMITTED}}",
				},
			},
			Address: "https://fictitious-scim.server:8080/identityiq/scim/v2",
			Id:      "SCIM",
			Type:    "SCIM2.0-1.0.0",
			Config: []byte(`{
			"scimProtocolVersion": 2
		}`),
		},
		Entity: &adapter_api_v1.EntityConfig{
			Id:         "SCIMUser",
			ExternalId: "Users",
			Ordered:    false,
			// all externalIds of attributes represented as a string
			Attributes: []*adapter_api_v1.AttributeConfig{
				{
					Id:         "id",
					ExternalId: "id",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "externalId",
					ExternalId: "externalId",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userName",
					ExternalId: "userName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "formattedName",
					ExternalId: "$.name.formatted",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "familyName",
					ExternalId: "$.name.familyName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "givenName",
					ExternalId: "$.name.givenName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "middleName",
					ExternalId: "$.name.middleName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "honorificPrefix",
					ExternalId: "$.name.honorificPrefix",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "honorificSuffix",
					ExternalId: "$.name.honorificSuffix",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "displayName",
					ExternalId: "displayName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "nickName",
					ExternalId: "nickName",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "profileUrl",
					ExternalId: "profileUrl",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "userType",
					ExternalId: "userType",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "title",
					ExternalId: "title",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "preferredLanguage",
					ExternalId: "preferredLanguage",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "locale",
					ExternalId: "locale",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "timezone",
					ExternalId: "timezone",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "active",
					ExternalId: "active",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
				},
				{
					Id:         "password",
					ExternalId: "password",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "employeeNumber",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].employeeNumber`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "costCenter",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].costCenter`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "organization",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].organization`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "division",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].division`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "department",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].department`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "managerId",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].manager.value`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "managerDisplayName",
					ExternalId: `$..["urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"].manager.displayName`,
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "resourceType",
					ExternalId: "$.meta.resourceType",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "created",
					ExternalId: "$.meta.created",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
				},
				{
					Id:         "lastModified",
					ExternalId: "$.meta.lastModified",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME,
				},
				{
					Id:         "version",
					ExternalId: "$.meta.version",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
				{
					Id:         "location",
					ExternalId: "$.meta.location",
					Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},

			// complex multivalued attributes
			// emails, addresses, phoneNumbers, ims, photos, groups, entitlements, roles, x509Certificates
			ChildEntities: []*adapter_api_v1.EntityConfig{
				{
					Id:         "SCIMUserEmail",
					ExternalId: "emails",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserAddress",
					ExternalId: "addresses",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "formatted",
							ExternalId: "formatted",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "streetAddress",
							ExternalId: "streetAddress",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "locality",
							ExternalId: "locality",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "region",
							ExternalId: "region",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "postalCode",
							ExternalId: "postalCode",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "country",
							ExternalId: "country",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				{
					Id:         "SCIMUserPhoneNumbers",
					ExternalId: "phoneNumbers",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserIMS",
					ExternalId: "ims",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserPhotos",
					ExternalId: "photos",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserGroups",
					ExternalId: "groups",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "ref",
							ExternalId: `..["$ref"]`,
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				{
					Id:         "SCIMUserEntitlements",
					ExternalId: "entitlements",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserRoles",
					ExternalId: "roles",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
				{
					Id:         "SCIMUserCertificates",
					ExternalId: "x509Certificates",
					Ordered:    false,
					Attributes: []*adapter_api_v1.AttributeConfig{
						{
							Id:         "value",
							ExternalId: "value",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "display",
							ExternalId: "display",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "type",
							ExternalId: "type",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
						{
							Id:         "primary",
							ExternalId: "primary",
							Type:       adapter_api_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
						},
					},
				},
			},
		},
		PageSize: 1,
		Cursor:   "",
	}

	gotResp, err := adapterClient.GetPage(ctx, &req)
	if err != nil {
		t.Fatal(err)
	}

	wantResp := new(adapter_api_v1.GetPageResponse)

	err = protojson.Unmarshal([]byte(`{
		"success": {
			"objects": [
				{
					"attributes": [
						{
							"id": "costCenter",
							"values": [
								{
									"string_value": "4130"
								}
							]
						},
						{
							"id": "department",
							"values": [
								{
									"string_value": "Tour Operations"
								}
							]
						},
						{
							"id": "division",
							"values": [
								{
									"string_value": "Theme Park"
								}
							]
						},
						{
							"id": "employeeNumber",
							"values": [
								{
									"string_value": "701984"
								}
							]
						},
						{
							"id": "managerDisplayName",
							"values": [
								{
									"string_value": "John Smith"
								}
							]
						},
						{
							"id": "managerId",
							"values": [
								{
									"string_value": "26118915-6090-4610-87e4-49d8ca9f808d"
								}
							]
						},
						{
							"id": "organization",
							"values": [
								{
									"string_value": "Universal Studios"
								}
							]
						},
						{
							"id": "created",
							"values": [
								{
									"datetime_value": {
										"timestamp": "2010-01-23T04:56:22.000000000Z",
										"timezone_offset": 0
									}
								}
							]
						},
						{
							"id": "lastModified",
							"values": [
								{
									"datetime_value": {
										"timestamp": "2011-05-13T04:42:34.000000000Z",
										"timezone_offset": 0
									}
								}
							]
						},
						{
							"id": "location",
							"values": [
								{
									"string_value": "https://example.com/v2/Users/2819c223-7f76-453a-919d-413861904646"
								}
							]
						},
						{
							"id": "resourceType",
							"values": [
								{
									"string_value": "User"
								}
							]
						},
						{
							"id": "version",
							"values": [
								{
									"string_value": "W/\"3694e05e9dff591\""
								}
							]
						},
						{
							"id": "id",
							"values": [
								{
									"string_value": "2819c223-7f76-453a-919d-413861904000"
								}
							]
						},
						{
							"id": "externalId",
							"values": [
								{
									"string_value": "701984"
								}
							]
						},
						{
							"id": "userName",
							"values": [
								{
									"string_value": "bjensen@example.com"
								}
							]
						},
						{
							"id": "formattedName",
							"values": [
								{
									"string_value": "Ms. Barbara J Jensen, III"
								}
							]
						},
						{
							"id": "familyName",
							"values": [
								{
									"string_value": "Jensen"
								}
							]
						},
						{
							"id": "givenName",
							"values": [
								{
									"string_value": "Barbara"
								}
							]
						},
						{
							"id": "middleName",
							"values": [
								{
									"string_value": "Jane"
								}
							]
						},
						{
							"id": "honorificPrefix",
							"values": [
								{
									"string_value": "Ms."
								}
							]
						},
						{
							"id": "honorificSuffix",
							"values": [
								{
									"string_value": "III"
								}
							]
						},
						{
							"id": "displayName",
							"values": [
								{
									"string_value": "Babs Jensen"
								}
							]
						},
						{
							"id": "nickName",
							"values": [
								{
									"string_value": "Babs"
								}
							]
						},
						{
							"id": "profileUrl",
							"values": [
								{
									"string_value": "https://login.example.com/bjensen"
								}
							]
						},
						{
							"id": "userType",
							"values": [
								{
									"string_value": "Employee"
								}
							]
						},
						{
							"id": "title",
							"values": [
								{
									"string_value": "Tour Guide"
								}
							]
						},
						{
							"id": "preferredLanguage",
							"values": [
								{
									"string_value": "en-US"
								}
							]
						},
						{
							"id": "locale",
							"values": [
								{
									"string_value": "en-US"
								}
							]
						},
						{
							"id": "timezone",
							"values": [
								{
									"string_value": "America/Los_Angeles"
								}
							]
						},
						{
							"id": "active",
							"values": [
								{
									"bool_value": true
								}
							]
						},
						{
							"id": "password",
							"values": [
								{
									"string_value": "OMITTED"
								}
							]
						}
					],
					"child_objects": [
						{
							"entity_id": "SCIMUserAddress",
							"objects": [
								{
									"attributes": [
										{
											"id": "country",
											"values": [
												{
													"string_value": "USA"
												}
											]
										},
										{
											"id": "formatted",
											"values": [
												{
													"string_value": "100 Universal City Plaza\nHollywood, CA 91608 USA"
												}
											]
										},
										{
											"id": "locality",
											"values": [
												{
													"string_value": "Hollywood"
												}
											]
										},
										{
											"id": "postalCode",
											"values": [
												{
													"string_value": "91608"
												}
											]
										},
										{
											"id": "region",
											"values": [
												{
													"string_value": "CA"
												}
											]
										},
										{
											"id": "streetAddress",
											"values": [
												{
													"string_value": "100 Universal City Plaza"
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "work"
												}
											]
										}
									]
								},
								{
									"attributes": [
										{
											"id": "country",
											"values": [
												{
													"string_value": "USA"
												}
											]
										},
										{
											"id": "formatted",
											"values": [
												{
													"string_value": "456 Hollywood Blvd\nHollywood, CA 91608 USA"
												}
											]
										},
										{
											"id": "locality",
											"values": [
												{
													"string_value": "Hollywood"
												}
											]
										},
										{
											"id": "postalCode",
											"values": [
												{
													"string_value": "91608"
												}
											]
										},
										{
											"id": "region",
											"values": [
												{
													"string_value": "CA"
												}
											]
										},
										{
											"id": "streetAddress",
											"values": [
												{
													"string_value": "456 Hollywood Blvd"
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "home"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserEmail",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "bjensen@example.com"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "work"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "bjensen@example.com"
												}
											]
										}
									]
								},
								{
									"attributes": [
										{
											"id": "type",
											"values": [
												{
													"string_value": "home"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "babs@jensen.org"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserEntitlements",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "E1"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "entitlement"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "e9e30dba-f08f-4109-8486-d5c6a331abc"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserPhoneNumbers",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "someaimhandle"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "aim"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "someaimhandle"
												}
											]
										}
									]
								},
								{
									"attributes": [
										{
											"id": "type",
											"values": [
												{
													"string_value": "mobile"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "555-555-4444"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserPhotos",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "https://photos.example.com/profilephoto/72930000000Ccne/F"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "photo"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "https://photos.example.com/profilephoto/72930000000Ccne/F"
												}
											]
										}
									]
								},
								{
									"attributes": [
										{
											"id": "type",
											"values": [
												{
													"string_value": "thumbnail"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "https://photos.example.com/profilephoto/72930000000Ccne/T"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserRoles",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "Role A"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "role"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "e9e30dba-f08f-4109-8486-d5c6a33role"
												}
											]
										}
									]
								}
							]
						},
						{
							"entity_id": "SCIMUserCertificates",
							"objects": [
								{
									"attributes": [
										{
											"id": "display",
											"values": [
												{
													"string_value": "SGNLCertificate"
												}
											]
										},
										{
											"id": "primary",
											"values": [
												{
													"bool_value": true
												}
											]
										},
										{
											"id": "type",
											"values": [
												{
													"string_value": "secret"
												}
											]
										},
										{
											"id": "value",
											"values": [
												{
													"string_value": "some_certificate_value"
												}
											]
										}
									]
								}
							]
						}
					]
				}
			],
			"nextCursor": ""
		}
	}`), wantResp)
	if err != nil {
		t.Fatal(err)
	}

	wantAttrs := wantResp.GetSuccess().GetObjects()[0].GetAttributes()
	gotAttrs := gotResp.GetSuccess().GetObjects()[0].GetAttributes()
	wantChildObjs := wantResp.GetSuccess().GetObjects()[0].GetChildObjects()
	gotChildObjs := wantResp.GetSuccess().GetObjects()[0].GetChildObjects()

	// The list of attributes is sorted by the externalID instead of the order specified in the entity config
	sort.Slice(wantAttrs, func(i, j int) bool {
		return wantAttrs[i].Id < wantAttrs[j].Id
	})
	sort.Slice(gotAttrs, func(i, j int) bool {
		return gotAttrs[i].Id < gotAttrs[j].Id
	})

	if diff := cmp.Diff(gotAttrs, wantAttrs, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	if diff := cmp.Diff(wantChildObjs, gotChildObjs, common.CmpOpts...); diff != "" {
		t.Fatal(diff)
	}

	close(stop)
}
