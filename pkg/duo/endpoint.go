// Copyright 2026 SGNL.ai, Inc.
package duo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

type EndpointInfo struct {
	URL  string
	Auth string
	Date string
}

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (*EndpointInfo, *framework.Error) {
	if request == nil {
		return nil, &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	var offset int64

	if request.Cursor != nil && request.Cursor.Cursor != nil {
		if *request.Cursor.Cursor <= 0 {
			return nil, &framework.Error{
				Message: "Cursor must be greater than 0.",
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}

		offset = *request.Cursor.Cursor
	}

	path := fmt.Sprintf("/admin/%s/%s", request.APIVersion, ValidEntityExternalIDs[request.EntityExternalID].path)
	params := fmt.Sprintf("limit=%d&offset=%d", request.PageSize, offset)
	auth, date := ConfigureAuth(request, path, params)
	baseURL := request.BaseURL
	endpoint := fmt.Sprintf("%s%s?%s", baseURL, path, params)

	return &EndpointInfo{
		URL:  endpoint,
		Auth: auth,
		Date: date,
	}, nil
}

// ConfigureAuth configures the auth and date headers for the request per the Duo Admin API standards
// Example Request Signature: https://duo.com/docs/adminapi#authentication
// Tue, 21 Aug 2012 17:29:18 -0000
// GET
// api-xxxxxxxx.duosecurity.com
// /admin/v1/users
// limit=1&offset=0.

func ConfigureAuth(request *Request, path, params string) (string, string) {
	date := time.Now().Format(RFC2822)
	hostname := strings.Replace(request.BaseURL, "https://", "", 1)

	hmac := hmac.New(sha1.New, []byte(request.Secret))
	hmac.Write([]byte(fmt.Sprintf("%s\nGET\n%s\n%s\n%s", date, hostname, path, params)))

	token := fmt.Sprintf("%s:%s", request.IntegrationKey, hex.EncodeToString(hmac.Sum(nil)))
	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(token)))

	return auth, date
}
