// Copyright 2025 SGNL.ai, Inc.

package ldap

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/go-objectsid"
	ldap_v3 "github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	grpc_proxy_v1 "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger"
	"github.com/sgnl-ai/adapters/pkg/logs/zaplogger/fields"
	"github.com/sgnl-ai/adapters/pkg/pagination"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// Proxy is an interface for LDAP proxy requests.
// It is used to send LDAP requests to a remote connector via the SGNL connector proxy.
type Proxy interface {
	ProxyRequest(ctx context.Context, ci *connector.ConnectorInfo, request *Request) (*Response, *framework.Error)
}

// Requester is an interface for LDAP requests.
// It is used to send LDAP requests directly to a publicly accessible LDAP server.
type Requester interface {
	Request(ctx context.Context, request *Request) (*Response, *framework.Error)
}

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct {
	Client Dispatcher
}

// Dispatcher is an interface that combines Proxy and Requester.
// It is used to determine if the LDAP request should be proxied or sent directly to the LDAP server.
// The IsProxied method checks if the LDAP request is proxied.
type Dispatcher interface {
	IsProxied() bool
	Proxy
	Requester
}

// NewLDAPRequester creates a new LDAP Requester instance.
// It is used to create a new LDAP client for making LDAP search requests.
// It also manages a session pool to reuse LDAP connections.
func NewLDAPRequester(ttl time.Duration, cleanupInterval time.Duration) Requester {
	client := &ldapClient{
		sessionPool: NewSessionPool(ttl, cleanupInterval),
	}

	return client
}

// ldapClient for making LDAP search request directly to a publicly accessible LDAP server,
// or for using a proxy client to access an on-premises LDAP server via the on-premises SGNL
// connector.
// It also manages a session pool to reuse LDAP connections.
type ldapClient struct {
	proxyClient grpc_proxy_v1.ProxyServiceClient
	sessionPool *SessionPool
}

func (c *ldapClient) IsProxied() bool {
	return c.proxyClient != nil
}

