// Copyright 2025 SGNL.ai, Inc.

package github_test

import (
	"testing"

	"github.com/sgnl-ai/adapters/pkg/github"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

func TestSetFilterParameter(t *testing.T) {
	tests := map[string]struct {
		filter *string
		want   string
	}{
		"nil_filter": {
			filter: nil,
			want:   "",
		},
		"empty_filter": {
			filter: testutil.GenPtr(""),
			want:   "",
		},
		"visibility_filter": {
			filter: testutil.GenPtr("visibility: PUBLIC"),
			want:   ", visibility: PUBLIC",
		},
		"states_filter": {
			filter: testutil.GenPtr("states: OPEN"),
			want:   ", states: OPEN",
		},
		"complex_filter": {
			filter: testutil.GenPtr("states: [OPEN, MERGED], labels: [\"bug\", \"enhancement\"]"),
			want:   ", states: [OPEN, MERGED], labels: [\"bug\", \"enhancement\"]",
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := github.SetFilterParameter(tc.filter)

			if got != tc.want {
				t.Errorf("SetFilterParameter() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSetOrderByParameter(t *testing.T) {
	tests := map[string]struct {
		orderBy *string
		want    string
	}{
		"nil_orderBy": {
			orderBy: nil,
			want:    "",
		},
		"empty_orderBy": {
			orderBy: testutil.GenPtr(""),
			want:    "",
		},
		"created_at_desc": {
			orderBy: testutil.GenPtr("orderBy: {field: CREATED_AT, direction: DESC}"),
			want:    ", orderBy: {field: CREATED_AT, direction: DESC}",
		},
		"updated_at_asc": {
			orderBy: testutil.GenPtr("orderBy: {field: UPDATED_AT, direction: ASC}"),
			want:    ", orderBy: {field: UPDATED_AT, direction: ASC}",
		},
		"name_asc": {
			orderBy: testutil.GenPtr("orderBy: {field: NAME, direction: ASC}"),
			want:    ", orderBy: {field: NAME, direction: ASC}",
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := github.SetOrderByParameter(tc.orderBy)

			if got != tc.want {
				t.Errorf("SetOrderByParameter() = %q, want %q", got, tc.want)
			}
		})
	}
}