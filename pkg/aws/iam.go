// Copyright 2025 SGNL.ai, Inc.

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// Implementation of EntityHandler for IAM User.
type UserHandler struct {
	Client *iam.Client
}

// Implementation of EntityHandler for IAM Group.
type GroupHandler struct {
	Client *iam.Client
}

// Implementation of EntityHandler for IAM Role.
type RoleHandler struct {
	Client *iam.Client
}

// Implementation of EntityHandler for IAM Policy.
type PolicyHandler struct {
	Client *iam.Client
}

// Implementation of EntityHandler for IAM Identity Providers.
type IDPHandler struct {
	Client *iam.Client
}

// Implementation of AttachedGroupPolicies.
type AttachedGroupPoliciesHandler struct {
	Client *iam.Client
}

// Implementation of AttachedRolePolicies.
type AttachedRolePoliciesHandler struct {
	Client *iam.Client
}

// Implementation of AttachedUserPolicies.
type AttachedUserPoliciesHandler struct {
	Client *iam.Client
}

// Implementation of GroupMembers.
type GroupMemberHandler struct {
	Client *iam.Client
}

var (
	// List + Get for IAM entities.
	_ EntityLister[types.User]   = (*UserHandler)(nil)
	_ EntityGetter[types.User]   = (*UserHandler)(nil)
	_ EntityLister[types.Group]  = (*GroupHandler)(nil)
	_ EntityGetter[types.Group]  = (*GroupHandler)(nil)
	_ EntityLister[types.Role]   = (*RoleHandler)(nil)
	_ EntityGetter[types.Role]   = (*RoleHandler)(nil)
	_ EntityLister[types.Policy] = (*PolicyHandler)(nil)
	_ EntityGetter[types.Policy] = (*PolicyHandler)(nil)

	// List for IAM entities.
	_ EntityLister[types.SAMLProviderListEntry] = (*IDPHandler)(nil)
	_ EntityLister[types.AttachedPolicy]        = (*AttachedGroupPoliciesHandler)(nil)
	_ EntityLister[types.AttachedPolicy]        = (*AttachedRolePoliciesHandler)(nil)
	_ EntityLister[types.AttachedPolicy]        = (*AttachedUserPoliciesHandler)(nil)
	_ EntityLister[types.User]                  = (*GroupMemberHandler)(nil)
)

func (h *UserHandler) List(ctx context.Context, opts *Options,
) ([]types.User, *string, error) {
	output, err := h.Client.ListUsers(ctx, &iam.ListUsersInput{
		MaxItems:   opts.MaxItems,
		PathPrefix: opts.PathPrefix,
		Marker:     opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Users, output.Marker, nil
}

func (h *UserHandler) Get(ctx context.Context, user types.User) (types.User, error) {
	output, err := h.Client.GetUser(ctx, &iam.GetUserInput{
		UserName: user.UserName,
	})
	if err != nil {
		return types.User{}, err
	}

	return *output.User, nil
}

func (h *GroupHandler) List(ctx context.Context, opts *Options,
) ([]types.Group, *string, error) {
	output, err := h.Client.ListGroups(ctx, &iam.ListGroupsInput{
		MaxItems:   opts.MaxItems,
		PathPrefix: opts.PathPrefix,
		Marker:     opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Groups, output.Marker, nil
}

func (h *GroupHandler) Get(ctx context.Context, group types.Group,
) (types.Group, error) {
	output, err := h.Client.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: group.GroupName,
	})
	if err != nil {
		return types.Group{}, err
	}

	return *output.Group, nil
}

func (h *RoleHandler) List(ctx context.Context, opts *Options,
) ([]types.Role, *string, error) {
	output, err := h.Client.ListRoles(ctx, &iam.ListRolesInput{
		MaxItems:   opts.MaxItems,
		PathPrefix: opts.PathPrefix,
		Marker:     opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Roles, output.Marker, nil
}

func (h *RoleHandler) Get(ctx context.Context, role types.Role) (types.Role, error) {
	output, err := h.Client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: role.RoleName,
	})
	if err != nil {
		return types.Role{}, err
	}

	return *output.Role, nil
}

func (h *PolicyHandler) List(ctx context.Context, opts *Options,
) ([]types.Policy, *string, error) {
	output, err := h.Client.ListPolicies(ctx, &iam.ListPoliciesInput{
		MaxItems:   opts.MaxItems,
		PathPrefix: opts.PathPrefix,
		Marker:     opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Policies, output.Marker, nil
}

func (h *PolicyHandler) Get(ctx context.Context, policy types.Policy,
) (types.Policy, error) {
	output, err := h.Client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: policy.Arn,
	})
	if err != nil {
		return types.Policy{}, err
	}

	return *output.Policy, nil
}

func (h *IDPHandler) List(ctx context.Context, _ *Options,
) ([]types.SAMLProviderListEntry, *string, error) {
	output, err := h.Client.ListSAMLProviders(ctx, &iam.ListSAMLProvidersInput{})
	if err != nil {
		return nil, nil, err
	}

	return output.SAMLProviderList, nil, nil
}

func (h *AttachedGroupPoliciesHandler) List(ctx context.Context, opts *Options,
) ([]types.AttachedPolicy, *string, error) {
	output, err := h.Client.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
		GroupName: opts.UniqueName,
		MaxItems:  opts.MaxItems,
		Marker:    opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.AttachedPolicies, output.Marker, nil
}

func (h *AttachedRolePoliciesHandler) List(ctx context.Context, opts *Options,
) ([]types.AttachedPolicy, *string, error) {
	output, err := h.Client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: opts.UniqueName,
		MaxItems: opts.MaxItems,
		Marker:   opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.AttachedPolicies, output.Marker, nil
}

func (h *AttachedUserPoliciesHandler) List(ctx context.Context, opts *Options,
) ([]types.AttachedPolicy, *string, error) {
	output, err := h.Client.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
		UserName: opts.UniqueName,
		MaxItems: opts.MaxItems,
		Marker:   opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.AttachedPolicies, output.Marker, nil
}

func (h *GroupMemberHandler) List(ctx context.Context, opts *Options,
) ([]types.User, *string, error) {
	output, err := h.Client.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: opts.UniqueName,
		MaxItems:  opts.MaxItems,
		Marker:    opts.Marker,
	})
	if err != nil {
		return nil, nil, err
	}

	return output.Users, output.Marker, nil
}
