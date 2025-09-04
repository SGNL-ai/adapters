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
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/pagination"
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

	return response, nil
}

// Request sends a paginated search query directly to an LDAP server and returns the
// processed response.
func (c *ldapClient) Request(_ context.Context, request *Request) (*Response, *framework.Error) {
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

	response := &Response{}

	// Perform search
	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		// Extract LDAP result code from the error
		if ldapErr, ok := err.(*ldap_v3.Error); ok {
			statusCode := ldapErrToHTTPStatusCode(ldapErr)
			response.StatusCode = statusCode

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

	isEmptyResult := searchResult == nil || len(searchResult.Entries) == 0
	isEmptyCursor := request.Cursor == nil || request.Cursor.CollectionID == nil

	if isEmptyResult && isEmptyCursor {
		response.StatusCode = http.StatusNotFound

		return response, nil
	}

	// Indicating a successful LDAP search operation.
	// In case of no error (err == nil), ldap_v3.Search is considered successful,
	// returning LDAP Result Code Success(0) equivalent to HTTP status code StatusOK.
	response.StatusCode = http.StatusOK

	if err := ProcessLDAPSearchResult(searchResult, response, request); err != nil {
		return nil, err
	}

	// Handle paging: store or cleanup session
	var nextCookie []byte

	if response.NextCursor != nil && response.NextCursor.Cursor != nil {
		pageInfo, decodeErr := DecodePageInfo(response.NextCursor.Cursor)
		if decodeErr == nil && pageInfo.NextPageCursor != nil {
			var err error

			nextCookie, err = OctetStringToBytes(*pageInfo.NextPageCursor)
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to encode next page cursor received from LDAP server: %v.", err),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}
		}
	}

	if len(nextCookie) > 0 {
		// Move the session to the new cookie key: remove old key, add new key (without closing conn)
		nextKey := sessionKey(address, nextCookie)
		c.sessionPool.UpdateKey(key, nextKey)
	} else {
		// Paging done, cleanup
		c.sessionPool.Delete(key)
	}

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

		cookie, err := OctetStringToBytes(*pageInfo.NextPageCursor)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to parse cursor value: %v.", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		pageControl.SetCookie(cookie)
	}

	return pageControl, nil
}

