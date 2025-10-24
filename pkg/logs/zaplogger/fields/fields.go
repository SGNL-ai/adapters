// Copyright 2025 SGNL.ai, Inc.
package fields

import (
	"encoding/json"
	"io"

	"go.uber.org/zap"
)

// Log fields which are commonly used throughout the adapter codebase.
const (
	FieldRequestEntityExternalID  = "requestEntityExternalId"
	FieldRequestPageSize          = "requestPageSize"
	FieldRequestURL               = "requestUrl"
	FieldResponseBody             = "responseBody"
	FieldResponseNextCursor       = "responseNextCursor"
	FieldResponseObjectCount      = "responseObjectCount"
	FieldResponseRetryAfterHeader = "responseRetryAfterHeader"
	FieldResponseStatusCode       = "responseStatusCode"

	// FieldSGNLEventType is a special field used by SGNL to identify the type of event being logged.
	FieldSGNLEventType      = "eventType"
	SGNLEventTypeErrorValue = "sgnl.adapterSvc.error"
)

func RequestEntityExternalID(entityExternalID string) zap.Field {
	return zap.String(FieldRequestEntityExternalID, entityExternalID)
}

func RequestPageSize(pageSize int64) zap.Field {
	return zap.Int64(FieldRequestPageSize, pageSize)
}

func RequestURL(url string) zap.Field {
	return zap.String(FieldRequestURL, url)
}

// ResponseBody either reads from an io.ReadCloser or takes a byte slice
// and returns a zap field containing the response body for logging purposes.
func ResponseBody(body any) zap.Field {
	var bodyBytes []byte

	switch body := body.(type) {
	case io.ReadCloser:
		// Best effort read of response body for logging purposes.
		// WARNING: This will consume the body, so it should only be used
		// in contexts where the body is not needed afterwards.
		bodyBytes, _ = io.ReadAll(body)
	case []byte:
		bodyBytes = body
	}

	if json.Valid(bodyBytes) {
		return zap.Any(FieldResponseBody, json.RawMessage(bodyBytes))
	}

	return zap.ByteString(FieldResponseBody, bodyBytes)
}

func ResponseNextCursor(cursor any) zap.Field {
	return zap.Any(FieldResponseNextCursor, cursor)
}

func ResponseObjectCount(count int) zap.Field {
	return zap.Int(FieldResponseObjectCount, count)
}

func ResponseRetryAfterHeader(retryAfter string) zap.Field {
	return zap.String(FieldResponseRetryAfterHeader, retryAfter)
}

func ResponseStatusCode(statusCode int) zap.Field {
	return zap.Int(FieldResponseStatusCode, statusCode)
}

func SGNLEventTypeError() zap.Field {
	return zap.String(FieldSGNLEventType, SGNLEventTypeErrorValue)
}
