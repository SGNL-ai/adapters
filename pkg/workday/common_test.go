// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package workday_test

import (
	"net/http"
)

// Define the endpoints and responses for the mock Workday server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer testtoken" {
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
	// Workers: total of 33 split between 4 requests of page size 10
	case "/api/wql/v1/SGNL/data?limit=5&offset=0&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers":
		w.Write([]byte(`{
			"total": 18,
			"data": [
				{
					"worker": {
						"descriptor": "user1",
						"id": "3aa5550b7fe348b98d7b5741afc65534"
					},
					"email_Work": [],
					"employeeID": "21001",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "4 Vice President",
						"id": "679d4d1ac6da40e19deb7d91e170431d"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00004",
					"jobTitle": "Vice President, Human Resources"
				},
				{
					"worker": {
						"descriptor": "user2",
						"id": "0e44c92412d34b01ace61e80a47aaf6d"
					},
					"email_Work": [
						{
							"descriptor": "user2@workdaySJTest.net",
							"id": "d7fef59db8e21001de457203a69e0001"
						}
					],
					"employeeID": "21002",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "2 Chief Executive Officer",
						"id": "3de1f2834f064394a40a40a727fb6c6d"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Not Declared",
						"id": "a14bf6afa9204ff48a8ea353dd71eb22"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00001",
					"jobTitle": "Chief Executive Officer"
				},
				{
					"worker": {
						"descriptor": "user3",
						"id": "3895af7993ff4c509cbea2e1817172e0"
					},
					"email_Work": [
						{
							"descriptor": "user3@workday.net",
							"id": "d7fef59db8e21001dddaa607a7d30001"
						}
					],
					"employeeID": "21003",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "3 Executive Vice President",
						"id": "0ceb3292987b474bbc40c751a1e22c69"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00002",
					"jobTitle": "Chief Information Officer"
				},
				{
					"worker": {
						"descriptor": "user4",
						"id": "3bf7df19491f4d039fd54decdd84e05c"
					},
					"email_Work": [
						{
							"descriptor": "user4@workday.net",
							"id": "2eab98c6070f4a609adf9ce702bfa9c3"
						}
					],
					"employeeID": "21004",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "3 Executive Vice President",
						"id": "0ceb3292987b474bbc40c751a1e22c69"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00005",
					"jobTitle": "Chief Operating Officer"
				},
				{
					"worker": {
						"descriptor": "user5",
						"id": "26c439a5deed4a7dbab76709e0d2d2ca"
					},
					"email_Work": [
						{
							"descriptor": "user5@workday.net",
							"id": "3aff08c6468b45998638dbbaeaaf4ab8"
						}
					],
					"employeeID": "21005",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00124",
					"jobTitle": "Director, Field Marketing"
				}
			]
		}`))
	case "/api/wql/v1/SGNL/data?limit=5&offset=5&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers":
		w.Write([]byte(`{
			"total": 18,
			"data": [
				{
					"worker": {
						"descriptor": "user6",
						"id": "cc7fb31eecd544e9ae8e03653c63bfab"
					},
					"email_Work": [
						{
							"descriptor": "user6@workday.net",
							"id": "d7fef59db8e21001de09700cef810002"
						}
					],
					"employeeID": "21006",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00011",
					"jobTitle": "Director, Employee Benefits"
				},
				{
					"worker": {
						"descriptor": "user7 (Terminated)",
						"id": "3a37558d68944bf394fad59ff267f4a1"
					},
					"email_Work": [
						{
							"descriptor": "user7@workday.net",
							"id": "4c4aa6815de541bfb24cf6144a0550cc"
						}
					],
					"employeeID": "21007",
					"workerActive": false,
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "0"
				},
				{
					"worker": {
						"descriptor": "user8",
						"id": "3bcc416214054db6911612ef25d51e9f"
					},
					"email_Work": [
						{
							"descriptor": "user8@workday.net",
							"id": "1d53eb9c5247461781f6a415bf94ad49"
						}
					],
					"employeeID": "21008",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Not Declared",
						"id": "a14bf6afa9204ff48a8ea353dd71eb22"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00010",
					"jobTitle": "Director, Payroll Operations"
				},
				{
					"worker": {
						"descriptor": "user9",
						"id": "d66d21e0b1c949b2b1a3decd2fad1375"
					},
					"email_Work": [
						{
							"descriptor": "user9@workday.net",
							"id": "8261477c74b748a2b03482bb9cdb7287"
						}
					],
					"employeeID": "21009",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00009",
					"jobTitle": "Director, Workforce Planning"
				},
				{
					"worker": {
						"descriptor": "user10",
						"id": "50ef79568a9b463a9c5fc431e074125b"
					},
					"email_Work": [
						{
							"descriptor": "user10@workday.net",
							"id": "60355b860cae4f7ea300e51594b8e610"
						}
					],
					"employeeID": "21012",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00508",
					"jobTitle": "Staff Payroll Specialist"
				}
			]
		}`))
	case "/api/wql/v1/SGNL/data?limit=5&offset=10&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers":
		w.Write([]byte(`{
			"total": 18,
			"data": [
				{
					"worker": {
						"descriptor": "user11 (On Leave)",
						"id": "cf9f717959444023b9bc9226a2556661"
					},
					"email_Work": [
						{
							"descriptor": "user11@workday.net",
							"id": "d80ef4c876e04e2fadffca124b944ce4"
						}
					],
					"employeeID": "21010",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00503",
					"jobTitle": "Senior Benefits Analyst"
				},
				{
					"worker": {
						"descriptor": "user12",
						"id": "f21231394b71433c8f75f6fe78264f33"
					},
					"email_Work": [
						{
							"descriptor": "user12@workday.net",
							"id": "68a02b2bff3a48afbfc4bd7c89750ee1"
						}
					],
					"employeeID": "21014",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00515",
					"jobTitle": "Staff Recruiter"
				},
				{
					"worker": {
						"descriptor": "user13",
						"id": "0a46063523fd469f96d4e81ed4d17812"
					},
					"email_Work": [
						{
							"descriptor": "user13@workday.net",
							"id": "6a1ae07ebe754bc19fb624d345fc6a68"
						}
					],
					"employeeID": "21011",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Not Declared",
						"id": "a14bf6afa9204ff48a8ea353dd71eb22"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00509",
					"jobTitle": "Staff Payroll Specialist"
				},
				{
					"worker": {
						"descriptor": "user14",
						"id": "cb625aa152344212970023a793f2c2ac"
					},
					"email_Work": [
						{
							"descriptor": "user14@workday.net",
							"id": "e26999c7731641b8a1c0f678aae7d385"
						}
					],
					"employeeID": "21013",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00512",
					"jobTitle": "Director, Payroll Operations"
				},
				{
					"worker": {
						"descriptor": "user15",
						"id": "2014150640fa42ebbafb6ab936b08073"
					},
					"email_Work": [
						{
							"descriptor": "user15@workday.net",
							"id": "06d86d5ac21343c5ac866179d320d27e"
						}
					],
					"employeeID": "21015",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00517",
					"jobTitle": "Senior Workforce Analyst"
				}
			]
		}`))
	case "/api/wql/v1/SGNL/data?limit=5&offset=15&query=SELECT+FTE%2C+company%2C+email_Work%2C+employeeID%2C+employeeType%2C+gender%2C+hireDate%2C+jobTitle%2C+managementLevel%2C+positionID%2C+worker%2C+workerActive+FROM+allWorkers":
		w.Write([]byte(`{
			"total": 18,
			"data": [
				{
					"worker": {
						"descriptor": "user16",
						"id": "16d87047a76a47b399b4a677058d629f"
					},
					"email_Work": [
						{
							"descriptor": "user16@workday.net",
							"id": "8e8aaddd60814dc693b89a938c192cba"
						}
					],
					"employeeID": "21016",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "8 Individual Contributor",
						"id": "7a379eea3b0c4a10a2b50663b2bd15e4"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services, Inc. (USA)",
						"id": "cb550da820584750aae8f807882fa79a"
					},
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00502",
					"jobTitle": "Senior Benefits Analyst"
				},
				{
					"worker": {
						"descriptor": "user17",
						"id": "1cf028c6f4484c248e8d7d573d7b8845"
					},
					"email_Work": [
						{
							"descriptor": "user17@workday.net",
							"id": "87dd6909b97e4fc6b2d77769ee1503ac"
						}
					],
					"employeeID": "21017",
					"workerActive": true,
					"managementLevel": {
						"descriptor": "5 Director",
						"id": "0b778018b3b44ca3959e498041865645"
					},
					"employeeType": [
						{
							"descriptor": "Regular",
							"id": "9459f5e6f1084433b767c7901ec04416"
						}
					],
					"company": {
						"descriptor": "Global Modern Services S.p.A (Italy)",
						"id": "e4859d59e6094f52a8f2e865cca82cef"
					},
					"gender": {
						"descriptor": "Female",
						"id": "9cce3bec2d0d420283f76f51b928d885"
					},
					"hireDate": "2000-01-01",
					"FTE": "1",
					"positionID": "P-00013",
					"jobTitle": "Director, Accounting"
				},
				{
					"worker": {
						"descriptor": "user18 (Terminated)",
						"id": "f2c673e5b73245889be3581d53187731"
					},
					"email_Work": [
						{
							"descriptor": "user18@workday.net",
							"id": "db497b4fba714fd6b457cc56b821c604"
						}
					],
					"employeeID": "21018",
					"workerActive": false,
					"gender": {
						"descriptor": "Male",
						"id": "d3afbf8074e549ffb070962128e1105a"
					},
					"hireDate": "2000-01-01",
					"FTE": "0"
				}
			]
		}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
