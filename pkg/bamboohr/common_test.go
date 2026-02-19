// Copyright 2026 SGNL.ai, Inc.

// nolint: lll, goconst
package bamboohr_test

import (
	"net/http"
)

// Define the endpoints and responses for the mock BambooHR server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Basic YXBpS2V5MTIzOnJhbmRvbVN0cmluZw==" {
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
	// Employees (no pagination support)
	case "/sgnltestdev/v1/reports/custom?format=JSON&onlyCurrent=true":
		w.Write([]byte(`{
			"title": "Report",
			"fields": [
				{
					"id": "id",
					"type": "int",
					"name": "EEID"
				},
				{
					"id": "bestEmail",
					"type": "email",
					"name": "Email"
				},
				{
					"id": "dateOfBirth",
					"type": "date",
					"name": "Birth Date"
				},
				{
					"id": "fullName1",
					"type": "text",
					"name": "First Name Last Name"
				},
				{
					"id": "isPhotoUploaded",
					"type": "bool",
					"name": "Is employee photo uploaded"
				},
				{
					"id": "customcustomBoolField",
					"type": "checkbox",
					"name": "customBoolField"
				},
				{
					"id": "supervisorEId",
					"type": "text",
					"name": "Supervisor EID"
				},
				{
					"id": "supervisorEmail",
					"type": "email",
					"name": "Manager's email"
				},
				{
					"id": "lastChanged",
					"type": "timestamp",
					"name": "Last changed"
				}
			],
			"employees": [
				{
					"id": "4",
					"bestEmail": "cabbott@efficientoffice.com",
					"dateOfBirth": "1996-09-02",
					"fullName1": "Charlotte Abbott",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "9",
					"supervisorEmail": "jcaldwell@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "5",
					"bestEmail": "aadams@efficientoffice.com",
					"dateOfBirth": "1983-06-30",
					"fullName1": "Ashley Adams",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "9",
					"supervisorEmail": "jcaldwell@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "6",
					"bestEmail": "cagluinda@efficientoffice.com",
					"dateOfBirth": "1996-08-27",
					"fullName1": "Christina Agluinda",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "9",
					"supervisorEmail": "jcaldwell@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "7",
					"bestEmail": "sanderson@efficientoffice.com",
					"dateOfBirth": "0000-00-00",
					"fullName1": "Shannon Anderson",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "9",
					"supervisorEmail": "jcaldwell@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "8",
					"bestEmail": "arvind@sgnl.ai",
					"dateOfBirth": "0000-00-00",
					"fullName1": "Arvind",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": null,
					"supervisorEmail": null,
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "9",
					"bestEmail": "jcaldwell@efficientoffice.com",
					"dateOfBirth": "1975-01-26",
					"fullName1": "Jennifer Caldwell",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "0000-00-00T00:00:00+00:00"
				},
				{
					"id": "10",
					"bestEmail": "rsaito@efficientoffice.com",
					"dateOfBirth": "1968-12-28",
					"fullName1": "Ryota Saito",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "0000-00-00T00:00:00"
				},
				{
					"id": "11",
					"bestEmail": "dvance@efficientoffice.com",
					"dateOfBirth": "1978-08-23",
					"fullName1": "Daniel Vance",
					"isPhotoUploaded": "",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "12",
					"bestEmail": "easture@efficientoffice.com",
					"dateOfBirth": "1990-07-01",
					"fullName1": "Eric Asture",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "null",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "13",
					"bestEmail": "cbarnet@efficientoffice.com",
					"dateOfBirth": "1987-06-16",
					"fullName1": "Cheryl Barnet",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "Null",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:50+00:00"
				},
				{
					"id": "14",
					"bestEmail": "mandev@efficientoffice.com",
					"dateOfBirth": "1987-06-05",
					"fullName1": "Maja Andev",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "15",
					"bestEmail": "twalsh@efficientoffice.com",
					"dateOfBirth": "1981-03-18",
					"fullName1": "Trent Walsh",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "16",
					"bestEmail": "jbryan@efficientoffice.com",
					"dateOfBirth": "1970-12-07",
					"fullName1": "Jake Bryan",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "17",
					"bestEmail": "dchou@efficientoffice.com",
					"dateOfBirth": "1987-05-08",
					"fullName1": "Dorothy Chou",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "8",
					"supervisorEmail": "arvind@sgnl.ai",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "18",
					"bestEmail": "javier@efficientoffice.com",
					"dateOfBirth": "1996-08-28",
					"fullName1": "Javier Cruz",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "15",
					"supervisorEmail": "twalsh@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "19",
					"bestEmail": "shelly@efficientoffice.com",
					"dateOfBirth": "1993-06-01",
					"fullName1": "Shelly Cluff",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "10",
					"supervisorEmail": "rsaito@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "22",
					"bestEmail": "dillon@efficientoffice.com",
					"dateOfBirth": "1972-06-06",
					"fullName1": "Dillon (Remote) Park",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "105",
					"supervisorEmail": "norma@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "23",
					"bestEmail": "darlene@efficientoffice.com",
					"dateOfBirth": "1975-09-16",
					"fullName1": "Darlene Handley",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "15",
					"supervisorEmail": "twalsh@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "24",
					"bestEmail": "zack@efficientoffice.com",
					"dateOfBirth": "2000-08-02",
					"fullName1": "Zack Miller",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "11",
					"supervisorEmail": "dvance@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "25",
					"bestEmail": "philip@efficientoffice.com",
					"dateOfBirth": "1975-11-26",
					"fullName1": "Philip Wagener",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "12",
					"supervisorEmail": "easture@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "26",
					"bestEmail": "agranger@efficientoffice.com",
					"dateOfBirth": "1998-11-26",
					"fullName1": "Amy Granger",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "13",
					"supervisorEmail": "cbarnet@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "27",
					"bestEmail": "debra@efficientoffice.com",
					"dateOfBirth": "1966-10-18",
					"fullName1": "Debra Tuescher",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "49",
					"supervisorEmail": "robert@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "28",
					"bestEmail": "andy@efficientoffice.com",
					"dateOfBirth": "2001-02-25",
					"fullName1": "Andy Graves",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "26",
					"supervisorEmail": "agranger@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "29",
					"bestEmail": "catherine@efficientoffice.com",
					"dateOfBirth": "1993-12-18",
					"fullName1": "Catherine Jones",
					"isPhotoUploaded": "false",
					"customcustomBoolField": "0",
					"supervisorEId": "4",
					"supervisorEmail": "cabbott@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "30",
					"bestEmail": "corey@efficientoffice.com",
					"dateOfBirth": "1995-05-01",
					"fullName1": "Corey Ross",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "15",
					"supervisorEmail": "twalsh@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "31",
					"bestEmail": "sally@efficientoffice.com",
					"dateOfBirth": "1984-05-31",
					"fullName1": "Sally Harmon",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "16",
					"supervisorEmail": "jbryan@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "32",
					"bestEmail": "carly@efficientoffice.com",
					"dateOfBirth": "1984-07-02",
					"fullName1": "Carly Seymour",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "52",
					"supervisorEmail": "nate@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "33",
					"bestEmail": "erin@efficientoffice.com",
					"dateOfBirth": "1993-02-26",
					"fullName1": "Erin Farr",
					"isPhotoUploaded": "false",
					"customcustomBoolField": "0",
					"supervisorEId": "16",
					"supervisorEmail": "jbryan@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "34",
					"bestEmail": "emily@efficientoffice.com",
					"dateOfBirth": "1994-04-30",
					"fullName1": "Emily Gomez",
					"isPhotoUploaded": "false",
					"customcustomBoolField": "0",
					"supervisorEId": "36",
					"supervisorEmail": "melany@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "35",
					"bestEmail": "aaron@efficientoffice.com",
					"dateOfBirth": "1998-08-16",
					"fullName1": "Aaron Eckerly",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "36",
					"supervisorEmail": "melany@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:48+00:00"
				},
				{
					"id": "36",
					"bestEmail": "melany@efficientoffice.com",
					"dateOfBirth": "1986-11-25",
					"fullName1": "Melany Olsen",
					"isPhotoUploaded": "false",
					"customcustomBoolField": "0",
					"supervisorEId": "13",
					"supervisorEmail": "cbarnet@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "37",
					"bestEmail": "whitney@efficientoffice.com",
					"dateOfBirth": "1992-12-30",
					"fullName1": "Whitney Webster",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "5",
					"supervisorEmail": "aadams@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "38",
					"bestEmail": "marrissa@efficientoffice.com",
					"dateOfBirth": "1995-01-31",
					"fullName1": "Marrissa Mellon",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "69",
					"supervisorEmail": "karin@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "39",
					"bestEmail": "paige@efficientoffice.com",
					"dateOfBirth": "1993-02-02",
					"fullName1": "Paige Rasmussen",
					"isPhotoUploaded": "true",
					"customcustomBoolField": "0",
					"supervisorEId": "57",
					"supervisorEmail": "liam@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				},
				{
					"id": "40",
					"bestEmail": "kelli@efficientoffice.com",
					"dateOfBirth": "1988-03-01",
					"fullName1": "Kelli Crandle",
					"isPhotoUploaded": "false",
					"customcustomBoolField": "0",
					"supervisorEId": "49",
					"supervisorEmail": "robert@efficientoffice.com",
					"lastChanged": "2024-04-12T19:33:47+00:00"
				}
			]
		}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})