// ProxyRequest proxies an LDAP adapter's request to a remote connector by sending a
// LdapSearchRequest containing serialized request data.
// The connector responds with the serialized, processed response.
func (c *ldapClient) ProxyRequest(
	ctx context.Context, ci *connector.ConnectorInfo, request *Request,
) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
		fields.BaseURL(request.BaseURL),
		fields.ConnectorID(ci.ID),
		fields.ConnectorSourceID(ci.SourceID),
		fields.ConnectorSourceType(int(ci.SourceType)),
	)

	logger.Info("Sending request to datasource")

	data, err := json.Marshal(request)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to proxy LDAP request: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	r := &grpc_proxy_v1.ProxyRequestMessage{
		ConnectorId: ci.ID,
		ClientId:    ci.ClientID,
		TenantId:    ci.TenantID,
		Request: &grpc_proxy_v1.Request{
			RequestType: &grpc_proxy_v1.Request_LdapSearchRequest{
				LdapSearchRequest: &grpc_proxy_v1.LDAPSearchRequest{
					Request: string(data),
				},
			},
		},
	}

	response := &Response{}

	proxyResp, err := c.proxyClient.ProxyRequest(ctx, r)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		if st, ok := status.FromError(err); ok {
			response.StatusCode = customerror.GRPCErrStatusToHTTPStatusCode(st, err)

			return response, nil
		}

		return nil, customerror.UpdateError(
			&framework.Error{
				Message: fmt.Sprintf("Error searching LDAP server: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
			customerror.WithRequestTimeoutMessage(
				err, request.RequestTimeoutSeconds,
			),
		)
	}

	ldapResp := proxyResp.GetLdapSearchResponse()

	if ldapResp == nil {
		return nil, &framework.Error{
			Message: "Error received nil response from the proxy",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	if ldapResp.Error != "" {
		var respErr framework.Error
		// Unmarshal the error response from the proxy.
		err = json.Unmarshal([]byte(ldapResp.Error), &respErr)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Error unmarshalling LDAP error response from the proxy: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		return nil, &respErr
	}

	if ldapResp.Response == "" {
		return nil, &framework.Error{
			Message: "Error received empty response from the proxy",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Unmarshal the response from the proxy.
	if err = json.Unmarshal([]byte(ldapResp.Response), response); err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Error unmarshalling LDAP response from the proxy: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

// Request sends a paginated search query directly to an LDAP server and returns the
// processed response.
func (c *ldapClient) Request(ctx context.Context, request *Request) (*Response, *framework.Error) {
	logger := zaplogger.FromContext(ctx).With(
		fields.RequestEntityExternalID(request.EntityExternalID),
		fields.RequestPageSize(request.PageSize),
	)

	logger.Info("Sending request to datasource")

	tlsConfig, configErr := GetTLSConfig(request)
	if configErr != nil {
		return nil, configErr
	}

	// Determine the paging cookie (if any)
	var cookie []byte

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		pageInfo, decodeErr := DecodePageInfo(request.Cursor.Cursor)
		if decodeErr != nil {
			return nil, decodeErr
		}

		if pageInfo.NextPageCursor != nil {
			var err error

			cookie, err = OctetStringToBytes(*pageInfo.NextPageCursor)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse cursor value: %v.", err),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}
		}
	}

	// get existing session or create a new one
	address := request.BaseURL
	key := sessionKey(address, cookie)

	session, found := c.sessionPool.Get(key)
	if !found {
		session = &Session{}
		c.sessionPool.Set(key, session)
	}

	conn, err := session.GetOrCreateConn(request.BaseURL, tlsConfig, request.BindDN, request.BindPassword)
	if err != nil {
		if ldapErr, ok := err.(*ldap_v3.Error); ok && ldapErr.ResultCode == 49 {
			return nil, &framework.Error{
				Message: "Failed to bind credentials: LDAP Result Code 49 \"Invalid Credentials\": .",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
			}
		}

		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to dial/bind LDAP server: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Set cursor from the request (if exists)
	pageControl, pageErr := setPageControl(request)
	if pageErr != nil {
		return nil, pageErr
	}

	filters, filterErr := SetFilters(request)
	if filterErr != nil {
		return nil, filterErr
	}

	attributes := make([]string, 0, len(request.Attributes))
	for _, attr := range request.Attributes {
		attributes = append(attributes, attr.ExternalId)
	}

	// Define LDAP search with filtering, attributes and paging.
	searchRequest := ldap_v3.NewSearchRequest(
		request.BaseDN,                 // BaseDN
		ldap_v3.ScopeWholeSubtree,      // Scope
		ldap_v3.DerefAlways,            // DeferAliases
		0,                              // SizeLimit
		request.RequestTimeoutSeconds,  // TimeLimit
		false,                          // TypesOnly
		filters,                        // Filters
		attributes,                     // Attributes
		[]ldap_v3.Control{pageControl}, // Controls
	)

	// Perform search
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		logger.Error("Request to datasource failed",
			fields.SGNLEventTypeError(),
			zap.Error(err),
		)

		// Extract LDAP result code from the error
		if ldapErr, ok := err.(*ldap_v3.Error); ok {
			return &Response{
				StatusCode: ldapErrToHTTPStatusCode(ldapErr),
			}, nil
		}

		return nil, customerror.UpdateError(
			&framework.Error{
				Message: fmt.Sprintf("Error searching LDAP server: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
			customerror.WithRequestTimeoutMessage(
				err, request.RequestTimeoutSeconds,
			),
		)
	}

	isEmptyResult := searchResult == nil || len(searchResult.Entries) == 0
	isEmptyCursor := request.Cursor == nil || request.Cursor.CollectionID == nil

	if isEmptyResult && isEmptyCursor {
		return &Response{
			StatusCode: http.StatusNotFound,
		}, nil
	}

	response, ferr := ProcessLDAPSearchResult(searchResult, request)
	if ferr != nil {
		return nil, ferr
	}

	// Handle paging: store or cleanup session
	cookie = cookie[:0]

	if response.NextCursor != nil && response.NextCursor.Cursor != nil {
		pageInfo, decodeErr := DecodePageInfo(response.NextCursor.Cursor)
		if decodeErr == nil && pageInfo.NextPageCursor != nil {
			var err error

			cookie, err = OctetStringToBytes(*pageInfo.NextPageCursor)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to encode next page cursor received from LDAP server: %v.", err),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}
		}
	}

	// If we have a cookie, update the session with the new cookie for the next request.
	if len(cookie) > 0 {
		// Move the session to the new cookie key: remove old key, add new key (without closing conn)
		nextKey := sessionKey(address, cookie)
		c.sessionPool.UpdateKey(key, nextKey)
	} else {
		// Paging done, cleanup
		c.sessionPool.Delete(key)
	}

	logger.Info("Datasource request completed successfully",
		fields.ResponseStatusCode(response.StatusCode),
		fields.ResponseObjectCount(len(response.Objects)),
		fields.ResponseNextCursor(response.NextCursor),
	)

	return response, nil
}

// GetTLSConfig creates a TLS config using certchain from the request.
func GetTLSConfig(request *Request) (*tls.Config, *framework.Error) {
	if !request.IsLDAPS {
		return &tls.Config{}, nil
	}

	decodedCertChain, err := base64.StdEncoding.DecodeString(request.CertificateChain)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to load certificates: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(decodedCertChain)

	// Use url.Hostname() for more reliable hostname extraction, especially for IPv6
	u, err := url.Parse(request.BaseURL)
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to parse URL: %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	return &tls.Config{
		RootCAs:    caCertPool,
		ServerName: u.Hostname(),
	}, nil
}

// setPageControl for setting up the page control cookie based upon the
// cursor value in the request.
func setPageControl(request *Request) (*ldap_v3.ControlPaging, *framework.Error) {
	// Set cursor from the request (if exists)
	pageControl := ldap_v3.NewControlPaging(uint32(request.PageSize))

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		pageInfo, decodeErr := DecodePageInfo(request.Cursor.Cursor)
		if decodeErr != nil {
			return nil, decodeErr
		}

		if pageInfo.NextPageCursor != nil {
			cookie, err := OctetStringToBytes(*pageInfo.NextPageCursor)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse cursor value: %v.", err),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			pageControl.SetCookie(cookie)
		}
	}

	return pageControl, nil
}

type PageInfo struct {
	// Collection is a map of the attributes of the collection entity.
	Collection map[string]any `json:"collection"`

	// NextPageCursor is the cursor to the next page of results.
	NextPageCursor *string `json:"nextPageCursor"`

	// NextGroupProcessed tracks the next group DN (for resuming)
	NextGroupProcessed string `json:"nextGroupProcessed,omitempty"`

	// NextMemberProcessed tracks the next member's offset within
	// the list of members for a group (for resuming)
	NextMemberProcessed int64 `json:"nextMemberProcessed,omitempty"`

	// RangeAttribute is the attribute used for range queries.
	RangeAttribute bool `json:"rangeAttribute,omitempty"`
}

func getPageInfo(req *Request) (*PageInfo, *framework.Error) {
	if req.Cursor != nil && req.Cursor.Cursor != nil {
		pageInfo, decodeErr := DecodePageInfo(req.Cursor.Cursor)
		if decodeErr == nil {
			return pageInfo, nil
		}

		return nil, decodeErr
	}

	return &PageInfo{}, nil
}

const (
	ErrorMsgAttributeTypeDoesNotMatchFmt = "Attribute '%s' was returned from the " +
		"configured datasource as type %s; wanted type %s"

	// This attribute specifies the unique identifier for an object.
	//
	// see: https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-ada3/937eb5c6-f6b3-4652-a276-5d6bb8979658
	objectGUID = "objectGUID"

	// Syntax: String(Sid)
	// An octet string that contains a security identifier (SID).
	//
	// see: https://learn.microsoft.com/en-us/windows/win32/adschema/s-string-sid
	objectSid          = "objectSid"
	sidHistory         = "SIDHistory"
	creatorSID         = "mS-DS-CreatorSID"
	securityIdentifier = "securityIdentifier"

	// Syntax: String(NT-Sec-Desc)
	// An octet string that contains a Windows NT or Windows 2000 security descriptor.
	//
	// See: https://learn.microsoft.com/en-us/windows/win32/adschema/s-string-nt-sec-desc
	nTSecurityDescriptor                    = "nTSecurityDescriptor"
	msDSAllowedToActOnBehalfOfOtherIdentity = "msDS-AllowedToActOnBehalfOfOtherIdentity"
	fRSRootSecurity                         = "fRSRootSecurity"
	pKIEnrollmentAccess                     = "pKIEnrollmentAccess"
	msDSGroupMSAMembership                  = "msDS-GroupMSAMembership"
	msDFSLinkSecurityDescriptorv2           = "msDFS-LinkSecurityDescriptorv2"
)

func (d *Datasource) getPage(ctx context.Context, request *Request, memberOf *string) (*Response, *framework.Error) {
	if validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		memberOf != nil,
	); validationErr != nil {
		return nil, validationErr
	}

	// Make sure if the connector context is set and client can proxy the request.
	if d.Client.IsProxied() {
		if ci, ok := connector.FromContext(ctx); ok {
			return d.Client.ProxyRequest(ctx, &ci, request)
		}
	}

	resp, err := d.Client.Request(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		// An adapter error message is generated if the response status code from the
		// collection API is not successful (i.e. if not statusCode >= 200 && statusCode < 300).
		if adapterErr := web.HTTPError(resp.StatusCode, resp.RetryAfterHeader); adapterErr != nil {
			return nil, adapterErr
		}
	}

	return resp, nil
}

// getMemberOfPage handles the special case of fetching members of groups.
// This involves first fetching groups in configured batch size, then for each group, fetching its members.
// The process continues until the requested page size is fulfilled or there are no more groups/members to
// process.
// In case of a group with large number of members, range queries are used to fetch members in chunks.
// The function also manages pagination state using a composite cursor that tracks both group and member
// processing state.
// Next group to process, next member offset within the group, and whether range queries are being used
// are all tracked as part of the pagination state.
func (d *Datasource) getMemberOfPage(
	ctx context.Context, request *Request, entityConfig *EntityConfig,
) (*Response, *framework.Error) {
	// This is the top-level cursor/pageinfo - only gets updated when the
	// entityConfig.MemberOfBatchSize is processed completely or entities equal to
	// page size are processed.
	pageInfo, ferr := getPageInfo(request)
	if ferr != nil {
		return nil, ferr
	}

	collectionAttribute := entityConfig.CollectionAttribute
	memberOfUniqueIDAttribute := entityConfig.MemberOfUniqueIDAttribute
	memberOfAttributes := []*framework.AttributeConfig{
		{
			ExternalId: *memberOfUniqueIDAttribute,
			Type:       framework.AttributeTypeString,
			UniqueId:   true,
		},
	}

	// For a member's entityConfig.Query, {{CollectionID}} is typically replaced by collectionAttribute
	// to filter the member entity. The collectionAttribute and memberOfUniqueIdAttribute can be identical.
	// If they are not identical, add the collectionAttribute request.
	if collectionAttribute != nil && *collectionAttribute != *memberOfUniqueIDAttribute {
		memberOfAttributes = append(memberOfAttributes, &framework.AttributeConfig{
			ExternalId: *collectionAttribute,
			Type:       framework.AttributeTypeString,
		})
	}

	memberOfReq := &Request{
		BaseURL:           request.BaseURL,
		Attributes:        memberOfAttributes,
		ConnectionParams:  request.ConnectionParams,
		UniqueIDAttribute: *memberOfUniqueIDAttribute,
		EntityExternalID:  *entityConfig.MemberOf,
		EntityConfigMap: map[string]*EntityConfig{
			*entityConfig.MemberOf: request.EntityConfigMap[*entityConfig.MemberOf],
		},
		PageSize:              entityConfig.MemberOfGroupBatchSize,
		RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		Cursor:                request.Cursor,
	}

	memberOfResp, err := d.getPage(ctx, memberOfReq, nil)
	if err != nil {
		return nil, err
	}

	// Process groups one by one, making individual range queries for each group
	var (
		allGroupMembers      []map[string]any
		targetPageSize       = int(request.PageSize)
		nextGroupProcessedID string
		nextMemberProcessed  int64
		isRangeAttribute     = false
		skipToGroup          = pageInfo != nil && pageInfo.NextGroupProcessed != ""
	)

	//nolint:nestif
	if len(memberOfResp.Objects) > 0 {
		for _, groupObj := range memberOfResp.Objects {
			// Get the group's unique ID value
			groupUniqueIDValue, ok := groupObj[*entityConfig.MemberOfUniqueIDAttribute].(string)
			if !ok {
				continue
			}

			// Skip groups until we reach the last processed group
			if skipToGroup {
				if groupUniqueIDValue != pageInfo.NextGroupProcessed {
					continue // Skip this group, haven't reached the current one yet
				}

				skipToGroup = false // Found the group, stop skipping
			}

			// Calculate how many members we still need
			remainingNeeded := targetPageSize - len(allGroupMembers)
			if remainingNeeded <= 0 {
				nextGroupProcessedID = groupUniqueIDValue

				break // We have enough members, don't process this group
			}

			// Calculate member offset within this specific group
			memberOffsetInGroup := int64(0)
			if pageInfo != nil && groupUniqueIDValue == pageInfo.NextGroupProcessed {
				// We're resuming from this group, use the stored offset
				memberOffsetInGroup = pageInfo.NextMemberProcessed
				isRangeAttribute = pageInfo.RangeAttribute
			}

			memberOffsetEnd := memberOffsetInGroup + int64(remainingNeeded)

			var memberAttribute string
			if isRangeAttribute {
				memberAttribute = fmt.Sprintf("member;range=%d-%d", memberOffsetInGroup, memberOffsetEnd)
			} else if entityConfig.MemberAttribute != nil {
				memberAttribute = fmt.Sprintf("%s", *entityConfig.MemberAttribute)
			} else {
				memberAttribute = defaultMemberAttribute
			}

			// Create a Group request for this specific group
			// Create a modified entity config that filters for this specific group
			modifiedGroupMemberQuery := strings.ReplaceAll(
				entityConfig.Query, "{{CollectionAttribute}}", *entityConfig.CollectionAttribute,
			)
			modifiedGroupMemberQuery = strings.ReplaceAll(modifiedGroupMemberQuery, "{{CollectionId}}", groupUniqueIDValue)

			modifiedEntityConfigMap := make(map[string]*EntityConfig)

			for k, v := range request.EntityConfigMap {
				entityConfigCopy := *v
				modifiedEntityConfigMap[k] = &entityConfigCopy
			}

			modifiedEntityConfigMap[request.EntityExternalID].Query = modifiedGroupMemberQuery

			memberRequest := &Request{
				ConnectionParams:      request.ConnectionParams,
				BaseURL:               request.BaseURL,
				PageSize:              1,                        // Only get this specific group
				EntityExternalID:      request.EntityExternalID, // Group entity, not GroupMember
				UniqueIDAttribute:     request.UniqueIDAttribute,
				EntityConfigMap:       modifiedEntityConfigMap,
				RequestTimeoutSeconds: request.RequestTimeoutSeconds,
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: *entityConfig.MemberOfUniqueIDAttribute,
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: memberAttribute, // Request the specific range query
						Type:       framework.AttributeTypeString,
						List:       true,
					},
				},
				Cursor: &pagination.CompositeCursor[string]{
					CollectionID: &groupUniqueIDValue,
				},
			}

			memberResp, err := d.getPage(ctx, memberRequest, entityConfig.MemberOf)
			if err != nil {
				return nil, err
			}

			if len(memberResp.Objects) > 0 {
				memberObj := memberResp.Objects[0]

				// Extract member DNs from this group
				membersDN, action := extractMembersDNFromGroup(memberObj, int(memberOffsetEnd)-int(memberOffsetInGroup))
				memberOffsetEnd = memberOffsetInGroup + int64(len(membersDN))

				// Convert member DNs to GroupMember objects
				objects := make([]map[string]any, len(membersDN))
				// id
				memberOfUniqueIDAttribute := fmt.Sprintf("group_%s", *entityConfig.MemberOfUniqueIDAttribute)

				// member_<memberUniqueIDAttribute>
				memberUniqueIDAttribute := fmt.Sprintf("member_%s", *entityConfig.MemberUniqueIDAttribute)

				for idx, memberDN := range membersDN {
					objects[idx] = map[string]any{
						request.UniqueIDAttribute: fmt.Sprintf("%s-%s", memberDN, groupUniqueIDValue),
						memberUniqueIDAttribute:   memberDN,
						memberOfUniqueIDAttribute: groupUniqueIDValue,
					}
				}

				allGroupMembers = append(allGroupMembers, objects...)

				remainingNeeded := targetPageSize - len(allGroupMembers)
				if remainingNeeded <= 0 && action != MemberExtractionActionDone {
					// Handle pagination for more members
					nextGroupProcessedID = groupUniqueIDValue
					nextMemberProcessed = memberOffsetEnd

					if action == MemberExtractionActionContinueRange {
						isRangeAttribute = true
					}

					break // We have enough members, don't process this group
				}

				nextGroupProcessedID = ""
				nextMemberProcessed = 0
			}
		}
	}

	// Create response
	response := &Response{
		StatusCode: http.StatusOK,
		Objects:    allGroupMembers,
	}

	// Handle pagination cursor if needed
	//nolint:nestif
	if memberOfResp.NextCursor != nil || nextGroupProcessedID != "" {
		// If we have more groups to process or hit the page limit, create/update cursor
		var nextPageInfo *PageInfo

		if memberOfResp.NextCursor != nil && memberOfResp.NextCursor.Cursor != nil {
			// Decode existing cursor
			if decodedPageInfo, decodeErr := DecodePageInfo(memberOfResp.NextCursor.Cursor); decodeErr == nil {
				nextPageInfo = decodedPageInfo
			}
		}

		if nextPageInfo == nil {
			nextPageInfo = &PageInfo{}
		}

		nextPageInfo.NextGroupProcessed = nextGroupProcessedID
		nextPageInfo.NextMemberProcessed = nextMemberProcessed

		nextPageInfo.RangeAttribute = isRangeAttribute

		if nextGroupProcessedID != "" {
			nextPageInfo.NextPageCursor = nil
			if pageInfo != nil {
				nextPageInfo.NextPageCursor = pageInfo.NextPageCursor
			}
		}

		// Encode updated cursor
		if b, marshalErr := json.Marshal(nextPageInfo); marshalErr == nil {
			encodedCursor := base64.StdEncoding.EncodeToString(b)
			response.NextCursor = &pagination.CompositeCursor[string]{
				Cursor: &encodedCursor,
			}
		} else {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", marshalErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}
	}

	return response, nil
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	entityConfig := request.EntityConfigMap[request.EntityExternalID]

	if entityConfig.MemberOf != nil {
		return d.getMemberOfPage(ctx, request, entityConfig)
	}

	return d.getPage(ctx, request, entityConfig.MemberOf)
}

type MemberExtractionAction string

const (
	MemberExtractionActionContinueRange   MemberExtractionAction = "ContinueRange"
	MemberExtractionActionContinueRegular MemberExtractionAction = "ContinueRegular"
	MemberExtractionActionDone            MemberExtractionAction = "Done"
)

// extractMembersDNFromGroup extracts member DNs from a group's member attributes
// Returns the DNs, whether range queries are needed, and any error.
func extractMembersDNFromGroup(memberObjs map[string]any, count int) ([]string, MemberExtractionAction) {
	var membersDN []string

	// Check for range attributes first (indicates large groups)
	for attrName, value := range memberObjs {
		memberList, ok := value.([]any)
		if !ok {
			continue
		}

		if count > len(memberList) {
			count = len(memberList)
		}

		for _, member := range memberList[:count] {
			if memberDN, ok := member.(string); ok {
				membersDN = append(membersDN, memberDN)
			} else {
				membersDN = membersDN[:0] // Clear the slice

				break
			}
		}

		if len(membersDN) == 0 {
			continue // Try next attribute
		}

		// Handle range attribute (e.g., member;range=0-1499)
		// Regex pattern: attrname;range=start-end or attrname;range=start-*
		// Example matches: member;range=0-1499, member;range=1500-*
		// Ref: https://learn.microsoft.com/en-us/windows/win32/adschema/attributes/member
		// Ref: https://learn.microsoft.com/en-us/windows/win32/adschema/attributes/memberof

		if matches := rangeAttributePattern.FindStringSubmatch(attrName); matches != nil {
			if matches[3] == "*" && count == len(memberList) {
				return membersDN, MemberExtractionActionDone
			}

			return membersDN, MemberExtractionActionContinueRange
		}

		if count == len(memberList) {
			return membersDN, MemberExtractionActionDone
		}

		return membersDN, MemberExtractionActionContinueRegular
	}

	return []string{}, MemberExtractionActionDone
}

func ProcessLDAPSearchResult(result *ldap_v3.SearchResult, request *Request) (*Response, *framework.Error) {
	objects, pageInfo, frameworkErr := ParseResponse(result, attrIDToConfig(request.Attributes))
	if frameworkErr != nil {
		return nil, frameworkErr
	}

	// Indicating a successful LDAP search operation.
	// In case of no error (err == nil), ldap_v3.Search is considered successful,
	// returning LDAP Result Code Success(0) equivalent to HTTP status code StatusOK.
	response := &Response{
		StatusCode: http.StatusOK,
		Objects:    objects,
	}

	if pageInfo != nil && pageInfo.NextPageCursor != nil {
		b, marshalErr := json.Marshal(pageInfo)
		if marshalErr != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", marshalErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		encodedCursor := base64.StdEncoding.EncodeToString(b)
		response.NextCursor = &pagination.CompositeCursor[string]{
			Cursor: &encodedCursor,
		}
	}

	return response, nil
}

func ParseResponse(searchResult *ldap_v3.SearchResult, attributes map[string]*framework.AttributeConfig) (
	objects []map[string]any, pageInfo *PageInfo, err *framework.Error) {
	objects = make([]map[string]any, 0, len(searchResult.Entries))

	for _, entry := range searchResult.Entries {
		if entry != nil {
			object, err := EntryToObject(entry, attributes)
			if err != nil {
				return nil, nil, err
			}

			objects = append(objects, object)
		}
	}

	if len(searchResult.Controls) != 0 {
		// Update Cursor
		pagingControl := ldap_v3.FindControl(searchResult.Controls, ldap_v3.ControlTypePaging)

		if ctrl, ok := pagingControl.(*ldap_v3.ControlPaging); ok && ctrl != nil && len(ctrl.Cookie) != 0 {
			// The Control.cookie is represented as an octet string.
			// Ref: https://www.ietf.org/rfc/rfc2696.txt.
			//
			// During the adapter flow, we transform the cookie into different types based on the specific
			// requirements at each level. To maintain the integrity of the octet string, we must
			// cautiously convert the []byte into a string, as any missing character could lead to unexpected query results.
			pageInfo = &PageInfo{
				NextPageCursor: BytesToOctetString(ctrl.Cookie),
			}

			return objects, pageInfo, nil
		}
	}

	return objects, nil, nil
}

func BytesToOctetString(data []byte) *string {
	octetStr := base64.StdEncoding.EncodeToString(data)

	return &octetStr
}

func OctetStringToBytes(octalString string) ([]byte, error) {
	octetStr, err := base64.StdEncoding.DecodeString(octalString)
	if err != nil {
		return nil, err
	}

	return octetStr, nil
}

const (
	rangeAttributePrefix = "member;range="
	rangeAttrList        = true
	rangeAttrType        = framework.AttributeTypeString
)

var (
	// Regex pattern for matching LDAP range attributes: attrname;range=start-end or attrname;range=start-*.
	rangeAttributePattern = regexp.MustCompile(`^(member);range=(\d+)-(\d+|\*)$`)
)

func EntryToObject(e *ldap_v3.Entry,
	attrConfig map[string]*framework.AttributeConfig,
) (map[string]any, *framework.Error) {
	result := make(map[string]any)
	result["dn"] = e.DN

	// Iterate over each attribute in the LDAP entry.
	// If the attribute is not in the config, it is skipped.
	// If the attribute is a range attribute, it is treated as a list of strings.
	// - Range attributes are used for large multi-valued attributes that exceed the LDAP server's limit.
	// If the attribute type conversion fails, an error is returned.
	for _, attribute := range e.Attributes {
		// Skip attributes that are not configured in the attribute config (could be due to user error)
		currAttrConfig, ok := attrConfig[attribute.Name]

		isRangeAttribute := strings.HasPrefix(attribute.Name, rangeAttributePrefix)
		if !ok && !isRangeAttribute {
			continue
		}

		var (
			isList   bool
			attrType framework.AttributeType
		)

		if isRangeAttribute {
			isList = rangeAttrList
			attrType = rangeAttrType
		}

		// Use the attribute config if it exists
		if ok {
			isList = currAttrConfig.List
			attrType = currAttrConfig.Type
		}

		// Convert string values to the requested type
		value, err := StringAttrValuesToRequestedType(attribute, isList, attrType)
		if err != nil {
			return nil, err
		}

		result[attribute.Name] = value
	}

	return result, nil
}

func getAttrType(attrType api_adapter_v1.AttributeType) framework.AttributeType {
	return framework.AttributeType(*api_adapter_v1.AttributeType.Enum(attrType))
}

func attrIDToConfig(attrConfig []*framework.AttributeConfig) map[string]*framework.AttributeConfig {
	result := make(map[string]*framework.AttributeConfig, len(attrConfig))
	for _, config := range attrConfig {
		result[config.ExternalId] = config
	}

	return result
}

func StringAttrValuesToRequestedType(
	attr *ldap_v3.EntryAttribute,
	isList bool,
	attrType framework.AttributeType,
) (any, *framework.Error) {
	if isList {
		if len(attr.Values) == 0 { // empty values
			return attr.Values, nil
		}

		values := make([]any, 0, len(attr.Values))

		for _, v := range attr.Values {
			listAttr := &ldap_v3.EntryAttribute{
				Name:   attr.Name,
				Values: []string{v},
			}

			value, err := StringAttrValuesToRequestedType(listAttr, false, attrType)
			if err != nil {
				return nil, err
			}

			values = append(values, value)
		}

		return values, nil
	}

	// Return an empty string in case of no values
	if len(attr.Values) == 0 {
		return "", nil
	}

	switch attrType {
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING):
		switch attr.Name {
		// Special AD syntaxes.
		case objectGUID:
			if len(attr.ByteValues) == 0 || attr.ByteValues[0] == nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Missing or nil GUID bytes for attribute: %s", attr.Name),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

			guid, err := uuid.Parse(hex.EncodeToString(attr.ByteValues[0]))
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf(
						ErrorMsgAttributeTypeDoesNotMatchFmt,
						attr.Name,
						reflect.TypeOf(attr.Values[0]),
						"string",
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

			return guid.String(), nil
		case objectSid, sidHistory, creatorSID, securityIdentifier:
			if len(attr.ByteValues) == 0 || attr.ByteValues[0] == nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Missing or nil SID bytes for attribute: %s", attr.Name),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

			sid := objectsid.Decode(attr.ByteValues[0])

			return sid.String(), nil
		case nTSecurityDescriptor, msDSAllowedToActOnBehalfOfOtherIdentity, fRSRootSecurity, pKIEnrollmentAccess,
			msDSGroupMSAMembership, msDFSLinkSecurityDescriptorv2:
			if len(attr.ByteValues) == 0 || attr.ByteValues[0] == nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Missing or nil security descriptor bytes for attribute: %s", attr.Name),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}
			// Convert the security descriptor bytes to base64 string
			sddl := base64.StdEncoding.EncodeToString(attr.ByteValues[0])

			return sddl, nil
		default:
			return attr.Values[0], nil
		}
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME):
		return attr.Values[0], nil
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL):
		value, err := strconv.ParseBool(attr.Values[0])
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(ErrorMsgAttributeTypeDoesNotMatchFmt,
					attr.Name, reflect.TypeOf(attr.Values[0]), "boolean"),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			}
		}

		return value, nil
	// TODO: optimize case of DOUBLE & INT64 type
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DOUBLE):
		value, err := strconv.ParseFloat(attr.Values[0], 64)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(ErrorMsgAttributeTypeDoesNotMatchFmt,
					attr.Name, reflect.TypeOf(attr.Values[0]), "float64"),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			}
		}

		return value, nil
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64):
		// All numbers are unmarshalled into float64. Further framework converts into int64.
		value, err := strconv.ParseFloat(attr.Values[0], 64)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(ErrorMsgAttributeTypeDoesNotMatchFmt,
					attr.Name, reflect.TypeOf(attr.Values[0]), "int64"),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			}
		}

		return value, nil
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DURATION):
		value, err := strconv.ParseInt(attr.Values[0], 10, 64)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf(ErrorMsgAttributeTypeDoesNotMatchFmt,
					attr.Name, reflect.TypeOf(attr.Values[0]), "int64"),
				Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
			}
		}

		return value, nil
	default:
		return nil, &framework.Error{
			Message: "Unsupported type requested for Active Directory SoR",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
		}
	}
}

