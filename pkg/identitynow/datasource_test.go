// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package identitynow_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/identitynow"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock IdentityNow server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer token" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{ "error": "JWT is required" }`))
	}

	switch r.URL.RequestURI() {
	// Accounts Page 1
	case "/v3/accounts?limit=2&offset=0":
		w.Write([]byte(`[
			{
				"authoritative": false,
				"systemAccount": false,
				"uncorrelated": true,
				"features": "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
				"uuid": null,
				"nativeIdentity": "0e826bf03710200044e0bfc8bcbe5d85",
				"description": null,
				"disabled": false,
				"locked": false,
				"manuallyCorrelated": false,
				"hasEntitlements": false,
				"sourceId": "1fb19cd2dcd440b09711aca31dabf616",
				"sourceName": "ServiceNow test-instance",
				"identityId": "0017c13884f348c2bca7433480b7a68a",
				"identity": {
					"type": "IDENTITY",
					"id": "0017c13884f348c2bca7433480b7a68a",
					"name": "victor.mcintosh"
				},
				"attributes": {
					"calendar_integration": "Outlook",
					"gender": "Male",
					"user_name": "victor.mcintosh",
					"sys_updated_on": "2022-09-10 12:46:42",
					"sys_class_name": "User",
					"notification": "Enable",
					"sys_id": "0e826bf03710200044e0bfc8bcbe5d85",
					"sys_updated_by": "system",
					"sys_created_on": "2012-02-17 19:04:50",
					"sys_domain": "global",
					"company": "ACME North America",
					"department": "Development",
					"vip": "false",
					"first_name": "Victor",
					"email": "victor.mcintosh@example.com",
					"idNowDescription": "a300f8ca5d6cdf45c4b75622602c886ca2f6d8fca4bf9dfd9284856fc7a32e6a",
					"sys_created_by": "admin",
					"locked_out": "false",
					"sys_mod_count": "5",
					"active": "true",
					"last_name": "Mcintosh",
					"cost_center": "Engineering",
					"name": "Victor Mcintosh",
					"location": "{{OMITTED}}",
					"password_needs_reset": "false"
				},
				"id": "1a1bb825eb7e4f76b72fbecb27699b31",
				"name": "victor.mcintosh",
				"created": "2023-09-22T16:46:54.250Z",
				"modified": "2023-09-22T16:46:54.325Z"
			},
			{
				"authoritative": false,
				"systemAccount": false,
				"uncorrelated": true,
				"features": "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
				"uuid": null,
				"nativeIdentity": "e46393fb1bd0b5509a55631e6e4bcbf7",
				"description": null,
				"disabled": false,
				"locked": false,
				"manuallyCorrelated": false,
				"hasEntitlements": true,
				"sourceId": "1fb19cd2dcd440b09711aca31dabf616",
				"sourceName": "ServiceNow test-instance",
				"identityId": "0017db102a20473ab350a537a4037009",
				"identity": {
					"type": "IDENTITY",
					"id": "0017db102a20473ab350a537a4037009",
					"name": "Cleo.Yoder"
				},
				"attributes": {
					"calendar_integration": "Outlook",
					"user_name": "Cleo.Yoder",
					"roles": [
						"e098ecf6c0a80165002aaec84d906014"
					],
					"sys_updated_on": "2023-08-03 18:33:00",
					"title": "ED Nurse",
					"sys_class_name": "User",
					"notification": "Enable",
					"sys_id": "e46393fb1bd0b5509a55631e6e4bcbf7",
					"sys_updated_by": "admin",
					"sys_created_on": "2023-08-03 18:33:00",
					"sys_domain": "global",
					"vip": "false",
					"first_name": "Cleo",
					"email": "Cleo.Yoder@260.sailpointtechnologies.com",
					"idNowDescription": "77dbd87c31d4c34978af9f2bb4ba0af6165dd9b419353c067af057d4766b3126",
					"sys_created_by": "admin",
					"locked_out": "false",
					"sys_mod_count": "0",
					"active": "true",
					"groups": [
						"5b3c2f56db45e01061a5a5bb1396197f"
					],
					"last_name": "Yoder",
					"phone": "+11111111111",
					"name": "Cleo Yoder",
					"password_needs_reset": "false"
				},
				"id": "ba699287e60b4014bcc4319f30e9b59e",
				"name": "Cleo.Yoder",
				"created": "2023-09-22T16:47:35.702Z",
				"modified": "2023-09-22T16:47:35.797Z"
			}
		]`))

	// Accounts Page 2 (last page)
	case "/v3/accounts?limit=2&offset=2":
		w.Write([]byte(`[
			{
				"authoritative": false,
				"systemAccount": false,
				"uncorrelated": false,
				"features": "GROUP_PROVISIONING, ENABLE, SEARCH, PASSWORD, NO_PERMISSIONS_PROVISIONING, NO_GROUP_PERMISSIONS_PROVISIONING, SYNC_PROVISIONING, PROVISIONING, CURRENT_PASSWORD, AUTHENTICATE, UNSTRUCTURED_TARGETS, GROUPS_HAVE_MEMBERS, UNLOCK, MANAGER_LOOKUP, PREFER_UUID",
				"uuid": "{8efc8cee-2e3d-46c6-a488-f7f000388e82}",
				"nativeIdentity": "CN=Cynthia Edwards,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
				"description": null,
				"disabled": false,
				"locked": false,
				"manuallyCorrelated": false,
				"hasEntitlements": true,
				"sourceId": "602dbeacc6eb429c9038a4bb2d776e28",
				"sourceName": "Active Directory",
				"identityId": "00263fd218d2487eac5ab27fd8b47f47",
				"identity": {
					"type": "IDENTITY",
					"id": "00263fd218d2487eac5ab27fd8b47f47",
					"name": "Cynthia.Edwards"
				},
				"attributes": {
					"mail": "Cynthia.Edwards@sailpointdemo.com",
					"displayName": "Cynthia Edwards",
					"distinguishedName": "CN=Cynthia,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
					"objectType": "user",
					"objectguid": "{8efc8cee-2e3d-46c6-a488-f7f000388e82}",
					"memberOf": [
						"CN=Development,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"CN=ENG_Internal,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"CN=Employees,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"CN=All_Users,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"CN=ENG_WestCoast,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com"
					],
					"sn": "Edwards",
					"department": "Engineering",
					"idNowDescription": "107d4be3f29b60ff19ced8ebd85924056d182170490cc584481ff68c8624c634",
					"userPrincipalName": "Cynthia.Edwards@sailpointdemo.com",
					"passwordLastSet": 1614901945829,
					"manager": "CN=Rahim Riddle,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
					"sAMAccountName": "Cynthia.Edwards",
					"msNPAllowDialin": "Not Set",
					"givenName": "Cynthia",
					"objectClass": [
						"top",
						"person",
						"organizationalPerson",
						"user"
					],
					"cn": "Cynthia Edwards",
					"accountFlags": [
						"Normal User Account",
						"Password Cannot Expire"
					],
					"NetBIOSName": null,
					"domain": "sailpointdemo.com",
					"primaryGroupID": "513",
					"objectSid": "S-1-5-21-2981491572-779881612-3979282638-32437",
					"msDS-PrincipalName": "SERI\\Cynthia.Edwards",
					"pwdLastSet": "132593755458298019"
				},
				"id": "28b1e9bf40ab458981067f4e4dc330b3",
				"name": "Cynthia.Edwards",
				"created": "2023-09-22T16:46:45.885Z",
				"modified": "2023-09-22T16:46:46.230Z"
			}
		]`))

	// Accounts with Filter (single page)
	case "/v3/accounts?limit=2&offset=0&filters=id%20eq%20%221a1bb825eb7e4f76b72fbecb27699b31%22":
		w.Write([]byte(`[
			{
				"authoritative": false,
				"systemAccount": false,
				"uncorrelated": true,
				"features": "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
				"uuid": null,
				"nativeIdentity": "0e826bf03710200044e0bfc8bcbe5d85",
				"description": null,
				"disabled": false,
				"locked": false,
				"manuallyCorrelated": false,
				"hasEntitlements": false,
				"sourceId": "1fb19cd2dcd440b09711aca31dabf616",
				"sourceName": "ServiceNow test-instance",
				"identityId": "0017c13884f348c2bca7433480b7a68a",
				"identity": {
					"type": "IDENTITY",
					"id": "0017c13884f348c2bca7433480b7a68a",
					"name": "victor.mcintosh"
				},
				"attributes": {
					"calendar_integration": "Outlook",
					"gender": "Male",
					"user_name": "victor.mcintosh",
					"sys_updated_on": "2022-09-10 12:46:42",
					"sys_class_name": "User",
					"notification": "Enable",
					"sys_id": "0e826bf03710200044e0bfc8bcbe5d85",
					"sys_updated_by": "system",
					"sys_created_on": "2012-02-17 19:04:50",
					"sys_domain": "global",
					"company": "ACME North America",
					"department": "Development",
					"vip": "false",
					"first_name": "Victor",
					"email": "victor.mcintosh@example.com",
					"idNowDescription": "a300f8ca5d6cdf45c4b75622602c886ca2f6d8fca4bf9dfd9284856fc7a32e6a",
					"sys_created_by": "admin",
					"locked_out": "false",
					"sys_mod_count": "5",
					"active": "true",
					"last_name": "Mcintosh",
					"cost_center": "Engineering",
					"name": "Victor Mcintosh",
					"location": "{{OMITTED}}",
					"password_needs_reset": "false"
				},
				"id": "1a1bb825eb7e4f76b72fbecb27699b31",
				"name": "victor.mcintosh",
				"created": "2023-09-22T16:46:54.250Z",
				"modified": "2023-09-22T16:46:54.325Z"
			}
		]`))

	// Accounts Page 97 (to test concatenating attributes)
	case "/v3/accounts?limit=2&offset=97":
		w.Write([]byte(`[
		{
			"authoritative": false,
			"systemAccount": false,
			"uncorrelated": true,
			"features": "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
			"uuid": null,
			"nativeIdentity": "0e826bf03710200044e0bfc8bcbe5d85",
			"description": null,
			"disabled": false,
			"locked": false,
			"manuallyCorrelated": false,
			"hasEntitlements": false,
			"sourceId": "1fb19cd2dcd440b09711aca31dabf616",
			"sourceName": "ServiceNow test-instance",
			"identityId": "0017c13884f348c2bca7433480b7a68a",
			"identity": {
				"type": "IDENTITY",
				"id": "0017c13884f348c2bca7433480b7a68a",
				"name": "victor.mcintosh"
			},
			"attributes": {
				"groups": ["GROUP1", "GROUP2", "GROUP3"],
				"memberOf": ["GROUP1", "GROUP2", "GROUP3"],
				"Groups": ["GROUP1", "GROUP2", "GROUP3"],
				"calendar_integration": "Outlook",
				"gender": "Male",
				"user_name": "victor.mcintosh",
				"sys_updated_on": "2022-09-10 12:46:42",
				"sys_class_name": "User",
				"notification": "Enable",
				"sys_id": "0e826bf03710200044e0bfc8bcbe5d85",
				"sys_updated_by": "system",
				"sys_created_on": "2012-02-17 19:04:50",
				"sys_domain": "global",
				"company": "ACME North America",
				"department": "Development",
				"vip": "false",
				"first_name": "Victor",
				"email": "victor.mcintosh@example.com",
				"idNowDescription": "a300f8ca5d6cdf45c4b75622602c886ca2f6d8fca4bf9dfd9284856fc7a32e6a",
				"sys_created_by": "admin",
				"locked_out": "false",
				"sys_mod_count": "5",
				"active": "true",
				"last_name": "Mcintosh",
				"cost_center": "Engineering",
				"name": "Victor Mcintosh",
				"location": "{{OMITTED}}",
				"password_needs_reset": "false"
			},
			"id": "1a1bb825eb7e4f76b72fbecb27699b31",
			"name": "victor.mcintosh",
			"created": "2023-09-22T16:46:54.250Z",
			"modified": "2023-09-22T16:46:54.325Z"
		}
	]`))

	// Accounts Page 98 (to test concatenating attributes with invalid attribute type, e.g. int)
	case "/v3/accounts?limit=2&offset=98":
		w.Write([]byte(`[
		{
			"authoritative": false,
			"systemAccount": false,
			"uncorrelated": true,
			"features": "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
			"uuid": null,
			"nativeIdentity": "0e826bf03710200044e0bfc8bcbe5d85",
			"description": null,
			"disabled": false,
			"locked": false,
			"manuallyCorrelated": false,
			"hasEntitlements": false,
			"sourceId": "1fb19cd2dcd440b09711aca31dabf616",
			"sourceName": "ServiceNow test-instance",
			"identityId": "0017c13884f348c2bca7433480b7a68a",
			"identity": {
				"type": "IDENTITY",
				"id": "0017c13884f348c2bca7433480b7a68a",
				"name": "victor.mcintosh"
			},
			"attributes": {
				"groups": ["GROUP1", 5, 10],
				"memberOf": "NOT_ARRAY",
				"Groups": ["GROUP1", "GROUP2", "GROUP3"],
				"calendar_integration": "Outlook",
				"gender": "Male",
				"user_name": "victor.mcintosh",
				"sys_updated_on": "2022-09-10 12:46:42",
				"sys_class_name": "User",
				"notification": "Enable",
				"sys_id": "0e826bf03710200044e0bfc8bcbe5d85",
				"sys_updated_by": "system",
				"sys_created_on": "2012-02-17 19:04:50",
				"sys_domain": "global",
				"company": "ACME North America",
				"department": "Development",
				"vip": "false",
				"first_name": "Victor",
				"email": "victor.mcintosh@example.com",
				"idNowDescription": "a300f8ca5d6cdf45c4b75622602c886ca2f6d8fca4bf9dfd9284856fc7a32e6a",
				"sys_created_by": "admin",
				"locked_out": "false",
				"sys_mod_count": "5",
				"active": "true",
				"last_name": "Mcintosh",
				"cost_center": "Engineering",
				"name": "Victor Mcintosh",
				"location": "{{OMITTED}}",
				"password_needs_reset": "false"
			},
			"id": "1a1bb825eb7e4f76b72fbecb27699b31",
			"name": "victor.mcintosh",
			"created": "2023-09-22T16:46:54.250Z",
			"modified": "2023-09-22T16:46:54.325Z"
		}
	]`))

	// Accounts Page 99 (contrived test case to return a response that is not http.StatusOK)
	case "/v3/accounts?limit=2&offset=99":
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))

	// Entitlements Page 1
	case "/beta/entitlements?limit=2&offset=0":
		w.Write([]byte(`[
		    {
				"sourceSchemaObjectType": "applicationRole",
				"attribute": "appRoleAssignments",
				"attributes": {
					"displayName": "Basic purchaser [on] Windows Store for Business"
				},
				"value": "AZURE_APP_ID_123:AZURE_RESOURCE_ID_456",
				"description": null,
				"privileged": true,
				"cloudGoverned": false,
				"requestable": true,
				"id": "ENTITLEMENT_ID_456",
				"created": "2023-09-22T16:50:12.053Z",
				"modified": "2023-09-22T16:56:12.896Z"
			},
			{
				"sourceSchemaObjectType": "applicationRole",
				"attribute": "appRoleAssignments",
				"attributes": {
					"displayName": "default access [on] TrustedPublishersProxyService"
				},
				"value": "efd1eb6f-44f3-4ffe-b4e4-eb68162ea4ae:00000000-0000-0000-0000-000000000000",
				"description": null,
				"privileged": true,
				"cloudGoverned": false,
				"requestable": true,
				"id": "00218206fe614e7da637f528accdf15e",
				"created": "2023-09-22T16:50:54.856Z",
				"modified": "2023-09-22T16:56:13.788Z"
			}
		]`))

	// New tests for faster fetching of account entitlements
	case "/beta/accounts?limit=5&offset=0&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId1",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId2",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId3",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId4",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId5",
				"hasEntitlements": false
			}
		]`))
	case "/beta/accounts?limit=5&offset=4&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId5",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId6",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId7",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId8",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId9",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=7&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId7",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId8",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId9",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId10",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId11",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=20&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId12",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId13",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId14",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId15",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId16",
				"hasEntitlements": false
			}
		]`))
	case "/beta/accounts?limit=5&offset=101&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId101",
				"hasEntitlements": "incorrect value"
			},
			{
				"id": "testaccountId102",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId103",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId104",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId105",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=110&sorters=id":
		w.Write([]byte(`[
			{
				"hasEntitlements": true
			},
			{
				"id": "testaccountId111",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId112",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId113",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId114",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=120&sorters=id":
		w.Write([]byte(`[
			{
				"id": 120,
				"hasEntitlements": true
			},
			{
				"id": "testaccountId111",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId112",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId113",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId114",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=250&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId250",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId251",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId252",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId253",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId254",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=300&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId300",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId301",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=302&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId302",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId303",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId304",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=80&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId80",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId81",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId82",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId83",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId84",
				"hasEntitlements": true
			}
		]`))
	case "/beta/accounts?limit=5&offset=400&sorters=id":
		w.Write([]byte(`[]`))
	case "/beta/accounts?limit=5&offset=450&sorters=id":
		w.Write([]byte(`[
			{
				"id": "testaccountId450",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId451",
				"hasEntitlements": true
			},
			{
				"id": "testaccountId452",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId453",
				"hasEntitlements": false
			},
			{
				"id": "testaccountId454",
				"hasEntitlements": false
			}
		]`))
	// AccountEntitlements for accountId: testaccountId1
	case "/beta/accounts/testaccountId1/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId1"
			},
			{
				"id": "entitlementId2"
			}
		]`))

	case "/beta/accounts/testaccountId2/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId3"
			},
			{
				"id": "entitlementId4"
			},
			{
				"id": "entitlementId5"
			}
		]`))

	case "/beta/accounts/testaccountId3/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId6"
			},
			{
				"id": "entitlementId7"
			},
			{
				"id": "entitlementId8"
			}
		]`))
	case "/beta/accounts/testaccountId4/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId9"
			},
			{
				"id": "entitlementId10"
			}
		]`))
	case "/beta/accounts/testaccountId6/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId11"
			},
			{
				"id": "entitlementId12"
			},
			{
				"id": "entitlementId13"
			},
			{
				"id": "entitlementId14"
			},
			{
				"id": "entitlementId15"
			}
		]`))
	case "/beta/accounts/testaccountId7/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId16"
			},
			{
				"id": "entitlementId17"
			},
			{
				"id": "entitlementId18"
			},
			{
				"id": "entitlementId19"
			},
			{
				"id": "entitlementId20"
			},
			{
				"id": "entitlementId21"
			}
		]`))
	case "/beta/accounts/testaccountId7/entitlements?limit=10&offset=5":
		w.Write([]byte(`[
			{
				"id": "entitlementId21"
			}
		]`))
	case "/beta/accounts/testaccountId8/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId22"
			},
			{
				"id": "entitlementId23"
			}
		]`))
	case "/beta/accounts/testaccountId9/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId24"
			},
			{
				"id": "entitlementId25"
			}
		]`))
	case "/beta/accounts/testaccountId10/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId26"
			}
		]`))
	case "/beta/accounts/testaccountId11/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId27"
			}
		]`))
	case "/beta/accounts/testaccountId101/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId27"
			}
		]`))
	case "/beta/accounts/testaccountId250/entitlements?limit=10&offset=99":
		w.Write([]byte(`[
			{
				"id": 250
			}
		]`))
	case "/beta/accounts/testaccountId300/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_One"
			},
			{
				"id": "entitlementId_Two"
			}
		]`))
	case "/beta/accounts/testaccountId301/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_Three"
			},
			{
				"id": "entitlementId_Four"
			}
		]`))
	case "/beta/accounts/testaccountId80/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId801"
			},
			{
				"id": "entitlementId802"
			},
			{
				"id": "entitlementId803"
			},
			{
				"id": "entitlementId804"
			},
			{
				"id": "entitlementId805"
			},
			{
				"id": "entitlementId806"
			},
			{
				"id": "entitlementId807"
			},
			{
				"id": "entitlementId808"
			},
			{
				"id": "entitlementId809"
			},
			{
				"id": "entitlementId810"
			},
			{
				"id": "entitlementId811"
			},
			{
				"id": "entitlementId812"
			}
		]`))
	case "/beta/accounts/testaccountId302/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_a"
			},
			{
				"id": "entitlementId_b"
			},
			{
				"id": "entitlementId_c"
			}
		]`))
	case "/beta/accounts/testaccountId303/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_d"
			},
			{
				"id": "entitlementId_e"
			},
			{
				"id": "entitlementId_f"
			},
			{
				"id": "entitlementId_g"
			}
		]`))
	case "/beta/accounts/testaccountId304/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_h"
			},
			{
				"id": "entitlementId_i"
			},
			{
				"id": "entitlementId_j"
			},
			{
				"id": "entitlementId_k"
			},
			{
				"id": "entitlementId_l"
			}
		]`))
	case "/beta/accounts/testaccountId450/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_a"
			},
			{
				"id": "entitlementId_b"
			}
		]`))
	case "/beta/accounts/testaccountId451/entitlements?limit=10&offset=0":
		w.Write([]byte(`[
			{
				"id": "entitlementId_c"
			},
			{
				"id": "entitlementId_d"
			}
		]`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestGetAccountsPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	identitynowClient := identitynow.NewClient(client, 0)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context      context.Context
		request      *identitynow.Request
		wantRes      *identitynow.Response
		wantErr      *framework.Error
		expectedLogs []map[string]any
	}{
		"first_page": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accounts",
				PageSize:              2,
				APIVersion:            "v3",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"authoritative":      false,
						"systemAccount":      false,
						"uncorrelated":       true,
						"features":           "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
						"uuid":               nil,
						"nativeIdentity":     "0e826bf03710200044e0bfc8bcbe5d85",
						"description":        nil,
						"disabled":           false,
						"locked":             false,
						"manuallyCorrelated": false,
						"hasEntitlements":    false,
						"sourceId":           "1fb19cd2dcd440b09711aca31dabf616",
						"sourceName":         "ServiceNow test-instance",
						"identityId":         "0017c13884f348c2bca7433480b7a68a",
						"identity": map[string]any{
							"type": "IDENTITY",
							"id":   "0017c13884f348c2bca7433480b7a68a",
							"name": "victor.mcintosh",
						},
						"attributes": map[string]any{
							"calendar_integration": "Outlook",
							"gender":               "Male",
							"user_name":            "victor.mcintosh",
							"sys_updated_on":       "2022-09-10 12:46:42",
							"sys_class_name":       "User",
							"notification":         "Enable",
							"sys_id":               "0e826bf03710200044e0bfc8bcbe5d85",
							"sys_updated_by":       "system",
							"sys_created_on":       "2012-02-17 19:04:50",
							"sys_domain":           "global",
							"company":              "ACME North America",
							"department":           "Development",
							"vip":                  "false",
							"first_name":           "Victor",
							"email":                "victor.mcintosh@example.com",
							"idNowDescription":     "a300f8ca5d6cdf45c4b75622602c886ca2f6d8fca4bf9dfd9284856fc7a32e6a",
							"sys_created_by":       "admin",
							"locked_out":           "false",
							"sys_mod_count":        "5",
							"active":               "true",
							"last_name":            "Mcintosh",
							"cost_center":          "Engineering",
							"name":                 "Victor Mcintosh",
							"location":             "{{OMITTED}}",
							"password_needs_reset": "false",
						},
						"id":       "1a1bb825eb7e4f76b72fbecb27699b31",
						"name":     "victor.mcintosh",
						"created":  "2023-09-22T16:46:54.250Z",
						"modified": "2023-09-22T16:46:54.325Z",
					},
					{
						"authoritative":      false,
						"systemAccount":      false,
						"uncorrelated":       true,
						"features":           "SEARCH, UNLOCK, SYNC_PROVISIONING, PASSWORD, GROUP_PROVISIONING, ENABLE, PROVISIONING",
						"uuid":               nil,
						"nativeIdentity":     "e46393fb1bd0b5509a55631e6e4bcbf7",
						"description":        nil,
						"disabled":           false,
						"locked":             false,
						"manuallyCorrelated": false,
						"hasEntitlements":    true,
						"sourceId":           "1fb19cd2dcd440b09711aca31dabf616",
						"sourceName":         "ServiceNow test-instance",
						"identityId":         "0017db102a20473ab350a537a4037009",
						"identity": map[string]any{
							"type": "IDENTITY",
							"id":   "0017db102a20473ab350a537a4037009",
							"name": "Cleo.Yoder",
						},
						"attributes": map[string]any{
							"calendar_integration": "Outlook",
							"user_name":            "Cleo.Yoder",
							"roles": []any{
								"e098ecf6c0a80165002aaec84d906014",
							},
							"sys_updated_on":   "2023-08-03 18:33:00",
							"title":            "ED Nurse",
							"sys_class_name":   "User",
							"notification":     "Enable",
							"sys_id":           "e46393fb1bd0b5509a55631e6e4bcbf7",
							"sys_updated_by":   "admin",
							"sys_created_on":   "2023-08-03 18:33:00",
							"sys_domain":       "global",
							"vip":              "false",
							"first_name":       "Cleo",
							"email":            "Cleo.Yoder@260.sailpointtechnologies.com",
							"idNowDescription": "77dbd87c31d4c34978af9f2bb4ba0af6165dd9b419353c067af057d4766b3126",
							"sys_created_by":   "admin",
							"locked_out":       "false",
							"sys_mod_count":    "0",
							"active":           "true",
							"groups": []any{
								"5b3c2f56db45e01061a5a5bb1396197f",
							},
							"last_name":            "Yoder",
							"phone":                "+11111111111",
							"name":                 "Cleo Yoder",
							"password_needs_reset": "false",
						},
						"groupsMembership": "5b3c2f56db45e01061a5a5bb1396197f",
						"id":               "ba699287e60b4014bcc4319f30e9b59e",
						"name":             "Cleo.Yoder",
						"created":          "2023-09-22T16:47:35.702Z",
						"modified":         "2023-09-22T16:47:35.797Z",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
			},
			wantErr: nil,
			expectedLogs: []map[string]any{
				{
					"level":                             "info",
					"msg":                               "Starting datasource request",
					fields.FieldRequestEntityExternalID: "accounts",
					fields.FieldRequestPageSize:         int64(2),
				},
				{
					"level":                             "info",
					"msg":                               "Sending HTTP request to datasource",
					fields.FieldRequestEntityExternalID: "accounts",
					fields.FieldRequestPageSize:         int64(2),
					fields.FieldRequestURL:              server.URL + "/v3/accounts?limit=2&offset=0",
				},
				{
					"level":                             "info",
					"msg":                               "Datasource request completed successfully",
					fields.FieldRequestEntityExternalID: "accounts",
					fields.FieldRequestPageSize:         int64(2),
					fields.FieldResponseStatusCode:      int64(200),
					fields.FieldResponseObjectCount:     int64(2),
					fields.FieldResponseNextCursor:      nil,
				},
			},
		},
		"last_page": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accounts",
				PageSize:              2,
				APIVersion:            "v3",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"authoritative":      false,
						"systemAccount":      false,
						"uncorrelated":       false,
						"features":           "GROUP_PROVISIONING, ENABLE, SEARCH, PASSWORD, NO_PERMISSIONS_PROVISIONING, NO_GROUP_PERMISSIONS_PROVISIONING, SYNC_PROVISIONING, PROVISIONING, CURRENT_PASSWORD, AUTHENTICATE, UNSTRUCTURED_TARGETS, GROUPS_HAVE_MEMBERS, UNLOCK, MANAGER_LOOKUP, PREFER_UUID",
						"uuid":               "{8efc8cee-2e3d-46c6-a488-f7f000388e82}",
						"nativeIdentity":     "CN=Cynthia Edwards,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"description":        nil,
						"disabled":           false,
						"locked":             false,
						"manuallyCorrelated": false,
						"hasEntitlements":    true,
						"sourceId":           "602dbeacc6eb429c9038a4bb2d776e28",
						"sourceName":         "Active Directory",
						"identityId":         "00263fd218d2487eac5ab27fd8b47f47",
						"identity": map[string]any{
							"type": "IDENTITY",
							"id":   "00263fd218d2487eac5ab27fd8b47f47",
							"name": "Cynthia.Edwards",
						},
						"attributes": map[string]any{
							"mail":              "Cynthia.Edwards@sailpointdemo.com",
							"displayName":       "Cynthia Edwards",
							"distinguishedName": "CN=Cynthia,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
							"objectType":        "user",
							"objectguid":        "{8efc8cee-2e3d-46c6-a488-f7f000388e82}",
							"memberOf": []any{
								"CN=Development,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
								"CN=ENG_Internal,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
								"CN=Employees,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
								"CN=All_Users,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
								"CN=ENG_WestCoast,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
							},
							"sn":                "Edwards",
							"department":        "Engineering",
							"idNowDescription":  "107d4be3f29b60ff19ced8ebd85924056d182170490cc584481ff68c8624c634",
							"userPrincipalName": "Cynthia.Edwards@sailpointdemo.com",
							"passwordLastSet":   float64(1614901945829),
							"manager":           "CN=Rahim Riddle,OU=Singapore,OU=Asia-Pacific,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
							"sAMAccountName":    "Cynthia.Edwards",
							"msNPAllowDialin":   "Not Set",
							"givenName":         "Cynthia",
							"objectClass": []any{
								"top",
								"person",
								"organizationalPerson",
								"user",
							},
							"cn": "Cynthia Edwards",
							"accountFlags": []any{
								"Normal User Account",
								"Password Cannot Expire",
							},
							"NetBIOSName":        nil,
							"domain":             "sailpointdemo.com",
							"primaryGroupID":     "513",
							"objectSid":          "S-1-5-21-2981491572-779881612-3979282638-32437",
							"msDS-PrincipalName": "SERI\\Cynthia.Edwards",
							"pwdLastSet":         "132593755458298019",
						},
						"memberOfMembership": "CN=Development,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com" +
							" | CN=ENG_Internal,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com" +
							" | CN=Employees,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com" +
							" | CN=All_Users,OU=BirthRight,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com" +
							" | CN=ENG_WestCoast,OU=Groups,OU=Demo,DC=seri,DC=sailpointdemo,DC=com",
						"id":       "28b1e9bf40ab458981067f4e4dc330b3",
						"name":     "Cynthia.Edwards",
						"created":  "2023-09-22T16:46:45.885Z",
						"modified": "2023-09-22T16:46:46.230Z",
					},
				},
			},
			wantErr: nil,
		},
		"not_found_response": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accounts",
				PageSize:              2,
				APIVersion:            "v3",
				RequestTimeoutSeconds: 5,
				// The test handler has been hardcoded to serve a 404 for a cursor value of 99.
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](99),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusNotFound,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctxWithLogger, observedLogs := testutil.NewContextWithObservableLogger(tt.context)

			gotRes, gotErr := identitynowClient.GetPage(ctxWithLogger, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			testutil.ValidateLogOutput(t, observedLogs, tt.expectedLogs)
		})
	}
}

func TestGetAccountEntitlements(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	accountCollectionPageSize := 5
	identitynowClient := identitynow.NewClient(client, accountCollectionPageSize)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *identitynow.Request
		wantRes *identitynow.Response
		wantErr *framework.Error
	}{
		"fetch_account_entitlements_for_first_5_accounts_by_4th_account_page_is_full": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId1-entitlementId1",
						"accountId":     "testaccountId1",
						"entitlementId": "entitlementId1",
					},
					{
						"id":            "testaccountId1-entitlementId2",
						"accountId":     "testaccountId1",
						"entitlementId": "entitlementId2",
					},
					{
						"id":            "testaccountId2-entitlementId3",
						"accountId":     "testaccountId2",
						"entitlementId": "entitlementId3",
					},
					{
						"id":            "testaccountId2-entitlementId4",
						"accountId":     "testaccountId2",
						"entitlementId": "entitlementId4",
					},
					{
						"id":            "testaccountId2-entitlementId5",
						"accountId":     "testaccountId2",
						"entitlementId": "entitlementId5",
					},
					{
						"id":            "testaccountId3-entitlementId6",
						"accountId":     "testaccountId3",
						"entitlementId": "entitlementId6",
					},
					{
						"id":            "testaccountId3-entitlementId7",
						"accountId":     "testaccountId3",
						"entitlementId": "entitlementId7",
					},
					{
						"id":            "testaccountId3-entitlementId8",
						"accountId":     "testaccountId3",
						"entitlementId": "entitlementId8",
					},
					{
						"id":            "testaccountId4-entitlementId9",
						"accountId":     "testaccountId4",
						"entitlementId": "entitlementId9",
					},
					{
						"id":            "testaccountId4-entitlementId10",
						"accountId":     "testaccountId4",
						"entitlementId": "entitlementId10",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr[string]("testaccountId4"),
					CollectionCursor: testutil.GenPtr[int64](4),
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_5_accounts_by_2nd_account_page_is_full_and_excess_entitlements_for_account": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           nil,
					CollectionID:     testutil.GenPtr[string]("testaccountId4"),
					CollectionCursor: testutil.GenPtr[int64](4),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId6-entitlementId11",
						"accountId":     "testaccountId6",
						"entitlementId": "entitlementId11",
					},
					{
						"id":            "testaccountId6-entitlementId12",
						"accountId":     "testaccountId6",
						"entitlementId": "entitlementId12",
					},
					{
						"id":            "testaccountId6-entitlementId13",
						"accountId":     "testaccountId6",
						"entitlementId": "entitlementId13",
					},
					{
						"id":            "testaccountId6-entitlementId14",
						"accountId":     "testaccountId6",
						"entitlementId": "entitlementId14",
					},
					{
						"id":            "testaccountId6-entitlementId15",
						"accountId":     "testaccountId6",
						"entitlementId": "entitlementId15",
					},
					{
						"id":            "testaccountId7-entitlementId16",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId16",
					},
					{
						"id":            "testaccountId7-entitlementId17",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId17",
					},
					{
						"id":            "testaccountId7-entitlementId18",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId18",
					},
					{
						"id":            "testaccountId7-entitlementId19",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId19",
					},
					{
						"id":            "testaccountId7-entitlementId20",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId20",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](5),
					CollectionID:     testutil.GenPtr[string]("testaccountId7"),
					CollectionCursor: testutil.GenPtr[int64](5),
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_5_accounts_grab_excess_entitlement_and_page_not_full": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](5),
					CollectionID:     testutil.GenPtr[string]("testaccountId7"),
					CollectionCursor: testutil.GenPtr[int64](7),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId7-entitlementId21",
						"accountId":     "testaccountId7",
						"entitlementId": "entitlementId21",
					},
					{
						"id":            "testaccountId8-entitlementId22",
						"accountId":     "testaccountId8",
						"entitlementId": "entitlementId22",
					},
					{
						"id":            "testaccountId8-entitlementId23",
						"accountId":     "testaccountId8",
						"entitlementId": "entitlementId23",
					},
					{
						"id":            "testaccountId9-entitlementId24",
						"accountId":     "testaccountId9",
						"entitlementId": "entitlementId24",
					},
					{
						"id":            "testaccountId9-entitlementId25",
						"accountId":     "testaccountId9",
						"entitlementId": "entitlementId25",
					},
					{
						"id":            "testaccountId10-entitlementId26",
						"accountId":     "testaccountId10",
						"entitlementId": "entitlementId26",
					},
					{
						"id":            "testaccountId11-entitlementId27",
						"accountId":     "testaccountId11",
						"entitlementId": "entitlementId27",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr[string]("testaccountId11"),
					CollectionCursor: testutil.GenPtr[int64](12),
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_5_accounts_and_none_of_the_accounts_have_entitlements": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr[string]("testaccountId21"),
					CollectionCursor: testutil.GenPtr[int64](20),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusNoContent,
				NextCursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](25),
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_for_last_few_accounts_and_entitlements_not_more_than_page_size": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](300),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId300-entitlementId_One",
						"accountId":     "testaccountId300",
						"entitlementId": "entitlementId_One",
					},
					{
						"id":            "testaccountId300-entitlementId_Two",
						"accountId":     "testaccountId300",
						"entitlementId": "entitlementId_Two",
					},
					{
						"id":            "testaccountId301-entitlementId_Three",
						"accountId":     "testaccountId301",
						"entitlementId": "entitlementId_Three",
					},
					{
						"id":            "testaccountId301-entitlementId_Four",
						"accountId":     "testaccountId301",
						"entitlementId": "entitlementId_Four",
					},
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_for_5_accounts_with_more_than_pagesize_entitlements_for_first_account": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](80),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId80-entitlementId801",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId801",
					},
					{
						"id":            "testaccountId80-entitlementId802",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId802",
					},
					{
						"id":            "testaccountId80-entitlementId803",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId803",
					},
					{
						"id":            "testaccountId80-entitlementId804",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId804",
					},
					{
						"id":            "testaccountId80-entitlementId805",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId805",
					},
					{
						"id":            "testaccountId80-entitlementId806",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId806",
					},
					{
						"id":            "testaccountId80-entitlementId807",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId807",
					},
					{
						"id":            "testaccountId80-entitlementId808",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId808",
					},
					{
						"id":            "testaccountId80-entitlementId809",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId809",
					},
					{
						"id":            "testaccountId80-entitlementId810",
						"accountId":     "testaccountId80",
						"entitlementId": "entitlementId810",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](10),
					CollectionID:     testutil.GenPtr[string]("testaccountId80"),
					CollectionCursor: testutil.GenPtr[int64](80),
				},
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_for_last_few_accounts_and_entitlements_more_than_page_size": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](302),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId302-entitlementId_a",
						"accountId":     "testaccountId302",
						"entitlementId": "entitlementId_a",
					},
					{
						"id":            "testaccountId302-entitlementId_b",
						"accountId":     "testaccountId302",
						"entitlementId": "entitlementId_b",
					},
					{
						"id":            "testaccountId302-entitlementId_c",
						"accountId":     "testaccountId302",
						"entitlementId": "entitlementId_c",
					},
					{
						"id":            "testaccountId303-entitlementId_d",
						"accountId":     "testaccountId303",
						"entitlementId": "entitlementId_d",
					},
					{
						"id":            "testaccountId303-entitlementId_e",
						"accountId":     "testaccountId303",
						"entitlementId": "entitlementId_e",
					},
					{
						"id":            "testaccountId303-entitlementId_f",
						"accountId":     "testaccountId303",
						"entitlementId": "entitlementId_f",
					},
					{
						"id":            "testaccountId303-entitlementId_g",
						"accountId":     "testaccountId303",
						"entitlementId": "entitlementId_g",
					},
					{
						"id":            "testaccountId304-entitlementId_h",
						"accountId":     "testaccountId304",
						"entitlementId": "entitlementId_h",
					},
					{
						"id":            "testaccountId304-entitlementId_i",
						"accountId":     "testaccountId304",
						"entitlementId": "entitlementId_i",
					},
					{
						"id":            "testaccountId304-entitlementId_j",
						"accountId":     "testaccountId304",
						"entitlementId": "entitlementId_j",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](3),
					CollectionID:     testutil.GenPtr[string]("testaccountId304"),
					CollectionCursor: testutil.GenPtr[int64](304),
				},
			},
			wantErr: nil,
		},
		// This would not be a reality but just to test out a certain scenario, the account endpoint
		// returns an empty body.
		"fetch_account_entitlements_for_accounts_where_no_accounts_are_returned": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](400),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusNoContent,
			},
			wantErr: nil,
		},
		"fetch_account_entitlements_5_accounts_and_only_first_two_accounts_have_entitlements": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					CollectionCursor: testutil.GenPtr[int64](450),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"id":            "testaccountId450-entitlementId_a",
						"accountId":     "testaccountId450",
						"entitlementId": "entitlementId_a",
					},
					{
						"id":            "testaccountId450-entitlementId_b",
						"accountId":     "testaccountId450",
						"entitlementId": "entitlementId_b",
					},
					{
						"id":            "testaccountId451-entitlementId_c",
						"accountId":     "testaccountId451",
						"entitlementId": "entitlementId_c",
					},
					{
						"id":            "testaccountId451-entitlementId_d",
						"accountId":     "testaccountId451",
						"entitlementId": "entitlementId_d",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					CollectionID:     testutil.GenPtr[string]("testaccountId451"),
					CollectionCursor: testutil.GenPtr[int64](455),
				},
			},
			wantErr: nil,
		},
		"account_entitlement_account_has_no_id": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           nil,
					CollectionID:     testutil.GenPtr[string]("testaccountId110"),
					CollectionCursor: testutil.GenPtr[int64](110),
				},
			},
			wantErr: &framework.Error{
				Message: "Failed to validate required fields for an account when fetching its entitlements.: Key: 'AccountObject.AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"account_entitlement_account_id_not_string": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           nil,
					CollectionID:     testutil.GenPtr[string]("testaccountId120"),
					CollectionCursor: testutil.GenPtr[int64](120),
				},
			},
			wantErr: &framework.Error{
				Message: "Failed to decode account data when fetching its entitlements.: 1 error(s) decoding:\n\n* 'id' expected type 'string', got unconvertible type 'float64', value: '120'.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"account_entitlement_id_not_string": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "accountEntitlements",
				PageSize:              10,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor:           testutil.GenPtr[int64](99),
					CollectionID:     testutil.GenPtr[string]("testaccountId250"),
					CollectionCursor: testutil.GenPtr[int64](250),
				},
			},
			wantErr: &framework.Error{
				Message: "Failed to convert IdentityNow account entitlement object id field to string: 250.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := identitynowClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetEntitlementsPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	identitynowClient := identitynow.NewClient(client, 0)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *identitynow.Request
		wantRes *identitynow.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "entitlements",
				PageSize:              2,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]any{
					{
						"attribute": "appRoleAssignments",
						"attributes": map[string]any{
							"displayName": "Basic purchaser [on] Windows Store for Business",
						},
						"cloudGoverned":          false,
						"created":                "2023-09-22T16:50:12.053Z",
						"description":            nil,
						"id":                     "ENTITLEMENT_ID_456",
						"modified":               "2023-09-22T16:56:12.896Z",
						"privileged":             true,
						"requestable":            true,
						"sourceSchemaObjectType": "applicationRole",
						"value":                  "AZURE_APP_ID_123:AZURE_RESOURCE_ID_456",
					},
					{
						"attribute": "appRoleAssignments",
						"attributes": map[string]any{
							"displayName": "default access [on] TrustedPublishersProxyService",
						},
						"cloudGoverned":          false,
						"created":                "2023-09-22T16:50:54.856Z",
						"description":            nil,
						"id":                     "00218206fe614e7da637f528accdf15e",
						"modified":               "2023-09-22T16:56:13.788Z",
						"privileged":             true,
						"requestable":            true,
						"sourceSchemaObjectType": "applicationRole",
						"value":                  "efd1eb6f-44f3-4ffe-b4e4-eb68162ea4ae:00000000-0000-0000-0000-000000000000",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
			},
			wantErr: nil,
		},
		"not_found_response": {
			context: context.Background(),
			request: &identitynow.Request{
				Token:                 "Bearer token",
				BaseURL:               server.URL,
				EntityExternalID:      "entitlements",
				PageSize:              2,
				APIVersion:            "beta",
				RequestTimeoutSeconds: 5,
				// The test handler has been hardcoded to serve a 404 for a cursor value of 99.
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](99),
				},
			},
			wantRes: &identitynow.Response{
				StatusCode: http.StatusNotFound,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := identitynowClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