type PageInfo struct {
	// Collection is a map of the attributes of the collection entity.
	Collection map[string]any `json:"collection"`

	// NextPageCursor is the cursor to the next page of results.
	NextPageCursor *string `json:"nextPageCursor"`
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

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	entityConfig := request.EntityConfigMap[request.EntityExternalID]
	memberOf := entityConfig.MemberOf

	// nolint: nestif
	if memberOf != nil {
		// Update required attribute for [Member] Entity
		request.Attributes = append(request.Attributes, &framework.AttributeConfig{
			ExternalId: *entityConfig.MemberUniqueIDAttribute,
			Type:       framework.AttributeTypeString,
		})

		collectionAttribute := entityConfig.CollectionAttribute
		memberOfUniqueIDAttribute := entityConfig.MemberOfUniqueIDAttribute
		memberOfAttributes := []*framework.AttributeConfig{
			{
				ExternalId: *memberOfUniqueIDAttribute,
				Type:       framework.AttributeTypeString,
				UniqueId:   true,
			},
		}

		// For a member's entity.query, {{CollectionID}} is typically replaced by collectionAttribute
		// to filter the member entity. The collectionAttribute and memberOfUniqueIdAttribute can be the identical.
		// If they are not identical, add the collectionAttribute request attributes.
		if collectionAttribute != memberOfUniqueIDAttribute {
			memberOfAttributes = append(memberOfAttributes, &framework.AttributeConfig{
				ExternalId: *collectionAttribute,
				Type:       framework.AttributeTypeString,
			})
		}

		memberOfReq := &Request{
			BaseURL:               request.BaseURL,
			Attributes:            memberOfAttributes,
			ConnectionParams:      request.ConnectionParams,
			UniqueIDAttribute:     *memberOfUniqueIDAttribute,
			EntityExternalID:      *memberOf,
			EntityConfigMap:       map[string]*EntityConfig{*memberOf: request.EntityConfigMap[*memberOf]},
			PageSize:              1,
			RequestTimeoutSeconds: request.RequestTimeoutSeconds,
		}

		// If the CollectionCursor is set, use that as the Cursor for the next call to `GetPage`.
		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			memberOfReq.Cursor = &pagination.CompositeCursor[string]{
				Cursor: request.Cursor.CollectionCursor,
			}
		}

		if request.Cursor == nil {
			request.Cursor = &pagination.CompositeCursor[string]{}
		}

		// Update filter for member entity if the cursor is set.
		if request.Cursor.Cursor != nil {
			query := entityConfig.Query
			entityConfig.Query = strings.Replace(query, "{{CollectionId}}", *request.Cursor.CollectionID, -1)
		}

		isEmptyLastPage, cursorErr := pagination.UpdateNextCursorFromCollectionAPI(
			ctx,
			request.Cursor,
			func(ctx context.Context, _ *Request) (
				int, string, []map[string]any, *pagination.CompositeCursor[string], *framework.Error,
			) {
				resp, err := d.GetPage(ctx, memberOfReq)
				if err != nil || resp == nil {
					return 0, "", nil, nil, err
				}

				if len(resp.Objects) > 0 {
					if collectionID, ok := resp.Objects[0][*collectionAttribute].(string); ok {
						query := entityConfig.Query
						entityConfig.Query = strings.Replace(query, "{{CollectionId}}", collectionID, -1)
					}
				}

				return resp.StatusCode, resp.RetryAfterHeader, resp.Objects, resp.NextCursor, nil
			},
			memberOfReq,
			*collectionAttribute,
		)

		if cursorErr != nil {
			return nil, cursorErr
		}

		if isEmptyLastPage {
			return &Response{
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	validationErr := pagination.ValidateCompositeCursor(
		request.Cursor,
		request.EntityExternalID,
		// Send a bool indicating if the entity is a member of a collection.
		memberOf != nil,
	)
	if validationErr != nil {
		return nil, validationErr
	}

	// Make sure if the connector context is set and client can proxy the request.
	if d.Client.IsProxied() {
		if ci, ok := connector.FromContext(ctx); ok {
			return d.Client.ProxyRequest(ctx, &ci, request)
		}
	}

	return d.Client.Request(ctx, request)
}

func ProcessLDAPSearchResult(result *ldap_v3.SearchResult, response *Response, request *Request) *framework.Error {
	requestAttributeMap := attrIDToConfig(request.Attributes)
	entityConfig := request.EntityConfigMap[request.EntityExternalID]
	memberOf := entityConfig.MemberOf

	objects, pageInfo, frameworkErr := ParseResponse(result, requestAttributeMap)
	if frameworkErr != nil {
		return frameworkErr
	}

	// we query only 1 page at a time for collections
	if pageInfo != nil && len(objects) == 1 && memberOf == nil {
		pageInfo.Collection = make(map[string]any, 1)
		if uniqueID, ok := objects[0][request.UniqueIDAttribute].(string); ok {
			pageInfo.Collection[request.UniqueIDAttribute] = uniqueID
		}
	}

	if pageInfo != nil && pageInfo.NextPageCursor != nil {
		b, marshalErr := json.Marshal(pageInfo)
		if marshalErr != nil {
			return &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", marshalErr),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		encodedCursor := base64.StdEncoding.EncodeToString(b)
		response.NextCursor = &pagination.CompositeCursor[string]{
			Cursor: &encodedCursor,
		}
	}

	// [MemberEntities] Set `id`, `memberId` and `memberType`.
	if memberOf != nil {
		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return &framework.Error{
				Message: "Cursor or CollectionID is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		memberUniqueIDAttribute := *entityConfig.MemberUniqueIDAttribute
		memberOfUniqueIDAttribute := *entityConfig.MemberOfUniqueIDAttribute

		memberOfUniqueIDValue, err := getMemberOfUniqueIDValue(request.Cursor, memberOfUniqueIDAttribute)
		if err != nil {
			return err
		}

		for idx, member := range objects {
			memberUniqueIDValue, ok := member[memberUniqueIDAttribute].(string)
			if !ok {
				return &framework.Error{
					Message: fmt.Sprintf(
						"Failed to parse %s field in Active Directory Member response as string.",
						memberUniqueIDAttribute,
					),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
				}
			}

			// Set Member details
			objects[idx][request.UniqueIDAttribute] = fmt.Sprintf("%s-%s", memberUniqueIDValue, memberOfUniqueIDValue)

			// format: member_attributeName
			objects[idx]["member_"+memberUniqueIDAttribute] = memberUniqueIDValue

			// format: entity_attributeName
			objects[idx][strings.ToLower(*memberOf)+"_"+memberOfUniqueIDAttribute] = memberOfUniqueIDValue
		}

		if response.NextCursor != nil && response.NextCursor.Cursor != nil {
			request.Cursor.Cursor = response.NextCursor.Cursor
		} else {
			request.Cursor.Cursor = nil
		}

		// If we have a next cursor for either the base collection (Groups) or members (Group Members),
		// encode the cursor for the next page. Otherwise, don't set a cursor as this sync is complete.
		if request.Cursor.Cursor != nil || request.Cursor.CollectionCursor != nil {
			response.NextCursor = request.Cursor
		}
	}

	response.Objects = objects

	return nil
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

func EntryToObject(e *ldap_v3.Entry, attrConfig map[string]*framework.AttributeConfig) (
	map[string]interface{}, *framework.Error) {
	result := make(map[string]interface{})
	result["dn"] = e.DN

	for _, attribute := range e.Attributes {
		// Skip attributes that are not configured in the attribute config (could be due to user error)
		currAttrConfig, ok := attrConfig[attribute.Name]
		if !ok {
			continue
		}

		value, err := StringAttrValuesToRequestedType(attribute, currAttrConfig.List, currAttrConfig.Type)
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

func getMemberOfUniqueIDValue(
	cursor *pagination.CompositeCursor[string],
	memberOfUniqueIDAttribute string,
) (any, *framework.Error) {
	var memberOfUniqueIDValue any

	if cursor.CollectionCursor != nil {
		memberOfPageInfo, pageInfoErr := DecodePageInfo(cursor.CollectionCursor)

		if pageInfoErr != nil {
			return nil, pageInfoErr
		}

		if memberOfPageInfo == nil {
			return nil, &framework.Error{
				Message: "MemberOfPageInfo is empty",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		memberOfUniqueIDValue = memberOfPageInfo.Collection[memberOfUniqueIDAttribute]
	} else {
		memberOfUniqueIDValue = *cursor.CollectionID
	}

	return memberOfUniqueIDValue, nil
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
