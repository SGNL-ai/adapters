// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package azuread_test

import (
	"net/http"
	"strings"
)

// CreateTestServerHandler creates a handler that replaces https://graph.microsoft.com with the test server URL.
func CreateTestServerHandler(serverURL string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the original handler
		TestServerHandler.ServeHTTP(&responseRewriter{ResponseWriter: w, serverURL: serverURL}, r)
	})
}

// responseRewriter wraps http.ResponseWriter to replace URLs in responses.
type responseRewriter struct {
	http.ResponseWriter
	serverURL string
}

func (rw *responseRewriter) Write(b []byte) (int, error) {
	// Replace https://graph.microsoft.com with the test server URL
	modified := strings.ReplaceAll(string(b), "https://graph.microsoft.com", rw.serverURL)

	return rw.ResponseWriter.Write([]byte(modified))
}

// Define the endpoints and responses for the mock Azure AD server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer Testtoken" {
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
	case "/v1.0/users?$select=id&$top=2":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA",
			"value": [
				{
					"businessPhones": [],
					"displayName": "Conf Room Adams",
					"givenName": null,
					"jobTitle": null,
					"mail": "Adams@M365x214355.onmicrosoft.com",
					"mobilePhone": null,
					"officeLocation": null,
					"preferredLanguage": null,
					"surname": null,
					"userPrincipalName": "Adams@M365x214355.onmicrosoft.com",
					"id": "6e7b768e-07e2-4810-8459-485f84f8f204"
				},
				{
					"businessPhones": [
						"+1 425 555 0109"
					],
					"displayName": "Adele Vance",
					"givenName": "Adele",
					"jobTitle": "Product Marketing Manager",
					"mail": "AdeleV@M365x214355.onmicrosoft.com",
					"mobilePhone": null,
					"officeLocation": "18/2111",
					"preferredLanguage": "en-US",
					"surname": "Vance",
					"userPrincipalName": "AdeleV@M365x214355.onmicrosoft.com",
					"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd"
				}
			]
		}`))

	// Users Page 2
	case "/v1.0/users?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"value": [
				{
					"businessPhones": [
						"8006427676"
					],
					"displayName": "MOD Administrator",
					"givenName": "MOD",
					"jobTitle": null,
					"mail": "admin@M365x214355.onmicrosoft.com",
					"mobilePhone": "5555555555",
					"officeLocation": null,
					"preferredLanguage": "en-US",
					"surname": "Administrator",
					"userPrincipalName": "admin@M365x214355.onmicrosoft.com",
					"id": "5bde3e51-d13b-4db1-9948-fe4b109d11a7"
				},
				{
					"businessPhones": [
						"+1 858 555 0110"
					],
					"displayName": "Alex Wilber",
					"givenName": "Alex",
					"jobTitle": "Marketing Assistant",
					"mail": "AlexW@M365x214355.onmicrosoft.com",
					"mobilePhone": null,
					"officeLocation": "131/1104",
					"preferredLanguage": "en-US",
					"surname": "Wilber",
					"userPrincipalName": "AlexW@M365x214355.onmicrosoft.com",
					"id": "4782e723-f4f4-4af3-a76e-25e3bab0d896"
				}
			]
		}`))

	// Groups Page 1:
	case "/v1.0/groups?$select=id&$top=2":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA",
			"value": [
				{
					"id": "02bd9fd6-8f93-4758-87c3-1fb73740a315",
					"deletedDateTime": null,
					"classification": null,
					"createdDateTime": "2017-07-31T18:56:16Z",
					"creationOptions": [
						"ExchangeProvisioningFlags:481"
					],
					"description": "Welcome to the HR Taskforce team.",
					"displayName": "HR Taskforce",
					"expirationDateTime": null,
					"groupTypes": [
						"Unified"
					],
					"isAssignableToRole": null,
					"mail": "HRTaskforce@M365x214355.onmicrosoft.com",
					"mailEnabled": true,
					"mailNickname": "HRTaskforce",
					"membershipRule": null,
					"membershipRuleProcessingState": null,
					"onPremisesDomainName": null,
					"onPremisesLastSyncDateTime": null,
					"onPremisesNetBiosName": null,
					"onPremisesSamAccountName": null,
					"onPremisesSecurityIdentifier": null,
					"onPremisesSyncEnabled": null,
					"preferredDataLocation": null,
					"preferredLanguage": null,
					"proxyAddresses": [
						"SMTP:HRTaskforce@M365x214355.onmicrosoft.com",
						"SPO:SPO_896cf652-b200-4b74-8111-c013f64406cf@SPO_dcd219dd-bc68-4b9b-bf0b-4a33a796be35"
					],
					"renewedDateTime": "2020-01-24T19:01:14Z",
					"resourceBehaviorOptions": [],
					"resourceProvisioningOptions": [
						"Team"
					],
					"securityEnabled": false,
					"securityIdentifier": "S-1-12-1-45981654-1196986259-3072312199-363020343",
					"theme": null,
					"visibility": "Private",
					"onPremisesProvisioningErrors": [],
					"serviceProvisioningErrors": []
				},
				{
					"id": "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23",
					"deletedDateTime": null,
					"classification": null,
					"createdDateTime": "2017-07-31T17:38:15Z",
					"creationOptions": [],
					"description": "Video Production",
					"displayName": "Video Production",
					"expirationDateTime": null,
					"groupTypes": [
						"Unified"
					],
					"isAssignableToRole": null,
					"mail": "VideoProduction@M365x214355.onmicrosoft.com",
					"mailEnabled": true,
					"mailNickname": "VideoProduction",
					"membershipRule": null,
					"membershipRuleProcessingState": null,
					"onPremisesDomainName": null,
					"onPremisesLastSyncDateTime": null,
					"onPremisesNetBiosName": null,
					"onPremisesSamAccountName": null,
					"onPremisesSecurityIdentifier": null,
					"onPremisesSyncEnabled": null,
					"preferredDataLocation": null,
					"preferredLanguage": null,
					"proxyAddresses": [
						"SMTP:VideoProduction@M365x214355.onmicrosoft.com",
						"SPO:SPO_16219fd2-fafd-4fea-8084-8b5eaa8c5ad2@SPO_dcd219dd-bc68-4b9b-bf0b-4a33a796be35"
					],
					"renewedDateTime": "2017-07-31T17:38:15Z",
					"resourceBehaviorOptions": [],
					"resourceProvisioningOptions": [],
					"securityEnabled": true,
					"securityIdentifier": "S-1-12-1-116797296-1315870759-261025683-595303213",
					"theme": null,
					"visibility": "Public",
					"onPremisesProvisioningErrors": [],
					"serviceProvisioningErrors": []
				}
			]
		}`))

	// Groups Page 2:
	case "/v1.0/groups?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups",
			"value": [
				{
					"id": "0a53828f-36c9-44c3-be3d-99a7fce977ac",
					"deletedDateTime": null,
					"classification": null,
					"createdDateTime": "2017-09-02T02:54:25Z",
					"creationOptions": [
						"YammerProvisioning"
					],
					"description": "Marketing Campaigns",
					"displayName": "Marketing Campaigns",
					"expirationDateTime": null,
					"groupTypes": [
						"Unified"
					],
					"isAssignableToRole": null,
					"mail": "marketingcampaigns@M365x214355.onmicrosoft.com",
					"mailEnabled": true,
					"mailNickname": "marketingcampaigns",
					"membershipRule": null,
					"membershipRuleProcessingState": null,
					"onPremisesDomainName": null,
					"onPremisesLastSyncDateTime": null,
					"onPremisesNetBiosName": null,
					"onPremisesSamAccountName": null,
					"onPremisesSecurityIdentifier": null,
					"onPremisesSyncEnabled": null,
					"preferredDataLocation": null,
					"preferredLanguage": null,
					"proxyAddresses": [
						"SMTP:marketingcampaigns@M365x214355.onmicrosoft.com",
						"SPO:SPO_8cfbec68-642c-4d90-a15e-68d5d55e1c1f@SPO_dcd219dd-bc68-4b9b-bf0b-4a33a796be35"
					],
					"renewedDateTime": "2017-09-02T02:54:25Z",
					"resourceBehaviorOptions": [
						"YammerProvisioning"
					],
					"resourceProvisioningOptions": [],
					"securityEnabled": false,
					"securityIdentifier": "S-1-12-1-173245071-1153644233-2811837886-2893539836",
					"theme": null,
					"visibility": "Public",
					"onPremisesProvisioningErrors": [],
					"serviceProvisioningErrors": []
				},
				{
					"id": "1381c058-2ee8-41ce-a005-aae7d91fe086",
					"deletedDateTime": null,
					"classification": null,
					"createdDateTime": "2017-09-15T01:03:59Z",
					"creationOptions": [
						"ProvisionGroupHomepage"
					],
					"description": "Where we share innovative ideas.",
					"displayName": "Ideas",
					"expirationDateTime": null,
					"groupTypes": [
						"Unified"
					],
					"isAssignableToRole": null,
					"mail": "Ideas@M365x214355.onmicrosoft.com",
					"mailEnabled": true,
					"mailNickname": "Ideas",
					"membershipRule": null,
					"membershipRuleProcessingState": null,
					"onPremisesDomainName": null,
					"onPremisesLastSyncDateTime": null,
					"onPremisesNetBiosName": null,
					"onPremisesSamAccountName": null,
					"onPremisesSecurityIdentifier": null,
					"onPremisesSyncEnabled": null,
					"preferredDataLocation": null,
					"preferredLanguage": null,
					"proxyAddresses": [
						"SMTP:Ideas@M365x214355.onmicrosoft.com",
						"SPO:SPO_877ab59c-6bed-4a01-9a7f-ced52a812f57@SPO_dcd219dd-bc68-4b9b-bf0b-4a33a796be35"
					],
					"renewedDateTime": "2017-09-15T01:03:59Z",
					"resourceBehaviorOptions": [],
					"resourceProvisioningOptions": [],
					"securityEnabled": false,
					"securityIdentifier": "S-1-12-1-327270488-1104031464-3886679456-2262835161",
					"theme": null,
					"visibility": "Public",
					"onPremisesProvisioningErrors": [],
					"serviceProvisioningErrors": []
				}
			]
		}`))

	// Applications Page 1
	case "/v1.0/applications?$select=id&$top=2":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#applications",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/applications?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA",
			"value": [
				{
					"id": "acc848e9-e8ec-4feb-a521-8d58b5482e09",
					"deletedDateTime": null,
					"appId": "05b10a2d-62db-420c-8626-55f3a5e7865b",
					"applicationTemplateId": null,
					"disabledByMicrosoftStatus": null,
					"createdDateTime": "2020-02-13T21:09:15Z",
					"displayName": "apisandboxproxy",
					"description": null,
					"groupMembershipClaims": null,
					"identifierUris": [
						"https://M365x214355.onmicrosoft.com/apisandboxproxy"
					],
					"isDeviceOnlyAuthSupported": null,
					"isFallbackPublicClient": true,
					"notes": null,
					"publisherDomain": "M365x214355.onmicrosoft.com",
					"serviceManagementReference": null,
					"signInAudience": "AzureADMyOrg",
					"tags": [],
					"tokenEncryptionKeyId": null,
					"samlMetadataUrl": null,
					"defaultRedirectUri": null,
					"certification": null,
					"optionalClaims": null,
					"servicePrincipalLockConfiguration": null,
					"requestSignatureVerification": null,
					"addIns": [],
					"api": {
						"acceptMappedClaims": null,
						"knownClientApplications": [],
						"requestedAccessTokenVersion": null,
						"oauth2PermissionScopes": [],
						"preAuthorizedApplications": []
					},
					"appRoles": [],
					"info": {
						"logoUrl": null,
						"marketingUrl": null,
						"privacyStatementUrl": null,
						"supportUrl": null,
						"termsOfServiceUrl": null
					},
					"keyCredentials": [],
					"parentalControlSettings": {
						"countriesBlockedForMinors": [],
						"legalAgeGroupRule": "Allow"
					}
				},
				{
					"id": "cfa98ac0-a32c-4b4c-a78b-94c9912ed7b2",
					"deletedDateTime": null,
					"appId": "c305b21c-fda6-4ecb-aa01-8a8141fdfd51",
					"applicationTemplateId": null,
					"disabledByMicrosoftStatus": null,
					"createdDateTime": "2018-03-27T02:45:04Z",
					"displayName": "EduPopulationHelper",
					"description": null,
					"groupMembershipClaims": null,
					"identifierUris": [
						"https://M365x214355.onmicrosoft.com/a60c216f-657f-4925-980a-d8ef69942167"
					],
					"isDeviceOnlyAuthSupported": null,
					"isFallbackPublicClient": false,
					"notes": null,
					"publisherDomain": null,
					"serviceManagementReference": null,
					"signInAudience": "AzureADMyOrg",
					"tags": [],
					"tokenEncryptionKeyId": null,
					"samlMetadataUrl": null,
					"defaultRedirectUri": null,
					"certification": null,
					"optionalClaims": null,
					"servicePrincipalLockConfiguration": null,
					"requestSignatureVerification": null,
					"addIns": [],
					"appRoles": [],
					"info": {
						"logoUrl": null,
						"marketingUrl": null,
						"privacyStatementUrl": null,
						"supportUrl": null,
						"termsOfServiceUrl": null
					},
					"keyCredentials": [],
					"parentalControlSettings": {
						"countriesBlockedForMinors": [],
						"legalAgeGroupRule": "Allow"
					},
					"publicClient": {
						"redirectUris": []
					}
				}
			]
		}`))

	// Applications Page 2
	case "/v1.0/applications?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#applications",
		"value": [
			{
				"id": "d6dbf9e0-98a4-4eea-b4c1-df8695277868",
				"deletedDateTime": null,
				"appId": "377a3df7-6ff0-42b2-ad55-194dcc7aacd9",
				"applicationTemplateId": null,
				"disabledByMicrosoftStatus": null,
				"createdDateTime": "2020-07-24T07:54:39Z",
				"displayName": "permissions-scraper-app",
				"description": null,
				"groupMembershipClaims": null,
				"identifierUris": [],
				"isDeviceOnlyAuthSupported": null,
				"isFallbackPublicClient": null,
				"notes": null,
				"publisherDomain": "M365x214355.onmicrosoft.com",
				"serviceManagementReference": null,
				"signInAudience": "AzureADandPersonalMicrosoftAccount",
				"tags": [],
				"tokenEncryptionKeyId": null,
				"samlMetadataUrl": null,
				"defaultRedirectUri": null,
				"certification": null,
				"optionalClaims": null,
				"servicePrincipalLockConfiguration": null,
				"requestSignatureVerification": null,
				"addIns": [],
				"api": {
					"acceptMappedClaims": null,
					"knownClientApplications": [],
					"requestedAccessTokenVersion": 2,
					"oauth2PermissionScopes": [],
					"preAuthorizedApplications": []
				},
				"appRoles": [],
				"info": {
					"logoUrl": null,
					"marketingUrl": null,
					"privacyStatementUrl": null,
					"supportUrl": null,
					"termsOfServiceUrl": null
				},
				"keyCredentials": [],
				"parentalControlSettings": {
					"countriesBlockedForMinors": [],
					"legalAgeGroupRule": "Allow"
				}
			}
		]
	}`))

	// Devices Page 1
	case "/v1.0/devices?$select=id&$top=2":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#devices",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/devices?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA",
		"value": [
			{
				"id":"0357837b-ca6e-402d-9429-9e54dd51d97a",
				"accountEnabled":true,
				"deviceId":"00000000-0000-0000-0000-000000000000",
				"deviceVersion":1,
				"displayName":"contoso_pixel",
				"Manufacturer":"Google",
				"Model":"Pixel 3a",
				"operatingSystemVersion":"10.0"
			},
			{
				"id":"4d1ed9a4-519e-421b-b9f6-158991feff5b",
				"accountEnabled":true,
				"deviceId":"00000000-0000-0000-0000-000000000001",
				"deviceVersion":1,
				"displayName":"contoso_galaxy",
				"Manufacturer":"Samsung",
				"Model":"Galaxy Note 7",
				"operatingSystemVersion":"8.2"
			}
		]
	}`))

	// Devices Page 2
	case "/v1.0/devices?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAADBBcHBsaWNhdGlvbl9jZmE5OGFjMC1hMzJjLTRiNGMtYTc4Yi05NGM5OTEyZWQ3YjIwQXBwbGljYXRpb25fY2ZhOThhYzAtYTMyYy00YjRjLWE3OGItOTRjOTkxMmVkN2IyAAAAAAAAAAAAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#devices",
		"value": [
			{
				"id":"6a59ea83-02bd-468f-a40b-f2c3d1821983",
				"accountEnabled":true,
				"deviceId":"00000000-0000-0000-0000-000000000002",
				"deviceVersion":1,
				"displayName":"contoso_iphone",
				"Manufacturer":"Apple",
				"Model":"iPhone 11 Pro Max",
				"operatingSystemVersion":"11.2"
			}
		]
	}`))

	// Group Members - Groups Page 1 (Page size 1):
	case "/v1.0/groups?$select=id&$top=1":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA",
			"value": [
				{
					"id": "02bd9fd6-8f93-4758-87c3-1fb73740a315"
				}
			]
		}`))

	// Group Members - Groups Page 2 (Page size 1):
	case "/v1.0/groups?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"value": [
				{
					"id": "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23"
				}
			]
		}`))

	// Group Members - Invalid Groups (Page size 1 requested, 2 returned):
	case "/v1.0/groups?$select=id&$top=1&$skiptoken=TOO-MANY-GROUPS":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"value": [
				{
					"id": "06f62f70-9827-4e6e-93ef-8e0f2d9b7b23"
				},
				{
					"id": "163edf70-9827-4e6e-93ef-8e0f2d9des82"
				}
			]
		}`))

	// Group Members - Invalid Groups (Page size 1 requested, 0 returned):
	case "/v1.0/groups?$select=id&$top=1&$skiptoken=NOT-ENOUGH-GROUPS":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"value": []
		}`))

	// Group Members - Groups Page 1 (Filtered id in ('02bd9fd6-8f93-4758-87c3-1fb73740a315','0a53828f-36c9-44c3-be3d-99a7fce977ac'))
	case "/v1.0/groups?$select=id&$top=1&$filter=id+in+%28%2702bd9fd6-8f93-4758-87c3-1fb73740a315%27%2C%270a53828f-36c9-44c3-be3d-99a7fce977ac%27%29":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$filter=id+in+%28%2702bd9fd6-8f93-4758-87c3-1fb73740a315%27%2C%270a53828f-36c9-44c3-be3d-99a7fce977ac%27%29&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA",
		"value": [
			{
				"id": "02bd9fd6-8f93-4758-87c3-1fb73740a315"
			}
		]
	}`))

	// Group Members - Groups Page 2 (Filtered id in ('02bd9fd6-8f93-4758-87c3-1fb73740a315','0a53828f-36c9-44c3-be3d-99a7fce977ac'))
	case "https://graph.microsoft.com/v1.0/groups?$select=id&$top=1&$filter=id+in+%28%2702bd9fd6-8f93-4758-87c3-1fb73740a315%27%2C%270a53828f-36c9-44c3-be3d-99a7fce977ac%27%29&$skiptoken=RFNwdAIAAQAAACpHcm91cF8wNmY2MmY3MC05ODI3LTRlNmUtOTNlZi04ZTBmMmQ5YjdiMjMqR3JvdXBfMDZmNjJmNzAtOTgyNy00ZTZlLTkzZWYtOGUwZjJkOWI3YjIzAAAAAAAAAAAAAAA":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#groups(id)",
			"value": [
				{
					"id": "0a53828f-36c9-44c3-be3d-99a7fce977ac"
				}
			]
		}`))

	// Group Members - 02bd9fd6-8f93-4758-87c3-1fb73740a315 - Members Page 1:
	case "/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA",
			"value": [
				{
					"id": "6e7b768e-07e2-4810-8459-485f84f8f204"
				},
				{
					"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd"
				}
			]
		}`))

	// Group Members - 02bd9fd6-8f93-4758-87c3-1fb73740a315 - Members Page 1 with 3 top value:
	case "/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=3":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA",
			"value": [
				{
					"id": "6e7b768e-07e2-4810-8459-485f84f8f204"
				},
				{
					"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd"
				}
			]
		}`))

	// Group Members - 02bd9fd6-8f93-4758-87c3-1fb73740a315 - Members Page 2:
	case "/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
		"value": [
			{
				"id": "5bde3e51-d13b-4db1-9948-fe4b109d11a7"
			},
			{
				"id": "4782e723-f4f4-4af3-a76e-25e3bab0d896"
			}
		]
	}`))

	// Group Members - 06f62f70-9827-4e6e-93ef-8e0f2d9b7b23 - Members Page 1:
	case "/v1.0/groups/06f62f70-9827-4e6e-93ef-8e0f2d9b7b23/members?$select=id&$top=2":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/groups/06f62f70-9827-4e6e-93ef-8e0f2d9b7b23/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA",
			"value": [
				{
					"id": "6e7b768e-07e2-4810-8459-485f84f8f204"
				},
				{
					"id": "87d349ed-44d7-43e1-9a83-5f2406dee5bd"
				}
			]
		}`))

	// Group Members - 06f62f70-9827-4e6e-93ef-8e0f2d9b7b23 - Members Page 2:
	case "/v1.0/groups/06f62f70-9827-4e6e-93ef-8e0f2d9b7b23/members?$select=id&$top=2&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"value": [
				{
					"id": "5bde3e51-d13b-4db1-9948-fe4b109d11a7"
				},
				{
					"id": "4782e723-f4f4-4af3-a76e-25e3bab0d896"
				}
			]
		}`))

	// Group Members - 02bd9fd6-8f93-4758-87c3-1fb73740a315 (Filtered id eq '6e7b768e-07e2-4810-8459-485f84f8f204') - Members Page 1
	case "/v1.0/groups/02bd9fd6-8f93-4758-87c3-1fb73740a315/members?$select=id&$top=1&$filter=id+eq+%276e7b768e-07e2-4810-8459-485f84f8f204%27":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAAQAAACM6QWRlbGVWQE0zNjV4MjE0MzU1Lm9ubWljcm9zb2Z0LmNvbSlVc2VyXzg3ZDM0OWVkLTQ0ZDctNDNlMS05YTgzLTVmMjQwNmRlZTViZLkAAAAAAAAAAAAA",
			"value": [
				{
					"id": "6e7b768e-07e2-4810-8459-485f84f8f204"
				}
			]
		}`))

	// Group Members - 0a53828f-36c9-44c3-be3d-99a7fce977ac (Filtered id eq '6e7b768e-07e2-4810-8459-485f84f8f204') - Members Page 1
	// Note: Filtering on a specific user that is not present in the group will return a 404 when interacting with the real AAD servers.
	// This is just being used as an example filter so we're not mimicking this behavior.
	case "/v1.0/groups/0a53828f-36c9-44c3-be3d-99a7fce977ac/members?$select=id&$top=2&$filter=id+eq+%276e7b768e-07e2-4810-8459-485f84f8f204%27":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"value": []
		}`))

	case "/v1.0/directoryRoles?$select=id,displayName,description,roleTemplateId":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
			"value": [
				{
					"id": "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"description": "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName": "Directory Readers",
					"roleTemplateId": "88d8e3e3-c189-46e8-94e1-9b9898b8876b"
				},
				{
					"id": "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"description": "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName": "Global Administrator",
					"roleTemplateId": "62e90394-3621-4004-a7cb-012177145e10"
				},
				{
					"id": "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"description": "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName": "Application Developer",
					"roleTemplateId": "cf1c38e5-69f5-4237-9190-879624dced7c"
				},
				{
					"id": "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"description": "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName": "Privileged Role Administrator",
					"roleTemplateId": "e8611ab8-8f55-4a1e-953a-60213ab1f814"
				},
				{
					"id": "62ceaa28-7382-48f9-j386-f8ed6e9a7c84",
					"description": "Can read and write basic directory information. For granting access to applications, not intended for users.",
					"displayName": "Directory Writers",
					"roleTemplateId": "9360feb5-f418-4baa-8175-e2a00bac4301"
				}
			]
		}`))

	case "/v1.0/directoryRoles?$select=id,displayName,description,roleTemplateId,deletedDateTime":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
			"value": [
				{
					"id": "0fea7f0d-dea1-4028-8ce8-a686ec639d75",
					"deletedDateTime": "2024-02-02T23:21:02Z",
					"description": "Can read basic directory information. Commonly used to grant directory read access to applications and guests.",
					"displayName": "Directory Readers",
					"roleTemplateId": "88d8e3e3-c189-46e8-94e1-9b9898b8876b"
				},
				{
					"id": "18eacdf7-8db3-458d-9099-69fcc2e3cd42",
					"deletedDateTime": null,
					"description": "Can manage all aspects of Microsoft Entra ID and Microsoft services that use Microsoft Entra identities.",
					"displayName": "Global Administrator",
					"roleTemplateId": "62e90394-3621-4004-a7cb-012177145e10"
				},
				{
					"id": "33a4c989-c3ff-4a77-bf46-ee0acd84476e",
					"deletedDateTime": null,
					"description": "Can create application registrations independent of the 'Users can register applications' setting.",
					"displayName": "Application Developer",
					"roleTemplateId": "cf1c38e5-69f5-4237-9190-879624dced7c"
				},
				{
					"id": "321fd63c-c37c-4597-81c4-81e0a93ffb6e",
					"deletedDateTime": null,
					"description": "Can manage role assignments in Microsoft Entra ID, and all aspects of Privileged Identity Management.",
					"displayName": "Privileged Role Administrator",
					"roleTemplateId": "e8611ab8-8f55-4a1e-953a-60213ab1f814"
				},
				{
					"id": "62ceaa28-7382-48f9-j386-f8ed6e9a7c84",
					"deletedDateTime": null,
					"description": "Can read and write basic directory information. For granting access to applications, not intended for users.",
					"displayName": "Directory Writers",
					"roleTemplateId": "9360feb5-f418-4baa-8175-e2a00bac4301"
				}
			]
		}`))

	// Role Members - Users Page 1 (Page size 1):
	case "/v1.0/users?$select=id&$top=1":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1",
			"value": [
				{
					"id": "65bb46a4-7d3j-9302-8a21-4d90f7a0efdb"
				}
			]
		}`))

	// Role Members - Users Page 2 (Page size 1):
	case "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_1":
		w.Write([]byte(`{
			"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
			"@odata.nextLink": "https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA",
			"value": [
				{

					"id": "df102bb2-2365-235g-a2g6-edb774169548"
				}
			]
	}`))

	// Role Members - Users Page 3 (Page size 1):
	case "/v1.0/users?$select=id&$top=1&$skiptoken=RFNwdAIAADA6YWFyb24uYXlhbGE4OTUwNzM0N0BzZ25sYWFkZGV2MS5vbm1pY3Jvc29mdC5jb20pVXNlcl9kZjEwMmJiMi0zNTMyLTQ1M2MtYTNiNC1lZGI3NzQxNjk1NDgAMDphYXJvbi5heWFsYTg5NTA3MzQ3QHNnbmxhYWRkZXYxLm9ubWljcm9zb2Z0LmNvbSlVc2VyX2RmMTAyYmIyLTM1MzItNDUzYy1hM2I0LWVkYjc3NDE2OTU0OLkAAAAAAAAAAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_2",
		"value": [
		{
			"id": "201d31c0-653d-43a6-adf0-aee89a79c805"
		}
	]
	}`))

	// Role Members - Users Page 4 (Page size 1):
	case "/v1.0/users?$select=id&$top=1&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_2":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users",
		"value": [
		{
			"id": "uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy"
		}
	]
	}`))

	// Role Members - 65bb46a4-7d3j-9302-8a21-4d90f7a0efdb - Members Page 1:
	case "/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4",
		"value": [
			{
				"id": "0fea7f0d-dea1-458d-9099-69fcc2e3cd42"
			},
			{
				"id": "795326a8-6eef-410e-9604-649ca68e1241"
			}
		]
	}`))

	// Role Members - 65bb46a4-7d3j-9302-8a21-4d90f7a0efdb - Members Page 2:

	case "/v1.0/users/65bb46a4-7d3j-9302-8a21-4d90f7a0efdb/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"value": [
			{
				"id": "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84"
			}
		]
	}`))

	// Role Members - df102bb2-2365-235g-a2g6-edb774169548 - Members Page 1:
	case "/v1.0/users/df102bb2-2365-235g-a2g6-edb774169548/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=3":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users/df102bb2-2365-235g-a2g6-edb774169548/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4",
		"value": [
			{
				"id": "62ceaa28-4794-48f9-9b54-f8ed6e9a7c84"
			},
			{
				"id": "795326a8-6eef-410e-9604-649ca68e1241"
			}
		]
	}`))

	// Role Members - df102bb2-2365-235g-a2g6-edb774169548 - Members Page 2:
	case "/v1.0/users/df102bb2-2365-235g-a2g6-edb774169548/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2&$skiptoken=NEXTLINK_TOKEN_PLACEHOLDER_4":
		w.Write([]byte(`{
	"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
	"value": [
		{
			"id": "e8d9279e-6883-4add-96e8-5f7c8df5637f"
		}
	]
	}`))

	// Role Members - 201d31c0-653d-43a6-adf0-aee89a79c805 - Members Page 1:
	case "/v1.0/users/201d31c0-653d-43a6-adf0-aee89a79c805/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"value": []
	}`))

	// Role Members - uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy - Members Page 1:
	case "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$select=id&$top=2":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$select=id%2cdisplayName&$top=2&$skiptoken=RFNwdAoAAQAAAAAAAAAAFAAAABkp8fswrv1Ls8cLjYDqBRABAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA",
		"value": [
			{
				"id": "0fea7f0d-dea1-458d-9099-69fcc2e3cd42"
			},
			{
				"id": "d973db57-eb50-4356-959e-f1ce19a22b98"
			}
		]
	}`))

	// Role Members - uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy - Members Page 2:
	// No records with nextLink
	case "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$select=id%2cdisplayName&$top=2&$skiptoken=RFNwdAoAAQAAAAAAAAAAFAAAABkp8fswrv1Ls8cLjYDqBRABAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAAPWE8iLxC5NNtqCdf_NZ8bcCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA",
		"value": []
	}`))

	// Role Members - uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy - Members Page 3:
	// less records with nextLink
	case "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAAPWE8iLxC5NNtqCdf_NZ8bcCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA":
		w.Write([]byte(`{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
		"@odata.nextLink": "https://graph.microsoft.com/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAABgFnxJuzI1NsFSV18Bt7PgCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA",
		"value": [
			{
				"id": "540b4b34-c25b-437d-8eee-329463952334"
			}
		]
	}`))

	// Role Members - uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy - Members Page 4:
	// last page
	case "/v1.0/users/uuuuuuuu-vvvv-wwww-xxxx-yyyyyyyyyyyy/transitiveMemberOf/microsoft.graph.directoryRole?$top=2&$select=id%2cdisplayName&$skiptoken=RFNwdAoAAAAAAAAAAAAAFAAAABgFnxJuzI1NsFSV18Bt7PgCAAAAAAAAAAAAAAAAAAAXMS4yLjg0MC4xMTM1NTYuMS40LjIzMzEGAAAAAY8MlBPpl2xBua2SNJARSM0AAfn9agujeJBOp41SpLihArMBzAAAAAEBAAAA":
		w.Write([]byte(`{
	"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#directoryRoles",
	"value": [
		{
			"id": "fc6c3c82-669c-4e24-b089-2a2847a43d14"
		}
	]
	}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
