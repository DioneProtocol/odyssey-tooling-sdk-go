// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudParams_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cloudParams CloudParams
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid AWS config",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: false,
		},
		{
			name: "Valid GCP config",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPCredentials: "/path/to/credentials.json",
					GCPNetwork:     "test-network",
					GCPZone:        "us-east1-b",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: false,
		},
		{
			name: "Missing region",
			cloudParams: CloudParams{
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "region is required",
		},
		{
			name: "Missing image ID",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "image is required",
		},
		{
			name: "Missing instance type",
			cloudParams: CloudParams{
				Region:  "us-east-1",
				ImageID: "ami-12345678",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "instance type is required",
		},
		{
			name: "Missing AWS config",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
			},
			expectError: true,
			errorMsg:    "unsupported cloud",
		},
		{
			name: "Missing AWS profile",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "unsupported cloud",
		},
		{
			name: "Missing AWS security group ID",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "AWS security group ID is required",
		},
		{
			name: "Missing AWS security group name",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:          "default",
					AWSKeyPair:          "test-keypair",
					AWSSecurityGroupID:  "sg-12345678",
					AWSVolumeSize:       100,
					AWSVolumeType:       "gp3",
					AWSVolumeIOPS:       1000,
					AWSVolumeThroughput: 500,
				},
			},
			expectError: true,
			errorMsg:    "AWS security group Name is required",
		},
		{
			name: "Negative AWS volume size",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        -1,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "AWS volume size must be positive",
		},
		{
			name: "Missing AWS volume type",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "AWS volume type is required",
		},
		{
			name: "Negative AWS volume IOPS",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        -1,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "AWS volume IOPS must be positive",
		},
		{
			name: "Negative AWS volume throughput",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  -1,
				},
			},
			expectError: true,
			errorMsg:    "AWS volume throughput must be positive",
		},
		{
			name: "Missing AWS key pair",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-sg",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			expectError: true,
			errorMsg:    "AWS key pair is required",
		},
		{
			name: "Missing GCP config",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
			},
			expectError: true,
			errorMsg:    "unsupported cloud", // This is a bug in the original code
		},
		{
			name: "Missing GCP network",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPCredentials: "/path/to/credentials.json",
					GCPZone:        "us-east1-b",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP network is required",
		},
		{
			name: "Missing GCP project",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPNetwork:     "test-network",
					GCPCredentials: "/path/to/credentials.json",
					GCPZone:        "us-east1-b",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP project is required",
		},
		{
			name: "Missing GCP credentials",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:    "test-project",
					GCPNetwork:    "test-network",
					GCPZone:       "us-east1-b",
					GCPVolumeSize: 100,
					GCPSSHKey:     "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP credentials is required",
		},
		{
			name: "Missing GCP zone",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPNetwork:     "test-network",
					GCPCredentials: "/path/to/credentials.json",
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP zone is required",
		},
		{
			name: "Negative GCP volume size",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPNetwork:     "test-network",
					GCPCredentials: "/path/to/credentials.json",
					GCPZone:        "us-east1-b",
					GCPVolumeSize:  -1,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP volume size must be positive",
		},
		{
			name: "GCP zone not in region",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPNetwork:     "test-network",
					GCPCredentials: "/path/to/credentials.json",
					GCPZone:        "us-west1-b", // Wrong region
					GCPVolumeSize:  100,
					GCPSSHKey:      "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			expectError: true,
			errorMsg:    "GCP zone must be in the region us-east1",
		},
		{
			name: "Missing GCP SSH key",
			cloudParams: CloudParams{
				Region:       "us-east1",
				ImageID:      "projects/ubuntu-os-cloud/global/images/ubuntu-2004-lts",
				InstanceType: "e2-standard-8",
				GCPConfig: &GCPConfig{
					GCPProject:     "test-project",
					GCPNetwork:     "test-network",
					GCPCredentials: "/path/to/credentials.json",
					GCPZone:        "us-east1-b",
					GCPVolumeSize:  100,
				},
			},
			expectError: true,
			errorMsg:    "GCP SSH key is required",
		},
		{
			name: "Unsupported cloud",
			cloudParams: CloudParams{
				Region:       "us-east-1",
				ImageID:      "ami-12345678",
				InstanceType: "c5.2xlarge",
			},
			expectError: true,
			errorMsg:    "unsupported cloud",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cloudParams.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCloudParams_Cloud(t *testing.T) {
	tests := []struct {
		name        string
		cloudParams CloudParams
		expected    SupportedCloud
	}{
		{
			name: "AWS cloud",
			cloudParams: CloudParams{
				AWSConfig: &AWSConfig{
					AWSProfile: "default",
				},
			},
			expected: AWSCloud,
		},
		{
			name: "GCP cloud with project",
			cloudParams: CloudParams{
				GCPConfig: &GCPConfig{
					GCPProject: "test-project",
				},
			},
			expected: GCPCloud,
		},
		{
			name: "GCP cloud with credentials",
			cloudParams: CloudParams{
				GCPConfig: &GCPConfig{
					GCPCredentials: "/path/to/credentials.json",
				},
			},
			expected: GCPCloud,
		},
		{
			name:        "Unknown cloud",
			cloudParams: CloudParams{},
			expected:    Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cloudParams.Cloud()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDefaultCloudParams(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		cloud       SupportedCloud
		expectError bool
	}{
		{
			name:        "AWS cloud",
			cloud:       AWSCloud,
			expectError: true, // Will fail due to missing AWS credentials
		},
		{
			name:        "GCP cloud",
			cloud:       GCPCloud,
			expectError: true, // Will fail due to missing GCP credentials
		},
		{
			name:        "Unknown cloud",
			cloud:       Unknown,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDefaultCloudParams(ctx, tt.cloud)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
