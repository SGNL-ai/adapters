// Copyright 2025 SGNL.ai, Inc.
package fields

import "go.uber.org/zap"

const (
	FieldRequestEntityExternalID  = "requestEntityExternalId"
	FieldRequestPageSize          = "requestPageSize"
	FieldResponseNextCursor       = "responseNextCursor"
	FieldResponseObjectCount      = "responseObjectCount"
	FieldResponseRetryAfterHeader = "responseRetryAfterHeader"
	FieldResponseStatusCode       = "responseStatusCode"
	FieldURL                      = "url"
)

func RequestEntityExternalID(entityExternalID string) zap.Field {
	return zap.String(FieldRequestEntityExternalID, entityExternalID)
}

func RequestPageSize(pageSize int64) zap.Field {
	return zap.Int64(FieldRequestPageSize, pageSize)
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

func URL(url string) zap.Field {
	return zap.String(FieldURL, url)
}
