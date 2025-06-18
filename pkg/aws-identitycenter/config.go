package awsidentitycenter

import (
	"context"
	"errors"

	"github.com/sgnl-ai/adapters/pkg/config"
)

// Config is the configuration passed in each GetPage calls to the adapter.
// AWS Identity Center Adapter configuration example:
//
// {
//   "region": "us-west-2",
//   "identityStoreID": "d-1234567890",
//   "instanceARN": "arn:aws:sso:::instance/ssoins-1234567890"
// }

type Config struct {
	*config.CommonConfig

	// Region is the AWS region to query.
	Region string `json:"region"`

	// IdentityStoreID is the AWS Identity Store identifier.
	IdentityStoreID string `json:"identityStoreID"`

	// InstanceARN is the AWS Identity Center instance ARN.
	InstanceARN string `json:"instanceARN"`
}

// ValidateConfig validates that a Config received in a GetPage call is valid.
func (c *Config) Validate(_ context.Context) error {
	switch {
	case c == nil:
		return errors.New("the request contains an empty configuration")
	case c.Region == "":
		return errors.New("the AWS Region is not set in the configuration")
	case c.IdentityStoreID == "":
		return errors.New("identityStoreID is not set in the configuration")
	case c.InstanceARN == "":
		return errors.New("instanceARN is not set in the configuration")
	default:
		return nil
	}
}
