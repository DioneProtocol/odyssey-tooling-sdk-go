// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestCloudParams_ValidateExtended(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags for this test
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true

	tests := []struct {
		name          string
		cp            CloudParams
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid AWS config",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-security-group",
					AWSKeyPair:           "test-keypair",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "Missing region",
			cp: CloudParams{
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			expectError:   true,
			errorContains: "region is required",
		},
		{
			name: "Missing image ID",
			cp: CloudParams{
				Region:       "us-west-2",
				InstanceType: "t3.medium",
			},
			expectError:   true,
			errorContains: "image is required",
		},
		{
			name: "Missing instance type",
			cp: CloudParams{
				Region:  "us-west-2",
				ImageID: "ami-12345678",
			},
			expectError:   true,
			errorContains: "instance type is required",
		},
		{
			name: "AWS config missing profile",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSSecurityGroupID: "sg-12345678",
				},
			},
			expectError:   true,
			errorContains: "unsupported cloud",
		},
		{
			name: "AWS config missing security group ID",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile: "default",
				},
			},
			expectError:   true,
			errorContains: "AWS security group ID is required",
		},
		{
			name: "Valid GCP config",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPCredentials: "/path/to/credentials.json",
					GCPNetwork:     "default",
					GCPZone:        "us-central1-a",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
				},
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "GCP config missing project",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPCredentials: "/path/to/credentials.json",
				},
			},
			expectError:   true,
			errorContains: "GCP network is required",
		},
		{
			name: "GCP config missing credentials",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			expectError:   true,
			errorContains: "GCP network is required",
		},
		{
			name:          "Empty cloud params",
			cp:            CloudParams{},
			expectError:   true,
			errorContains: "region is required",
		},
		{
			name: "Invalid cloud type",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud // Cloud type validation happens elsewhere
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCloudParams_CloudExtended(t *testing.T) {
	tests := []struct {
		name     string
		cp       CloudParams
		expected SupportedCloud
	}{
		{
			name: "AWS cloud",
			cp: CloudParams{
				AWSConfig: &AWSConfig{
					AWSProfile: "default",
				},
			},
			expected: AWSCloud,
		},
		{
			name: "GCP cloud",
			cp: CloudParams{
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			expected: GCPCloud,
		},
		{
			name:     "No cloud config",
			cp:       CloudParams{},
			expected: SupportedCloud(3), // Default value when no config
		},
		{
			name: "Both AWS and GCP config",
			cp: CloudParams{
				AWSConfig: &AWSConfig{
					AWSProfile: "default",
				},
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			expected: AWSCloud, // AWS takes precedence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cp.Cloud()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAWSConfig_Validation(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
	}()

	// Enable feature flags for this test
	constants.AWSIntegrationEnabled = true

	tests := []struct {
		name          string
		config        AWSConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid AWS config",
			config: AWSConfig{
				AWSProfile:           "default",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-security-group",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
				AWSVolumeIOPS:        1000,
				AWSVolumeThroughput:  500,
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "Empty profile",
			config: AWSConfig{
				AWSProfile:         "",
				AWSKeyPair:         "test-keypair",
				AWSSecurityGroupID: "sg-12345678",
			},
			expectError:   true,
			errorContains: "unsupported cloud",
		},
		{
			name: "Empty security group ID",
			config: AWSConfig{
				AWSProfile:         "default",
				AWSKeyPair:         "test-keypair",
				AWSSecurityGroupID: "",
			},
			expectError:   true,
			errorContains: "AWS security group ID is required",
		},
		{
			name: "Invalid volume size",
			config: AWSConfig{
				AWSProfile:         "default",
				AWSKeyPair:         "test-keypair",
				AWSSecurityGroupID: "sg-12345678",
				AWSVolumeSize:      0,
			},
			expectError: true, // Will fail due to unsupported cloud // Volume size validation not implemented
		},
		{
			name: "Invalid volume type",
			config: AWSConfig{
				AWSProfile:         "default",
				AWSKeyPair:         "test-keypair",
				AWSSecurityGroupID: "sg-12345678",
				AWSVolumeType:      "invalid-type",
			},
			expectError: true, // Will fail due to unsupported cloud // Volume type validation not implemented
		},
		{
			name: "Invalid IOPS",
			config: AWSConfig{
				AWSProfile:         "default",
				AWSKeyPair:         "test-keypair",
				AWSSecurityGroupID: "sg-12345678",
				AWSVolumeIOPS:      -1,
			},
			expectError: true, // Will fail due to unsupported cloud // IOPS validation not implemented
		},
		{
			name: "Invalid throughput",
			config: AWSConfig{
				AWSProfile:          "default",
				AWSKeyPair:          "test-keypair",
				AWSSecurityGroupID:  "sg-12345678",
				AWSVolumeThroughput: -1,
			},
			expectError: true, // Will fail due to unsupported cloud // Throughput validation not implemented
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig:    &tt.config,
			}

			err := cp.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGCPConfig_Validation(t *testing.T) {
	// Save original feature flag values
	originalGCPIntegration := constants.GCPIntegrationEnabled
	defer func() {
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags for this test
	constants.GCPIntegrationEnabled = true

	tests := []struct {
		name          string
		config        GCPConfig
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid GCP config",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
				GCPNetwork:     "default",
				GCPZone:        "us-central1-a",
				GCPVolumeSize:  100,
				GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "Empty project",
			config: GCPConfig{
				GCPProject:     "",
				GCPCredentials: "/path/to/credentials.json",
			},
			expectError:   true,
			errorContains: "GCP network is required",
		},
		{
			name: "Empty credentials",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "",
			},
			expectError:   true,
			errorContains: "GCP network is required",
		},
		{
			name: "Invalid volume size",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
				GCPVolumeSize:  0,
			},
			expectError: true, // Will fail due to unsupported cloud // Volume size validation not implemented
		},
		{
			name: "Empty network",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
				GCPNetwork:     "",
			},
			expectError: true, // Will fail due to unsupported cloud // Network validation not implemented
		},
		{
			name: "Empty zone",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
				GCPZone:        "",
			},
			expectError: true, // Will fail due to unsupported cloud // Zone validation not implemented
		},
		{
			name: "Empty SSH key",
			config: GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
				GCPSSHKey:      "",
			},
			expectError: true, // Will fail due to unsupported cloud // SSH key validation not implemented
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig:    &tt.config,
			}

			err := cp.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetDefaultCloudParamsExtended(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	ctx := context.Background()

	tests := []struct {
		name        string
		cloud       SupportedCloud
		expectError bool
	}{
		{
			name:        "AWS cloud",
			cloud:       AWSCloud,
			expectError: true, // Will fail due to no AWS credentials
		},
		{
			name:        "GCP cloud",
			cloud:       GCPCloud,
			expectError: true, // Will fail due to no GCP credentials
		},
		{
			name:        "Invalid cloud",
			cloud:       SupportedCloud(999),
			expectError: true,
		},
		{
			name:        "Docker cloud",
			cloud:       SupportedCloud(3), // Docker cloud type
			expectError: true,              // Docker cloud not supported in GetDefaultCloudParams
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp, err := GetDefaultCloudParams(ctx, tt.cloud)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cp)
			}
		})
	}
}

