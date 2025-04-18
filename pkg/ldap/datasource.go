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
	"os"
	"reflect"
	"strconv"
	"strings"

	parser "github.com/Azure/azure-storage-azcopy/v10/sddl"
	"github.com/bwmarrin/go-objectsid"
	ldap_v3 "github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapters/pkg/errors"
	"github.com/sgnl-ai/adapters/pkg/pagination"
)

// Datasource directly implements a Client interface to allow querying an external datasource.
type Datasource struct{}

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

// NewClient returns a Client to query the datasource.
func NewClient() Client {
	return &Datasource{}
}

func (d *Datasource) GetPage(ctx context.Context, request *Request) (*Response, *framework.Error) {
	entityConfig := request.EntityConfigMap[request.EntityExternalID]
	memberOf := entityConfig.MemberOf

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
				if err != nil {
					return 0, "", nil, nil, err
				}

				if collectionID, ok := resp.Objects[0][*collectionAttribute].(string); ok {
					query := entityConfig.Query
					entityConfig.Query = strings.Replace(query, "{{CollectionId}}", collectionID, -1)
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

	tlsConfig := &tls.Config{}

	if request.IsLDAPS {
		decodedCertChain, err := base64.StdEncoding.DecodeString(request.CertificateChain)
		if err != nil {
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to load certificates - %v", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
			}
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(decodedCertChain)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.ServerName = request.Host
	}

	ldapConn, err := ldap_v3.DialURL(request.BaseURL, ldap_v3.DialWithTLSConfig(tlsConfig))
	if err != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to dial ldap server - %v", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	// Bind credentials with connection
	bindErr := ldapConn.Bind(request.BindDN, request.BindPassword)
	if bindErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to bind credentials - %v", bindErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_AUTHENTICATION_FAILED,
		}
	}
	defer ldapConn.Close()

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
				Message: fmt.Sprintf("Failed to parse cursor value - %v", err),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		pageControl.SetCookie(cookie)
	}

	filters, filterErr := SetFilters(request)
	if filterErr != nil {
		return nil, filterErr
	}

	attributes := []string{}

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
	searchResult, err := ldapConn.Search(searchRequest)
	if err != nil {
		// Extract LDAP result code from the error
		if ldapErr, ok := err.(*ldap_v3.Error); ok {
			statusCode := ResultCodeToHTTPStatusCode(ldapErr)
			response.StatusCode = statusCode

			return response, nil
		}

		return nil, customerror.UpdateError(&framework.Error{
			Message: fmt.Sprintf("Error searching LDAP server: - %v.", err),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		},
			customerror.WithRequestTimeoutMessage(err, request.RequestTimeoutSeconds),
		)
	} else if (searchResult == nil || len(searchResult.Entries) == 0) &&
		(request.Cursor == nil || request.Cursor.CollectionID == nil) {
		response.StatusCode = http.StatusNotFound

		return response, nil
	}

	// Indicating a successful LDAP search operation.
	// In case of no error (err == nil), ldap_v3.Search is considered successful,
	// returning LDAP Result Code Success(0) equivalent to HTTP status code StatusOK.
	response.StatusCode = http.StatusOK

	requestAttributeMap := attrIDToConfig(request.Attributes)

	objects, pageInfo, frameworkErr := ParseResponse(searchResult, requestAttributeMap)
	if frameworkErr != nil {
		return nil, frameworkErr
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
			return nil, &framework.Error{
				Message: fmt.Sprintf("Failed to create updated cursor: %v.", err),
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
		var memberOfUniqueIDValue any

		if request.Cursor == nil || request.Cursor.CollectionID == nil {
			return nil, &framework.Error{
				Message: "Cursor or CollectionID is nil",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			}
		}

		memberUniqueIDAttribute := *entityConfig.MemberUniqueIDAttribute
		memberOfUniqueIDAttribute := *entityConfig.MemberOfUniqueIDAttribute

		if request.Cursor != nil && request.Cursor.CollectionCursor != nil {
			memberOfPageInfo, _ := DecodePageInfo(request.Cursor.CollectionCursor)
			memberOfUniqueIDValue = memberOfPageInfo.Collection[memberOfUniqueIDAttribute]
		}

		for idx, member := range objects {
			memberUniqueIDValue, ok := member[memberUniqueIDAttribute].(string)
			if !ok {
				return nil, &framework.Error{
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

func EntryToObject(e *ldap_v3.Entry, attrConfig map[string]*framework.AttributeConfig) (
	map[string]interface{}, *framework.Error) {
	result := make(map[string]interface{})
	result["dn"] = e.DN

	for _, attribute := range e.Attributes {
		currAttrConfig := attrConfig[attribute.Name]

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

func StringAttrValuesToRequestedType(attr *ldap_v3.EntryAttribute, isList bool,
	attrType framework.AttributeType) (any, *framework.Error) {
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

	switch attrType {
	case getAttrType(api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING):
		switch attr.Name {
		// Special AD syntaxes.
		case objectGUID:
			guid, err := uuid.Parse(hex.EncodeToString(attr.ByteValues[0]))
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf(ErrorMsgAttributeTypeDoesNotMatchFmt,
						attr.Name, reflect.TypeOf(attr.Values[0]), "string"),
					Code: api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

			return guid.String(), nil
		case objectSid, sidHistory, creatorSID, securityIdentifier:
			sid := objectsid.Decode(attr.ByteValues[0])

			return sid.String(), nil
		case nTSecurityDescriptor, msDSAllowedToActOnBehalfOfOtherIdentity, fRSRootSecurity, pKIEnrollmentAccess,
			msDSGroupMSAMembership, msDFSLinkSecurityDescriptorv2:
			sddl, err := parser.SecurityDescriptorToString(attr.ByteValues[0])
			if err != nil {
				return nil, &framework.Error{
					Message: fmt.Sprintf("Failed to parse a String(NT-Sec-Desc) syntax attribute: %s", attr.Name),
					Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ATTRIBUTE_TYPE,
				}
			}

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

func ResultCodeToHTTPStatusCode(ldapError *ldap_v3.Error) int {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	ldapToHTTP := map[uint16]int{
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
		999: http.StatusInternalServerError, // Default value for unknown LDAP result codes
	}
	if httpStatusCode, ok := ldapToHTTP[ldapError.ResultCode]; ok {
		return httpStatusCode
	}

	logger.Printf("Unknown LDAP result code received: %v \t %v\n", ldapError.ResultCode, ldapError.Err.Error())

	return ldapToHTTP[999]
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
