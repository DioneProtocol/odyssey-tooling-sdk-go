// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCloudInstances_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	ctx := context.Background()
	cp := &CloudParams{
		Region:       "us-west-2",
		ImageID:      "ami-12345678",
		InstanceType: "t3.micro",
		AWSConfig: &AWSConfig{
			AWSProfile: "test-profile",
		},
	}

	tests := []struct {
		name            string
		awsEnabled      bool
		gcpEnabled      bool
		instanceEnabled bool
		sshEnabled      bool
		cloudType       SupportedCloud
		expectError     bool
		errorContains   string
	}{
		{
			name:            "AWS integration disabled",
			awsEnabled:      false,
			gcpEnabled:      true,
			instanceEnabled: true,
			sshEnabled:      true,
			cloudType:       AWSCloud,
			expectError:     true,
			errorContains:   "AWS integration functionality is disabled",
		},
		{
			name:            "GCP integration disabled",
			awsEnabled:      true,
			gcpEnabled:      false,
			instanceEnabled: true,
			sshEnabled:      true,
			cloudType:       GCPCloud,
			expectError:     true,
			errorContains:   "GCP integration functionality is disabled",
		},
		{
			name:            "Instance management disabled",
			awsEnabled:      true,
			gcpEnabled:      true,
			instanceEnabled: false,
			sshEnabled:      true,
			cloudType:       AWSCloud,
			expectError:     true,
			errorContains:   "instance management functionality is disabled",
		},
		{
			name:            "SSH key management disabled",
			awsEnabled:      true,
			gcpEnabled:      true,
			instanceEnabled: true,
			sshEnabled:      false,
			cloudType:       AWSCloud,
			expectError:     true,
			errorContains:   "SSH key management functionality is disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flags
			constants.AWSIntegrationEnabled = tt.awsEnabled
			constants.GCPIntegrationEnabled = tt.gcpEnabled
			constants.InstanceManagementEnabled = tt.instanceEnabled
			constants.SSHKeyManagementEnabled = tt.sshEnabled

			// Set cloud type by configuring the appropriate config
			if tt.cloudType == AWSCloud {
				cp.AWSConfig = &AWSConfig{AWSProfile: "test-profile"}
			} else if tt.cloudType == GCPCloud {
				cp.GCPConfig = &GCPConfig{GCPProject: "test-project"}
			}

			nodes, err := createCloudInstances(ctx, *cp, 1, false, "/path/to/key")

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, nodes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, nodes)
			}
		})
	}
}

func TestCreateCloudInstances_ValidationErrors(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags for validation tests
	constants.AWSIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name          string
		cp            CloudParams
		count         int
		sshKeyPath    string
		expectError   bool
		errorContains string
	}{
		{
			name: "Invalid cloud params",
			cp: CloudParams{
				Region: "", // Invalid region
			},
			count:         1,
			sshKeyPath:    "/path/to/key",
			expectError:   true,
			errorContains: "region is required",
		},
		{
			name: "Zero count",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			},
			count:         0,
			sshKeyPath:    "/path/to/key",
			expectError:   true,
			errorContains: "count must be greater than 0",
		},
		{
			name: "Empty SSH key path",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			},
			count:         1,
			sshKeyPath:    "",
			expectError:   true,
			errorContains: "SSH private key path is required",
		},
		{
			name: "Negative count",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			},
			count:         -1,
			sshKeyPath:    "/path/to/key",
			expectError:   true,
			errorContains: "count must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := createCloudInstances(ctx, tt.cp, tt.count, false, tt.sshKeyPath)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, nodes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, nodes)
			}
		})
	}
}

func TestCreateCloudInstances_CloudTypes(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name          string
		cp            CloudParams
		expectError   bool
		errorContains string
	}{
		{
			name: "AWS cloud type",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			},
			expectError:   true, // Will fail due to AWS API calls
			errorContains: "failed to create all instances",
		},
		{
			name: "GCP cloud type",
			cp: CloudParams{
				Region:       "us-west1",
				ImageID:      "projects/debian-cloud/global/images/family/debian-11",
				InstanceType: "e2-micro",
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			expectError:   true, // Will fail due to GCP API calls
			errorContains: "failed to create all instances",
		},
		{
			name: "Unknown cloud type",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				// No AWS or GCP config - will result in Unknown cloud type
			},
			expectError:   true,
			errorContains: "unsupported cloud type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := createCloudInstances(ctx, tt.cp, 1, false, "/path/to/key")

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, nodes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, nodes)
			}
		})
	}
}

func TestCreateCloudInstances_StaticIP(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags
	constants.AWSIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name        string
		useStaticIP bool
		expectError bool
	}{
		{
			name:        "Without static IP",
			useStaticIP: false,
			expectError: true, // Will fail due to AWS API calls
		},
		{
			name:        "With static IP",
			useStaticIP: true,
			expectError: true, // Will fail due to AWS API calls
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			}

			nodes, err := createCloudInstances(ctx, cp, 1, tt.useStaticIP, "/path/to/key")

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, nodes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, nodes)
			}
		})
	}
}

func TestCreateCloudInstances_EdgeCasesExtended(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags
	constants.AWSIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name        string
		count       int
		expectError bool
	}{
		{
			name:        "Single instance",
			count:       1,
			expectError: true, // Will fail due to AWS API calls
		},
		{
			name:        "Multiple instances",
			count:       3,
			expectError: true, // Will fail due to AWS API calls
		},
		{
			name:        "Large count",
			count:       10,
			expectError: true, // Will fail due to AWS API calls
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.micro",
				AWSConfig: &AWSConfig{
					AWSProfile: "test-profile",
				},
			}

			nodes, err := createCloudInstances(ctx, cp, tt.count, false, "/path/to/key")

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, nodes)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, nodes)
				assert.Len(t, nodes, tt.count)
			}
		})
	}
}
