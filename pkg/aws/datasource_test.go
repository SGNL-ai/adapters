// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"testing"

	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
	"github.com/stretchr/testify/assert"
)

func TestArnToAccountID(t *testing.T) {
	testCases := []struct {
		Arn        string
		EntityType string
		ExpectedID string
	}{
		{
			Arn:        "arn:aws:iam::123456789012:user/test-user",
			EntityType: aws_adapter.User,
			ExpectedID: "123456789012",
		},
		{
			Arn:        "arn:aws:iam::000000000000:group/Group1",
			EntityType: aws_adapter.Role,
			ExpectedID: "000000000000",
		},
		{
			Arn:        "arn:aws:iam::762319060234:role/sso.amazonaws.com/AWSReservedSSO_de25js739eef1832",
			EntityType: aws_adapter.Group,
			ExpectedID: "762319060234",
		},
		{
			Arn:        "arn:aws:iam::652319060462:policy/ExampleEngPolicy",
			EntityType: aws_adapter.Policy,
			ExpectedID: "652319060462",
		},
		{
			Arn:        "arn:aws:iam::123456789012:saml-provider/Provider1",
			EntityType: aws_adapter.IdentityProvider,
			ExpectedID: "123456789012",
		},
	}

	for _, tc := range testCases {
		entity := map[string]interface{}{
			"Arn": tc.Arn,
		}

		err := aws_adapter.ArnToAccountID(&entity, tc.EntityType)
		assert.NoError(t, err)
		assert.Equal(t, tc.ExpectedID, entity["AccountId"])
	}

	negativeTestCases := []struct {
		Entity        map[string]interface{}
		EntityType    string
		ExpectedError string
	}{
		{
			Entity:        nil,
			EntityType:    aws_adapter.User,
			ExpectedError: "Unable to find Arn in entity",
		},
		{
			Entity:        map[string]interface{}{},
			EntityType:    aws_adapter.User,
			ExpectedError: "Unable to find Arn in entity",
		},
		{
			Entity: map[string]interface{}{
				"Arn": 12345, // Arn is not a string
			},
			EntityType:    aws_adapter.User,
			ExpectedError: "failed to convert Arn to string",
		},
		{
			Entity: map[string]interface{}{
				"Arn": "invalid-arn-format",
			},
			EntityType:    aws_adapter.User,
			ExpectedError: "failed to parse Arn",
		},
	}

	for _, tc := range negativeTestCases {
		err := aws_adapter.ArnToAccountID(&tc.Entity, tc.EntityType)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), tc.ExpectedError)
	}
}
