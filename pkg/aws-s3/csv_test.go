// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	framework "github.com/sgnl-ai/adapter-framework"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
)

func TestCSVHeaders(t *testing.T) {
	tests := map[string]struct {
		inputReaderFn     func() *bufio.Reader
		expectedHeaders   []string
		expectedBytesRead int64
		expectedError     bool
		errorContains     string
	}{
		"valid_headers_newline_terminated": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("name,age,city\nJohn,25,NYC"))
			},
			expectedHeaders:   []string{"name", "age", "city"},
			expectedBytesRead: int64(len("name,age,city\n")),
		},
		"valid_headers_eof_terminated": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("name,age,city"))
			},
			expectedHeaders:   []string{"name", "age", "city"},
			expectedBytesRead: int64(len("name,age,city")),
		},
		"headers_with_spaces": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("First Name, Last Name, Email Address\nJohn,Doe,john@example.com"))
			},
			expectedHeaders:   []string{"First Name", " Last Name", " Email Address"},
			expectedBytesRead: int64(len("First Name, Last Name, Email Address\n")),
		},
		"single_header_newline": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("email\ntest@example.com"))
			},
			expectedHeaders:   []string{"email"},
			expectedBytesRead: int64(len("email\n")),
		},
		"single_header_eof": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("email"))
			},
			expectedHeaders:   []string{"email"},
			expectedBytesRead: int64(len("email")),
		},
		"empty_input_reader": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader(""))
			},
			expectedError: true,
			errorContains: "CSV header is empty or missing",
		},
		"header_just_newline": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("\n"))
			},
			expectedError: true,
			errorContains: "CSV file format is invalid or corrupted",
		},
		"invalid_csv_format_unclosed_quote": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("\"unclosed quote field\n"))
			},
			expectedError: true,
			errorContains: `parse error on line 1, column 23: extraneous or missing " in quoted-field`,
		},
		"header_exceeds_max_size": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader(strings.Repeat("a,", s3_adapter.MaxHeaderSizeBytes/2+1) + "last\n"))
			},
			expectedError: true,
			errorContains: fmt.Sprintf("CSV header line exceeds %dKB size limit", s3_adapter.MaxHeaderSizeBytes/1024),
		},
		"header_with_quoted_newline": {
			inputReaderFn: func() *bufio.Reader {
				return bufio.NewReader(strings.NewReader("name,\"multi\nline\nheader\",status\nvalue1,value2,value3"))
			},
			expectedHeaders:   []string{"name", "multi\nline\nheader", "status"},
			expectedBytesRead: int64(len("name,\"multi\nline\nheader\",status\n")),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			reader := tt.inputReaderFn()
			headers, bytesRead, err := s3_adapter.CSVHeaders(reader)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}

				if headers != nil {
					t.Errorf("Expected nil headers on error, got: %v", headers)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if diff := cmp.Diff(tt.expectedHeaders, headers); diff != "" {
					t.Errorf("Headers mismatch (-want +got):\n%s", diff)
				}

				if bytesRead != tt.expectedBytesRead {
					t.Errorf("Expected bytesRead %d, got %d", tt.expectedBytesRead, bytesRead)
				}
			}
		})
	}
}

