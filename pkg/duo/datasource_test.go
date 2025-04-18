// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package duo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapters/pkg/duo"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock Duo server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.RequestURI() {
	// Users Page 1
	case "/admin/v1/users?limit=4&offset=0":
		w.Write([]byte(`{
			"metadata": {
			  "next_offset": 4,
			  "total_objects": 10
			},
			"response": [
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user1@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 1",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUYC8O4O953VBGGKLHAL",
				"username": "user1",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user2@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  },
				  {
					"desc": "random group 2",
					"group_id": "DGIB125DJLJKYZ9W257F",
					"mobile_otp_enabled": false,
					"name": "group2",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 2",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUHUTX7KGB6D15WTD3VY",
				"username": "user2",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user3@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 3",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUEL2SL4CWLP04CE71SL",
				"username": "user3",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user4@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  },
				  {
					"desc": "random group 2",
					"group_id": "DGIB125DJLJKYZ9W257F",
					"mobile_otp_enabled": false,
					"name": "group2",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 4",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUB3BH17CE2V7B744RLI",
				"username": "user4",
				"webauthncredentials": []
			  }
			],
			"stat": "OK"
		  }`))

	// Users Page 2
	case "/admin/v1/users?limit=4&offset=4":
		w.Write([]byte(`{
			"metadata": {
			  "next_offset": 8,
			  "prev_offset": 0,
			  "total_objects": 10
			},
			"response": [
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user5@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 5",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUWC7NXJX7IM9I7J26AT",
				"username": "user5",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041056,
				"desktoptokens": [],
				"email": "user6@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": false,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [],
				"realname": "Test User 6",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUQ87KL4A6OU5VYMWWLT",
				"username": "user6",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041057,
				"desktoptokens": [],
				"email": "user7@example.com",
				"firstname": null,
				"groups": [
				  {
					"desc": "random group 1",
					"group_id": "DGKUKQSTG7ZFDN2N1XID",
					"mobile_otp_enabled": false,
					"name": "group1",
					"push_enabled": false,
					"sms_enabled": false,
					"status": "Active"
				  }
				],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 7",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DU9SRK429IRM2J7OEDCP",
				"username": "user7",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041057,
				"desktoptokens": [],
				"email": "user8@example.com",
				"firstname": null,
				"groups": [],
				"is_enrolled": false,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [],
				"realname": "Test User 8",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DU1E2CSOI6I5HEO043WN",
				"username": "user8",
				"webauthncredentials": []
			  }
			],
			"stat": "OK"
		  }`))

	// Users Page 3
	case "/admin/v1/users?limit=4&offset=8":
		w.Write([]byte(`{
			"metadata": {
			  "prev_offset": 4,
			  "total_objects": 10
			},
			"response": [
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041057,
				"desktoptokens": [],
				"email": "user9@example.com",
				"firstname": null,
				"groups": [],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  },
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPX0H7ZWQLSB735FEHVY",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 9",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DU2T7B5VIC0RSCN1A13W",
				"username": "user9",
				"webauthncredentials": []
			  },
			  {
				"alias1": null,
				"alias2": null,
				"alias3": null,
				"alias4": null,
				"aliases": {},
				"created": 1706041057,
				"desktoptokens": [],
				"email": "user10@example.com",
				"firstname": null,
				"groups": [],
				"is_enrolled": true,
				"last_directory_sync": null,
				"last_login": null,
				"lastname": null,
				"lockout_reason": null,
				"notes": "",
				"phones": [
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DPFL36P8Z8LZANN1FFEZ",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  },
				  {
					"activated": false,
					"capabilities": [
					  "sms"
					],
					"encrypted": "",
					"extension": "",
					"fingerprint": "",
					"last_seen": "",
					"model": "Unknown",
					"name": "",
					"number": "+11111111111",
					"phone_id": "DP7MW6K4G1OVMP8DTI08",
					"platform": "Generic Smartphone",
					"screenlock": "",
					"sms_passcodes_sent": false,
					"tampered": "",
					"type": "Mobile"
				  }
				],
				"realname": "Test User 10",
				"status": "active",
				"tokens": [],
				"u2ftokens": [],
				"user_id": "DUG1B8MRABMVKYVCFO8H",
				"username": "user10",
				"webauthncredentials": []
			  }
			],
			"stat": "OK"
		  }`))

	// Groups Page 1:
	case "/admin/v1/groups?limit=3&offset=0":
		w.Write([]byte(`{
			"metadata": {
			  "next_offset": 3,
			  "total_objects": 5
			},
			"response": [
			  {
				"desc": "random group 1",
				"group_id": "DGKUKQSTG7ZFDN2N1XID",
				"mobile_otp_enabled": false,
				"name": "group1",
				"push_enabled": false,
				"sms_enabled": false,
				"status": "Active"
			  },
			  {
				"desc": "random group 2",
				"group_id": "DGIB125DJLJKYZ9W257F",
				"mobile_otp_enabled": false,
				"name": "group2",
				"push_enabled": false,
				"sms_enabled": false,
				"status": "Active"
			  },
			  {
				"desc": "random group 3",
				"group_id": "DG36ABPJ1T3RZDL7ISLC",
				"mobile_otp_enabled": false,
				"name": "group3",
				"push_enabled": false,
				"sms_enabled": false,
				"status": "Active"
			  }
			],
			"stat": "OK"
		  }`))

	// Groups Page 2:
	case "/admin/v1/groups?limit=3&offset=3":
		w.Write([]byte(`{
			"metadata": {
			  "prev_offset": 0,
			  "total_objects": 5
			},
			"response": [
			  {
				"desc": "random group 4",
				"group_id": "DG6IHDSWM72IJJNXBA82",
				"mobile_otp_enabled": false,
				"name": "group4",
				"push_enabled": false,
				"sms_enabled": false,
				"status": "Active"
			  },
			  {
				"desc": "empty group 5",
				"group_id": "DGKQMVO91JT365VY36MU",
				"mobile_otp_enabled": false,
				"name": "group5",
				"push_enabled": false,
				"sms_enabled": false,
				"status": "Active"
			  }
			],
			"stat": "OK"
		  }`))

	// Phones Page 1:
	case "/admin/v1/phones?limit=2&offset=0":
		w.Write([]byte(`{
			"metadata": {
			  "next_offset": 2,
			  "total_objects": 3
			},
			"response": [
			  {
				"activated": false,
				"capabilities": [
				  "sms"
				],
				"encrypted": "",
				"extension": "",
				"fingerprint": "",
				"last_seen": "",
				"model": "Unknown",
				"name": "",
				"number": "+11111111111",
				"phone_id": "DPFL36P8Z8LZANN1FFEZ",
				"platform": "Generic Smartphone",
				"screenlock": "",
				"sms_passcodes_sent": false,
				"tampered": "",
				"type": "Mobile",
				"users": [
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041056,
					"email": "user1@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 1",
					"status": "active",
					"user_id": "DUYC8O4O953VBGGKLHAL",
					"username": "user1"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041057,
					"email": "user10@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 10",
					"status": "active",
					"user_id": "DUG1B8MRABMVKYVCFO8H",
					"username": "user10"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041056,
					"email": "user2@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 2",
					"status": "active",
					"user_id": "DUHUTX7KGB6D15WTD3VY",
					"username": "user2"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041056,
					"email": "user3@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 3",
					"status": "active",
					"user_id": "DUEL2SL4CWLP04CE71SL",
					"username": "user3"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041056,
					"email": "user4@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 4",
					"status": "active",
					"user_id": "DUB3BH17CE2V7B744RLI",
					"username": "user4"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041056,
					"email": "user5@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 5",
					"status": "active",
					"user_id": "DUWC7NXJX7IM9I7J26AT",
					"username": "user5"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041057,
					"email": "user7@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 7",
					"status": "active",
					"user_id": "DU9SRK429IRM2J7OEDCP",
					"username": "user7"
				  },
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041057,
					"email": "user9@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 9",
					"status": "active",
					"user_id": "DU2T7B5VIC0RSCN1A13W",
					"username": "user9"
				  }
				]
			  },
			  {
				"activated": false,
				"capabilities": [
				  "sms"
				],
				"encrypted": "",
				"extension": "",
				"fingerprint": "",
				"last_seen": "",
				"model": "Unknown",
				"name": "",
				"number": "+11111111111",
				"phone_id": "DPX0H7ZWQLSB735FEHVY",
				"platform": "Generic Smartphone",
				"screenlock": "",
				"sms_passcodes_sent": false,
				"tampered": "",
				"type": "Mobile",
				"users": [
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041057,
					"email": "user9@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 9",
					"status": "active",
					"user_id": "DU2T7B5VIC0RSCN1A13W",
					"username": "user9"
				  }
				]
			  }
			],
			"stat": "OK"
		  }`))

	// Phones Page 2:
	case "/admin/v1/phones?limit=2&offset=2":
		w.Write([]byte(`{
			"metadata": {
			  "prev_offset": 0,
			  "total_objects": 3
			},
			"response": [
			  {
				"activated": false,
				"capabilities": [
				  "sms"
				],
				"encrypted": "",
				"extension": "",
				"fingerprint": "",
				"last_seen": "",
				"model": "Unknown",
				"name": "",
				"number": "+11111111111",
				"phone_id": "DP7MW6K4G1OVMP8DTI08",
				"platform": "Generic Smartphone",
				"screenlock": "",
				"sms_passcodes_sent": false,
				"tampered": "",
				"type": "Mobile",
				"users": [
				  {
					"alias1": null,
					"alias2": null,
					"alias3": null,
					"alias4": null,
					"aliases": {},
					"created": 1706041057,
					"email": "user10@example.com",
					"firstname": null,
					"is_enrolled": false,
					"last_directory_sync": null,
					"last_login": null,
					"lastname": null,
					"notes": "",
					"realname": "Test User 10",
					"status": "active",
					"user_id": "DUG1B8MRABMVKYVCFO8H",
					"username": "user10"
				  }
				]
			  }
			],
			"stat": "OK"
		  }`))

	// Endpoints Page 1:
	case "/admin/v1/endpoints?limit=2&offset=0":
		w.Write([]byte(`{
			"stat": "OK",
			"response": [
			  {
				"browsers": [
				  {
					"browser_family": "Chrome",
					"browser_version": "91.0.4472.77",
					"flash_version": "uninstalled",
					"java_version": "uninstalled",
					"last_used": 1624451420
				  },
				  {
					"browser_family": "Safari",
					"browser_version": "14.1",
					"flash_version": "uninstalled",
					"java_version": "uninstalled",
					"last_used": 1624457297
				  }
				],
				"computer_sid": "",
				"cpu_id": "",
				"device_id": "",
				"device_identifier": "3FA47335-1976-3BED-8286-D3F1ABCDEA13",
				"device_identifier_type": "hardware_uuid",
				"device_name": "ejmac",
				"device_udid": "",
				"device_username": "mba22915’s MacBook Air/mba22915",
				"device_username_type": "os_username",
				"disk_encryption_status": "On",
				"domain_sid": "",
				"email": "ejennings@example.com",
				"epkey": "EP18JX1A10AB102M2T2X",
				"firewall_status": "On",
				"hardware_uuid": "3FA47335-1976-3BED-8286-D3F1ABCDEA13",
				"health_app_client_version": "2.13.1.0",
				"health_data_last_collected": 1624451421,
				"last_updated": 1624451420,
				"machine_guid": "",
				"model": "",
				"os_build": "19H1030",
				"os_family": "Mac OS X",
				"os_version": "10.11.7",
				"password_status": "Set",
				"security_agents": [
				  {
					"security_agent": "Cisco AMP for Endpoints",
					"version": "10.1.2.3"
				  }
				],
				"trusted_endpoint": "yes",
				"type": "",
				"username": "ejennings"
			  },
			  {
				"browsers": [
				  {
					"browser_family": "Mobile Safari",
					"browser_version": "9.0",
					"flash_version": "uninstalled",
					"java_version": "uninstalled"
				  }
				],
				"computer_sid": "",
				"cpu_id": "",
				"device_id": "",
				"device_identifier": "",
				"device_identifier_type": "",
				"device_name": "",
				"device_udid": "",
				"device_username": "",
				"device_username_type": "",
				"disk_encryption_status": "Unknown",
				"domain_sid": "",
				"email": "mhanson@example.com",
				"epkey": "EP65MWZWXA10AB1027TQ",
				"firewall_status": "Unknown",
				"hardware_uuid": "",
				"health_app_client_version": "",
				"health_data_last_collected": "",
				"last_updated": 1622036309,
				"machine_guid": "",
				"model": "iPhone",
				"os_build": "",
				"os_family": "iOS",
				"os_version": "14.5.1",
				"password_status": "Unknown",
				"security_agents": [],
				"trusted_endpoint": "unknown",
				"type": "",
				"username": "mhanson"
			  }
			]
		  }`))

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(``))
	}
})