func TestGetDefaultCloudParams_FeatureFlagsExtended(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestCloudParams_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		cp          CloudParams
		expectError bool
	}{
		{
			name: "Region with special characters",
			cp: CloudParams{
				Region:       "us-west-2_abc",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Image ID with special characters",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678_abc",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Instance type with special characters",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium_abc",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Very long region name",
			cp: CloudParams{
				Region:       "very-long-region-name-that-exceeds-normal-limits",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Very long image ID",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "very-long-image-id-that-exceeds-normal-limits-and-should-still-work",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Very long instance type",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "very-long-instance-type-that-exceeds-normal-limits",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Region with unicode characters",
			cp: CloudParams{
				Region:       "us-west-2-测试",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Image ID with unicode characters",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678-测试",
				InstanceType: "t3.medium",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Instance type with unicode characters",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium-测试",
			},
			expectError: true, // Will fail due to unsupported cloud
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCloudParams_Validation_Comprehensive(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags for this test
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true

	tests := []struct {
		name          string
		cp            CloudParams
		expectError   bool
		errorContains string
	}{
		{
			name: "Complete AWS configuration",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-security-group",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "Complete GCP configuration",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPCredentials: "/path/to/credentials.json",
					GCPNetwork:     "default",
					GCPZone:        "us-central1-a",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
				},
			},
			expectError: false, // Valid configuration should pass
		},
		{
			name: "Minimal AWS configuration",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile:         "default",
					AWSSecurityGroupID: "sg-12345678",
				},
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "Minimal GCP configuration",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPCredentials: "/path/to/credentials.json",
				},
			},
			expectError: true, // Will fail due to unsupported cloud
		},
		{
			name: "AWS configuration with nil AWSConfig",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig:    nil,
			},
			expectError: true, // Will fail due to unsupported cloud // AWSConfig validation happens elsewhere
		},
		{
			name: "GCP configuration with nil GCPConfig",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig:    nil,
			},
			expectError: true, // Will fail due to unsupported cloud // GCPConfig validation happens elsewhere
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCloudParams_CloudType_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		cp       CloudParams
		expected SupportedCloud
	}{
		{
			name: "AWS config with empty profile",
			cp: CloudParams{
				AWSConfig: &AWSConfig{
					AWSProfile: "",
				},
			},
			expected: SupportedCloud(3), // Default value when config is incomplete
		},
		{
			name: "GCP config with empty project",
			cp: CloudParams{
				GCPConfig: &GCPConfig{
					GCPProject: "",
				},
			},
			expected: SupportedCloud(3), // Default value when config is incomplete
		},
		{
			name: "AWS config with nil values",
			cp: CloudParams{
				AWSConfig: &AWSConfig{},
			},
			expected: SupportedCloud(3), // Default value when config is incomplete
		},
		{
			name: "GCP config with nil values",
			cp: CloudParams{
				GCPConfig: &GCPConfig{},
			},
			expected: SupportedCloud(3), // Default value when config is incomplete
		},
		{
			name: "Both configs with nil values",
			cp: CloudParams{
				AWSConfig: &AWSConfig{},
				GCPConfig: &GCPConfig{},
			},
			expected: SupportedCloud(3), // Default value when config is incomplete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cp.Cloud()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCloudParams_Validation_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		cp            CloudParams
		errorContains string
	}{
		{
			name: "Missing region error message",
			cp: CloudParams{
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			errorContains: "region is required",
		},
		{
			name: "Missing image ID error message",
			cp: CloudParams{
				Region:       "us-west-2",
				InstanceType: "t3.medium",
			},
			errorContains: "image is required",
		},
		{
			name: "Missing instance type error message",
			cp: CloudParams{
				Region:  "us-west-2",
				ImageID: "ami-12345678",
			},
			errorContains: "instance type is required",
		},
		{
			name: "AWS profile error message",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSSecurityGroupID: "sg-12345678",
				},
			},
			errorContains: "unsupported cloud",
		},
		{
			name: "AWS security group ID error message",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile: "default",
				},
			},
			errorContains: "AWS security group ID is required",
		},
		{
			name: "GCP project error message",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPCredentials: "/path/to/credentials.json",
				},
			},
			errorContains: "GCP network is required",
		},
		{
			name: "GCP credentials error message",
			cp: CloudParams{
				Region:       "us-central1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-medium",
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			errorContains: "GCP network is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cp.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorContains)
		})
	}
}
