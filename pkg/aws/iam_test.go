// Copyright 2026 SGNL.ai, Inc.

// nolint: goconst

package aws_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go/middleware"
	aws_adapter "github.com/sgnl-ai/adapters/pkg/aws"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

var (
	// Dummy Data Entities and Their Relationships
	//
	// This section contains mock data for AWS IAM entities, including Users, Groups, Roles, Policies,
	// SAML Providers, RolePolicies, GroupPolicies, UserPolicies and GroupMember along with their relationship.
	//
	// Users:
	// There are 6 users (user1 through user6).
	//
	// Groups:
	// There are 4 groups (Group1 through Group4).
	//
	// Roles:
	// There are 3 roles.
	//
	// Policies:
	// There are 2 policies (ExampleEngPolicy and Policy2).
	//
	// SAML Providers:
	// There are 2 SAML providers (Provider1 and Provider2).
	//
	// Relationships:
	// - Group1 members: (user1, user2)
	// - Group2 members: (user1, user2)
	// - Group1 has 2 policies attached: (ExampleEngPolicy, Policy2)
	// - Group4 has 1 policies attached: (ExampleEngPolicy)
	// - User1 has 1 policy attached: (Policy104)
	// - User6 has 1 policy attached: (Policy105)
	// - Role1 has 2 policies attached: (Policy106, Policy102)
	// - Role2 has 1 policy attached: (Policy107).

	mockUsers = []types.User{
		{
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user1"),
			UserName:         testutil.GenPtr("user1"),
			UserId:           testutil.GenPtr("user1"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user2"),
			UserName:         testutil.GenPtr("user2"),
			UserId:           testutil.GenPtr("user2"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user3"),
			UserName:         testutil.GenPtr("user3"),
			UserId:           testutil.GenPtr("user3"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user4"),
			UserName:         testutil.GenPtr("user4"),
			UserId:           testutil.GenPtr("user4"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user5"),
			UserName:         testutil.GenPtr("user5"),
			UserId:           testutil.GenPtr("user5"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		}, {
			Arn:              testutil.GenPtr("arn:aws:iam::000000000000:user/user6"),
			UserName:         testutil.GenPtr("user6"),
			UserId:           testutil.GenPtr("user6"),
			Path:             testutil.GenPtr("/"),
			CreateDate:       aws.Time(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			PasswordLastUsed: aws.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	mockGroups = []types.Group{
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:group/Group1"),
			GroupName:  testutil.GenPtr("Group1"),
			GroupId:    testutil.GenPtr("Group1"),
			Path:       testutil.GenPtr("/"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:group/Group2"),
			GroupName:  testutil.GenPtr("Group2"),
			GroupId:    testutil.GenPtr("Group2"),
			Path:       testutil.GenPtr("/"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:group/Group3"),
			GroupName:  testutil.GenPtr("Group3"),
			GroupId:    testutil.GenPtr("Group3"),
			Path:       testutil.GenPtr("/"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:group/Group4"),
			GroupName:  testutil.GenPtr("Group4"),
			GroupId:    testutil.GenPtr("Group4"),
			Path:       testutil.GenPtr("/"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	mockRoles = []types.Role{
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_de25js739eef1832"),
			RoleName:   testutil.GenPtr("AWSReservedSSO_de25js739eef1832"),
			RoleId:     testutil.GenPtr("AROAXXXXXXXXXXXXXXXX2"),
			Path:       testutil.GenPtr("/sso.amazonaws.com"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:role/sso.amazonaws.com/AWSReservedSSO_abcdef1234567890"),
			RoleName:   testutil.GenPtr("AWSReservedSSO_abcdef1234567890"),
			RoleId:     testutil.GenPtr("AROA3C7OBZZZCD433N4DQ"),
			Path:       testutil.GenPtr("/sso.amazonaws.com"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::000000000000:role/sso.amazonaws.com/role_3"),
			RoleName:   testutil.GenPtr("role_3"),
			RoleId:     testutil.GenPtr("AROA3C7OBYYYCD433N4DQ"),
			Path:       testutil.GenPtr("/sso.amazonaws.com"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	mockPolicies = []types.Policy{
		{
			Arn:                           testutil.GenPtr("arn:aws:iam::000000000000:policy/ExampleEngPolicy"),
			PolicyName:                    testutil.GenPtr("ExampleEngPolicy"),
			PolicyId:                      testutil.GenPtr("ANPA3C7OBZZZCD411N4DQ"),
			AttachmentCount:               testutil.GenPtr(int32(1)),
			DefaultVersionId:              testutil.GenPtr("v1"),
			CreateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			IsAttachable:                  true,
			Path:                          testutil.GenPtr("/"),
			UpdateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			PermissionsBoundaryUsageCount: testutil.GenPtr(int32(0)),
		},
		{
			Arn:                           testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy2"),
			PolicyName:                    testutil.GenPtr("Policy2"),
			PolicyId:                      testutil.GenPtr("ANPA3C7OBZZZCD433N4DQ"),
			AttachmentCount:               testutil.GenPtr(int32(1)),
			DefaultVersionId:              testutil.GenPtr("v1"),
			CreateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			IsAttachable:                  true,
			Path:                          testutil.GenPtr("/"),
			UpdateDate:                    aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			PermissionsBoundaryUsageCount: testutil.GenPtr(int32(0)),
		},
	}

	mockSAMLProviders = []types.SAMLProviderListEntry{
		{
			Arn:        testutil.GenPtr("arn:aws:iam::123456789012:saml-provider/Provider1"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			ValidUntil: aws.Time(time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			Arn:        testutil.GenPtr("arn:aws:iam::123456789012:saml-provider/Provider2"),
			CreateDate: aws.Time(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)),
			ValidUntil: aws.Time(time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	mockAttachedGroup1Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("ExampleEngPolicy"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/ExampleEngPolicy"),
		},
		{
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy2"),
			PolicyName: testutil.GenPtr("Policy2"),
		},
	}

	mockAttachedGroup4Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("ExampleEngPolicy"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/ExampleEngPolicy"),
		},
	}

	mockAttachedUser1Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("Policy104"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy104"),
		},
	}

	mockAttachedUser6Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("Policy105"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy105"),
		},
	}

	mockAttachedRole1Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("Policy106"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy106"),
		},
		{
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy102"),
			PolicyName: testutil.GenPtr("Policy102"),
		},
	}

	mockAttachedRole2Policies = []types.AttachedPolicy{
		{
			PolicyName: testutil.GenPtr("Policy107"),
			PolicyArn:  testutil.GenPtr("arn:aws:iam::000000000000:policy/Policy107"),
		},
	}
)

// getInputOptions is a middleware.InitializeMiddlewareFunc that helps to grab
// appropriate options for the IAM API calls.
func getInputOptions(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	switch v := in.Parameters.(type) {
	case *iam.ListUsersInput: // At the time of a particular API call
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListRolesInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListGroupsInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListPoliciesInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListSAMLProvidersInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{})
	case *iam.ListAttachedGroupPoliciesInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.GroupName,
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListAttachedRolePoliciesInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.RoleName,
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.ListAttachedUserPoliciesInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.UserName,
			InputParams: aws_adapter.InputParams{
				Marker:     v.Marker,
				MaxItems:   v.MaxItems,
				PathPrefix: v.PathPrefix,
			},
		})
	case *iam.GetUserInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.UserName,
		})
	case *iam.GetRoleInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.RoleName,
		})
	case *iam.GetGroupInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.GroupName,
		})
	case *iam.GetPolicyInput:
		ctx = middleware.WithStackValue(ctx, aws_adapter.Options{}, &aws_adapter.Options{
			UniqueName: v.PolicyArn,
		})
	}

	return next.HandleInitialize(ctx, in)
}

// setupMocks sets up the mocks for the AWS API calls based on the operation name and options.
// It returns the mocked output, metadata, and error.
//
// Reference:
//   - https://dev.to/aws-builders/testing-with-aws-sdk-for-go-v2-without-interface-mocks-55de
func setupMocks(operationName string, options aws_adapter.Options,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	switch operationName {
	case "ListUsers":
		if options.PathPrefix != nil {
			// If the PathPrefix is set to "some-internal-error", return an error.
			// this is just a test case to simulate an internal error.
			if *options.PathPrefix == "some-internal-error" {
				return middleware.FinalizeOutput{
					Result: nil,
				}, middleware.Metadata{}, fmt.Errorf("InternalFailure")
			}

			// If the PathPrefix is set to "/not-found", return an empty list of users.
			if *options.PathPrefix == "/not-found" {
				return middleware.FinalizeOutput{
					Result: &iam.ListUsersOutput{},
				}, middleware.Metadata{}, nil
			}
		}

		users, marker := paginate(mockUsers, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListUsersOutput{
				Users:  users,
				Marker: marker,
			},
		}, middleware.Metadata{}, nil
	case "GetUser":
		return middleware.FinalizeOutput{
			Result: &iam.GetUserOutput{
				User: func() *types.User {
					for _, u := range mockUsers {
						if u.UserName == options.UniqueName {
							return &u
						}
					}

					return nil
				}(),
			},
		}, middleware.Metadata{}, nil
	case "ListRoles":
		if options.PathPrefix != nil && *options.PathPrefix == internalError {
			return middleware.FinalizeOutput{
				Result: &iam.ListRolesOutput{},
			}, middleware.Metadata{}, fmt.Errorf("InternalFailure")
		}

		roles, marker := paginate(mockRoles, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListRolesOutput{
				Roles:  roles,
				Marker: marker,
			},
		}, middleware.Metadata{}, nil
	case "GetRole":
		return middleware.FinalizeOutput{
			Result: &iam.GetRoleOutput{
				Role: func() *types.Role {
					for _, r := range mockRoles {
						if r.RoleName == options.UniqueName {
							return &r
						}
					}

					return nil
				}(),
			},
		}, middleware.Metadata{}, nil
	case "ListGroups":
		if options.PathPrefix != nil && *options.PathPrefix == internalError {
			return middleware.FinalizeOutput{
				Result: &iam.ListGroupsOutput{},
			}, middleware.Metadata{}, fmt.Errorf("InternalFailure")
		}

		groups, marker := paginate(mockGroups, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListGroupsOutput{
				Groups: groups,
				Marker: marker,
			},
		}, middleware.Metadata{}, nil
	case "GetGroup":
		return middleware.FinalizeOutput{
			Result: &iam.GetGroupOutput{
				Group: func() *types.Group {
					for _, g := range mockGroups {
						if g.GroupName == options.UniqueName {
							return &g
						}
					}

					return nil
				}(),
				Users: func() []types.User {
					switch *options.UniqueName {
					case "Group1", "Group2":
						return mockUsers[:2]
					default:
						return nil
					}
				}(),
			},
		}, middleware.Metadata{}, nil
	case "ListPolicies":
		if options.PathPrefix != nil && *options.PathPrefix == internalError {
			return middleware.FinalizeOutput{
				Result: &iam.ListPoliciesOutput{},
			}, middleware.Metadata{}, fmt.Errorf("InternalFailure")
		}

		policies, marker := paginate(mockPolicies, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListPoliciesOutput{
				Policies: policies,
				Marker:   marker,
			},
		}, middleware.Metadata{}, nil
	case "GetPolicy":
		return middleware.FinalizeOutput{
			Result: &iam.GetPolicyOutput{
				Policy: func() *types.Policy {
					for _, p := range mockPolicies {
						if p.Arn == options.UniqueName {
							return &p
						}
					}

					return nil
				}(),
			},
		}, middleware.Metadata{}, nil
	case "ListSAMLProviders":
		return middleware.FinalizeOutput{
			Result: &iam.ListSAMLProvidersOutput{
				SAMLProviderList: mockSAMLProviders,
			},
		}, middleware.Metadata{}, nil
	case "ListAttachedGroupPolicies":
		var attachedPolicies []types.AttachedPolicy

		switch *options.UniqueName {
		case "Group1":
			attachedPolicies = mockAttachedGroup1Policies
		case "Group4":
			attachedPolicies = mockAttachedGroup4Policies
		default:
			attachedPolicies = []types.AttachedPolicy{}
		}

		policies, marker := paginate(attachedPolicies, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListAttachedGroupPoliciesOutput{
				AttachedPolicies: policies,
				Marker:           marker,
			},
		}, middleware.Metadata{}, nil

	case "ListAttachedRolePolicies":
		var attachedPolicies []types.AttachedPolicy

		switch *options.UniqueName {
		case "AWSReservedSSO_de25js739eef1832":
			attachedPolicies = mockAttachedRole1Policies
		case "AWSReservedSSO_abcdef1234567890":
			attachedPolicies = mockAttachedRole2Policies
		default:
			attachedPolicies = nil
		}

		policies, marker := paginate(attachedPolicies, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListAttachedRolePoliciesOutput{
				AttachedPolicies: policies,
				Marker:           marker,
			},
		}, middleware.Metadata{}, nil

	case "ListAttachedUserPolicies":
		var attachedPolicies []types.AttachedPolicy

		switch *options.UniqueName {
		case "user1":
			attachedPolicies = mockAttachedUser1Policies
		case "user6":
			attachedPolicies = mockAttachedUser6Policies
		default:
			attachedPolicies = nil
		}

		policies, marker := paginate(attachedPolicies, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListAttachedUserPoliciesOutput{
				AttachedPolicies: policies,
				Marker:           marker,
			},
		}, middleware.Metadata{}, nil

	default:
		return middleware.FinalizeOutput{}, middleware.Metadata{}, fmt.Errorf("operation not found")
	}
}

func setupEmptyMocks(operationName string,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	switch operationName {
	case "ListSAMLProviders":
		return middleware.FinalizeOutput{
			Result: &iam.ListSAMLProvidersOutput{
				SAMLProviderList: nil,
			},
		}, middleware.Metadata{}, nil

	default:
		return middleware.FinalizeOutput{}, middleware.Metadata{}, fmt.Errorf("operation not found")
	}
}

// This middleware is added purely for benchmarking purposes. It facilitates a `List`
// operation followed by several `Get` operations with minimal amount of `sleep` added
// to mimic a network roundtrip cost.
func setupPolicyDataset(operationName string, options aws_adapter.Options,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	mocks := policyDataset()

	switch operationName {
	case "ListPolicies":
		// Simulate a delay in the response.
		time.Sleep(100 * time.Millisecond)

		policies, marker := paginate(mocks, options.Marker, int(*options.MaxItems))

		return middleware.FinalizeOutput{
			Result: &iam.ListPoliciesOutput{
				Policies: policies,
				Marker:   marker,
			},
		}, middleware.Metadata{}, nil
	case "GetPolicy":
		// Simulate a delay in the response.
		time.Sleep(100 * time.Millisecond)

		return middleware.FinalizeOutput{
			Result: &iam.GetPolicyOutput{
				Policy: func() *types.Policy {
					for _, p := range mocks {
						if *p.Arn == *options.UniqueName {
							return &p
						}
					}

					return nil
				}(),
			},
		}, middleware.Metadata{}, nil
	default:
		return middleware.FinalizeOutput{}, middleware.Metadata{}, fmt.Errorf("operation not found")
	}
}
