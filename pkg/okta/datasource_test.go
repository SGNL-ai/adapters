// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package okta_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/okta"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock Okta server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "SSWS testtoken" && r.Header.Get("Authorization") != "Bearer Testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"errorCode": "E0000011",
			"errorSummary": "Invalid token provided",
			"errorLink": "E0000011",
			"errorId": "oaefW5oDjyLRLKVkrmTlp0Thg",
			"errorCauses": []
		}}`))
	}

	switch r.URL.RequestURI() {
	// Users Page 1
	case "/api/v1/users?limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/users?limit=2>; rel="self"`)
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>; rel="next"`)
		w.Write([]byte(`[
				{
					"id": "00ub0oNGTSWTBKOLGLNR",
					"status": "ACTIVE",
					"created": "2013-06-24T16:39:18.000Z",
					"activated": "2013-06-24T16:39:19.000Z",
					"statusChanged": "2013-06-24T16:39:19.000Z",
					"lastLogin": "2013-06-24T17:39:19.000Z",
					"lastUpdated": "2013-07-02T21:36:25.344Z",
					"passwordChanged": "2013-07-02T21:36:25.344Z",
					"profile": {
						"firstName": "Isaac",
						"lastName": "Brock",
						"email": "isaac.brock@example.com",
						"login": "isaac.brock@example.com",
						"mobilePhone": "111-111-1111"
					},
					"credentials": {
						"password": {},
						"recovery_question": {
							"question": "What was your school's mascot?"
						},
						"provider": {
							"type": "OKTA",
							"name": "OKTA"
						}
					},
					"_links": {
						"self": {
							"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR"
						}
					}
				},
				{
					"id": "00ub0oNGTSWTBKOCNDJI",
					"status": "ACTIVE",
					"created": "2013-06-24T16:42:20.000Z",
					"activated": "2013-06-24T16:42:20.000Z",
					"statusChanged": "2013-06-24T16:42:20.000Z",
					"lastLogin": "2013-06-24T16:43:12.000Z",
					"lastUpdated": "2013-06-24T16:42:20.000Z",
					"passwordChanged": "2013-06-24T16:42:20.000Z",
					"profile": {
						"firstName": "John",
						"lastName": "Smith",
						"email": "john.smith@example.com",
						"login": "john.smith@example.com",
						"mobilePhone": "111-111-1111"
					},
					"credentials": {
						"password": {},
						"recovery_question": {
							"question": "What is your mother's maiden name?"
						},
						"provider": {
							"type": "OKTA",
							"name": "OKTA"
						}
					},
					"_links": {
						"self": {
							"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI"
						}
					}
				}
			]`))

	// Users Page 2
	case "/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>; rel="self"`)
		w.Write([]byte(`[
				{
					"id": "00ub0oNGTSWTBKOMSUFE",
					"status": "ACTIVE",
					"created": "2013-06-24T18:02:12.000Z",
					"activated": "2013-06-24T18:02:12.000Z",
					"statusChanged": "2013-06-24T18:02:12.000Z",
					"lastLogin": "2013-06-24T19:14:58.000Z",
					"lastUpdated": "2013-06-24T18:02:12.000Z",
					"passwordChanged": "2013-06-24T18:02:12.000Z",
					"profile": {
						"firstName": "Brooke",
						"lastName": "Pearson",
						"email": "brooke.pearson@example.com",
						"login": "brooke.pearson@example.com",
						"mobilePhone": "111-111-1111"
					},
					"credentials": {
						"password": {},
						"recovery_question": {
							"question": "What is your middle name?"
						},
						"provider": {
							"type": "OKTA",
							"name": "OKTA"
						}
					},
					"_links": {
						"self": {
							"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE"
						}
					}
				}
			]`))

	// Groups Page 1:
	case "/api/v1/groups?filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22&limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="next"`)
		w.Write([]byte(`[
			{
				"id": "00g1emaKYZTWRYYRRTSK",
				"created": "2015-02-06T10:11:28.000Z",
				"lastUpdated": "2015-10-05T19:16:43.000Z",
				"lastMembershipUpdated": "2015-11-28T19:15:32.000Z",
				"objectClass": [
					"okta:user_group"
				],
				"type": "OKTA_GROUP",
				"profile": {
					"name": "West Coast Users",
					"description": "All Users West of The Rockies"
				}
			},
			{
				"id": "00garwpuyxHaWOkdV0g4",
				"created": "2015-08-15T19:15:17.000Z",
				"lastUpdated": "2015-11-18T04:02:19.000Z",
				"lastMembershipUpdated": "2015-08-15T19:15:17.000Z",
				"objectClass": [
					"okta:windows_security_principal"
				],
				"type": "APP_GROUP",
				"profile": {
					"name": "Engineering Users",
					"description": "corp.example.com/Engineering/Engineering Users",
					"groupType": "Security",
					"samAccountName": "Engineering Users",
					"objectSid": "S-1-5-21-717838489-685202119-709183397-1177",
					"groupScope": "Global",
					"dn": "CN=Engineering Users,OU=Engineering,DC=corp,DC=example,DC=com",
					"windowsDomainQualifiedName": "CORP\\Engineering Users",
					"externalId": "OZJdWdONCU6h7WjQKp+LPA=="
				},
				"source": {
					"id": "0oa2v0el0gP90aqjJ0g7"
				}
			}
		]`))

	// Groups Page 2:
	case "/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Write([]byte(`[
			{
				"id": "00g1emaKYZTWRCCNDHEU",
				"created": "2015-02-06T10:11:28.000Z",
				"lastUpdated": "2015-10-05T19:16:43.000Z",
				"lastMembershipUpdated": "2015-11-28T19:15:32.000Z",
				"objectClass": [
					"okta:user_group"
				],
				"type": "OKTA_GROUP",
				"profile": {
					"name": "East Coast Users",
					"description": "All Users East of The Rockies"
				}
			}
		]`))

	// Group Members - Groups Page 1 (Page size 1):
	case "/api/v1/groups?filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22&limit=1":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="next"`)
		w.Write([]byte(`[
		{
			"id": "00g1emaKYZTWRYYRRTSK",
			"created": "2015-02-06T10:11:28.000Z",
			"lastUpdated": "2015-10-05T19:16:43.000Z",
			"lastMembershipUpdated": "2015-11-28T19:15:32.000Z",
			"objectClass": [
				"okta:user_group"
			],
			"type": "OKTA_GROUP",
			"profile": {
				"name": "West Coast Users",
				"description": "All Users West of The Rockies"
			}
		}
	]`))

	// Group Members - Groups Page 2 (Page size 1):
	case "/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Write([]byte(`[
			{
				"id": "00garwpuyxHaWOkdV0g4",
				"created": "2015-08-15T19:15:17.000Z",
				"lastUpdated": "2015-11-18T04:02:19.000Z",
				"lastMembershipUpdated": "2015-08-15T19:15:17.000Z",
				"objectClass": [
					"okta:windows_security_principal"
				],
				"type": "APP_GROUP",
				"profile": {
					"name": "Engineering Users",
					"description": "corp.example.com/Engineering/Engineering Users",
					"groupType": "Security",
					"samAccountName": "Engineering Users",
					"objectSid": "S-1-5-21-717838489-685202119-709183397-1177",
					"groupScope": "Global",
					"dn": "CN=Engineering Users,OU=Engineering,DC=corp,DC=example,DC=com",
					"windowsDomainQualifiedName": "CORP\\Engineering Users",
					"externalId": "OZJdWdONCU6h7WjQKp+LPA=="
				},
				"source": {
					"id": "0oa2v0el0gP90aqjJ0g7"
				}
			}
		]`))

	// Group Members - Invalid Groups (Page size 1 requested, 2 returned):
	case "/api/v1/groups?after=00garwpuyxHaWOkdV0g4&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00garwpuyxHaWOkdV0g4&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Write([]byte(`[
			{
				"id": "00garwpuyxHaWOkdV0g4"
			},
			{
				"id": "00garwpuyxHaWOkcndD8"
			}
		]`))

	// Group Members - Invalid Groups (Page size 1 requested, 0 returned):
	case "/api/v1/groups?after=00garwpuyxHaWOkcndD8&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22>; rel="self">`)
		w.Write([]byte(`[]`))

	// Group Members - 00g1emaKYZTWRYYRRTSK - Members Page 1:
	case "/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?limit=2>; rel="self">`)
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2>; rel="next"`)
		w.Write([]byte(`[
			{
				"id": "00ub0oNGTSWTBKOLGLNR",
				"status": "ACTIVE",
				"created": "2013-06-24T16:39:18.000Z",
				"activated": "2013-06-24T16:39:19.000Z",
				"statusChanged": "2013-06-24T16:39:19.000Z",
				"lastLogin": "2013-06-24T17:39:19.000Z",
				"lastUpdated": "2013-07-02T21:36:25.344Z",
				"passwordChanged": "2013-07-02T21:36:25.344Z",
				"profile": {
					"firstName": "Isaac",
					"lastName": "Brock",
					"email": "isaac.brock@example.com",
					"login": "isaac.brock@example.com",
					"mobilePhone": "111-111-1111"
				},
				"credentials": {
					"password": {},
					"recovery_question": {
						"question": "What was your school's mascot?"
					},
					"provider": {
						"type": "OKTA",
						"name": "OKTA"
					}
				},
				"_links": {
					"self": {
						"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR"
					}
				}
			},
			{
				"id": "00ub0oNGTSWTBKOCNDJI",
				"status": "ACTIVE",
				"created": "2013-06-24T16:42:20.000Z",
				"activated": "2013-06-24T16:42:20.000Z",
				"statusChanged": "2013-06-24T16:42:20.000Z",
				"lastLogin": "2013-06-24T16:43:12.000Z",
				"lastUpdated": "2013-06-24T16:42:20.000Z",
				"passwordChanged": "2013-06-24T16:42:20.000Z",
				"profile": {
					"firstName": "John",
					"lastName": "Smith",
					"email": "john.smith@example.com",
					"login": "john.smith@example.com",
					"mobilePhone": "111-111-1111"
				},
				"credentials": {
					"password": {},
					"recovery_question": {
						"question": "What is your mother's maiden name?"
					},
					"provider": {
						"type": "OKTA",
						"name": "OKTA"
					}
				},
				"_links": {
					"self": {
						"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI"
					}
				}
			}
		]`))

	// Group Members - 00g1emaKYZTWRYYRRTSK - Members Page 2:
	case "/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2>; rel="self">`)
		w.Write([]byte(`[
			{
				"id": "00ub0oNGTSWTBKOMSUFE",
				"status": "ACTIVE",
				"created": "2013-06-24T18:02:12.000Z",
				"activated": "2013-06-24T18:02:12.000Z",
				"statusChanged": "2013-06-24T18:02:12.000Z",
				"lastLogin": "2013-06-24T19:14:58.000Z",
				"lastUpdated": "2013-06-24T18:02:12.000Z",
				"passwordChanged": "2013-06-24T18:02:12.000Z",
				"profile": {
					"firstName": "Brooke",
					"lastName": "Pearson",
					"email": "brooke.pearson@example.com",
					"login": "brooke.pearson@example.com",
					"mobilePhone": "111-111-1111"
				},
				"credentials": {
					"password": {},
					"recovery_question": {
						"question": "What is your middle name?"
					},
					"provider": {
						"type": "OKTA",
						"name": "OKTA"
					}
				},
				"_links": {
					"self": {
						"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE"
					}
				}
			}
		]`))

	// Group Members - 00garwpuyxHaWOkdV0g4 - Members Page 1:
	case "/api/v1/groups/00garwpuyxHaWOkdV0g4/users?limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?limit=2>; rel="self">`)
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOCNDJI&limit=2>; rel="next"`)
		w.Write([]byte(`[
		{
			"id": "00ub0oNGTSWTBKOLGLNR",
			"status": "ACTIVE",
			"created": "2013-06-24T16:39:18.000Z",
			"activated": "2013-06-24T16:39:19.000Z",
			"statusChanged": "2013-06-24T16:39:19.000Z",
			"lastLogin": "2013-06-24T17:39:19.000Z",
			"lastUpdated": "2013-07-02T21:36:25.344Z",
			"passwordChanged": "2013-07-02T21:36:25.344Z",
			"profile": {
				"firstName": "Isaac",
				"lastName": "Brock",
				"email": "isaac.brock@example.com",
				"login": "isaac.brock@example.com",
				"mobilePhone": "111-111-1111"
			},
			"credentials": {
				"password": {},
				"recovery_question": {
					"question": "What was your school's mascot?"
				},
				"provider": {
					"type": "OKTA",
					"name": "OKTA"
				}
			},
			"_links": {
				"self": {
					"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR"
				}
			}
		},
		{
			"id": "00ub0oNGTSWTBKOCNDJI",
			"status": "ACTIVE",
			"created": "2013-06-24T16:42:20.000Z",
			"activated": "2013-06-24T16:42:20.000Z",
			"statusChanged": "2013-06-24T16:42:20.000Z",
			"lastLogin": "2013-06-24T16:43:12.000Z",
			"lastUpdated": "2013-06-24T16:42:20.000Z",
			"passwordChanged": "2013-06-24T16:42:20.000Z",
			"profile": {
				"firstName": "John",
				"lastName": "Smith",
				"email": "john.smith@example.com",
				"login": "john.smith@example.com",
				"mobilePhone": "111-111-1111"
			},
			"credentials": {
				"password": {},
				"recovery_question": {
					"question": "What is your mother's maiden name?"
				},
				"provider": {
					"type": "OKTA",
					"name": "OKTA"
				}
			},
			"_links": {
				"self": {
					"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI"
				}
			}
		}
	]`))

	// Group Members - 00garwpuyxHaWOkdV0g4 - Members Page 2:
	case "/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOCNDJI&limit=2":
		w.Header().Add("link", `<https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOCNDJI&limit=2>; rel="self">`)
		w.Write([]byte(`[
		{
			"id": "00ub0oNGTSWTBKOMSUFE",
			"status": "ACTIVE",
			"created": "2013-06-24T18:02:12.000Z",
			"activated": "2013-06-24T18:02:12.000Z",
			"statusChanged": "2013-06-24T18:02:12.000Z",
			"lastLogin": "2013-06-24T19:14:58.000Z",
			"lastUpdated": "2013-06-24T18:02:12.000Z",
			"passwordChanged": "2013-06-24T18:02:12.000Z",
			"profile": {
				"firstName": "Brooke",
				"lastName": "Pearson",
				"email": "brooke.pearson@example.com",
				"login": "brooke.pearson@example.com",
				"mobilePhone": "111-111-1111"
			},
			"credentials": {
				"password": {},
				"recovery_question": {
					"question": "What is your middle name?"
				},
				"provider": {
					"type": "OKTA",
					"name": "OKTA"
				}
			},
			"_links": {
				"self": {
					"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE"
				}
			}
		}
	]`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body        []byte
		wantObjects []map[string]interface{}
		wantErr     *framework.Error
	}{
		"single_page": {
			body: []byte(`[{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]`),
			wantObjects: []map[string]interface{}{
				{"id": "00ub0oNGTSWTBKOLGLNR", "status": "ACTIVE"},
				{"id": "00ub0oNGTSWTBKOCHDKE", "status": "ACTIVE"},
			},
		},
		"invalid_object_structure": {
			body: []byte(`{"result": [{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go value of type []map[string]interface {}.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			body: []byte(`[{"00ub0oNGTSWTBKOLGLNR"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: invalid character '}' after object key.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotErr := okta.ParseResponse(tt.body)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetUsersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	oktaClient := okta.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *okta.Request
		wantRes *okta.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOLGLNR",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:39:18.000Z",
						"activated":       "2013-06-24T16:39:19.000Z",
						"statusChanged":   "2013-06-24T16:39:19.000Z",
						"lastLogin":       "2013-06-24T17:39:19.000Z",
						"lastUpdated":     "2013-07-02T21:36:25.344Z",
						"passwordChanged": "2013-07-02T21:36:25.344Z",
						"profile": map[string]any{
							"firstName":   "Isaac",
							"lastName":    "Brock",
							"email":       "isaac.brock@example.com",
							"login":       "isaac.brock@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What was your school's mascot?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR",
							},
						},
					},
					{
						"id":              "00ub0oNGTSWTBKOCNDJI",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:42:20.000Z",
						"activated":       "2013-06-24T16:42:20.000Z",
						"statusChanged":   "2013-06-24T16:42:20.000Z",
						"lastLogin":       "2013-06-24T16:43:12.000Z",
						"lastUpdated":     "2013-06-24T16:42:20.000Z",
						"passwordChanged": "2013-06-24T16:42:20.000Z",
						"profile": map[string]any{
							"firstName":   "John",
							"lastName":    "Smith",
							"email":       "john.smith@example.com",
							"login":       "john.smith@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your mother's maiden name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI",
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOMSUFE",
						"status":          "ACTIVE",
						"created":         "2013-06-24T18:02:12.000Z",
						"activated":       "2013-06-24T18:02:12.000Z",
						"statusChanged":   "2013-06-24T18:02:12.000Z",
						"lastLogin":       "2013-06-24T19:14:58.000Z",
						"lastUpdated":     "2013-06-24T18:02:12.000Z",
						"passwordChanged": "2013-06-24T18:02:12.000Z",
						"profile": map[string]any{
							"firstName":   "Brooke",
							"lastName":    "Pearson",
							"email":       "brooke.pearson@example.com",
							"login":       "brooke.pearson@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your middle name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE",
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := oktaClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetGroupsPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	oktaClient := okta.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *okta.Request
		wantRes *okta.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "Group",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    "00g1emaKYZTWRYYRRTSK",
						"created":               "2015-02-06T10:11:28.000Z",
						"lastUpdated":           "2015-10-05T19:16:43.000Z",
						"lastMembershipUpdated": "2015-11-28T19:15:32.000Z",
						"objectClass": []any{
							"okta:user_group",
						},
						"type": "OKTA_GROUP",
						"profile": map[string]any{
							"name":        "West Coast Users",
							"description": "All Users West of The Rockies",
						},
					},
					{
						"id":                    "00garwpuyxHaWOkdV0g4",
						"created":               "2015-08-15T19:15:17.000Z",
						"lastUpdated":           "2015-11-18T04:02:19.000Z",
						"lastMembershipUpdated": "2015-08-15T19:15:17.000Z",
						"objectClass": []any{
							"okta:windows_security_principal",
						},
						"type": "APP_GROUP",
						"profile": map[string]any{
							"name":                       "Engineering Users",
							"description":                "corp.example.com/Engineering/Engineering Users",
							"groupType":                  "Security",
							"samAccountName":             "Engineering Users",
							"objectSid":                  "S-1-5-21-717838489-685202119-709183397-1177",
							"groupScope":                 "Global",
							"dn":                         "CN=Engineering Users,OU=Engineering,DC=corp,DC=example,DC=com",
							"windowsDomainQualifiedName": "CORP\\Engineering Users",
							"externalId":                 "OZJdWdONCU6h7WjQKp+LPA==",
						},
						"source": map[string]any{
							"id": "0oa2v0el0gP90aqjJ0g7",
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "User",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00g3zvuhepAwReSDo1d7&limit=2&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":                    "00g1emaKYZTWRCCNDHEU",
						"created":               "2015-02-06T10:11:28.000Z",
						"lastUpdated":           "2015-10-05T19:16:43.000Z",
						"lastMembershipUpdated": "2015-11-28T19:15:32.000Z",
						"objectClass": []any{
							"okta:user_group",
						},
						"type": "OKTA_GROUP",
						"profile": map[string]any{
							"name":        "East Coast Users",
							"description": "All Users East of The Rockies",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := oktaClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetGroupMembersPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	oktaClient := okta.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *okta.Request
		wantRes *okta.Response
		wantErr *framework.Error
	}{
		"group_page_1_of_2_member_page_1_of_2": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOLGLNR-00g1emaKYZTWRYYRRTSK",
						"userId":          "00ub0oNGTSWTBKOLGLNR",
						"groupId":         "00g1emaKYZTWRYYRRTSK",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:39:18.000Z",
						"activated":       "2013-06-24T16:39:19.000Z",
						"statusChanged":   "2013-06-24T16:39:19.000Z",
						"lastLogin":       "2013-06-24T17:39:19.000Z",
						"lastUpdated":     "2013-07-02T21:36:25.344Z",
						"passwordChanged": "2013-07-02T21:36:25.344Z",
						"profile": map[string]any{
							"firstName":   "Isaac",
							"lastName":    "Brock",
							"email":       "isaac.brock@example.com",
							"login":       "isaac.brock@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What was your school's mascot?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR",
							},
						},
					},
					{
						"id":              "00ub0oNGTSWTBKOCNDJI-00g1emaKYZTWRYYRRTSK",
						"userId":          "00ub0oNGTSWTBKOCNDJI",
						"groupId":         "00g1emaKYZTWRYYRRTSK",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:42:20.000Z",
						"activated":       "2013-06-24T16:42:20.000Z",
						"statusChanged":   "2013-06-24T16:42:20.000Z",
						"lastLogin":       "2013-06-24T16:43:12.000Z",
						"lastUpdated":     "2013-06-24T16:42:20.000Z",
						"passwordChanged": "2013-06-24T16:42:20.000Z",
						"profile": map[string]any{
							"firstName":   "John",
							"lastName":    "Smith",
							"email":       "john.smith@example.com",
							"login":       "john.smith@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your mother's maiden name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI",
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// We're of syncing Members for a specific Groups, so this cursor is to the next page of Members.
					Cursor:       testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2"),
					CollectionID: testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantErr: nil,
		},
		"group_page_1_of_2_member_page_2_of_2": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:           testutil.GenPtr(server.URL + "/api/v1/groups/00g1emaKYZTWRYYRRTSK/users?after=00ub0oNGTSWTBKOCNDJI&limit=2"),
					CollectionCursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
					CollectionID:     testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOMSUFE-00g1emaKYZTWRYYRRTSK",
						"userId":          "00ub0oNGTSWTBKOMSUFE",
						"groupId":         "00g1emaKYZTWRYYRRTSK",
						"status":          "ACTIVE",
						"created":         "2013-06-24T18:02:12.000Z",
						"activated":       "2013-06-24T18:02:12.000Z",
						"statusChanged":   "2013-06-24T18:02:12.000Z",
						"lastLogin":       "2013-06-24T19:14:58.000Z",
						"lastUpdated":     "2013-06-24T18:02:12.000Z",
						"passwordChanged": "2013-06-24T18:02:12.000Z",
						"profile": map[string]any{
							"firstName":   "Brooke",
							"lastName":    "Pearson",
							"email":       "brooke.pearson@example.com",
							"login":       "brooke.pearson@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your middle name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE",
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// There is no Cursor since we've finished all pages of Members for the current Group.
					CollectionID: testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
					// GroupCursor to the next page of Groups.
					CollectionCursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantErr: nil,
		},
		"group_page_2_of_2_member_page_1_of_2": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00g1emaKYZTWRYYRRTSK&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
					CollectionID:     testutil.GenPtr("00g1emaKYZTWRYYRRTSK"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOLGLNR-00garwpuyxHaWOkdV0g4",
						"userId":          "00ub0oNGTSWTBKOLGLNR",
						"groupId":         "00garwpuyxHaWOkdV0g4",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:39:18.000Z",
						"activated":       "2013-06-24T16:39:19.000Z",
						"statusChanged":   "2013-06-24T16:39:19.000Z",
						"lastLogin":       "2013-06-24T17:39:19.000Z",
						"lastUpdated":     "2013-07-02T21:36:25.344Z",
						"passwordChanged": "2013-07-02T21:36:25.344Z",
						"profile": map[string]any{
							"firstName":   "Isaac",
							"lastName":    "Brock",
							"email":       "isaac.brock@example.com",
							"login":       "isaac.brock@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What was your school's mascot?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOLGLNR",
							},
						},
					},
					{
						"id":              "00ub0oNGTSWTBKOCNDJI-00garwpuyxHaWOkdV0g4",
						"userId":          "00ub0oNGTSWTBKOCNDJI",
						"groupId":         "00garwpuyxHaWOkdV0g4",
						"status":          "ACTIVE",
						"created":         "2013-06-24T16:42:20.000Z",
						"activated":       "2013-06-24T16:42:20.000Z",
						"statusChanged":   "2013-06-24T16:42:20.000Z",
						"lastLogin":       "2013-06-24T16:43:12.000Z",
						"lastUpdated":     "2013-06-24T16:42:20.000Z",
						"passwordChanged": "2013-06-24T16:42:20.000Z",
						"profile": map[string]any{
							"firstName":   "John",
							"lastName":    "Smith",
							"email":       "john.smith@example.com",
							"login":       "john.smith@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your mother's maiden name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOCNDJI",
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[string]{
					// Cursor to the next page of Members for the current Group.
					Cursor:       testutil.GenPtr("https://test-instance.oktapreview.com/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOCNDJI&limit=2"),
					CollectionID: testutil.GenPtr("00garwpuyxHaWOkdV0g4"),
					// There is no CollectionCursor since we're currently processing the last Group.
				},
			},
			wantErr: nil,
		},
		"group_page_2_of_2_member_page_2_of_2": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					Cursor:       testutil.GenPtr(server.URL + "/api/v1/groups/00garwpuyxHaWOkdV0g4/users?after=00ub0oNGTSWTBKOCNDJI&limit=2"),
					CollectionID: testutil.GenPtr("00garwpuyxHaWOkdV0g4"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":              "00ub0oNGTSWTBKOMSUFE-00garwpuyxHaWOkdV0g4",
						"userId":          "00ub0oNGTSWTBKOMSUFE",
						"groupId":         "00garwpuyxHaWOkdV0g4",
						"status":          "ACTIVE",
						"created":         "2013-06-24T18:02:12.000Z",
						"activated":       "2013-06-24T18:02:12.000Z",
						"statusChanged":   "2013-06-24T18:02:12.000Z",
						"lastLogin":       "2013-06-24T19:14:58.000Z",
						"lastUpdated":     "2013-06-24T18:02:12.000Z",
						"passwordChanged": "2013-06-24T18:02:12.000Z",
						"profile": map[string]any{
							"firstName":   "Brooke",
							"lastName":    "Pearson",
							"email":       "brooke.pearson@example.com",
							"login":       "brooke.pearson@example.com",
							"mobilePhone": "111-111-1111",
						},
						"credentials": map[string]any{
							"password": map[string]any{},
							"recovery_question": map[string]any{
								"question": "What is your middle name?",
							},
							"provider": map[string]any{
								"type": "OKTA",
								"name": "OKTA",
							},
						},
						"_links": map[string]any{
							"self": map[string]any{
								"href": "https://test-instance.oktapreview.com/api/v1/users/00ub0oNGTSWTBKOMSUFE",
							},
						},
					},
				},
				// Cursor and CollectionCursor are both nil, so no cursor is set as this is the last page for the sync.
			},
			wantErr: nil,
		},
		"group_members_too_many_groups_returned": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00garwpuyxHaWOkdV0g4&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantErr: &framework.Error{
				Message: "Too many collection objects returned in response; expected 1, got 2.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		// If no groups are present in the current page and there is no next group link, exit early.
		"group_members_no_groups_returned": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/api/v1/groups?after=00garwpuyxHaWOkcndD8&limit=1&filter=type+eq+%22OKTA_GROUP%22+or+type+eq+%22APP_GROUP%22"),
				},
			},
			wantRes: &okta.Response{
				StatusCode: http.StatusOK,
			},
		},
		"group_member_invalid_group_cursor": {
			context: context.Background(),
			request: &okta.Request{
				Token:                 "SSWS testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "GroupMember",
				PageSize:              2,
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[string]{
					CollectionCursor: testutil.GenPtr(server.URL + "/invalid"),
				},
			},
			wantErr: &framework.Error{
				Message: "Datasource rejected request, returned status code: 404.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := oktaClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
