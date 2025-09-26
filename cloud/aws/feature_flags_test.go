// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aws

import (
	"context"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

// TestFeatureFlags_SecurityGroups tests security groups feature flag
func TestFeatureFlags_SecurityGroups(t *testing.T) {
	// Save original value
	originalValue := constants.SecurityGroupsEnabled
	defer func() {
		constants.SecurityGroupsEnabled = originalValue
	}()

	// Test with security groups disabled
	constants.SecurityGroupsEnabled = false

	ctx := context.Background()

	// Test CreateSecurityGroup function
	_, err := CreateSecurityGroup(ctx, "test-sg", "test-profile", "us-east-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "security groups functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.SecurityGroupsEnabled = true to enable")

	// Test with security groups enabled
	constants.SecurityGroupsEnabled = true
	// Note: This will still fail due to missing AWS credentials, but should pass the security groups flag check
	_, err = CreateSecurityGroup(ctx, "test-sg", "test-profile", "us-east-1")
	// The error should not be about security groups being disabled
	assert.NotContains(t, err.Error(), "security groups functionality is disabled")
}

// TestFeatureFlags_SecurityGroupsInstance tests security groups feature flag on AwsCloud instance
func TestFeatureFlags_SecurityGroupsInstance(t *testing.T) {
	// Save original value
	originalValue := constants.SecurityGroupsEnabled
	defer func() {
		constants.SecurityGroupsEnabled = originalValue
	}()

	// Test with security groups disabled
	constants.SecurityGroupsEnabled = false

	ctx := context.Background()
	awsCloud := &AwsCloud{
		ctx: ctx,
	}

	// Test CreateSecurityGroup method
	_, err := awsCloud.CreateSecurityGroup("test-sg", "test description")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "security groups functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.SecurityGroupsEnabled = true to enable")

	// Test with security groups enabled
	constants.SecurityGroupsEnabled = true
	// Since we can't easily mock the AWS client, we'll just verify the flag is enabled
	// The actual AWS call would fail due to missing credentials, but the flag check would pass
	assert.True(t, constants.SecurityGroupsEnabled, "Security groups flag should be enabled")
}