func TestParseResponse(t *testing.T) {
	tests := map[string]struct {
		body           []byte
		wantObjects    []map[string]interface{}
		wantNextCursor *pagination.CompositeCursor[int64]
		wantErr        *framework.Error
	}{
		"first_page": {
			body: []byte(`{"metadata": {"next_offset": 2, "prev_offset": 0, "total_objects": 3}, "response": [{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}], "stat": "OK"}`),
			wantObjects: []map[string]interface{}{
				{"id": "00ub0oNGTSWTBKOLGLNR", "status": "ACTIVE"},
				{"id": "00ub0oNGTSWTBKOCHDKE", "status": "ACTIVE"},
			},
			wantNextCursor: &pagination.CompositeCursor[int64]{
				Cursor: testutil.GenPtr[int64](2),
			},
		},
		"last_page": {
			body: []byte(`{"metadata": {"prev_offset": 2, "total_objects": 3}, "response": [{"id": "00ub0oNGTSWTBKsdasss","status": "ACTIVE"}], "stat": "OK"}`),
			wantObjects: []map[string]interface{}{
				{"id": "00ub0oNGTSWTBKsdasss", "status": "ACTIVE"},
			},
			wantNextCursor: nil,
		},
		"invalid_object_structure": {
			body: []byte(`[{"id": "00ub0oNGTSWTBKOLGLNR","status": "ACTIVE"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: "Failed to unmarshal the datasource response: json: cannot unmarshal array into Go value of type duo.DatasourceResponse.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
		"invalid_objects": {
			body: []byte(`{"response": [{"00ub0oNGTSWTBKOLGLNR"}, {"id": "00ub0oNGTSWTBKOCHDKE","status": "ACTIVE"}]}`),
			wantErr: testutil.GenPtr(framework.Error{
				Message: `Failed to unmarshal the datasource response: invalid character '}' after object key.`,
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotObjects, gotCursor, gotErr := duo.ParseResponse(tt.body)

			if !reflect.DeepEqual(gotObjects, tt.wantObjects) {
				t.Errorf("gotObjects: %v, wantObjects: %v", gotObjects, tt.wantObjects)
			}

			if !reflect.DeepEqual(gotCursor, tt.wantNextCursor) {
				t.Errorf("gotNextLink: %v, wantNextLink: %v", gotCursor, tt.wantNextCursor)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetUserPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	duoClient := duo.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *duo.Request
		wantRes *duo.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:               server.URL,
				IntegrationKey:        "test key",
				Secret:                "test secret",
				PageSize:              4,
				EntityExternalID:      "User",
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user1@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 1",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUYC8O4O953VBGGKLHAL",
						"username":            "user1",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user2@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
							map[string]any{
								"desc":               "random group 2",
								"group_id":           "DGIB125DJLJKYZ9W257F",
								"mobile_otp_enabled": false,
								"name":               "group2",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 2",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUHUTX7KGB6D15WTD3VY",
						"username":            "user2",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user3@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 3",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUEL2SL4CWLP04CE71SL",
						"username":            "user3",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user4@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
							map[string]any{
								"desc":               "random group 2",
								"group_id":           "DGIB125DJLJKYZ9W257F",
								"mobile_otp_enabled": false,
								"name":               "group2",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 4",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUB3BH17CE2V7B744RLI",
						"username":            "user4",
						"webauthncredentials": make([]any, 0),
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](4),
				},
			},
			wantErr: nil,
		},
		"second_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:          server.URL,
				IntegrationKey:   "test key",
				Secret:           "test secret",
				PageSize:         4,
				EntityExternalID: "User",
				APIVersion:       "v1",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](4),
				},
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user5@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 5",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUWC7NXJX7IM9I7J26AT",
						"username":            "user5",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041056.0),
						"desktoptokens": make([]any, 0),
						"email":         "user6@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         false,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones":              make([]any, 0),
						"realname":            "Test User 6",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUQ87KL4A6OU5VYMWWLT",
						"username":            "user6",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":        nil,
						"alias2":        nil,
						"alias3":        nil,
						"alias4":        nil,
						"aliases":       map[string]interface{}{},
						"created":       float64(1706041057.0),
						"desktoptokens": make([]any, 0),
						"email":         "user7@example.com",
						"firstname":     nil,
						"groups": []any{
							map[string]any{
								"desc":               "random group 1",
								"group_id":           "DGKUKQSTG7ZFDN2N1XID",
								"mobile_otp_enabled": false,
								"name":               "group1",
								"push_enabled":       false,
								"sms_enabled":        false,
								"status":             "Active",
							},
						},
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 7",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DU9SRK429IRM2J7OEDCP",
						"username":            "user7",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":              nil,
						"alias2":              nil,
						"alias3":              nil,
						"alias4":              nil,
						"aliases":             map[string]interface{}{},
						"created":             float64(1706041057.0),
						"desktoptokens":       make([]any, 0),
						"email":               "user8@example.com",
						"firstname":           nil,
						"groups":              make([]any, 0),
						"is_enrolled":         false,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones":              make([]any, 0),
						"realname":            "Test User 8",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DU1E2CSOI6I5HEO043WN",
						"username":            "user8",
						"webauthncredentials": make([]any, 0),
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](8),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:          server.URL,
				IntegrationKey:   "test key",
				Secret:           "test secret",
				PageSize:         4,
				EntityExternalID: "User",
				APIVersion:       "v1",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](8),
				},
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"alias1":              nil,
						"alias2":              nil,
						"alias3":              nil,
						"alias4":              nil,
						"aliases":             map[string]interface{}{},
						"created":             float64(1706041057.0),
						"desktoptokens":       make([]any, 0),
						"email":               "user9@example.com",
						"firstname":           nil,
						"groups":              make([]any, 0),
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPX0H7ZWQLSB735FEHVY",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 9",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DU2T7B5VIC0RSCN1A13W",
						"username":            "user9",
						"webauthncredentials": make([]any, 0),
					},
					{
						"alias1":              nil,
						"alias2":              nil,
						"alias3":              nil,
						"alias4":              nil,
						"aliases":             map[string]interface{}{},
						"created":             float64(1706041057.0),
						"desktoptokens":       make([]any, 0),
						"email":               "user10@example.com",
						"firstname":           nil,
						"groups":              make([]any, 0),
						"is_enrolled":         true,
						"last_directory_sync": nil,
						"last_login":          nil,
						"lastname":            nil,
						"lockout_reason":      nil,
						"notes":               "",
						"phones": []any{
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
							map[string]any{
								"activated":          false,
								"capabilities":       []any{"sms"},
								"encrypted":          "",
								"extension":          "",
								"fingerprint":        "",
								"last_seen":          "",
								"model":              "Unknown",
								"name":               "",
								"number":             "+11111111111",
								"phone_id":           "DP7MW6K4G1OVMP8DTI08",
								"platform":           "Generic Smartphone",
								"screenlock":         "",
								"sms_passcodes_sent": false,
								"tampered":           "",
								"type":               "Mobile",
							},
						},
						"realname":            "Test User 10",
						"status":              "active",
						"tokens":              make([]any, 0),
						"u2ftokens":           make([]any, 0),
						"user_id":             "DUG1B8MRABMVKYVCFO8H",
						"username":            "user10",
						"webauthncredentials": make([]any, 0),
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := duoClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetGroupPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	duoClient := duo.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *duo.Request
		wantRes *duo.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:               server.URL,
				IntegrationKey:        "test key",
				Secret:                "test secret",
				PageSize:              3,
				EntityExternalID:      "Group",
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"desc":               "random group 1",
						"group_id":           "DGKUKQSTG7ZFDN2N1XID",
						"mobile_otp_enabled": false,
						"name":               "group1",
						"push_enabled":       false,
						"sms_enabled":        false,
						"status":             "Active",
					},
					{
						"desc":               "random group 2",
						"group_id":           "DGIB125DJLJKYZ9W257F",
						"mobile_otp_enabled": false,
						"name":               "group2",
						"push_enabled":       false,
						"sms_enabled":        false,
						"status":             "Active",
					},
					{
						"desc":               "random group 3",
						"group_id":           "DG36ABPJ1T3RZDL7ISLC",
						"mobile_otp_enabled": false,
						"name":               "group3",
						"push_enabled":       false,
						"sms_enabled":        false,
						"status":             "Active",
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](3),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:          server.URL,
				IntegrationKey:   "test key",
				Secret:           "test secret",
				PageSize:         3,
				EntityExternalID: "Group",
				APIVersion:       "v1",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](3),
				},
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"desc":               "random group 4",
						"group_id":           "DG6IHDSWM72IJJNXBA82",
						"mobile_otp_enabled": false,
						"name":               "group4",
						"push_enabled":       false,
						"sms_enabled":        false,
						"status":             "Active",
					},
					{
						"desc":               "empty group 5",
						"group_id":           "DGKQMVO91JT365VY36MU",
						"mobile_otp_enabled": false,
						"name":               "group5",
						"push_enabled":       false,
						"sms_enabled":        false,
						"status":             "Active",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := duoClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetPhonePage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	duoClient := duo.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *duo.Request
		wantRes *duo.Response
		wantErr *framework.Error
	}{
		"first_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:               server.URL,
				IntegrationKey:        "test key",
				Secret:                "test secret",
				PageSize:              2,
				EntityExternalID:      "Phone",
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"activated":          false,
						"capabilities":       []any{"sms"},
						"encrypted":          "",
						"extension":          "",
						"fingerprint":        "",
						"last_seen":          "",
						"model":              "Unknown",
						"name":               "",
						"number":             "+11111111111",
						"phone_id":           "DPFL36P8Z8LZANN1FFEZ",
						"platform":           "Generic Smartphone",
						"screenlock":         "",
						"sms_passcodes_sent": false,
						"tampered":           "",
						"type":               "Mobile",
						"users": []any{
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041056.0),
								"email":               "user1@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 1",
								"status":              "active",
								"user_id":             "DUYC8O4O953VBGGKLHAL",
								"username":            "user1",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041057.0),
								"email":               "user10@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 10",
								"status":              "active",
								"user_id":             "DUG1B8MRABMVKYVCFO8H",
								"username":            "user10",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041056.0),
								"email":               "user2@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 2",
								"status":              "active",
								"user_id":             "DUHUTX7KGB6D15WTD3VY",
								"username":            "user2",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041056.0),
								"email":               "user3@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 3",
								"status":              "active",
								"user_id":             "DUEL2SL4CWLP04CE71SL",
								"username":            "user3",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041056.0),
								"email":               "user4@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 4",
								"status":              "active",
								"user_id":             "DUB3BH17CE2V7B744RLI",
								"username":            "user4",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041056.0),
								"email":               "user5@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 5",
								"status":              "active",
								"user_id":             "DUWC7NXJX7IM9I7J26AT",
								"username":            "user5",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041057.0),
								"email":               "user7@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 7",
								"status":              "active",
								"user_id":             "DU9SRK429IRM2J7OEDCP",
								"username":            "user7",
							},
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041057.0),
								"email":               "user9@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 9",
								"status":              "active",
								"user_id":             "DU2T7B5VIC0RSCN1A13W",
								"username":            "user9",
							},
						},
					},
					{
						"activated":          false,
						"capabilities":       []any{"sms"},
						"encrypted":          "",
						"extension":          "",
						"fingerprint":        "",
						"last_seen":          "",
						"model":              "Unknown",
						"name":               "",
						"number":             "+11111111111",
						"phone_id":           "DPX0H7ZWQLSB735FEHVY",
						"platform":           "Generic Smartphone",
						"screenlock":         "",
						"sms_passcodes_sent": false,
						"tampered":           "",
						"type":               "Mobile",
						"users": []any{
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041057.0),
								"email":               "user9@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 9",
								"status":              "active",
								"user_id":             "DU2T7B5VIC0RSCN1A13W",
								"username":            "user9",
							},
						},
					},
				},
				NextCursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
			},
			wantErr: nil,
		},
		"last_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:          server.URL,
				IntegrationKey:   "test key",
				Secret:           "test secret",
				PageSize:         2,
				EntityExternalID: "Phone",
				APIVersion:       "v1",
				Cursor: &pagination.CompositeCursor[int64]{
					Cursor: testutil.GenPtr[int64](2),
				},
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"activated":          false,
						"capabilities":       []any{"sms"},
						"encrypted":          "",
						"extension":          "",
						"fingerprint":        "",
						"last_seen":          "",
						"model":              "Unknown",
						"name":               "",
						"number":             "+11111111111",
						"phone_id":           "DP7MW6K4G1OVMP8DTI08",
						"platform":           "Generic Smartphone",
						"screenlock":         "",
						"sms_passcodes_sent": false,
						"tampered":           "",
						"type":               "Mobile",
						"users": []any{
							map[string]any{
								"alias1":              nil,
								"alias2":              nil,
								"alias3":              nil,
								"alias4":              nil,
								"aliases":             map[string]interface{}{},
								"created":             float64(1706041057.0),
								"email":               "user10@example.com",
								"firstname":           nil,
								"is_enrolled":         false,
								"last_directory_sync": nil,
								"last_login":          nil,
								"lastname":            nil,
								"notes":               "",
								"realname":            "Test User 10",
								"status":              "active",
								"user_id":             "DUG1B8MRABMVKYVCFO8H",
								"username":            "user10",
							},
						},
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := duoClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetEndpointPage(t *testing.T) {
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	duoClient := duo.NewClient(client)
	server := httptest.NewServer(TestServerHandler)
	tests := map[string]struct {
		context context.Context
		request *duo.Request
		wantRes *duo.Response
		wantErr *framework.Error
	}{
		"first_and_last_page": {
			context: context.Background(),
			request: &duo.Request{
				BaseURL:               server.URL,
				IntegrationKey:        "test key",
				Secret:                "test secret",
				PageSize:              2,
				EntityExternalID:      "Endpoint",
				APIVersion:            "v1",
				RequestTimeoutSeconds: 5,
			},
			wantRes: &duo.Response{
				StatusCode: http.StatusOK,
				Objects: []map[string]interface{}{
					{
						"browsers": []any{
							map[string]any{
								"browser_family":  "Chrome",
								"browser_version": "91.0.4472.77",
								"flash_version":   "uninstalled",
								"java_version":    "uninstalled",
								"last_used":       float64(1.62445142e+09),
							},
							map[string]any{
								"browser_family":  "Safari",
								"browser_version": "14.1",
								"flash_version":   "uninstalled",
								"java_version":    "uninstalled",
								"last_used":       float64(1.624457297e+09),
							},
						},
						"computer_sid":               "",
						"cpu_id":                     "",
						"device_id":                  "",
						"device_identifier":          "3FA47335-1976-3BED-8286-D3F1ABCDEA13",
						"device_identifier_type":     "hardware_uuid",
						"device_name":                "ejmac",
						"device_udid":                "",
						"device_username":            "mba22915’s MacBook Air/mba22915",
						"device_username_type":       "os_username",
						"disk_encryption_status":     "On",
						"domain_sid":                 "",
						"email":                      "ejennings@example.com",
						"epkey":                      "EP18JX1A10AB102M2T2X",
						"firewall_status":            "On",
						"hardware_uuid":              "3FA47335-1976-3BED-8286-D3F1ABCDEA13",
						"health_app_client_version":  "2.13.1.0",
						"health_data_last_collected": float64(1.624451421e+09),
						"last_updated":               float64(1.62445142e+09),
						"machine_guid":               "",
						"model":                      "",
						"os_build":                   "19H1030",
						"os_family":                  "Mac OS X",
						"os_version":                 "10.11.7",
						"password_status":            "Set",
						"security_agents": []any{
							map[string]any{
								"security_agent": "Cisco AMP for Endpoints",
								"version":        "10.1.2.3",
							},
						},
						"trusted_endpoint": "yes",
						"type":             "",
						"username":         "ejennings",
					},
					{
						"browsers": []any{
							map[string]any{
								"browser_family":  "Mobile Safari",
								"browser_version": "9.0",
								"flash_version":   "uninstalled",
								"java_version":    "uninstalled",
							},
						},
						"computer_sid":               "",
						"cpu_id":                     "",
						"device_id":                  "",
						"device_identifier":          "",
						"device_identifier_type":     "",
						"device_name":                "",
						"device_udid":                "",
						"device_username":            "",
						"device_username_type":       "",
						"disk_encryption_status":     "Unknown",
						"domain_sid":                 "",
						"email":                      "mhanson@example.com",
						"epkey":                      "EP65MWZWXA10AB1027TQ",
						"firewall_status":            "Unknown",
						"hardware_uuid":              "",
						"health_app_client_version":  "",
						"health_data_last_collected": "",
						"last_updated":               float64(1.622036309e+09),
						"machine_guid":               "",
						"model":                      "iPhone",
						"os_build":                   "",
						"os_family":                  "iOS",
						"os_version":                 "14.5.1",
						"password_status":            "Unknown",
						"security_agents":            []interface{}{},
						"trusted_endpoint":           "unknown",
						"type":                       "",
						"username":                   "mhanson",
					},
				},
				NextCursor: nil,
			},
			wantErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotRes, gotErr := duoClient.GetPage(tt.context, tt.request)

			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("gotRes: %v, wantRes: %v", gotRes, tt.wantRes)
			}

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
