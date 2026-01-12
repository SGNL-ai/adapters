// Copyright 2025 SGNL.ai, Inc.

// nolint:lll
/*
# Creating smoke tests
To create smoke tests please follow the guidelines documented [here](https://eng.playbooks.sgnl.host/playbooks/adapters/development/#creating-smoke-tests-built-in-adapters-only)

## Setup a recorder and skip TLS verify
To skip TLS verify with the test SoR server and generate a fixture, setup the recorder as follows.

```go
r, err := recorder.NewWithOptions(

	&recorder.Options{
		CassetteName:       "fixtures/scim/user",
		Mode:               recorder.ModeRecordOnce,
		SkipRequestLatency: false,
		RealTransport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	},

)
```
*/
package smoketests
