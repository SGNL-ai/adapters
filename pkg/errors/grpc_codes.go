// Copyright 2025 SGNL.ai, Inc.

package customerror

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	GRPCStatusCodeToHTTP = map[codes.Code]int{
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusGatewayTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
	}
)

func GRPCErrStatusToHTTPStatusCode(s *status.Status, err error) int {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	if httpStatusCode, ok := GRPCStatusCodeToHTTP[s.Code()]; ok {
		return httpStatusCode
	}

	logger.Printf("Unknown gRPC status code received: %v \t %v\n", s.Code(), err)

	return http.StatusInternalServerError
}
