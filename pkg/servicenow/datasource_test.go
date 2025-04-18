// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package servicenow_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/auth"
	"github.com/sgnl-ai/adapters/pkg/servicenow"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

var (
	userIDs = []string{
		"9a826bf03710200044e0bfc8bcbe5dd1",
		"a2826bf03710200044e0bfc8bcbe5ddb",
		"aa826bf03710200044e0bfc8bcbe5ddf",
		"cf1ec0b4530360100999ddeeff7b129f",
	}

	groupIDs = []string{
		"c38f00f4530360100999ddeexxgroup1",
		"c38f00f4530360100999ddeexxgroup2",
		"c38f00f4530360100999ddeexxgroup3",
	}

	groupMemberIDs = []string{
		"002750f8530360100999ddeeff7b1206",
		"0c2750f8530360100999ddeeff7b1208",
		"195ebb573b331300ad3cc9bb34efc4ad",
		"1e7836bac0a8018b439ffb87b1fcc39d",
		"1e7836f0c0a8018b439ffb87529a3c00",
	}

	caseIDs = []string{
		"f6911038530360100999ddeeff7case1",
		"f6911038530360100999ddeeff7case2",
		"f6911038530360100999ddeeff7case3",
		"f6911038530360100999ddeeff7case4",
	}

	accountIDs = []string{
		"f6911038530360100999ddeeff77acc1",
		"f6911038530360100999ddeeff77acc2",
		"f6911038530360100999ddeeff77acc3",
		"f6911038530360100999ddeeff77acc4",
	}
)

