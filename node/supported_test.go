// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSupportedCloud(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SupportedCloud
	}{
		{
			name:     "AWS cloud",
			input:    "aws",
			expected: AWSCloud,
		},
		{
			name:     "GCP cloud",
			input:    "gcp",
			expected: GCPCloud,
		},
		{
			name:     "Docker cloud",
			input:    "docker",
			expected: Docker,
		},
		{
			name:     "Unknown cloud",
			input:    "unknown",
			expected: Unknown,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: Unknown,
		},
		{
			name:     "Case sensitive",
			input:    "AWS",
			expected: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSupportedCloud(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSupportedCloud_String(t *testing.T) {
	tests := []struct {
		name     string
		cloud    SupportedCloud
		expected string
	}{
		{
			name:     "AWS cloud",
			cloud:    AWSCloud,
			expected: "aws",
		},
		{
			name:     "GCP cloud",
			cloud:    GCPCloud,
			expected: "gcp",
		},
		{
			name:     "Docker cloud",
			cloud:    Docker,
			expected: "docker",
		},
		{
			name:     "Unknown cloud",
			cloud:    Unknown,
			expected: "unknown",
		},
		{
			name:     "Invalid cloud value",
			cloud:    SupportedCloud(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cloud.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewSupportedRole(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SupportedRole
	}{
		{
			name:     "Validator role",
			input:    "validator",
			expected: Validator,
		},
		{
			name:     "API role",
			input:    "api",
			expected: API,
		},
		{
			name:     "AWM Relayer role",
			input:    "awm-relayer",
			expected: AWMRelayer,
		},
		{
			name:     "Load test role",
			input:    "loadtest",
			expected: Loadtest,
		},
		{
			name:     "Monitor role",
			input:    "monitor",
			expected: Monitor,
		},
		{
			name:     "Unknown role defaults to Monitor",
			input:    "unknown",
			expected: Monitor,
		},
		{
			name:     "Empty string defaults to Monitor",
			input:    "",
			expected: Monitor,
		},
		{
			name:     "Case sensitive",
			input:    "Validator",
			expected: Monitor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSupportedRole(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSupportedRole_String(t *testing.T) {
	tests := []struct {
		name     string
		role     SupportedRole
		expected string
	}{
		{
			name:     "Validator role",
			role:     Validator,
			expected: "validator",
		},
		{
			name:     "API role",
			role:     API,
			expected: "api",
		},
		{
			name:     "AWM Relayer role",
			role:     AWMRelayer,
			expected: "awm-relayer",
		},
		{
			name:     "Load test role",
			role:     Loadtest,
			expected: "loadtest",
		},
		{
			name:     "Monitor role",
			role:     Monitor,
			expected: "monitor",
		},
		{
			name:     "Unknown role",
			role:     SupportedRole(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckRoles(t *testing.T) {
	tests := []struct {
		name        string
		roles       []SupportedRole
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Single validator role",
			roles:       []SupportedRole{Validator},
			expectError: false,
		},
		{
			name:        "Single API role",
			roles:       []SupportedRole{API},
			expectError: false,
		},
		{
			name:        "Single monitor role",
			roles:       []SupportedRole{Monitor},
			expectError: false,
		},
		{
			name:        "Single load test role",
			roles:       []SupportedRole{Loadtest},
			expectError: false,
		},
		{
			name:        "Single AWM relayer role",
			roles:       []SupportedRole{AWMRelayer},
			expectError: false,
		},
		{
			name:        "Empty roles",
			roles:       []SupportedRole{},
			expectError: false,
		},
		{
			name:        "Validator and API roles - should fail",
			roles:       []SupportedRole{Validator, API},
			expectError: true,
			errorMsg:    "cannot have both validator and api roles",
		},
		{
			name:        "Load test with other roles - should fail",
			roles:       []SupportedRole{Loadtest, Monitor},
			expectError: true,
			errorMsg:    "role cannot be combined with other roles",
		},
		{
			name:        "Load test with validator - should fail",
			roles:       []SupportedRole{Loadtest, Validator},
			expectError: true,
			errorMsg:    "role cannot be combined with other roles",
		},
		{
			name:        "Monitor with other roles - should fail",
			roles:       []SupportedRole{Monitor, Validator},
			expectError: true,
			errorMsg:    "role cannot be combined with other roles",
		},
		{
			name:        "Monitor with API - should fail",
			roles:       []SupportedRole{Monitor, API},
			expectError: true,
			errorMsg:    "role cannot be combined with other roles",
		},
		{
			name:        "AWM relayer with other roles - should be allowed",
			roles:       []SupportedRole{AWMRelayer, Validator},
			expectError: false,
		},
		{
			name:        "AWM relayer with API - should be allowed",
			roles:       []SupportedRole{AWMRelayer, API},
			expectError: false,
		},
		{
			name:        "Multiple valid roles",
			roles:       []SupportedRole{Validator, AWMRelayer},
			expectError: false,
		},
		{
			name:        "API with AWM relayer",
			roles:       []SupportedRole{API, AWMRelayer},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRoles(tt.roles)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSupportedCloud_Constants(t *testing.T) {
	// Test that the constants have the expected values
	assert.Equal(t, SupportedCloud(0), AWSCloud)
	assert.Equal(t, SupportedCloud(1), GCPCloud)
	assert.Equal(t, SupportedCloud(2), Docker)
	assert.Equal(t, SupportedCloud(3), Unknown)
}

func TestSupportedRole_Constants(t *testing.T) {
	// Test that the constants have the expected values
	assert.Equal(t, SupportedRole(0), Validator)
	assert.Equal(t, SupportedRole(1), API)
	assert.Equal(t, SupportedRole(2), AWMRelayer)
	assert.Equal(t, SupportedRole(3), Loadtest)
	assert.Equal(t, SupportedRole(4), Monitor)
}

func TestSupportedCloud_EdgeCases(t *testing.T) {
	// Test with negative values
	negativeCloud := SupportedCloud(-1)
	assert.Equal(t, "unknown", negativeCloud.String())

	// Test with very large values
	largeCloud := SupportedCloud(999999)
	assert.Equal(t, "unknown", largeCloud.String())
}

func TestSupportedRole_EdgeCases(t *testing.T) {
	// Test with negative values
	negativeRole := SupportedRole(-1)
	assert.Equal(t, "unknown", negativeRole.String())

	// Test with very large values
	largeRole := SupportedRole(999999)
	assert.Equal(t, "unknown", largeRole.String())
}

func TestCheckRoles_EdgeCases(t *testing.T) {
	// Test with nil slice
	err := CheckRoles(nil)
	assert.NoError(t, err)

	// Test with duplicate roles
	duplicateRoles := []SupportedRole{Validator, Validator}
	err = CheckRoles(duplicateRoles)
	assert.NoError(t, err) // Duplicates are allowed

	// Test with all roles (should fail due to conflicts)
	allRoles := []SupportedRole{Validator, API, AWMRelayer, Loadtest, Monitor}
	err = CheckRoles(allRoles)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot have both validator and api roles")
}

func TestRoleCombinations(t *testing.T) {
	// Test all valid single role combinations
	singleRoles := []SupportedRole{Validator, API, AWMRelayer, Loadtest, Monitor}
	for _, role := range singleRoles {
		t.Run("Single_"+role.String(), func(t *testing.T) {
			err := CheckRoles([]SupportedRole{role})
			assert.NoError(t, err)
		})
	}

	// Test valid multi-role combinations
	validCombinations := [][]SupportedRole{
		{Validator, AWMRelayer},
		{API, AWMRelayer},
	}

	for i, roles := range validCombinations {
		t.Run("Valid_Combination_"+string(rune(i+'0')), func(t *testing.T) {
			err := CheckRoles(roles)
			assert.NoError(t, err)
		})
	}

	// Test invalid multi-role combinations
	invalidCombinations := [][]SupportedRole{
		{Validator, API},
		{Loadtest, Validator},
		{Loadtest, API},
		{Loadtest, Monitor},
		{Loadtest, AWMRelayer},
		{Monitor, Validator},
		{Monitor, API},
		{Monitor, AWMRelayer},
	}

	for i, roles := range invalidCombinations {
		t.Run("Invalid_Combination_"+string(rune(i+'0')), func(t *testing.T) {
			err := CheckRoles(roles)
			assert.Error(t, err)
		})
	}
}
