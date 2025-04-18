// Copyright 2025 SGNL.ai, Inc.

// TODO: The contents of this file are unused at the moment. Do not remove. We need it for future improvements.
package crowdstrike

import (
	"fmt"
	"sort"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

type FragmentType struct {
	Parent     string
	Name       string
	Attributes []string
}

// AttributeNode stores the metadata required to build the inner part of the query for an entity.
type AttributeNode struct {
	Name         string
	Children     map[string]*AttributeNode
	IsFragment   bool
	FragmentType *FragmentType
}

var fragmentTypes = map[string][]FragmentType{
	"accounts": {
		{
			Parent:     "accounts",
			Name:       "ActiveDirectoryAccountDescriptor",
			Attributes: []string{"enabled", "ou", "domain", "creationTime", "datasource", "upn"},
		},
		{
			Parent:     "accounts",
			Name:       "ActivityParticipatingAccountDescriptor",
			Attributes: []string{"enabled", "mostRecentActivity", "datasource"},
		},
		{
			Parent:     "accounts",
			Name:       "SsoUserAccountDescriptor",
			Attributes: []string{"enabled", "creationTime", "mostRecentActivity", "datasource"},
		},
	},
	"user": {
		{
			Parent:     "user",
			Name:       "UserEntity",
			Attributes: []string{"riskScoreSeverity", "mostRecentActivity", "riskScore", "riskScoreSeverity"},
		},
	},
}

// AddChild adds a child to the current node and returns the child node.
func (node *AttributeNode) AddChild(path []string, isFragment bool, fragmentType *FragmentType) *AttributeNode {
	if node.Children == nil {
		node.Children = make(map[string]*AttributeNode)
	}

	if len(path) == 0 {
		return node
	}

	childName := path[0]
	if _, exists := node.Children[childName]; !exists {
		node.Children[childName] = &AttributeNode{
			Name:         childName,
			Children:     make(map[string]*AttributeNode),
			IsFragment:   isFragment,
			FragmentType: fragmentType,
		}
	}

	return node.Children[childName].AddChild(path[1:], isFragment, fragmentType)
}

func GetAttributePath(input string) []string {
	if strings.HasPrefix(input, "$") {
		parts := strings.Split(input, ".")
		if len(parts) > 1 {
			return parts[1:]
		}
	}

	return []string{input}
}

func AttributeQueryBuilder(
	entityConfig *framework.EntityConfig,
	rootName string,
) (*AttributeNode, *framework.Error) {
	rootParts := GetAttributePath(rootName)
	if len(rootParts) == 0 {
		return nil, &framework.Error{
			Message: "Root name is empty.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	rootNode := &AttributeNode{
		Name:     rootParts[0],
		Children: make(map[string]*AttributeNode),
	}
	// baseNode and rootNode are different when the rootName is a JSON path. This happens when
	// child entities with JSON path externalIDs are handled.
	baseNode := rootNode.AddChild(rootParts[1:], false, nil)

	for _, attr := range entityConfig.Attributes {
		attrParts := GetAttributePath(attr.ExternalId)
		isFragment := strings.Contains(attr.ExternalId, "[")

		var fragmentType *FragmentType

		var ok bool

		if isFragment {
			parts := strings.Split(attrParts[0], "[")
			if len(parts) > 1 {
				fragmentType, ok = findFragmentContainingAttributeForKey(parts[0], attrParts[1])
				if !ok {
					return nil, &framework.Error{
						Message: fmt.Sprintf("Fragment type not found for attribute: %s", attr.ExternalId),
						Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
					}
				}
			}
		}

		baseNode.AddChild(attrParts, isFragment, fragmentType)
	}

	for _, child := range entityConfig.ChildEntities {
		childParts := GetAttributePath(child.ExternalId)
		if len(childParts) == 0 {
			return nil, &framework.Error{
				Message: "Child name is empty.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
			}
		}

		childNode, err := AttributeQueryBuilder(child, child.ExternalId)
		if err != nil {
			return nil, err
		}

		baseNode.Children[childParts[0]] = childNode
	}

	return rootNode, nil
}

func findFragmentContainingAttributeForKey(key, attribute string) (*FragmentType, bool) {
	fragments, exists := fragmentTypes[key]
	if !exists {
		return nil, false
	}

	for _, fragment := range fragments {
		if contains(fragment.Attributes, attribute) {
			return &fragment, true
		}
	}

	return nil, false
}

// TODO: Remove this and replaces usages with slices.Contains.
func contains(attributes []string, attribute string) bool {
	for _, attr := range attributes {
		if attr == attribute {
			return true
		}
	}

	return false
}

func (node *AttributeNode) BuildQuery() string {
	if len(node.Children) == 0 {
		return node.Name
	}

	childrenQueries := make([]string, 0, len(node.Children))
	keys := make([]string, 0, len(node.Children))

	for key := range node.Children {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		child := node.Children[key]
		childrenQueries = append(childrenQueries, child.BuildQuery())
	}

	if node.IsFragment {
		return fmt.Sprintf("... on %s { %s }", node.FragmentType, strings.Join(childrenQueries, ", "))
	}

	return fmt.Sprintf("%s { %s }", node.Name, strings.Join(childrenQueries, ", "))
}

// This returns the PageInfo struct of the n deep layer, where 0 is the outermost layer.
// If n > number of layers or the pageInfo is nil, the function returns nil.
// If n < 0, the function returns the PageInfo of the outermost layer.
func GetPageInfoAfter(pageInfo *PageInfo, n int) *string {
	if pageInfo == nil {
		return nil
	}

	if n <= 0 {
		return &pageInfo.EndCursor
	}

	return GetPageInfoAfter(pageInfo.InnerPageInfo, n-1)
}