// Define the endpoints and responses for the mock Servicenow server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer testtoken" && r.Header.Get("Authorization") != auth.BasicAuthHeader("username", "password") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
				"errorCode": "E0000005",
				"errorSummary": "Invalid session",
				"errorLink": "E0000005",
				"errorId": "...",
				"errorCauses": []
			}`))
	}

	switch r.URL.Path {
	case "/api/now/v2/table/sys_user":
		// Edge case to test error response.
		switch r.URL.Query().Get("sysparm_limit") {
		case "999":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"detail": "The requested query is too long to build the response pagination header URLs. Please do one of the following: shorten the sysparm_query, or query without pagination by setting the parameter 'sysparm_suppress_pagination_header' to true, or set 'sysparm_limit' with a value larger then 110 to bypass the need for pagination.",
					"message": "Pagination not supported"
				},
				"status": "failure"
			}`),
			)

			return
		}

		switch r.URL.Query().Get("sysparm_offset") {
		case "":
			w.Header().Set("Link", `<https://localhost/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3>;rel="next"`)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + userIDs[0] + `",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"email": "freeman.soula@example.com",
							"manager": "` + userIDs[1] + `"
						},
						{
							"sys_id": "` + userIDs[1] + `",
							"sys_created_on": "2012-02-18 03:04:52",
							"active": "true",
							"email": "junior.wadlinger@example.com",
							"manager": "` + userIDs[2] + `"
						},
						{
							"sys_id": "` + userIDs[2] + `",
							"sys_created_on": "2012-02-18 03:04:52",
							"active": "true",
							"email": "curt.menedez@example.com",
							"manager": ""
						}
					]
				}`))
		case "3":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + userIDs[3] + `",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"email": "john.doe@example.com",
							"manager": "` + userIDs[1] + `"
						}
					]
				}`))

		// Edge case for testing invalid datetime format, so there is no link from the previous page
		case "4":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + userIDs[3] + `",
							"sys_created_on": "2021/01/01 00:00:00.000Z",
							"active": "true",
							"email": "john.doe@example.com",
							"manager": "` + userIDs[1] + `"
						}
					]
				}`))
		}

	case "/api/now/v2/table/sys_user_group":
		switch r.URL.Query().Get("sysparm_offset") {

		case "":
			w.Header().Set("Link", `<https://localhost/api/now/v2/table/sys_group?sysparm_fields=sys_id,parent,manager,default_assignee,sys_created_on,active,description&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3>;rel="next"`)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + groupIDs[0] + `",
							"parent": "` + groupIDs[2] + `",
							"manager": "Tom",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"defaultAssignee": "Tom",
							"description": "Development"
						},
						{
							"sys_id": "` + groupIDs[1] + `",
							"parent": "` + groupIDs[2] + `",
							"manager": "Jim",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"defaultAssignee": "Jim",
							"description": "Product Marketing"
						},
						{
							"sys_id": "` + groupIDs[2] + `",
							"parent": "",
							"manager": "Ben",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"defaultAssignee": "Ben",
							"description": "Management"
						}
					]
				}`))
		}

	case "/api/now/v2/table/sys_user_grmember":
		switch r.URL.Query().Get("sysparm_offset") {

		case "":
			w.Header().Set("Link", `https://localhost/api/now/v2/table/sys_user_grmember?sysparm_fields=sys_id,user,group&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=5>;rel="next"`)
			w.WriteHeader(http.StatusOK)

			/*
				| User        | Group        |
				| ----------- | ------------ |
				| userIDs[0]  | groupIDs[0]  |
				| userIDs[0]  | groupIDs[1]  |
				| userIDs[1]  | groupIDs[0]  |
				| userIDs[2]  | groupIDs[0]  |
				| userIDs[3]  | groupIDs[1]  |
			*/
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + groupMemberIDs[0] + `",
							"sys_created_on": "2021-03-19 16:05:36",
							"user": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
								"value": "` + userIDs[0] + `"
							},
							"group": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user_group/` + groupIDs[0] + `",
								"value": "` + groupIDs[0] + `"
							}
						},
						{
							"sys_id": "` + groupMemberIDs[1] + `",
							"sys_created_on": "2021-03-19 16:05:36",
							"user": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
								"value": "` + userIDs[0] + `"
							},
							"group": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user_group/` + groupIDs[1] + `",
								"value": "` + groupIDs[1] + `"
							}
						},
						{
							"sys_id": "` + groupMemberIDs[2] + `",
							"sys_created_on": "2021-03-19 16:05:36",
							"user": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user/` + userIDs[1] + `",
								"value": "` + userIDs[1] + `"
							},
							"group": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user_group/` + groupIDs[0] + `",
								"value": "` + groupIDs[0] + `"
							}
						},
						{
							"sys_id": "` + groupMemberIDs[3] + `",
							"sys_created_on": "2021-03-19 16:05:36",
							"user": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user/` + userIDs[2] + `",
								"value": "` + userIDs[2] + `"
							},
							"group": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user_group/` + groupIDs[0] + `",
								"value": "` + groupIDs[0] + `"
							}
						},
						{
							"sys_id": "` + groupMemberIDs[4] + `",
							"sys_created_on": "2021-03-19 16:05:36",
							"user": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user/` + userIDs[3] + `",
								"value": "` + userIDs[3] + `"
							},
							"group": {
								"link": "https://dev122280.service-now.com/api/now/v2/table/sys_user_group/` + groupIDs[1] + `",
								"value": "` + groupIDs[1] + `"
							}
						}
					]
				}`))
		}

	case "/api/now/v2/table/sn_customerservice_case":
		switch r.URL.Query().Get("sysparm_query") {
		case "ORDERBYsys_id":
			w.Header().Set("Link", `https://localhost/api/now/v2/table/sn_customerservice_case?sysparm_fields=sys_id,case,parent,assigned_to,account,description,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=4>;rel="next"`)
			w.WriteHeader(http.StatusOK)

			/*
				| Case        | Assigned To | Account       | Parent Case |
				| ----------- | ----------- | ------------- | ----------- |
				| caseIDs[0]  | userIDs[0]  | accountIDs[0] | caseIDs[1]  |
				| caseIDs[1]  | userIDs[0]  | accountIDs[0] | caseIDs[2]  |
				| caseIDs[2]  | userIDs[0]  | accountIDs[0] | ""          |
				| caseIDs[3]  | userIDs[1]  | accountIDs[1] | ""          |
			*/
			w.Write([]byte(`{
				"result": [
					{
						"sys_id": "` + caseIDs[0] + `",
						"case": "` + caseIDs[0] + `",
						"parent": "` + caseIDs[1] + `",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
							"value": "` + userIDs[0] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[0] + `",
							"value": "` + accountIDs[0] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					},
					{
						"sys_id": "` + caseIDs[1] + `",
						"case": "` + caseIDs[1] + `",
						"parent": "` + caseIDs[2] + `",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
							"value": "` + userIDs[0] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[0] + `",
							"value": "` + accountIDs[0] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					},
					{
						"sys_id": "` + caseIDs[2] + `",
						"case": "` + caseIDs[2] + `",
						"parent": "",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
							"value": "` + userIDs[0] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[0] + `",
							"value": "` + accountIDs[0] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					},
					{
						"sys_id": "` + caseIDs[3] + `",
						"case": "` + caseIDs[3] + `",
						"parent": "",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[1] + `",
							"value": "` + userIDs[1] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[1] + `",
							"value": "` + accountIDs[1] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					}
				]
			}`))

		case "active=true^ORDERBYsys_id":
			w.Header().Set("Link", `https://localhost/api/now/v2/table/sn_customerservice_case?sysparm_fields=sys_id,case,parent,assigned_to,account,description,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=active%3Dtrue%5EORDERBYsys_id&sysparm_offset=4>;rel="next"`)
			w.WriteHeader(http.StatusOK)

			/*
				| Case        | Assigned To | Account       | Parent Case |
				| ----------- | ----------- | ------------- | ----------- |
				| caseIDs[0]  | userIDs[0]  | accountIDs[0] | caseIDs[1]  |
				| caseIDs[3]  | userIDs[1]  | accountIDs[1] | ""          |
			*/
			w.Write([]byte(`{
				"result": [
					{
						"sys_id": "` + caseIDs[0] + `",
						"case": "` + caseIDs[0] + `",
						"parent": "` + caseIDs[1] + `",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[0] + `",
							"value": "` + userIDs[0] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[0] + `",
							"value": "` + accountIDs[0] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					},
					{
						"sys_id": "` + caseIDs[3] + `",
						"case": "` + caseIDs[3] + `",
						"parent": "",
						"assigned_to": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/sys_user/` + userIDs[1] + `",
							"value": "` + userIDs[1] + `"
						},
						"account": {
							"link": "https://test-instance.service-now.com/api/now/v2/table/customer_account/` + accountIDs[1] + `",
							"value": "` + accountIDs[1] + `"
						},
						"description": "",
						"sys_created_on": "2022-05-26 21:59:03",
						"active": true
					}
				]
			}`))
		}

	case "/api/now/v2/table/customer_account":
		switch r.URL.Query().Get("sysparm_offset") {

		case "":
			w.Header().Set("Link", `https://localhost/api/now/v2/table/customer_account?sysparm_fields=sys_id,number,account_parent,parent,sys_created_on,primary&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=4>;rel="next"`)
			w.WriteHeader(http.StatusOK)

			/*
				| Account       | Parent Account | Parent     |
				| ------------- | -------------- | ---------- |
				| accountIDs[0] | accountIDs[1]  | userIDs[0] |
				| accountIDs[1] | accountIDs[2]  | userIDs[0] |
				| accountIDs[2] | ""             | userIDs[0] |
				| accountIDs[3] | ""             | userIDs[1] |
			*/
			w.Write([]byte(`{
				"result": [
					{
						"parent": "` + userIDs[0] + `",
						"number": "ACCT0010001",
						"sys_created_on": "2022-10-20 20:11:49",
						"primary": "true",
						"account_parent": "` + accountIDs[1] + `",
						"sys_id": "` + accountIDs[0] + `"
					},
					{
						"parent": "` + userIDs[0] + `",
						"number": "ACCT0010002",
						"sys_created_on": "2022-10-20 20:11:49",
						"primary": "true",
						"account_parent": "` + accountIDs[2] + `",
						"sys_id": "` + accountIDs[1] + `"
					},
					{
						"parent": "` + userIDs[0] + `",
						"number": "ACCT0010003",
						"sys_created_on": "2022-10-20 20:11:49",
						"primary": "true",
						"account_parent": "",
						"sys_id": "` + accountIDs[2] + `"
					},
					{
						"parent": "` + userIDs[1] + `",
						"number": "ACCT0010004",
						"sys_created_on": "2022-10-20 20:11:49",
						"primary": "true",
						"account_parent": "",
						"sys_id": "` + accountIDs[3] + `"
					}
				]
			}`))
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

var TestServerAdvancedFiltersHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer testtoken" && r.Header.Get("Authorization") != auth.BasicAuthHeader("username", "password") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
				"errorCode": "E0000005",
				"errorSummary": "Invalid session",
				"errorLink": "E0000005",
				"errorId": "...",
				"errorCauses": []
			}`))
	}

	sysparmQuery := r.URL.Query().Get("sysparm_query")
	sysparmLimit := r.URL.Query().Get("sysparm_limit")
	sysarmOffset := r.URL.Query().Get("sysparm_offset")

	switch r.URL.Path {
	case "/api/now/v2/table/sys_user_group":
		// Handlers for pagination testing.
		if sysparmLimit == "1" {
			switch sysparmQuery {
			case "sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2^ORDERBYsys_id":
				switch sysarmOffset {
				case "":
					w.Header().Set("Link", `https://127.0.0.1:8443/api/now/v2/table/sys_user_group?sysparm_fields=sys_id,default_assignee&sysparm_exclude_reference_link=true&sysparm_limit=1&sysparm_query=sys_idINc38f00f4530360100999ddeexxgroup1,c38f00f4530360100999ddeexxgroup2^ORDERBYsys_id&sysparm_offset=1>;rel="next"`)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
							"result": [
								{
									"sys_id": "` + groupIDs[0] + `",
									"defaultAssignee": "Tom"
								}
							]
						}`,
					))

					return
				case "1":
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
							"result": [
								{
									"sys_id": "` + groupIDs[1] + `",
									"defaultAssignee": "John"
								}
							]
						}`,
					))

					return
				}
			}
		}

		switch sysparmQuery {
		case "active=true^ORDERBYsys_id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + groupIDs[0] + `",
							"manager": "Tom",
							"sys_created_on": "2012-02-18 03:04:51",
							"active": "true",
							"defaultAssignee": "Tom",
							"description": "Development"
						}
					]
				}`,
			))
		}

	case "/api/now/v2/table/sys_user_grmember":
		// Handlers for pagination testing.
		if sysparmLimit == "1" {
			switch sysparmQuery {
			case "groupINc38f00f4530360100999ddeexxgroup1^user.active=true^ORDERBYsys_id":
				switch sysarmOffset {
				case "":
					w.Header().Set("Link", `https://127.0.0.1:8443/api/now/v2/table/sys_user_grmember?sysparm_fields=sys_id,user.sys_id&sysparm_exclude_reference_link=true&sysparm_limit=1&sysparm_query=groupINc38f00f4530360100999ddeexxgroup1^user.active=true^ORDERBYsys_id&sysparm_offset=1>;rel="next"`)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
								"result": [
									{
										"sys_id": "` + userIDs[0] + `",
										"user.active": "true"
									}
								]
							}`,
					))

					return
				case "1":
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
								"result": [
									{
										"sys_id": "` + userIDs[1] + `",
										"user.active": "true"
									}
								]
							}`,
					))

					return
				}
			case "groupINc38f00f4530360100999ddeexxgroup2^user.active=true^ORDERBYsys_id":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(
					`{
						"result": [
							{
								"sys_id": "` + userIDs[2] + `",
								"user.active": "true"
							}
						]
					}`,
				))

				return
			}
		}

		switch sysparmQuery {
		case "groupIN" + groupIDs[0] + "^user.active=true^ORDERBYsys_id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + groupMemberIDs[0] + `",
							"user.sys_id": "` + userIDs[0] + `",
							"user.email": "user2@example.com"
						}
					]
				}`,
			))
		}

	case "/api/now/v2/table/sn_customerservice_case":
		// Handlers for pagination testing.
		if sysparmLimit == "1" {
			switch sysparmQuery {
			case "assignment_groupINc38f00f4530360100999ddeexxgroup1^ORDERBYsys_id":
				switch sysarmOffset {
				case "":
					w.Header().Set("Link", `https://127.0.0.1:8443/api/now/v2/table/sn_customerservice_case?sysparm_fields=sys_id,assignment_group&sysparm_exclude_reference_link=true&sysparm_limit=1&sysparm_query=assignment_groupINc38f00f4530360100999ddeexxgroup1^ORDERBYsys_id&sysparm_offset=1>;rel="next"`)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
							"result": [
								{
									"sys_id": "` + caseIDs[0] + `",
									"assignment_group": "` + groupIDs[0] + `"
								}
							]
						}`,
					))

					return
				case "1":
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
							"result": [
								{
									"sys_id": "` + caseIDs[1] + `",
									"assignment_group": "` + groupIDs[0] + `"
								}
							]
						}`,
					))

					return
				}
			case "assignment_groupINc38f00f4530360100999ddeexxgroup2^ORDERBYsys_id":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
						"result": [
							{
								"sys_id": "` + caseIDs[2] + `",
								"assignment_group": "` + groupIDs[1] + `"
							}
						]
					}`,
				))

				return
			}
		}

		switch sysparmQuery {
		case "assigned_toIN" + userIDs[0] + "^ORDERBYsys_id":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
					"result": [
						{
							"sys_id": "` + caseIDs[0] + `",
							"assigned_to": "` + userIDs[0] + `"
						}
					]
				}`,
			))
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body             []byte
		entityExternalID string
		cursor           string
		wantObjects      []map[string]interface{}
		wantNextCursor   *string
		wantErr          *framework.Error
	}{
		"single_page": {
			body:             []byte(`{"result": [{"sys_id": "9a826bf03710200044e0bfc8bcbe5dd1"}, {"sys_id": "a2826bf03710200044e0bfc8bcbe5ddb"}]}`),
			entityExternalID: "sys_user",
			wantObjects: []map[string]interface{}{
				{"sys_id": "9a826bf03710200044e0bfc8bcbe5dd1"},
				{"sys_id": "a2826bf03710200044e0bfc8bcbe5ddb"},
			},
		},
		"result_invalid_type": {
			body:             []byte(`{"result": {"Id": "500Hu000020yLuHIAU"}}`),
			entityExternalID: "sys_user",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal object into Go struct field DatasourceResponse.result of type []map[string]interface {}.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			body:             []byte(`{"result": ["500Hu000020yLuHIAU", "500Hu000020yLuMIAU"]}`),
			entityExternalID: "sys_user",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: json: cannot unmarshal string into Go struct field DatasourceResponse.result of type map[string]interface {}.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_object_structure": {
			body:             []byte(`{"result": [{"500Hu000020yLuHIAU"}, {"500Hu000020yLuMIAU"}]}`),
			entityExternalID: "sys_user",
			wantNextCursor:   nil,
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: invalid character '}' after object key.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotErr := servicenow.ParseResponse(tt.body)

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

	servicenowClient := servicenow.NewClient(client)
	server := httptest.NewServer(TestServerHandler)

	tests := map[string]struct {
		context context.Context
		request *servicenow.Request
		wantRes *servicenow.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &servicenow.Request{
				RequestTimeoutSeconds: 5,
				AuthorizationHeader:   "Bearer testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "sys_user",
				PageSize:              200,
				APIVersion:            "v2",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "email",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "sys_created_on",
						Type:       framework.AttributeTypeDateTime,
					},
				},
			},
			wantRes: &servicenow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"sys_id":         "9a826bf03710200044e0bfc8bcbe5dd1",
						"sys_created_on": "2012-02-18 03:04:51",
						"active":         "true",
						"email":          "freeman.soula@example.com",
						"manager":        "a2826bf03710200044e0bfc8bcbe5ddb",
					},
					{
						"sys_id":         "a2826bf03710200044e0bfc8bcbe5ddb",
						"sys_created_on": "2012-02-18 03:04:52",
						"active":         "true",
						"email":          "junior.wadlinger@example.com",
						"manager":        "aa826bf03710200044e0bfc8bcbe5ddf",
					},
					{
						"sys_id":         "aa826bf03710200044e0bfc8bcbe5ddf",
						"sys_created_on": "2012-02-18 03:04:52",
						"active":         "true",
						"email":          "curt.menedez@example.com",
						"manager":        "",
					},
				},
				NextCursor: testutil.GenPtr("https://localhost/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3"),
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &servicenow.Request{
				RequestTimeoutSeconds: 5,
				AuthorizationHeader:   "Bearer testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "sys_user",
				PageSize:              200,
				APIVersion:            "v2",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "email",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "sys_created_on",
						Type:       framework.AttributeTypeDateTime,
					},
				},
				Cursor: testutil.GenPtr(server.URL + "/api/now/v2/table/sys_user?sysparm_fields=sys_id,manager,email,sys_created_on,active&sysparm_exclude_reference_link=true&sysparm_limit=0&sysparm_query=ORDERBYsys_id&sysparm_offset=3"),
			},
			wantRes: &servicenow.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"sys_id":         "cf1ec0b4530360100999ddeeff7b129f",
						"sys_created_on": "2012-02-18 03:04:51",
						"active":         "true",
						"email":          "john.doe@example.com",
						"manager":        "a2826bf03710200044e0bfc8bcbe5ddb",
					},
				},
			},
			wantErr: nil,
		},
		"invalid_auth": {
			context: context.Background(),
			request: &servicenow.Request{
				RequestTimeoutSeconds: 5,
				AuthorizationHeader:   "testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "sys_user",
				PageSize:              200,
				APIVersion:            "v2",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "email",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "sys_created_on",
						Type:       framework.AttributeTypeDateTime,
					},
				},
			},
			wantRes: &servicenow.Response{
				StatusCode: http.StatusUnauthorized,
			},
			wantErr: nil,
		},
		"datasource_returned_400": {
			context: context.Background(),
			request: &servicenow.Request{
				RequestTimeoutSeconds: 5,
				AuthorizationHeader:   "Bearer testtoken",
				BaseURL:               server.URL,
				EntityExternalID:      "sys_user",
				PageSize:              999, // Dummy page size to trigger the error response.
				APIVersion:            "v2",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "sys_id",
						Type:       framework.AttributeTypeString,
					},
				},
			},
			wantRes: nil,
			wantErr: &framework.Error{
				Message: "Failed to get page from datasource: 400. Message: `Pagination not supported`. Details: `The requested query is too long to build the response pagination header URLs. Please do one of the following: shorten the sysparm_query, or query without pagination by setting the parameter 'sysparm_suppress_pagination_header' to true, or set 'sysparm_limit' with a value larger then 110 to bypass the need for pagination.`.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := servicenowClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
