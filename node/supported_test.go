// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name:     "Monitor role",
			input:    "monitor",
			expected: Monitor,
		},
		{
			name:     "Loadtest role",
			input:    "loadtest",
			expected: Loadtest,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: SupportedRole(999), // Invalid role
		},
		{
			name:     "Case sensitive",
			input:    "VALIDATOR",
			expected: SupportedRole(999), // Invalid role
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
			name:     "Monitor role",
			role:     Monitor,
			expected: "monitor",
		},
		{
			name:     "Loadtest role",
			role:     Loadtest,
			expected: "loadtest",
		},
		{
			name:     "Invalid role",
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

func TestSupportedRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     SupportedRole
		expected bool
	}{
		{
			name:     "Valid validator role",
			role:     Validator,
			expected: true,
		},
		{
			name:     "Valid API role",
			role:     API,
			expected: true,
		},
		{
			name:     "Valid monitor role",
			role:     Monitor,
			expected: true,
		},
		{
			name:     "Valid loadtest role",
			role:     Loadtest,
			expected: true,
		},
		{
			name:     "Invalid role",
			role:     SupportedRole(999),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test role validation through CheckRoles function
			err := CheckRoles([]SupportedRole{tt.role})
			result := err == nil
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCheckRoles(t *testing.T) {
	tests := []struct {
		name        string
		roles       []SupportedRole
		expectError bool
	}{
		{
			name:        "Empty roles",
			roles:       []SupportedRole{},
			expectError: true,
		},
		{
			name:        "Single valid role",
			roles:       []SupportedRole{Validator},
			expectError: false,
		},
		{
			name:        "Multiple valid roles",
			roles:       []SupportedRole{Validator, API},
			expectError: false,
		},
		{
			name:        "Invalid role",
			roles:       []SupportedRole{SupportedRole(999)},
			expectError: true,
		},
		{
			name:        "Mixed valid and invalid roles",
			roles:       []SupportedRole{Validator, SupportedRole(999)},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRoles(tt.roles)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
