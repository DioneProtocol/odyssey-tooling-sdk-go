// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package odyssey

import (
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestVMType_RepoName(t *testing.T) {
	tests := []struct {
		name     string
		vmType   VMType
		expected string
	}{
		{
			name:     "SubnetEvm",
			vmType:   SubnetEvm,
			expected: constants.SubnetEVMRepoName,
		},
		{
			name:     "Unknown VM type",
			vmType:   VMType("UnknownVM"),
			expected: "unknown",
		},
		{
			name:     "Empty VM type",
			vmType:   VMType(""),
			expected: "unknown",
		},
		{
			name:     "Custom VM type",
			vmType:   VMType("CustomVM"),
			expected: "unknown",
		},
		{
			name:     "EVM VM type",
			vmType:   VMType("EVM"),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.vmType.RepoName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVMTypeConstants(t *testing.T) {
	// Test that the VMType constants are properly defined
	assert.Equal(t, "Subnet-EVM", SubnetEvm)
	assert.NotEmpty(t, SubnetEvm)
}

func TestVMType_String(t *testing.T) {
	// Test that VMType can be converted to string
	tests := []struct {
		name     string
		vmType   VMType
		expected string
	}{
		{
			name:     "SubnetEvm to string",
			vmType:   SubnetEvm,
			expected: "Subnet-EVM",
		},
		{
			name:     "Custom VM type to string",
			vmType:   VMType("CustomVM"),
			expected: "CustomVM",
		},
		{
			name:     "Empty VM type to string",
			vmType:   VMType(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tt.vmType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVMType_Comparison(t *testing.T) {
	// Test VMType comparison operations
	t.Run("Equal VM types", func(t *testing.T) {
		vm1 := VMType(SubnetEvm)
		vm2 := VMType(SubnetEvm)
		assert.Equal(t, vm1, vm2)
		assert.True(t, vm1 == vm2)
	})

	t.Run("Different VM types", func(t *testing.T) {
		vm1 := VMType(SubnetEvm)
		vm2 := VMType("DifferentVM")
		assert.NotEqual(t, vm1, vm2)
		assert.False(t, vm1 == vm2)
	})

	t.Run("Case sensitive comparison", func(t *testing.T) {
		vm1 := VMType(SubnetEvm)
		vm2 := VMType("subnet-evm") // lowercase
		assert.NotEqual(t, vm1, vm2)
		assert.False(t, vm1 == vm2)
	})
}

func TestVMType_RepoNameIntegration(t *testing.T) {
	// Test that RepoName returns the expected constant value
	expectedRepoName := constants.SubnetEVMRepoName
	actualRepoName := VMType(SubnetEvm).RepoName()

	assert.Equal(t, expectedRepoName, actualRepoName)
	assert.NotEmpty(t, actualRepoName)
}

func TestVMType_EdgeCases(t *testing.T) {
	t.Run("VMType with special characters", func(t *testing.T) {
		vmType := VMType("VM-With-Special_Chars.123")
		result := vmType.RepoName()
		assert.Equal(t, "unknown", result)
	})

	t.Run("VMType with spaces", func(t *testing.T) {
		vmType := VMType("VM With Spaces")
		result := vmType.RepoName()
		assert.Equal(t, "unknown", result)
	})

	t.Run("VMType with unicode characters", func(t *testing.T) {
		vmType := VMType("VM-With-Ünicode")
		result := vmType.RepoName()
		assert.Equal(t, "unknown", result)
	})

	t.Run("Very long VMType", func(t *testing.T) {
		longName := "VeryLongVMTypeNameThatExceedsNormalLengthAndTestsEdgeCases"
		vmType := VMType(longName)
		result := vmType.RepoName()
		assert.Equal(t, "unknown", result)
	})
}

func TestVMType_TypeAssertion(t *testing.T) {
	// Test that VMType can be used in type assertions and interfaces
	var vmType interface{} = VMType(SubnetEvm)

	// Test type assertion
	if v, ok := vmType.(VMType); ok {
		assert.Equal(t, VMType(SubnetEvm), v)
		assert.Equal(t, constants.SubnetEVMRepoName, v.RepoName())
	} else {
		t.Fatal("Type assertion failed")
	}
}

func TestVMType_ZeroValue(t *testing.T) {
	// Test zero value of VMType
	var zeroVMType VMType
	assert.Equal(t, VMType(""), zeroVMType)
	assert.Equal(t, "unknown", zeroVMType.RepoName())
	assert.Equal(t, "", string(zeroVMType))
}

func TestVMType_ConstantsConsistency(t *testing.T) {
	// Test that the constant values are consistent
	assert.Equal(t, "Subnet-EVM", SubnetEvm)

	// Test that the constant can be used in comparisons
	if SubnetEvm == "Subnet-EVM" {
		// This should be true
		assert.True(t, true)
	} else {
		t.Fatal("Constant value mismatch")
	}
}
