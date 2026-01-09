// Copyright 2026 SGNL.ai, Inc.
package rootly

import (
	"math"
	"testing"
)

// TestCastToBool tests the castToBool function with various input types.
func TestCastToBool(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  bool
		expectErr bool
	}{
		// Nil input
		{
			name:      "nil input",
			input:     nil,
			expected:  false,
			expectErr: false,
		},

		// Bool inputs
		{
			name:      "bool true",
			input:     true,
			expected:  true,
			expectErr: false,
		},
		{
			name:      "bool false",
			input:     false,
			expected:  false,
			expectErr: false,
		},

		// String inputs
		{
			name:      "string 'true'",
			input:     "true",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "string 'false'",
			input:     "false",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "string '1'",
			input:     "1",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "string '0'",
			input:     "0",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "string 'True' (case insensitive)",
			input:     "True",
			expected:  true,
			expectErr: false,
		},
		{
			name:      "string 'FALSE' (case insensitive)",
			input:     "FALSE",
			expected:  false,
			expectErr: false,
		},
		{
			name:      "invalid string",
			input:     "not a bool",
			expected:  false,
			expectErr: true,
		},

		// Integer inputs
		{
			name:      "int 0",
			input:     int(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "int 1",
			input:     int(1),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "int -1",
			input:     int(-1),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "int 42",
			input:     int(42),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "int8 0",
			input:     int8(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "int16 1",
			input:     int16(1),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "int32 0",
			input:     int32(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "int64 1",
			input:     int64(1),
			expected:  true,
			expectErr: false,
		},

		// Unsigned integer inputs
		{
			name:      "uint 0",
			input:     uint(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "uint 1",
			input:     uint(1),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "uint8 0",
			input:     uint8(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "uint16 1",
			input:     uint16(1),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "uint32 0",
			input:     uint32(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "uint64 1",
			input:     uint64(1),
			expected:  true,
			expectErr: false,
		},

		// Float inputs
		{
			name:      "float32 0",
			input:     float32(0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "float32 1.5",
			input:     float32(1.5),
			expected:  true,
			expectErr: false,
		},
		{
			name:      "float64 0.0",
			input:     float64(0.0),
			expected:  false,
			expectErr: false,
		},
		{
			name:      "float64 -1.5",
			input:     float64(-1.5),
			expected:  true,
			expectErr: false,
		},

		// Invalid types
		{
			name:      "slice",
			input:     []int{1, 2, 3},
			expected:  false,
			expectErr: true,
		},
		{
			name:      "map",
			input:     map[string]int{"a": 1},
			expected:  false,
			expectErr: true,
		},
		{
			name:      "struct",
			input:     struct{ value int }{value: 1},
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			result, err := castToBool(tt.input)

			// Assert
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

// TestCastToFloat64 tests the castToFloat64 function with various input types.
func TestCastToFloat64(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  float64
		expectErr bool
	}{
		// Nil input
		{
			name:      "nil input",
			input:     nil,
			expected:  0,
			expectErr: false,
		},

		// Float inputs
		{
			name:      "float64 42.5",
			input:     float64(42.5),
			expected:  42.5,
			expectErr: false,
		},
		{
			name:      "float64 0.0",
			input:     float64(0.0),
			expected:  0.0,
			expectErr: false,
		},
		{
			name:      "float64 -3.14",
			input:     float64(-3.14),
			expected:  -3.14,
			expectErr: false,
		},
		{
			name:      "float32 1.5",
			input:     float32(1.5),
			expected:  1.5,
			expectErr: false,
		},

		// Integer inputs
		{
			name:      "int 42",
			input:     int(42),
			expected:  42.0,
			expectErr: false,
		},
		{
			name:      "int 0",
			input:     int(0),
			expected:  0.0,
			expectErr: false,
		},
		{
			name:      "int -10",
			input:     int(-10),
			expected:  -10.0,
			expectErr: false,
		},
		{
			name:      "int8 127",
			input:     int8(127),
			expected:  127.0,
			expectErr: false,
		},
		{
			name:      "int16 -500",
			input:     int16(-500),
			expected:  -500.0,
			expectErr: false,
		},
		{
			name:      "int32 1000",
			input:     int32(1000),
			expected:  1000.0,
			expectErr: false,
		},
		{
			name:      "int64 999999",
			input:     int64(999999),
			expected:  999999.0,
			expectErr: false,
		},

		// Unsigned integer inputs
		{
			name:      "uint 42",
			input:     uint(42),
			expected:  42.0,
			expectErr: false,
		},
		{
			name:      "uint8 255",
			input:     uint8(255),
			expected:  255.0,
			expectErr: false,
		},
		{
			name:      "uint16 1000",
			input:     uint16(1000),
			expected:  1000.0,
			expectErr: false,
		},
		{
			name:      "uint32 50000",
			input:     uint32(50000),
			expected:  50000.0,
			expectErr: false,
		},
		{
			name:      "uint64 999999",
			input:     uint64(999999),
			expected:  999999.0,
			expectErr: false,
		},

		// String inputs
		{
			name:      "string '42.5'",
			input:     "42.5",
			expected:  42.5,
			expectErr: false,
		},
		{
			name:      "string '0'",
			input:     "0",
			expected:  0.0,
			expectErr: false,
		},
		{
			name:      "string '-3.14'",
			input:     "-3.14",
			expected:  -3.14,
			expectErr: false,
		},
		{
			name:      "string '1e3'",
			input:     "1e3",
			expected:  1000.0,
			expectErr: false,
		},
		{
			name:      "string '1.23e-4'",
			input:     "1.23e-4",
			expected:  0.000123,
			expectErr: false,
		},
		{
			name:      "invalid string",
			input:     "not a number",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expected:  0,
			expectErr: true,
		},

		// Bool inputs
		{
			name:      "bool true",
			input:     true,
			expected:  1.0,
			expectErr: false,
		},
		{
			name:      "bool false",
			input:     false,
			expected:  0.0,
			expectErr: false,
		},

		// Special float values
		{
			name:      "max float64",
			input:     math.MaxFloat64,
			expected:  math.MaxFloat64,
			expectErr: false,
		},
		{
			name:      "smallest positive float64",
			input:     math.SmallestNonzeroFloat64,
			expected:  math.SmallestNonzeroFloat64,
			expectErr: false,
		},

		// Invalid types
		{
			name:      "slice",
			input:     []int{1, 2, 3},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "map",
			input:     map[string]int{"a": 1},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "struct",
			input:     struct{ value int }{value: 1},
			expected:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			result, err := castToFloat64(tt.input)

			// Assert
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

// TestCastToString tests the castToString function with various input types.
func TestCastToString(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  string
		expectErr bool
	}{
		// Nil input
		{
			name:      "nil input",
			input:     nil,
			expected:  "",
			expectErr: false,
		},

		// String inputs
		{
			name:      "string 'hello'",
			input:     "hello",
			expected:  "hello",
			expectErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			expected:  "",
			expectErr: false,
		},
		{
			name:      "string with spaces",
			input:     "hello world",
			expected:  "hello world",
			expectErr: false,
		},
		{
			name:      "string with special chars",
			input:     "hello\nworld\t!",
			expected:  "hello\nworld\t!",
			expectErr: false,
		},

		// Bool inputs
		{
			name:      "bool true",
			input:     true,
			expected:  "true",
			expectErr: false,
		},
		{
			name:      "bool false",
			input:     false,
			expected:  "false",
			expectErr: false,
		},

		// Integer inputs
		{
			name:      "int 42",
			input:     int(42),
			expected:  "42",
			expectErr: false,
		},
		{
			name:      "int 0",
			input:     int(0),
			expected:  "0",
			expectErr: false,
		},
		{
			name:      "int -10",
			input:     int(-10),
			expected:  "-10",
			expectErr: false,
		},
		{
			name:      "int8 127",
			input:     int8(127),
			expected:  "127",
			expectErr: false,
		},
		{
			name:      "int16 -500",
			input:     int16(-500),
			expected:  "-500",
			expectErr: false,
		},
		{
			name:      "int32 1000",
			input:     int32(1000),
			expected:  "1000",
			expectErr: false,
		},
		{
			name:      "int64 999999",
			input:     int64(999999),
			expected:  "999999",
			expectErr: false,
		},

		// Unsigned integer inputs
		{
			name:      "uint 42",
			input:     uint(42),
			expected:  "42",
			expectErr: false,
		},
		{
			name:      "uint8 255",
			input:     uint8(255),
			expected:  "255",
			expectErr: false,
		},
		{
			name:      "uint16 1000",
			input:     uint16(1000),
			expected:  "1000",
			expectErr: false,
		},
		{
			name:      "uint32 50000",
			input:     uint32(50000),
			expected:  "50000",
			expectErr: false,
		},
		{
			name:      "uint64 999999",
			input:     uint64(999999),
			expected:  "999999",
			expectErr: false,
		},

		// Float inputs
		{
			name:      "float32 1.5",
			input:     float32(1.5),
			expected:  "1.5",
			expectErr: false,
		},
		{
			name:      "float32 0.0",
			input:     float32(0.0),
			expected:  "0",
			expectErr: false,
		},
		{
			name:      "float64 42.5",
			input:     float64(42.5),
			expected:  "42.5",
			expectErr: false,
		},
		{
			name:      "float64 -3.14",
			input:     float64(-3.14),
			expected:  "-3.14",
			expectErr: false,
		},
		{
			name:      "float64 0.0",
			input:     float64(0.0),
			expected:  "0",
			expectErr: false,
		},
		{
			name:      "float64 1.23456789",
			input:     float64(1.23456789),
			expected:  "1.23456789",
			expectErr: false,
		},

		// Complex types (should use fmt.Sprintf fallback)
		{
			name:      "slice",
			input:     []int{1, 2, 3},
			expected:  "[1 2 3]",
			expectErr: false,
		},
		{
			name:      "map",
			input:     map[string]int{"a": 1, "b": 2},
			expected:  "map[a:1 b:2]",
			expectErr: false,
		},
		{
			name: "pointer",
			input: func() *int {
				i := 42

				return &i
			}(),
			expected:  "42",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			result, err := castToString(tt.input)

			// Assert
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

// TestCastToStringComplexTypes tests castToString with complex types that use fmt.Sprintf.
func TestCastToStringComplexTypes(t *testing.T) {
	// Test that complex types are converted using fmt.Sprintf without error
	tests := []struct {
		name  string
		input interface{}
	}{
		{"slice of strings", []string{"a", "b", "c"}},
		{"empty slice", []int{}},
		{"map", map[string]string{"key": "value"}},
		{"struct", struct{ Name string }{"test"}},
		{"nested struct", struct{ Inner struct{ Value int } }{Inner: struct{ Value int }{42}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange is in test table (tt.input)

			// Act
			result, err := castToString(tt.input)

			// Assert
			if err != nil {
				t.Errorf("unexpected error for complex type: %v", err)
			}
			if result == "" {
				t.Errorf("expected non-empty string for complex type, got empty string")
			}
		})
	}
}

// BenchmarkCastToBool benchmarks the castToBool function.
func BenchmarkCastToBool(b *testing.B) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{"bool", true},
		{"string", "true"},
		{"int", int(1)},
		{"float64", float64(1.0)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = castToBool(tc.input)
			}
		})
	}
}

// BenchmarkCastToFloat64 benchmarks the castToFloat64 function.
func BenchmarkCastToFloat64(b *testing.B) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{"float64", float64(42.5)},
		{"int", int(42)},
		{"string", "42.5"},
		{"bool", true},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = castToFloat64(tc.input)
			}
		})
	}
}

// BenchmarkCastToString benchmarks the castToString function.
func BenchmarkCastToString(b *testing.B) {
	testCases := []struct {
		name  string
		input interface{}
	}{
		{"string", "hello"},
		{"int", int(42)},
		{"float64", float64(42.5)},
		{"bool", true},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = castToString(tc.input)
			}
		})
	}
}