// See https://github.com/go-ldap/ldap/blob/e80f029a818542002267960a6b8dae32d79d0994/v3/error.go#L10-L93
// for all the LDAP releated error codes.
var (
	ldapToHTTPErrorCodes = map[uint16]int{
		2:   http.StatusBadRequest, // invalid page result cookie
		32:  http.StatusNotFound,
		34:  http.StatusBadRequest,
		33:  http.StatusBadRequest,
		81:  http.StatusServiceUnavailable,
		85:  http.StatusGatewayTimeout,
		49:  http.StatusUnauthorized,
		50:  http.StatusForbidden,
		12:  http.StatusNotImplemented,
		201: http.StatusBadRequest,
	}
)

func ldapErrToHTTPStatusCode(ldapError *ldap_v3.Error) int {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	if httpStatusCode, ok := ldapToHTTPErrorCodes[ldapError.ResultCode]; ok {
		return httpStatusCode
	}

	logger.Printf("Unknown LDAP result code received: %v \t %v\n", ldapError.ResultCode, ldapError.Err.Error())

	return http.StatusInternalServerError // default error code
}

func DecodePageInfo(cursor *string) (*PageInfo, *framework.Error) {
	b, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, &framework.Error{
			Message: "Cursor.Cursor base64 decoding failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	var pageInfo PageInfo

	err = json.Unmarshal(b, &pageInfo)
	if err != nil {
		return nil, &framework.Error{
			Message: "PageInfo unmarshalling failed.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return &pageInfo, nil
}

func sessionKey(address string, cookie []byte) string {
	cookieStr := base64.StdEncoding.EncodeToString(cookie)

	return address + "|" + cookieStr
}
