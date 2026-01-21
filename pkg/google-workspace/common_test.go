// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package googleworkspace_test

import (
	"net/http"
)

// Define the endpoints and responses for the mock Google Workspace server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer Testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"errorCode": "E0000011",
			"errorSummary": "Invalid auth provided",
			"errorLink": "E0000011",
			"errorId": "oaefW5oDjyLRLKVkrmTlp0Thg",
			"errorCauses": []
		}}`))

		return
	}

	switch r.URL.RequestURI() {
	// Users: Total of 3, 1 per request
	case "/admin/directory/v1/users?domain=sgnldemos.com&maxResults=1":
		w.Write([]byte(`{
			"kind": "admin#directory#users",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/Ero2jiIUqpkIHHjIvNRtJmBi99k\"",
			"users": [
				{
					"kind": "admin#directory#user",
					"id": "USER987654321",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/HFAgBzXOF0pogUtDLHSgZQ9lvls\"",
					"primaryEmail": "user1@sgnldemos.com",
					"name": {
						"givenName": "user1",
						"familyName": "user1",
						"fullName": "user1 user1"
					},
					"isAdmin": true,
					"isDelegatedAdmin": false,
					"lastLoginTime": "2024-02-02T23:30:53.000Z",
					"creationTime": "2024-02-02T23:30:06.000Z",
					"agreedToTerms": true,
					"suspended": false,
					"archived": false,
					"changePasswordAtNextLogin": false,
					"ipWhitelisted": false,
					"emails": [
						{
							"address": "user1@sgnldemos.com",
							"primary": true
						},
						{
							"address": "user1@sgnldemos.com.test-google-a.com"
						}
					],
					"languages": [
						{
							"languageCode": "en",
							"preference": "preferred"
						}
					],
					"nonEditableAliases": [
						"user1@sgnldemos.com.test-google-a.com"
					],
					"customerId": "CUST123456",
					"orgUnitPath": "/",
					"isMailboxSetup": true,
					"isEnrolledIn2Sv": false,
					"isEnforcedIn2Sv": false,
					"includeInGlobalAddressList": true,
					"recoveryEmail": "user1@sgnl.ai"
				}
			],
			"nextPageToken": "Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ=="
		}`))
	case "/admin/directory/v1/users?domain=sgnldemos.com&maxResults=1&pageToken=Q0FFUzl3SUJrUHpWQUhTWWhPa2VQc0dabG1wVlBxdld0N2NtenZPTDVMM29ZNTZKK2d6ME5NczdVNTFqWXRUZnVWclQzd0c3cVZWZk8rR0M0V1FhTWpyV2R3cnZzdS9ld3BUamNkSm1BMmhEZjFIcGw5QUhyWG8yazV1WjJONzdnNHhaMjU5R2pSenIveWpOeDZOT093bk1NY20zYW9DVitJMmc0WUlEajRDYW93RDk4T0M2MjM4aGlxdUE5U05BT2lnemtlYnpJSXorVGU5SUFSYzRBd2M5OUlJdngrQ3YzMkRLekFTWGpuSlR2aERqTHMyNnpxNXZLZG9yekxKSUZDdWFUdjd1OVpJZWNMWDBXejhwVWdUMDNSbVhUSzJ1eGoxeVBaanJXalkrUzZWWlJ2MnZJeUUrM3lkU3hqUHFsL2w1S0FwUG9lZ21TbUh0MHZiR3YvV04raUQxczh3NnFTNnJCNFhLS20wZExVR21zSkNJR09Vc2lXQjhoVjFkZUkrZmlDNW85aW5lUWFjYkZsTXVacS90ejlWM0JEbnk1Uk5mUDk1cm1Ed2hsRXRjMUZTeGNocU02QUw1OFhHVFVDTTVnOWZZYVZJN0E4THE5WHNITE9Wd3gyQ1FPbWk5cURMYy9pdXp1S1NJeE8xbTlFaXhyRlFFZDBRPQ%3D%3D":
		w.Write([]byte(`{
			"kind": "admin#directory#users",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/cmjch8MWvWZkyzHS25bLHspxSPg\"",
			"users": [
				{
					"kind": "admin#directory#user",
					"id": "102475661842232156723",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/zEsIEnG_BTp0IWmiowHAH5dGdhI\"",
					"primaryEmail": "user2@sgnldemos.com",
					"name": {
						"givenName": "user2",
						"familyName": "user2",
						"fullName": "user2 user2"
					},
					"isAdmin": true,
					"isDelegatedAdmin": false,
					"lastLoginTime": "2024-04-18T23:55:53.000Z",
					"creationTime": "2024-02-02T23:55:44.000Z",
					"agreedToTerms": true,
					"suspended": false,
					"archived": false,
					"changePasswordAtNextLogin": false,
					"ipWhitelisted": false,
					"emails": [
						{
							"address": "user2@sgnldemos.com",
							"primary": true
						},
						{
							"address": "user2@sgnldemos.com.test-google-a.com"
						}
					],
					"languages": [
						{
							"languageCode": "en",
							"preference": "preferred"
						}
					],
					"nonEditableAliases": [
						"user2@sgnldemos.com.test-google-a.com"
					],
					"customerId": "CUST123456",
					"orgUnitPath": "/",
					"isMailboxSetup": true,
					"isEnrolledIn2Sv": false,
					"isEnforcedIn2Sv": false,
					"includeInGlobalAddressList": true
				}
			],
			"nextPageToken": "Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0="
		}`))
	case "/admin/directory/v1/users?domain=sgnldemos.com&maxResults=1&pageToken=Q0FFUzhRSUJrUHpWQUdiNG04Q1IxMklwSjFJemdWcFhvODg4ZFRjV1RSWTBWK2JoeDNXWUQyQUpLeFkzN1NuQWI4ZDFHUWszMmpESGpxR1I3Um5EcXd4V2REbi9Xc0NJTUYyVXVYR2xZcEgwdUVNRk5ZWCtlVlhzYTRXeXA3MFJ2NUxqT25vM1hCeUMzZ0wvdDRwUXZHa3pnd21QcnVZSm9udFEzMk9zMlcyaEhZOUJ5OGd0UzZmU3BZdHBpeE1uUUtOUWJ6ZlYrTUI0WjVnNFBYVVB4ZjRDZTJVc0pXQ05GSG1FZnYzQkMreU9BRWNYZWRkWCt2U3REMFR0ZjI0SElMY1Z3VHB3SHh3WURzbk84d0N5eTFsNDAwUFNVVGJrNW9BUkFwajJEL3dYcFo4bmRIa3FRdmRGK3Z3b0EwWSt6ZTB6Y3ZkcUlMd3pwVjIzL25GL0tIN2JPcTVqMWFSVVYrRDN0NE4zZzNxaU44clJ5c1dxMWhadkxyT1R2TGJjdWNVdVEwMVgwcHp5cTZlSG5vTUVWeWttMHJUZEFBNGdVOVFtVjUvZStBMWlhVlVrOEQyTVlpMmJJUUtaRnRoUk5jK3lhM1dScyt0ZDFSZWFHd09MVlNnQVNpQWZrQlZYaDJqbFVqblRvL3Y3OUZYWlp1TT0%3D":
		w.Write([]byte(`{
			"kind": "admin#directory#users",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/jKQlemYRE_4Xx1psPHzJ4ldVbfE\"",
			"users": [
				{
					"kind": "admin#directory#user",
					"id": "114211199695002816249",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/0VBDYIsSDmoa_zjvxKsXSmzSk2I\"",
					"primaryEmail": "sor-dev@sgnldemos.com",
					"name": {
						"givenName": "SoR",
						"familyName": "Development",
						"fullName": "SoR Development"
					},
					"isAdmin": false,
					"isDelegatedAdmin": true,
					"lastLoginTime": "2024-04-20T05:56:11.000Z",
					"creationTime": "2024-04-18T00:30:35.000Z",
					"agreedToTerms": true,
					"suspended": false,
					"archived": false,
					"changePasswordAtNextLogin": false,
					"ipWhitelisted": false,
					"emails": [
						{
							"address": "sgnl-demos@sgnl.ai",
							"type": "work"
						},
						{
							"address": "sor-dev@sgnldemos.com",
							"primary": true
						},
						{
							"address": "sor-dev@sgnldemos.com.test-google-a.com"
						}
					],
					"languages": [
						{
							"languageCode": "en",
							"preference": "preferred"
						}
					],
					"nonEditableAliases": [
						"sor-dev@sgnldemos.com.test-google-a.com"
					],
					"customerId": "CUST123456",
					"orgUnitPath": "/serviceAccounts",
					"isMailboxSetup": true,
					"isEnrolledIn2Sv": false,
					"isEnforcedIn2Sv": false,
					"includeInGlobalAddressList": true
				}
			]
		}`))

	// Groups: Total of 3 groups, 1 per request
	// The Groups endpoint is returning a pageToken on the last page of the request.
	// The fourth mock request is to simulate this.
	case "/admin/directory/v1/groups?domain=sgnldemos.com&maxResults=1":
		w.Write([]byte(`{
			"kind": "admin#directory#groups",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/62cbvvHzBBknFDn9M-zq-AUUQVo\"",
			"groups": [
				{
					"kind": "admin#directory#group",
					"id": "01qoc8b13vgdlqb",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/jegUmo9HAPpclSc2UeC_v9oHm2E\"",
					"email": "emptygroup@sgnldemos.com",
					"name": "Empty Group",
					"directMembersCount": "0",
					"description": "",
					"adminCreated": true,
					"nonEditableAliases": [
						"emptygroup@sgnldemos.com.test-google-a.com"
					]
				}
			],
			"nextPageToken": "Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09"
		}`))
	case "/admin/directory/v1/groups?domain=sgnldemos.com&maxResults=1&pageToken=Q2lvd0xDSmxiWEIwZVdkeWIzVndRSE5uYm14a1pXMXZjeTVqYjIwaUxERXdOVGMyTVRnek1EWXhNemxJQTJDOG11U2lCQT09":
		w.Write([]byte(`{
			"kind": "admin#directory#groups",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/6N4IK3UYxu0dY4kKeb-nhAS-UJ0\"",
			"groups": [
				{
					"kind": "admin#directory#group",
					"id": "048pi1tg0qf1f8g",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/n5phxV9NHrbJaK4QhiAp3pDR21k\"",
					"email": "group2@sgnldemos.com",
					"name": "Group2",
					"directMembersCount": "3",
					"description": "",
					"adminCreated": true,
					"nonEditableAliases": [
						"group2@sgnldemos.com.test-google-a.com"
					]
				}
			],
			"nextPageToken": "Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09"
		}`))
	case "/admin/directory/v1/groups?domain=sgnldemos.com&maxResults=1&pageToken=Q2lVd0xDSm5jbTkxY0RKQWMyZHViR1JsYlc5ekxtTnZiU0lzTnpVM05EWTFNVFExTVRnMFNBTmdoNWJnNmY3X19fX19BUT09":
		w.Write([]byte(`{
			"kind": "admin#directory#groups",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/mp27A9nnTqU_uz5BRPRMvhbkMh8\"",
			"groups": [
				{
					"kind": "admin#directory#group",
					"id": "030j0zll41obyxg",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/knD5dVuTcKRMQdMc1Qr7e37h_n4\"",
					"email": "hello@sgnldemos.com",
					"name": "SGNLDemos",
					"directMembersCount": "2",
					"description": "",
					"adminCreated": true,
					"nonEditableAliases": [
						"hello@sgnldemos.com.test-google-a.com"
					]
				}
			],
			"nextPageToken": "Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0="
		}`))
	case "/admin/directory/v1/groups?domain=sgnldemos.com&maxResults=1&pageToken=Q2lNd0xDSm9aV3hzYjBCeloyNXNaR1Z0YjNNdVkyOXRJaXd4TURVek56STVOVFF4TWtnRFlQYVZ5TGY5X19fX193RT0%3D":
		w.Write([]byte(`{
			"kind": "admin#directory#groups",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/HMFwD2wLX237BRmKZQUJYB5ZE7U\""
		}`))

	// Members: Total of 3 requests, 1 per request
	case "/admin/directory/v1/groups/01qoc8b13vgdlqb/members?domain=sgnldemos.com&maxResults=2":
		w.Write([]byte(`{
			"kind": "admin#directory#members",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/HMFwD2wLX237BRmKZQUJYB5ZE7U\""
		}`))
	case "/admin/directory/v1/groups/048pi1tg0qf1f8g/members?domain=sgnldemos.com&maxResults=2":
		w.Write([]byte(`{
			"kind": "admin#directory#members",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/Vk5BaXyT_vbk_bqt6oxXydC4hGU\"",
			"members": [
				{
					"kind": "admin#directory#member",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/1pVaKo8kmZjl-_gkcKlCYvYBPjA\"",
					"id": "USER987654321",
					"email": "user1@sgnldemos.com",
					"role": "MEMBER",
					"type": "USER",
					"status": "ACTIVE"
				},
				{
					"kind": "admin#directory#member",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/4j9E51Yhqt_d3t65g3OEza89nps\"",
					"id": "102475661842232156723",
					"email": "user2@sgnldemos.com",
					"role": "OWNER",
					"type": "USER",
					"status": "ACTIVE"
				}
			],
			"nextPageToken": "CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E"
		}`))
	case "/admin/directory/v1/groups/048pi1tg0qf1f8g/members?domain=sgnldemos.com&maxResults=2&pageToken=CjRJaDhLSFFqSXZNaWZ6Z0VTRW0xaGNtTkFjMmR1YkdSbGJXOXpMbU52YlJnQllKeUppTThFIh8KHQjIvMifzgESEm1hcmNAc2dubGRlbW9zLmNvbRgBYJyJiM8E":
		w.Write([]byte(`{
			"kind": "admin#directory#members",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/ANvWLDdSVCGI5Or1KNF_qrrxIEM\"",
			"members": [
				{
					"kind": "admin#directory#member",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/d566bqA1WOuRhrsagV1MWF5o2f8\"",
					"id": "114211199695002816249",
					"email": "sor-dev@sgnldemos.com",
					"role": "OWNER",
					"type": "USER",
					"status": "ACTIVE"
				}
			]
		}`))
	case "/admin/directory/v1/groups/030j0zll41obyxg/members?domain=sgnldemos.com&maxResults=2":
		w.Write([]byte(`{
			"kind": "admin#directory#members",
			"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/PsQMBqkd8nUPEnlr8Ca0MaOK7fU\"",
			"members": [
				{
					"kind": "admin#directory#member",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/Z1pIME_BX7RFRQE_z0f1Gfewx70\"",
					"id": "048pi1tg0qf1f8g",
					"email": "group2@sgnldemos.com",
					"role": "MEMBER",
					"type": "GROUP",
					"status": "ACTIVE"
				},
				{
					"kind": "admin#directory#member",
					"etag": "\"OgGr5k3K9Y5ATXvRGH2u3OKSYpWBthUtvh8XRXQHGKs/J-EKKLlSn8MJS_aGpigqCNOJvYs\"",
					"id": "102475661842232156723",
					"email": "user2@sgnldemos.com",
					"role": "OWNER",
					"type": "USER",
					"status": "ACTIVE"
				}
			]
		}`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
