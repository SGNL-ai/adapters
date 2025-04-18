// Copyright 2025 SGNL.ai, Inc.

package awss3_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
)

func TestGetObjectKeyFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *s3_adapter.Request
		want    string
	}{
		{
			name: "simple",
			request: &s3_adapter.Request{
				PathPrefix:       "data/internal",
				FileType:         "csv",
				EntityExternalID: "users",
			},
			want: "data/internal/users.csv",
		},
		{
			name: "simple_with_trailing_slash",
			request: &s3_adapter.Request{
				PathPrefix:       "data/internal/", // trailing slash here
				FileType:         "csv",
				EntityExternalID: "users",
			},
			want: "data/internal/users.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s3_adapter.GetObjectKeyFromRequest(tt.request)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
