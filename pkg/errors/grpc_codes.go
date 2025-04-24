// Copyright 2025 SGNL.ai, Inc.

package customerror

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/status"
)

var (
	GRPCStatusCodeToHTTP = map[code.Code]int{
		//		code.Code_CANCELLED:           , // Client cancelled the request
		code.Code_UNKNOWN:             http.StatusInternalServerError,
		code.Code_INVALID_ARGUMENT:    http.StatusBadRequest,
		code.Code_DEADLINE_EXCEEDED:   http.StatusGatewayTimeout,
		code.Code_NOT_FOUND:           http.StatusNotFound,
		code.Code_ALREADY_EXISTS:      http.StatusConflict,
		code.Code_PERMISSION_DENIED:   http.StatusForbidden,
		code.Code_RESOURCE_EXHAUSTED:  http.StatusTooManyRequests,
		code.Code_FAILED_PRECONDITION: http.StatusBadRequest,
		code.Code_ABORTED:             http.StatusConflict,
		code.Code_OUT_OF_RANGE:        http.StatusBadRequest,
		code.Code_UNIMPLEMENTED:       http.StatusNotImplemented,
		code.Code_INTERNAL:            http.StatusInternalServerError,
		code.Code_UNAVAILABLE:         http.StatusServiceUnavailable,
		code.Code_DATA_LOSS:           http.StatusInternalServerError,
	}
)

func GRPCStatusCodeToHTTPStatusCode(s *status.Status, err error) int {
	logger := log.New(os.Stdout, "adapter", log.Lmicroseconds|log.LUTC|log.Lshortfile)

	if httpStatusCode, ok := GRPCStatusCodeToHTTP[code.Code(s.Code())]; ok {
		return httpStatusCode
	}

	logger.Printf("Unknown gRPC status code received: %v \t %v\n", s.Code(), err)

	return http.StatusInternalServerError
}
