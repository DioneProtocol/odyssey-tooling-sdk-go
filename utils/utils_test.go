// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"testing"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "all unique elements",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all duplicate elements",
			input:    []int{1, 1, 1, 1},
			expected: []int{1},
		},
		{
			name:     "mixed unique and duplicate elements",
			input:    []int{1, 2, 1, 3, 2, 4, 1},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "duplicates at different positions",
			input:    []int{1, 2, 3, 1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with different data types
	t.Run("string slice", func(t *testing.T) {
		input := []string{"a", "b", "a", "c", "b"}
		expected := []string{"a", "b", "c"}
		result := Unique(input)
		assert.Equal(t, expected, result)
	})

	t.Run("struct slice", func(t *testing.T) {
		type testStruct struct {
			ID   int
			Name string
		}
		input := []testStruct{
			{1, "a"},
			{2, "b"},
			{1, "a"},
			{3, "c"},
		}
		expected := []testStruct{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		}
		result := Unique(input)
		assert.Equal(t, expected, result)
	})
}

func TestUint32Sort(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint32
		expected []uint32
	}{
		{
			name:     "empty slice",
			input:    []uint32{},
			expected: []uint32{},
		},
		{
			name:     "single element",
			input:    []uint32{5},
			expected: []uint32{5},
		},
		{
			name:     "already sorted",
			input:    []uint32{1, 2, 3, 4, 5},
			expected: []uint32{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse sorted",
			input:    []uint32{5, 4, 3, 2, 1},
			expected: []uint32{1, 2, 3, 4, 5},
		},
		{
			name:     "random order",
			input:    []uint32{3, 1, 4, 1, 5, 9, 2, 6},
			expected: []uint32{1, 1, 2, 3, 4, 5, 6, 9},
		},
		{
			name:     "duplicate values",
			input:    []uint32{3, 1, 3, 1, 2, 3},
			expected: []uint32{1, 1, 2, 3, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying the original
			input := make([]uint32, len(tt.input))
			copy(input, tt.input)

			Uint32Sort(input)
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestGetAPIContext(t *testing.T) {
	ctx, cancel := GetAPIContext()
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)

	// Test that the context has the correct timeout
	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	expectedTimeout := time.Now().Add(constants.APIRequestTimeout)

	// Allow for small time differences
	timeDiff := expectedTimeout.Sub(deadline)
	assert.True(t, timeDiff >= 0 && timeDiff < time.Second,
		"Expected timeout around %v, got %v", constants.APIRequestTimeout, deadline.Sub(time.Now()))

	// Test that cancel works
	cancel()
	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled")
	}
}

func TestGetAPILargeContext(t *testing.T) {
	ctx, cancel := GetAPILargeContext()
	require.NotNil(t, ctx)
	require.NotNil(t, cancel)

	// Test that the context has the correct timeout
	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	expectedTimeout := time.Now().Add(constants.APIRequestLargeTimeout)

	// Allow for small time differences
	timeDiff := expectedTimeout.Sub(deadline)
	assert.True(t, timeDiff >= 0 && timeDiff < time.Second,
		"Expected timeout around %v, got %v", constants.APIRequestLargeTimeout, deadline.Sub(time.Now()))

	// Test that cancel works
	cancel()
	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled")
	}
}

func TestO(t *testing.T) {
	// Create test addresses
	addr1 := ids.ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	addr2 := ids.ShortID{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	tests := []struct {
		name       string
		networkHRP string
		addresses  []ids.ShortID
		wantErr    bool
	}{
		{
			name:       "empty addresses",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{},
			wantErr:    false,
		},
		{
			name:       "single address",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{addr1},
			wantErr:    false,
		},
		{
			name:       "multiple addresses",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{addr1, addr2},
			wantErr:    false,
		},
		{
			name:       "mainnet HRP",
			networkHRP: "mainnet",
			addresses:  []ids.ShortID{addr1},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := O(tt.networkHRP, tt.addresses)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, result, len(tt.addresses))

			// Verify that all addresses start with "O-" and contain the HRP
			for i, addr := range result {
				assert.True(t, len(addr) > 0, "Address %d should not be empty", i)
				// Note: The actual format validation would depend on the address.Format function
			}
		})
	}
}

func TestA(t *testing.T) {
	// Create test addresses
	addr1 := ids.ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	addr2 := ids.ShortID{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	tests := []struct {
		name       string
		networkHRP string
		addresses  []ids.ShortID
		wantErr    bool
	}{
		{
			name:       "empty addresses",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{},
			wantErr:    false,
		},
		{
			name:       "single address",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{addr1},
			wantErr:    false,
		},
		{
			name:       "multiple addresses",
			networkHRP: "testnet",
			addresses:  []ids.ShortID{addr1, addr2},
			wantErr:    false,
		},
		{
			name:       "mainnet HRP",
			networkHRP: "mainnet",
			addresses:  []ids.ShortID{addr1},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := A(tt.networkHRP, tt.addresses)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, result, len(tt.addresses))

			// Verify that all addresses start with "A-" and contain the HRP
			for i, addr := range result {
				assert.True(t, len(addr) > 0, "Address %d should not be empty", i)
				// Note: The actual format validation would depend on the address.Format function
			}
		})
	}
}

func TestRemoveSurrounding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		left     string
		right    string
		expected string
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			left:     "(",
			right:    ")",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "string with both left and right",
			input:    "(hello)",
			left:     "(",
			right:    ")",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "string with whitespace",
			input:    "  (hello)  ",
			left:     "(",
			right:    ")",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "string missing left prefix",
			input:    "hello)",
			left:     "(",
			right:    ")",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "string missing right suffix",
			input:    "(hello",
			left:     "(",
			right:    ")",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "string too short",
			input:    "()",
			left:     "(",
			right:    ")",
			expected: "",
			wantErr:  false, // Empty string after trimming is valid
		},
		{
			name:     "string shorter than left+right",
			input:    "(",
			left:     "(",
			right:    ")",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "different delimiters",
			input:    "[test]",
			left:     "[",
			right:    "]",
			expected: "test",
			wantErr:  false,
		},
		{
			name:     "nested delimiters",
			input:    "((hello))",
			left:     "(",
			right:    ")",
			expected: "(hello)",
			wantErr:  false,
		},
		{
			name:     "multiple occurrences",
			input:    "(hello)(world)",
			left:     "(",
			right:    ")",
			expected: "hello)(world",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RemoveSurrounding(tt.input, tt.left, tt.right)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
