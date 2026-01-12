// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	framework "github.com/sgnl-ai/adapter-framework"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func largePolicyObjects() []framework.Object {
	var objs []framework.Object

	for i := 1; i <= 1000; i++ {
		policy := framework.Object{
			"Arn":                           fmt.Sprintf("arn:aws:iam::000000000000:policy/Policy-%v", i),
			"PolicyName":                    fmt.Sprintf("Policy-%v", i),
			"PolicyId":                      fmt.Sprintf("ANPA3C7OBZZZCD411N4D%v", i),
			"AttachmentCount":               int64(1),
			"DefaultVersionId":              "v1",
			"AccountId":                     "000000000000",
			"CreateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			"IsAttachable":                  true,
			"Path":                          "/",
			"UpdateDate":                    time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			"PermissionsBoundaryUsageCount": int64(0),
		}
		objs = append(objs, policy)
	}

	return objs
}

func policyDataset() []types.Policy {
	var policies []types.Policy

	for i := 1; i <= 1000; i++ {
		policy := types.Policy{
			Arn:                           testutil.GenPtr(fmt.Sprintf("arn:aws:iam::000000000000:policy/Policy-%v", i)),
			PolicyName:                    testutil.GenPtr(fmt.Sprintf("Policy-%v", i)),
			PolicyId:                      testutil.GenPtr(fmt.Sprintf("ANPA3C7OBZZZCD411N4D%v", i)),
			AttachmentCount:               testutil.GenPtr(int32(1)),
			DefaultVersionId:              testutil.GenPtr("v1"),
			CreateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			IsAttachable:                  true,
			Path:                          testutil.GenPtr("/"),
			UpdateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			PermissionsBoundaryUsageCount: testutil.GenPtr(int32(0)),
		}
		policies = append(policies, policy)
	}

	return policies
}
