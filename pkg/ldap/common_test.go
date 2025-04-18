// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package ldap_test

import (
	framework "github.com/sgnl-ai/adapter-framework"
	ldap_adapter "github.com/sgnl-ai/adapters/pkg/ldap"
)

var (
	mockLDAPAddr     = "ldap://127.0.0.1:389"
	mockLDAPSAddr    = "ldaps://127.0.0.1:636"
	mockLDAPUser     = "cn=admin,dc=example,dc=org"
	mockLDAPPassword = "admin"

	validAuthCredentials = &framework.DatasourceAuthCredentials{
		Basic: &framework.BasicAuthCredentials{
			Username: mockLDAPUser,
			Password: mockLDAPPassword,
		},
	}

	validCommonConfig = &ldap_adapter.Config{
		BaseDN: "dc=corp,dc=example,dc=io",
		EntityConfigMap: map[string]*ldap_adapter.EntityConfig{
			"Person": {
				Query: "(&(objectClass=person))",
			},
		},
		CertificateChain: `
		LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNzRENDQWhtZ0F3SUJBZ0lKQUx3enJKRUlCT2FlTUEwR0NTcUd
		TSWIzRFFFQkJRVUFNRVV4Q3pBSkJnTlYKQkFZVEFrRlZNUk13RVFZRFZRUUlFd3BUYjIxbExWTjBZWFJsTVNFd0h3WU
		RWUVFLRXhoSmJuUmxjbTVsZENCWAphV1JuYVhSeklGQjBlU0JNZEdRd0hoY05NVEV3T1RNd01UVXlOak0yV2hjTk1qR
		XdPVEkzTVRVeU5qTTJXakJGCk1Rc3dDUVlEVlFRR0V3SkJWVEVUTUJFR0ExVUVDQk1LVTI5dFpTMVRkR0YwWlRFaE1C
		OEdBMVVFQ2hNWVNXNTAKWlhKdVpYUWdWMmxrWjJsMGN5QlFkSGtnVEhSa01JR2ZNQTBHQ1NxR1NJYjNEUUVCQVFVQUE
		0R05BRENCaVFLQgpnUUM4OENrd3J1OVZSMnAyS0oxV1F5cWVzTHpyOTV0YU5iaGtZZnNkMGo4VGwwTUdZNWgrZGN6Q2
		FNUXowWVkzCnhIWHVVNXlBUVFUWmppa3MrRDNLQTNjeCtpS0RmMnAxcTc3b1h4UWN4NUNrclhCV1R hWDJvcVZ0SG0z
		YVgyM0IKQUlPUkd1UGswMGI0clQzY2xkN1ZoY0VGbXpSTmJ5STBFcUxNQXhJd2NlVUtTUUlEQVFBQm80R25NSUdrTUI
		wRwpBMVVkRGdRV0JCU0dtT2R2U1hLWGNsaWM1VU9LUFczNUpMTUVFakIxQmdOVkhTTUViakJzZ0JTR21PZHZTWEtYCm
		NsaWM1VU9LUFczNUpMTUVFcUZKcEVjd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdUQ2xOdmJXVXQKVTNSa
		GRHVXhJVEFmQmdOVkJBb1RHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpJSUpBTHd6ckpFSQpCT2FlTUF3
		R0ExVWRFd1FGTUFNQkFmOHdEUVlKS29aSWh2Y05BUUVGQlFBRGdZRUFjUGZXbjQ5cGdBWDU0amk1ClNpVVBGRk5DdVF
		HU1NUSGgySStUTXJzMUcxTWIzYTBYMWRWNUNOTFJ5WHl1VnhzcWhpTS9IMnZlRm5UejJRNFUKd2RZL2tQeEUxOUF1d2
		N6OUF2Q2t3N29sMUxJbExmSnZCemp6T2pFcFpKTnRrWFR4OFJPU29vTnJEZUpsM0h5TgpjY2lTNWhmODBYeklGcXdoe
		mFWUzlnbWl5TTg9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0=
		`,
	}
)