func TestStreamingCSVToPage(t *testing.T) {
	sampleHeaders := []string{"name", "age", "city", "aliases"}
	attrConfigDefault := []*framework.AttributeConfig{
		{ExternalId: "name", Type: framework.AttributeTypeString},
		{ExternalId: "age", Type: framework.AttributeTypeInt64}, // Will be parsed as float64
		{ExternalId: "city", Type: framework.AttributeTypeString},
		// aliases will be handled by default JSON detection
	}
	attrConfigAllString := []*framework.AttributeConfig{
		{ExternalId: "name", Type: framework.AttributeTypeString},
		{ExternalId: "age", Type: framework.AttributeTypeString},
		{ExternalId: "city", Type: framework.AttributeTypeString},
		{ExternalId: "aliases", Type: framework.AttributeTypeString},
	}

	csvDataBasic := `John,25,NYC,"[{""alias"":""Johnny""}]"
Jane,30,LA,"[{""alias"":""Janey""}]"
Bob,35,SF,"[{""alias"":""Bobby""}]"`
	csvDataWithEmptyLine := "John,25,NYC,\n\nJane,30,LA," // Empty line, then valid line
	csvDataOneLine := "Alice,40,BOS,"

	// Generate a row that would exceed MaxCSVRowSizeBytes
	longString := strings.Repeat("longdata", s3_adapter.MaxCSVRowSizeBytes/(len("longdata")))
	csvDataExceedsRowLimit := fmt.Sprintf("name1,%s,city1,\nname2,short,city2,", longString)

	csvDataForProcessingLimit := "r1,1\nr2,22\nr3,333\n"
	headersForProcessingLimit := []string{"colA", "colB"}
	attrConfigForProcessingLimit := []*framework.AttributeConfig{
		{ExternalId: "colA", Type: framework.AttributeTypeString},
		{ExternalId: "colB", Type: framework.AttributeTypeString},
	}

	tests := map[string]struct {
		csvData                 string
		headers                 []string
		pageSize                int64
		attrConfig              []*framework.AttributeConfig
		maxProcessingBytesTotal int64
		expectedObjects         []map[string]any
		expectedHasNext         bool
		expectedError           bool
		errorContains           string
	}{
		"success_basic_csv_page1": {
			csvData:    csvDataBasic,
			headers:    sampleHeaders,
			pageSize:   2,
			attrConfig: attrConfigDefault,
			expectedObjects: []map[string]any{
				{"name": "John", "age": float64(25), "city": "NYC", "aliases": []any{map[string]any{"alias": "Johnny"}}},
				{"name": "Jane", "age": float64(30), "city": "LA", "aliases": []any{map[string]any{"alias": "Janey"}}},
			},
			expectedHasNext: true,
		},
		"success_basic_csv_page2_and_last": {
			csvData:    `Bob,35,SF,"[{""alias"":""Bobby""}]"`,
			headers:    sampleHeaders,
			pageSize:   2,
			attrConfig: attrConfigDefault,
			expectedObjects: []map[string]any{
				{"name": "Bob", "age": float64(35), "city": "SF", "aliases": []any{map[string]any{"alias": "Bobby"}}},
			},
			expectedHasNext: false,
		},
		"success_last_page_exact_size": {
			csvData:    csvDataOneLine,
			headers:    sampleHeaders,
			pageSize:   1,
			attrConfig: attrConfigDefault,
			expectedObjects: []map[string]any{
				{"name": "Alice", "age": float64(40), "city": "BOS", "aliases": ""},
			},
			expectedHasNext: false,
		},
		"success_page_size_larger_than_data": {
			csvData:    csvDataOneLine,
			headers:    sampleHeaders,
			pageSize:   5,
			attrConfig: attrConfigDefault,
			expectedObjects: []map[string]any{
				{"name": "Alice", "age": float64(40), "city": "BOS", "aliases": ""},
			},
			expectedHasNext: false,
		},
		"success_with_json_field_auto_detected": {
			csvData:  `John,"[{""alias"": ""Johnny"", ""primary"": true}]"`,
			headers:  []string{"name", "details"},
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{ // No specific config for "details"
				{ExternalId: "name", Type: framework.AttributeTypeString},
			},
			expectedObjects: []map[string]any{
				{"name": "John", "details": []any{map[string]any{"alias": "Johnny", "primary": true}}},
			},
			expectedHasNext: false,
		},
		"success_numeric_conversion_int_and_double": {
			csvData:  "John,85.5,4",
			headers:  []string{"name", "score", "rating"},
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{ExternalId: "name", Type: framework.AttributeTypeString},
				{ExternalId: "score", Type: framework.AttributeTypeDouble},
				{ExternalId: "rating", Type: framework.AttributeTypeInt64}, // Becomes float64
			},
			expectedObjects: []map[string]any{
				{"name": "John", "score": 85.5, "rating": float64(4)},
			},
			expectedHasNext: false,
		},
		"success_empty_csv_data_after_headers": {
			csvData:         "",
			headers:         sampleHeaders,
			pageSize:        2,
			attrConfig:      attrConfigDefault,
			expectedObjects: []map[string]any{},
			expectedHasNext: false,
		},
		"success_skips_empty_lines": {
			csvData:    csvDataWithEmptyLine,
			headers:    sampleHeaders,
			pageSize:   2,
			attrConfig: attrConfigDefault,
			expectedObjects: []map[string]any{
				{"name": "John", "age": float64(25), "city": "NYC", "aliases": ""},
				{"name": "Jane", "age": float64(30), "city": "LA", "aliases": ""},
			},
			expectedHasNext: false, // Assuming the data ends after Jane
		},
		"error_invalid_number_in_data": {
			csvData:       "John,not_a_number,NYC,",
			headers:       sampleHeaders,
			pageSize:      1,
			attrConfig:    attrConfigDefault, // age is Int64
			expectedError: true,
			errorContains: `CSV contains invalid numeric value "not_a_number" in column "age"`,
		},
		"error_invalid_json_in_data": {
			csvData:       `John,"[{invalid json}]"`,
			headers:       []string{"name", "data"},
			pageSize:      1,
			attrConfig:    []*framework.AttributeConfig{{ExternalId: "name", Type: framework.AttributeTypeString}},
			expectedError: true,
			errorContains: `failed to unmarshal the value: "[{invalid json}]" in column: data`,
		},
		"error_row_exceeds_max_size": {
			csvData:       csvDataExceedsRowLimit,
			headers:       sampleHeaders,
			pageSize:      2,
			attrConfig:    attrConfigAllString,
			expectedError: true,
			errorContains: fmt.Sprintf("CSV file contains a single row larger than %d MB",
				s3_adapter.MaxCSVRowSizeBytes/(1024*1024)),
		},
		"success_max_processing_bytes_total_exact_one_row": {
			csvData:                 csvDataForProcessingLimit, // "r1,1\n" (5b), "r2,22\n" (6b), "r3,333\n" (7b)
			headers:                 headersForProcessingLimit,
			pageSize:                3,
			attrConfig:              attrConfigForProcessingLimit,
			maxProcessingBytesTotal: 5, // Allows only first row (5 bytes)
			expectedObjects: []map[string]any{
				{"colA": "r1", "colB": "1"},
			},
			expectedHasNext: true, // because total data is larger
		},
		"success_max_processing_bytes_total_allows_two_rows": {
			csvData:                 csvDataForProcessingLimit, // r1 (5b), r2 (6b) = 11b total
			headers:                 headersForProcessingLimit,
			pageSize:                3,
			attrConfig:              attrConfigForProcessingLimit,
			maxProcessingBytesTotal: 11, // Allows first two rows
			expectedObjects: []map[string]any{
				{"colA": "r1", "colB": "1"},
				{"colA": "r2", "colB": "22"},
			},
			expectedHasNext: true,
		},
		"success_max_processing_bytes_total_mid_row_allowance": {
			// Total bytes read for r1=5. Next check: 5 < 8 is true. Read r2 (6 bytes). Total=11. Process r2.
			// Next check: 11 < 8 is false. Break.
			csvData:                 csvDataForProcessingLimit, // r1 (5b), r2 (6b), r3 (7b)
			headers:                 headersForProcessingLimit,
			pageSize:                3,
			attrConfig:              attrConfigForProcessingLimit,
			maxProcessingBytesTotal: 8, // Allows r1. Then r1+r2 (11b) > 8, so r2 is read & processed, then stop.
			expectedObjects: []map[string]any{
				{"colA": "r1", "colB": "1"},
				{"colA": "r2", "colB": "22"}, // r2 is processed as 5(r1) < 8, then 5+6=11.
			},
			expectedHasNext: true,
		},
		"success_max_processing_bytes_total_unlimited": {
			csvData:                 csvDataForProcessingLimit,
			headers:                 headersForProcessingLimit,
			pageSize:                3,
			attrConfig:              attrConfigForProcessingLimit,
			maxProcessingBytesTotal: 0, // Unlimited
			expectedObjects: []map[string]any{
				{"colA": "r1", "colB": "1"},
				{"colA": "r2", "colB": "22"},
				{"colA": "r3", "colB": "333"},
			},
			expectedHasNext: false,
		},
		"success_max_processing_bytes_total_less_than_first_row_but_first_row_read": {
			csvData:                 csvDataForProcessingLimit, // r1 is 5 bytes
			headers:                 headersForProcessingLimit,
			pageSize:                3,
			attrConfig:              attrConfigForProcessingLimit,
			maxProcessingBytesTotal: 3, // Less than first row
			expectedObjects: []map[string]any{
				// First row (5 bytes) is read because initial totalBytes (0) < 3 is true.
				// Then 5 >= 3, so loop breaks after processing it.
				{"colA": "r1", "colB": "1"},
			},
			expectedHasNext: true,
		},
		"error_on_record_parse_after_first_row": {
			// First row OK, second row has unclosed quote causing csv.Reader.Read() to error.
			csvData:  "good,data\n\"bad,data",
			headers:  []string{"f1", "f2"},
			pageSize: 2,
			attrConfig: []*framework.AttributeConfig{{ExternalId: "f1", Type: framework.AttributeTypeString},
				{ExternalId: "f2", Type: framework.AttributeTypeString}},
			expectedError: true,
			// The error message includes the problematic row.
			errorContains: `CSV file format is invalid or corrupted (record parse error): parse error on line 1,
			column 10: extraneous or missing " in quoted-field. Row: '"bad,data'`,
		},
		"header_name_not_in_attr_config": {
			csvData:  "valX,valY",
			headers:  []string{"HeaderX", "HeaderY"}, // HeaderY not in attrConfig
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{ExternalId: "HeaderX", Type: framework.AttributeTypeString},
			},
			expectedObjects: []map[string]any{
				{"HeaderX": "valX", "HeaderY": "valY"}, // HeaderY included as string
			},
			expectedHasNext: false,
		},
		"attr_config_for_non_existent_header": {
			csvData:  "valA",
			headers:  []string{"HeaderA"},
			pageSize: 1,
			attrConfig: []*framework.AttributeConfig{
				{ExternalId: "HeaderA", Type: framework.AttributeTypeString},
				{ExternalId: "NonExistentHeader", Type: framework.AttributeTypeInt64},
			},
			expectedObjects: []map[string]any{
				{"HeaderA": "valA"},
			},
			expectedHasNext: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			streamReader := bufio.NewReader(strings.NewReader(tt.csvData))
			objects, _, hasNext, err := s3_adapter.StreamingCSVToPage(
				streamReader,
				tt.headers,
				tt.pageSize,
				tt.attrConfig,
				tt.maxProcessingBytesTotal,
			)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}

				if diff := cmp.Diff(tt.expectedObjects, objects); diff != "" {
					t.Errorf("Objects mismatch: %s", diff)
				}

				if hasNext != tt.expectedHasNext {
					t.Errorf("HasNext mismatch: got %v, want %v", hasNext, tt.expectedHasNext)
				}
			}
		})
	}
}
