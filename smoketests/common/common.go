// Copyright 2026 SGNL.ai, Inc.

package common

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	adapter_api_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	CmpOpts = []cmp.Option{
		cmpopts.IgnoreUnexported(adapter_api_v1.GetPageResponse{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.GetPageResponse_Success{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.Page{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.Object{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.Attribute{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.AttributeValue{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.AttributeValue_DatetimeValue{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.DateTime{}),
		cmpopts.IgnoreUnexported(timestamppb.Timestamp{}),
		cmpopts.IgnoreUnexported(adapter_api_v1.EntityObjects{}),
	}
)

// StartRecorder starts a recorder with a default mode of ModeRecordOnce. This function returns a http client
// and the underlying recorder, which needs to be stopped after recording any interactions in order for
// them to be saved.
func StartRecorder(t *testing.T, path string) (*http.Client, *recorder.Recorder) {
	r, err := recorder.New(path)
	if err != nil {
		t.Fatal(err)
	}

	return r.GetDefaultClient(), r
}

func GetAdapterCtx() (context.Context, context.CancelFunc) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "token", "dGhpc2lzYXRlc3R0b2tlbg==")

	return ctx, cancelCtx
}
