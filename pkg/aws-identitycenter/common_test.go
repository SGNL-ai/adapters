package awsidentitycenter_test

import (
	framework "github.com/sgnl-ai/adapter-framework"
	adapter "github.com/sgnl-ai/adapters/pkg/aws-identitycenter"
)

var (
	validAuthCredentials = &framework.DatasourceAuthCredentials{
		Basic: &framework.BasicAuthCredentials{
			Username: "access",
			Password: "secret",
		},
	}

	validConfig = &adapter.Config{
		Region:          "us-west-2",
		IdentityStoreID: "d-1234567890",
		InstanceARN:     "arn:aws:sso:::instance/ssoins-1234567890",
	}
)
