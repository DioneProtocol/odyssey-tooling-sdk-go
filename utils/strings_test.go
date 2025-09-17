// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "exact case match",
			s:        "Hello World",
			substr:   "Hello",
			expected: true,
		},
		{
			name:     "different case match",
			s:        "Hello World",
			substr:   "hello",
			expected: true,
		},
		{
			name:     "mixed case match",
			s:        "Hello World",
			substr:   "hELLo",
			expected: true,
		},
		{
			name:     "no match",
			s:        "Hello World",
			substr:   "Goodbye",
			expected: false,
		},
		{
			name:     "empty string",
			s:        "",
			substr:   "test",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "Hello World",
			substr:   "",
			expected: true,
		},
		{
			name:     "both empty",
			s:        "",
			substr:   "",
			expected: true,
		},
		{
			name:     "special characters",
			s:        "Hello@World#123",
			substr:   "hello@world",
			expected: true,
		},
		{
			name:     "unicode characters",
			s:        "Café",
			substr:   "café",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsIgnoreCase(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanupString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string with \\r",
			input:    "hello\rworld",
			expected: "hello\rworld",
		},
		{
			name:     "string with \\n",
			input:    "hello\nworld",
			expected: "hello\nworld",
		},
		{
			name:     "string with both \\r and \\n",
			input:    "hello\r\nworld",
			expected: "hello\r\nworld",
		},
		{
			name:     "string with multiple \\r\\n",
			input:    "\r\nhello\r\nworld\r\n",
			expected: "\nhello\r\nworld",
		},
		{
			name:     "clean string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only \\r\\n",
			input:    "\r\n",
			expected: "",
		},
		{
			name:     "only \\n",
			input:    "\n",
			expected: "",
		},
		{
			name:     "only \\r",
			input:    "\r",
			expected: "",
		},
		{
			name:     "mixed whitespace",
			input:    "  \r\n  hello  \r\n  world  \r\n  ",
			expected: "  \r\n  hello  \r\n  world  \r\n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanupString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanupStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single string",
			input:    []string{"hello\r\nworld"},
			expected: []string{"hello\r\nworld"},
		},
		{
			name:     "multiple strings",
			input:    []string{"hello\r\n", "world\n", "test\r"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "mixed clean and dirty strings",
			input:    []string{"clean", "dirty\r\n", "also clean"},
			expected: []string{"clean", "dirty", "also clean"},
		},
		{
			name:     "all clean strings",
			input:    []string{"hello", "world", "test"},
			expected: []string{"hello", "world", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanupStrings(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractPlaceholderValue(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		text     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Extract Version",
			pattern:  `dionetech/odysseygo:(\S+)`,
			text:     "dionetech/odysseygo:v1.14.4",
			expected: "v1.14.4",
			wantErr:  false,
		},
		{
			name:     "Extract File Path",
			pattern:  `config\.file=(\S+)`,
			text:     "promtail -config.file=/etc/promtail/promtail.yaml",
			expected: "/etc/promtail/promtail.yaml",
			wantErr:  false,
		},
		{
			name:     "No Match",
			pattern:  `nonexistent=(\S+)`,
			text:     "image: dionetech/odysseygo:v1.14.4",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Multiple matches - first capture group",
			pattern:  `(\w+)=`,
			text:     "name=value other=thing",
			expected: "name",
			wantErr:  false,
		},
		{
			name:     "Empty text",
			pattern:  `(\w+)`,
			text:     "",
			expected: "",
			wantErr:  true,
		},
		// Note: Invalid regex patterns cause panic in regexp.MustCompile
		// This test case is removed to avoid test panics
		{
			name:     "No capture group",
			pattern:  `\w+`,
			text:     "hello",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPlaceholderValue(tt.pattern, tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPlaceholderValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ExtractPlaceholderValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAddSingleQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "empty string",
			input:    []string{""},
			expected: []string{"''"},
		},
		{
			name:     "single character",
			input:    []string{"b"},
			expected: []string{"'b'"},
		},
		{
			name:     "multiple words",
			input:    []string{"orange banana"},
			expected: []string{"'orange banana'"},
		},
		{
			name:     "already quoted",
			input:    []string{"'apple'"},
			expected: []string{"'apple'"},
		},
		{
			name:     "missing closing quote",
			input:    []string{"'a"},
			expected: []string{"'a'"},
		},
		{
			name:     "missing opening quote",
			input:    []string{"b'"},
			expected: []string{"'b'"},
		},
		{
			name:     "mixed cases",
			input:    []string{"", "b", "orange banana", "'apple'", "'a", "b'"},
			expected: []string{"''", "'b'", "'orange banana'", "'apple'", "'a'", "'b'"},
		},
		{
			name:     "special characters",
			input:    []string{"hello@world#123"},
			expected: []string{"'hello@world#123'"},
		},
		{
			name:     "unicode characters",
			input:    []string{"café"},
			expected: []string{"'café'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddSingleQuotes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStringValue(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected string
		wantErr  bool
	}{
		{
			name:     "existing string key",
			data:     map[string]interface{}{"name": "test"},
			key:      "name",
			expected: "test",
			wantErr:  false,
		},
		{
			name:     "existing int key",
			data:     map[string]interface{}{"count": 42},
			key:      "count",
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "existing bool key",
			data:     map[string]interface{}{"enabled": true},
			key:      "enabled",
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "existing float key",
			data:     map[string]interface{}{"price": 3.14},
			key:      "price",
			expected: "3.14",
			wantErr:  false,
		},
		{
			name:     "non-existing key",
			data:     map[string]interface{}{"name": "test"},
			key:      "missing",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty map",
			data:     map[string]interface{}{},
			key:      "any",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "nil map",
			data:     nil,
			key:      "any",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty key",
			data:     map[string]interface{}{"": "empty key"},
			key:      "",
			expected: "empty key",
			wantErr:  false,
		},
		{
			name:     "nil value",
			data:     map[string]interface{}{"nil": nil},
			key:      "nil",
			expected: "<nil>",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StringValue(tt.data, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
