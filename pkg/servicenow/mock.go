// Copyright 2026 SGNL.ai, Inc.
package servicenow

import (
	"context"
	"fmt"
	"net/http"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/web"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// If the request address is equal to MockServiceNowAddress, we return a mock recorded response.
const MockServiceNowAddress = "mock.servicenow.sgnl.ai"

// MockServiceNowGetPage returns a mock response for a GetPage request by using `go-vcr` and a recorded
// fixture file. The only request parameter used is the entity ID. The other parameters (e.g. page size, cursor)
// are ignored and do not affect the response.
func MockServiceNowGetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	fixtureName := fmt.Sprintf("pkg/mock/servicenow/fixtures/%s", request.Entity.ExternalId)

	r, err := recorder.New(
		fixtureName,
		recorder.WithMode(recorder.ModeReplayOnly),
		// The fixture file is for a specific URL (i.e. list of query params, etc.).
		// If the request slightly changes, the recorder will return an error.
		// Therefore, use a custom matcher which returns true to match any request URL
		// and always use the same fixture file for a given entity.
		recorder.WithMatcher(func(_ *http.Request, _ cassette.Request) bool { return true }),
	)
	if err != nil || r == nil {
		return framework.NewGetPageResponseError(&framework.Error{
			Message: fmt.Sprintf("Fixture file %s not found for entity: %s.", fixtureName, request.Entity.ExternalId),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		})
	}
	defer r.Stop()

	mockHTTPClient := r.GetDefaultClient()
	if mockHTTPClient == nil {
		return framework.NewGetPageResponseError(&framework.Error{
			Message: fmt.Sprintf("Failed to create mock client for entity: %s.", request.Entity.ExternalId),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		})
	}

	mockResponse, frameworkErr := NewClient(mockHTTPClient).GetPage(ctx, &Request{
		BaseURL:               request.Address,
		EntityExternalID:      request.Entity.ExternalId,
		RequestTimeoutSeconds: 10,
		// The other request fields are not required in a mock request.
	})
	if frameworkErr != nil {
		return framework.NewGetPageResponseError(frameworkErr)
	}

	objects, parserErr := web.ConvertJSONObjectList(
		&request.Entity,
		mockResponse.Objects,
		// TODO [sc-14078]: Remove support for complex attribute names.
		web.WithComplexAttributeNameDelimiter("__"),
		web.WithJSONPathAttributeNames(),
		// nolint: lll
		// The below formats are the defaults specified by Servicenow, however users are able to override the
		// global date or time format with a personal preference. TODO [sc-16472].
		// https://docs.servicenow.com/bundle/vancouver-platform-administration/page/administer/time/reference/r_FormatDateAndTimeFields.html
		web.WithDateTimeFormats(
			[]web.DateTimeFormatWithTimeZone{
				{Format: "2006-01-02 15:04:05", HasTimeZone: false},
				{Format: time.DateOnly, HasTimeZone: false},
			}...,
		),
	)
	if parserErr != nil {
		return framework.NewGetPageResponseError(&framework.Error{
			Message: fmt.Sprintf("Failed to convert datasource response objects: %v.", parserErr),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		})
	}

	return framework.NewGetPageResponseSuccess(&framework.Page{Objects: objects})
}
