// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAny(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "empty slice",
			input:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  false,
		},
		{
			name:      "predicate always true",
			input:     []int{1, 2, 3},
			predicate: func(x int) bool { return true },
			expected:  true,
		},
		{
			name:      "predicate always false",
			input:     []int{1, 2, 3},
			predicate: func(x int) bool { return false },
			expected:  false,
		},
		{
			name:      "predicate true for some elements",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 3 },
			expected:  true,
		},
		{
			name:      "predicate true for first element",
			input:     []int{5, 1, 2, 3},
			predicate: func(x int) bool { return x > 3 },
			expected:  true,
		},
		{
			name:      "predicate true for last element",
			input:     []int{1, 2, 3, 5},
			predicate: func(x int) bool { return x > 3 },
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Any(tt.input, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with different data types
	t.Run("string slice", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		predicate := func(s string) bool { return len(s) > 4 }
		result := Any(input, predicate)
		assert.True(t, result)
	})
}

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  *int
	}{
		{
			name:      "empty slice",
			input:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  nil,
		},
		{
			name:      "find first element",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 2 },
			expected:  func() *int { v := 3; return &v }(),
		},
		{
			name:      "find middle element",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x == 3 },
			expected:  func() *int { v := 3; return &v }(),
		},
		{
			name:      "find last element",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x == 5 },
			expected:  func() *int { v := 5; return &v }(),
		},
		{
			name:      "element not found",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 10 },
			expected:  nil,
		},
		{
			name:      "predicate always false",
			input:     []int{1, 2, 3},
			predicate: func(x int) bool { return false },
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Find(tt.input, tt.predicate)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestBelongs(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		elem     int
		expected bool
	}{
		{
			name:     "empty slice",
			input:    []int{},
			elem:     1,
			expected: false,
		},
		{
			name:     "element exists",
			input:    []int{1, 2, 3, 4, 5},
			elem:     3,
			expected: true,
		},
		{
			name:     "element does not exist",
			input:    []int{1, 2, 3, 4, 5},
			elem:     6,
			expected: false,
		},
		{
			name:     "element is first",
			input:    []int{1, 2, 3, 4, 5},
			elem:     1,
			expected: true,
		},
		{
			name:     "element is last",
			input:    []int{1, 2, 3, 4, 5},
			elem:     5,
			expected: true,
		},
		{
			name:     "duplicate elements",
			input:    []int{1, 2, 2, 3, 4},
			elem:     2,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Belongs(tt.input, tt.elem)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with different data types
	t.Run("string slice", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		result := Belongs(input, "world")
		assert.True(t, result)

		result = Belongs(input, "missing")
		assert.False(t, result)
	})
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "empty slice",
			input:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
		{
			name:      "filter all elements",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{1, 2, 3, 4, 5},
		},
		{
			name:      "filter no elements",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 10 },
			expected:  []int{},
		},
		{
			name:      "filter some elements",
			input:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 3 },
			expected:  []int{4, 5},
		},
		{
			name:      "filter even numbers",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with different data types
	t.Run("string slice", func(t *testing.T) {
		input := []string{"hello", "world", "test", "go"}
		predicate := func(s string) bool { return len(s) > 3 }
		expected := []string{"hello", "world", "test"}
		result := Filter(input, predicate)
		assert.Equal(t, expected, result)
	})
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapper   func(int) string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []int{},
			mapper:   func(x int) string { return string(rune(x)) },
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []int{65},
			mapper:   func(x int) string { return string(rune(x)) },
			expected: []string{"A"},
		},
		{
			name:     "multiple elements",
			input:    []int{65, 66, 67},
			mapper:   func(x int) string { return string(rune(x)) },
			expected: []string{"A", "B", "C"},
		},
		{
			name:     "convert to string",
			input:    []int{1, 2, 3, 4, 5},
			mapper:   func(x int) string { return string(rune('0' + x)) },
			expected: []string{"1", "2", "3", "4", "5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.mapper)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test with different input/output types
	t.Run("string to int", func(t *testing.T) {
		input := []string{"1", "2", "3", "4", "5"}
		mapper := func(s string) int { return len(s) }
		expected := []int{1, 1, 1, 1, 1}
		result := Map(input, mapper)
		assert.Equal(t, expected, result)
	})
}

func TestMapWithError(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		mapper   func(int) (string, error)
		expected []string
		wantErr  bool
	}{
		{
			name:     "empty slice",
			input:    []int{},
			mapper:   func(x int) (string, error) { return string(rune(x)), nil },
			expected: []string{},
			wantErr:  false,
		},
		{
			name:     "successful mapping",
			input:    []int{65, 66, 67},
			mapper:   func(x int) (string, error) { return string(rune(x)), nil },
			expected: []string{"A", "B", "C"},
			wantErr:  false,
		},
		{
			name:  "error in mapping",
			input: []int{65, 66, 67},
			mapper: func(x int) (string, error) {
				if x == 66 {
					return "", errors.New("mapping error")
				}
				return string(rune(x)), nil
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "error on first element",
			input: []int{65, 66, 67},
			mapper: func(x int) (string, error) {
				if x == 65 {
					return "", errors.New("mapping error")
				}
				return string(rune(x)), nil
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapWithError(tt.input, tt.mapper)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAppendSlices(t *testing.T) {
	tests := []struct {
		name   string
		slices [][]interface{}
		want   []interface{}
	}{
		{
			name:   "AppendSlices with strings",
			slices: [][]interface{}{{"a", "b", "c"}, {"d", "e", "f"}, {"g", "h", "i"}},
			want:   []interface{}{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
		},
		{
			name:   "AppendSlices with ints",
			slices: [][]interface{}{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			want:   []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:   "AppendSlices with empty slices",
			slices: [][]interface{}{{}, {}, {}},
			want:   []interface{}{},
		},
		{
			name:   "Append identical slices",
			slices: [][]interface{}{{"a", "b", "c"}, {"a", "b", "c"}},
			want:   []interface{}{"a", "b", "c", "a", "b", "c"},
		},
		{
			name:   "Single slice",
			slices: [][]interface{}{{"a", "b", "c"}},
			want:   []interface{}{"a", "b", "c"},
		},
		{
			name:   "No slices",
			slices: [][]interface{}{},
			want:   []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AppendSlices(tt.slices...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Mock function for testing retries.
func mockFunction() (interface{}, error) {
	return nil, errors.New("error occurred")
}

func mockSuccessFunction() (interface{}, error) {
	return "success", nil
}

func TestRetry(t *testing.T) {
	success := "success"
	// Test with a function that always returns an error.
	result, err := Retry(WrapContext(mockFunction), 100*time.Millisecond, 3, "test error")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Test with a function that succeeds on the first attempt.
	result, err = Retry(WrapContext(mockSuccessFunction), 100*time.Millisecond, 3, "test error")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != success {
		t.Errorf("Expected 'success' result, got %v", result)
	}

	// Test with a function that succeeds after multiple attempts.
	count := 0
	fn := func() (interface{}, error) {
		count++
		if count < 3 {
			return nil, errors.New("error occurred")
		}
		return success, nil
	}
	result, err = Retry(WrapContext(fn), 100*time.Millisecond, 5, "test error")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != success {
		t.Errorf("Expected 'success' result, got %v", result)
	}

	// Test with invalid retry interval (should use default).
	result, err = Retry(WrapContext(mockFunction), 0, 3, "test error")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Test error message format
	result, err = Retry(WrapContext(mockFunction), 100*time.Millisecond, 2, "custom error")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "custom error")
	assert.Contains(t, err.Error(), "maximum retry attempts 2 reached")
}

func TestWrapContext(t *testing.T) {
	// Test with function that completes before timeout
	fn := func() (string, error) {
		return "success", nil
	}
	wrapped := WrapContext(fn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := wrapped(ctx)
	require.NoError(t, err)
	assert.Equal(t, "success", result)

	// Test with function that times out
	slowFn := func() (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "success", nil
	}
	wrappedSlow := WrapContext(slowFn)

	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result, err = wrappedSlow(ctx)
	require.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Empty(t, result)

	// Test with function that returns error
	errorFn := func() (string, error) {
		return "", errors.New("function error")
	}
	wrappedError := WrapContext(errorFn)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err = wrappedError(ctx)
	require.Error(t, err)
	assert.Equal(t, "function error", err.Error())
	assert.Empty(t, result)
}

func TestCallWithTimeout(t *testing.T) {
	// Test with function that completes quickly
	fn := func() (string, error) {
		return "success", nil
	}

	result, err := CallWithTimeout("test", fn, 100*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "success", result)

	// Test with function that times out
	slowFn := func() (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "success", nil
	}

	result, err = CallWithTimeout("slow test", slowFn, 50*time.Millisecond)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slow test timeout of 0 seconds")
	assert.Empty(t, result)

	// Test with function that returns error
	errorFn := func() (string, error) {
		return "", errors.New("function error")
	}

	result, err = CallWithTimeout("error test", errorFn, 100*time.Millisecond)
	require.Error(t, err)
	assert.Equal(t, "function error", err.Error())
	assert.Empty(t, result)
}

func TestRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 0",
			length: 0,
		},
		{
			name:   "length 1",
			length: 1,
		},
		{
			name:   "length 10",
			length: 10,
		},
		{
			name:   "length 100",
			length: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.length)
			assert.Len(t, result, tt.length)

			// Test that generated strings contain only expected characters
			for _, char := range result {
				assert.True(t, char >= 'a' && char <= 'z',
					"Character %c should be between 'a' and 'z'", char)
			}
		})
	}

	// Test that generated strings are different
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result := RandomString(10)
		assert.False(t, results[result], "Generated duplicate string: %s", result)
		results[result] = true
	}
}

func TestSupportedOdysseyGoArch(t *testing.T) {
	archs := SupportedOdysseyGoArch()
	assert.NotEmpty(t, archs)
	assert.Contains(t, archs, string(types.ArchitectureTypeArm64))
	assert.Contains(t, archs, string(types.ArchitectureTypeX8664))
}

func TestArchSupported(t *testing.T) {
	tests := []struct {
		name     string
		arch     string
		expected bool
	}{
		{
			name:     "supported arm64",
			arch:     string(types.ArchitectureTypeArm64),
			expected: true,
		},
		{
			name:     "supported x86_64",
			arch:     string(types.ArchitectureTypeX8664),
			expected: true,
		},
		{
			name:     "unsupported architecture",
			arch:     "unsupported",
			expected: false,
		},
		{
			name:     "empty architecture",
			arch:     "",
			expected: false,
		},
		{
			name:     "case sensitive",
			arch:     "ARM64",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ArchSupported(tt.arch)
			assert.Equal(t, tt.expected, result)
		})
	}
}
