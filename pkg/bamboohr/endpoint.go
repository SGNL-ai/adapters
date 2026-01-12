// Copyright 2025 SGNL.ai, Inc.
package bamboohr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

type RequestBody struct {
	Fields []string `json:"fields"`
}

type EndpointInfo struct {
	URL  string
	Body string
}

// ConstructEndpoint constructs and returns the endpoint to query the datasource.
func ConstructEndpoint(request *Request) (*EndpointInfo, *framework.Error) {
	if request == nil {
		return nil, &framework.Error{
			Message: "Request is nil.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	var sb strings.Builder

	sb.Grow(len(request.BaseURL) + len(request.APIVersion) + len(request.CompanyDomain) + 15)
	sb.WriteString(request.BaseURL)
	sb.WriteString("/")
	sb.WriteString(request.CompanyDomain)
	sb.WriteString("/")
	sb.WriteString(request.APIVersion)
	sb.WriteString(ValidEntityExternalIDs[request.EntityConfig.ExternalId].path)

	params := url.Values{}
	params.Add("format", "JSON")
	params.Add("onlyCurrent", strconv.FormatBool(request.OnlyCurrent))

	paramString := params.Encode()

	sb.Grow(len(paramString) + 1)
	sb.WriteString("?")
	sb.WriteString(paramString)

	requestBody := RequestBody{
		Fields: make([]string, len(request.EntityConfig.Attributes)),
	}

	for idx, attr := range request.EntityConfig.Attributes {
		requestBody.Fields[idx] = attr.ExternalId
	}

	jsonData, marshalErr := json.Marshal(requestBody)
	if marshalErr != nil {
		return nil, &framework.Error{
			Message: fmt.Sprintf("Failed to marshal requestBody into JSON: %v.", marshalErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return &EndpointInfo{
		URL:  sb.String(),
		Body: string(jsonData),
	}, nil
}
