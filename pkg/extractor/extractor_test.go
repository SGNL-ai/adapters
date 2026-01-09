// Copyright 2026 SGNL.ai, Inc.
package extractor_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sgnl-ai/adapters/pkg/extractor"
)

func TestValueFromList(t *testing.T) {
	tests := map[string]struct {
		inputValues         []string
		inputIncludedPrefix string
		inputExcludedSuffix string
		want                string
	}{
		"simple": {
			inputValues: []string{
				"<https://test-instance.oktapreview.com/api/v1/users?limit=2>; rel=\"self\"",
				"<https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>; rel=\"next\"",
			},
			inputIncludedPrefix: "https://",
			inputExcludedSuffix: ">; rel=\"next\"",
			want:                "https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2",
		},
		"missing_value": {
			inputValues: []string{
				"<https://test-instance.oktapreview.com/api/v1/users?limit=2>; rel=\"self\"",
			},
			inputIncludedPrefix: "https://",
			inputExcludedSuffix: ">; rel=\"next\"",
			want:                "",
		},
		"missing_suffix_value_present": {
			inputValues: []string{
				"<https://test-instance.oktapreview.com/api/v1/users?limit=2>; rel=\"self\"",
				"<https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>",
			},
			inputIncludedPrefix: "https://",
			inputExcludedSuffix: ">; rel=\"next\"",
			want:                "",
		},
		"missing_prefix_value_present": {
			inputValues: []string{
				"<https://test-instance.oktapreview.com/api/v1/users?limit=2>; rel=\"self\"",
				"</test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>; rel=\"next\"",
			},
			inputIncludedPrefix: "https://",
			inputExcludedSuffix: ">; rel=\"next\"",
			want:                "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractor.ValueFromList(tt.inputValues, tt.inputIncludedPrefix, tt.inputExcludedSuffix)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestAttributesFromJSONPath(t *testing.T) {
	tests := map[string]struct {
		inputExpression string
		wantAttributes  []string
		wantError       error
	}{
		"empty_expressions": {
			inputExpression: "$.",
		},
		"invalid_prefix": {
			inputExpression: "@.store.book.author.firstName",
			wantError:       errors.New("expression missing required '$.' prefix"),
		},
		"missing_prefix": {
			inputExpression: "store.book.author.firstName",
			wantError:       errors.New("expression missing required '$.' prefix"),
		},
		"simple_dot_notation": {
			inputExpression: "$.store",
			wantAttributes:  []string{"store"},
		},
		"dot_notation": {
			inputExpression: "$.store.book.author.firstName",
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"dot_notation_with_dash": {
			inputExpression: "$.flagship-store.book",
			wantAttributes:  []string{"flagship-store", "book"},
		},
		"dot_bracket_notation_single_quote": {
			inputExpression: "$.['store'].['book'].['author'].['firstName']",
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"bracket_notation_single_quote": {
			inputExpression: "$.['store']['book']['author']['firstName']",
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"bracket_and_dot_bracket_notation_single_quote": {
			inputExpression: "$.['store']['book'].['author'].['firstName']['initial']",
			wantAttributes:  []string{"store", "book", "author", "firstName", "initial"},
		},
		"dot_bracket_notation_escaped_single_quote": {
			inputExpression: `$.[\'store\'].[\'book\'].[\'author\'].[\'firstName\']`,
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"dot_bracket_notation_double_quote": {
			inputExpression: `$.["store"].["book"].["author"].["firstName"]`,
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"dot_bracket_notation_escaped_double_quote": {
			inputExpression: `$.[\"store\"].[\"book\"].[\"author\"].[\"firstName\"]`,
			wantAttributes:  []string{"store", "book", "author", "firstName"},
		},
		"bracket_notation_after_recursive_operator": {
			inputExpression: "$..['$ref']",
			wantAttributes:  []string{"$ref"},
		},
		"bracket_notation_empty_string": {
			inputExpression: "$.['']",
			wantAttributes:  []string{""},
		},
		"bracket_notation_quoted_array_slice_literal": {
			inputExpression: "$..[':']",
			wantAttributes:  []string{":"},
		},
		"bracket_notation_quoted_join_literal": {
			inputExpression: "$..[',']",
			wantAttributes:  []string{","},
		},
		"bracket_notation_quoted_closing_bracket_literal": {
			inputExpression: "$..[']']",
			wantAttributes:  []string{"]"},
		},
		"bracket_notation_quoted_opening_bracket_literal": {
			inputExpression: "$..['[']",
			wantAttributes:  []string{"["},
		},
		"dot_bracket_notation_quoted_closing_bracket_literal_no_closing_bracket": {
			inputExpression: "$.[']'",
			wantError:       errors.New("invalid expression provided: missing closing bracket"),
		},
		"bracket_notation_quoted_dot_literal": {
			inputExpression: "$..['.']",
			wantAttributes:  []string{"."},
		},
		"bracket_notation_quoted_closing_bracket_and_dot_literals": {
			inputExpression: "$..['].']",
			wantAttributes:  []string{"]."},
		},
		"bracket_notation_quoted_opening_bracket_and_dot_literals": {
			inputExpression: "$..['[.']",
			wantAttributes:  []string{"[."},
		},
		"bracket_notation_quoted_number": {
			inputExpression: "$..['1']",
			wantAttributes:  []string{"1"},
		},
		"bracket_notation_quoted_wildcard_literal": {
			inputExpression: "$..['*']",
			wantAttributes:  []string{"*"},
		},
		"dot_notation_after_recursive_operator": {
			inputExpression: "$..book",
			wantAttributes:  []string{"book"},
		},
		"recursive_operator_after_dot_notation": {
			inputExpression: "$.store..book",
			wantAttributes:  []string{"store", "book"},
		},
		"recursive_operator_after_bracket_notation": {
			inputExpression: "$.['store']..book",
			wantAttributes:  []string{"store", "book"},
		},
		"dot_bracket_notation_with_dot": {
			inputExpression: "$.['store.book.author']",
			wantAttributes:  []string{"store.book.author"},
		},
		"ignored_subscript_operator_after_dot_notation": {
			inputExpression: "$.store.[1].book",
			wantAttributes:  []string{"store", "book"},
		},
		"ignored_filter_after_dot_notation": {
			inputExpression: "$.store[?(@.key==$.value)].book",
			wantAttributes:  []string{"store", "book"},
		},
		"ignored_subscript_operator": {
			inputExpression: "$.[0]",
		},
		"dot_notation_after_ignored_subscript_operator": {
			inputExpression: "$.[0].book",
			wantAttributes:  []string{"book"},
		},
		"ignored_wildcard": {
			inputExpression: "$.*",
		},
		"ignored_wildcard_in_brackets": {
			inputExpression: "$.[*]",
		},
		"ignored_range_operator": {
			inputExpression: "$.stores..[1:4].book",
			wantAttributes:  []string{"stores", "book"},
		},
		"ignored_join_operator": {
			inputExpression: "$.stores..[1,4].book",
			wantAttributes:  []string{"stores", "book"},
		},
		"ignored_script_expression": {
			inputExpression: "$.stores..[(@.length-1)].book",
			wantAttributes:  []string{"stores", "book"},
		},
		"ignored_filter_expression": {
			inputExpression: "$.stores..[?(@.isbn)].book",
			wantAttributes:  []string{"stores", "book"},
		},
		"ignored_filter_and_subscript": {
			inputExpression: "$.book[?(@.price<10)].authors[1]",
			wantAttributes:  []string{"book", "authors"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := extractor.AttributesFromJSONPath(tt.inputExpression)

			if !reflect.DeepEqual(got, tt.wantAttributes) {
				t.Errorf("got: %v, want: %v", got, tt.wantAttributes)
			}

			if (tt.wantError != nil && gotErr == nil) || (tt.wantError == nil && gotErr != nil) {
				t.Errorf("gotErr: %v, wantError: %v", gotErr, tt.wantError)
			}

			if tt.wantError != nil && gotErr != nil {
				if !reflect.DeepEqual(gotErr.Error(), tt.wantError.Error()) {
					t.Errorf("gotErr: %v, wantError: %v", gotErr, tt.wantError)
				}
			}
		})
	}
}
