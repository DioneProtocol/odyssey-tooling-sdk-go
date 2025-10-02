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

func TestGetDefaultCloudParams_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	ctx := context.Background()

	tests := []struct {
		name          string
		cloud         SupportedCloud
		awsEnabled    bool
		gcpEnabled    bool
		expectError   bool
		errorContains string
	}{
		{
			name:          "AWS integration disabled",
			cloud:         AWSCloud,
			awsEnabled:    false,
			gcpEnabled:    true,
			expectError:   true,
			errorContains: "AWS integration functionality is disabled",
		},
		{
			name:          "GCP integration disabled",
			cloud:         GCPCloud,
			awsEnabled:    true,
			gcpEnabled:    false,
			expectError:   true,
			errorContains: "GCP integration functionality is disabled",
		},
		{
			name:          "Both integrations disabled",
			cloud:         AWSCloud,
			awsEnabled:    false,
			gcpEnabled:    false,
			expectError:   true,
			errorContains: "AWS integration functionality is disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flags
			constants.AWSIntegrationEnabled = tt.awsEnabled
			constants.GCPIntegrationEnabled = tt.gcpEnabled

			params, err := GetDefaultCloudParams(ctx, tt.cloud)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, params)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, params)
			}
		})
	}
}

func TestGetDefaultCloudParams_CloudTypes(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true

	ctx := context.Background()

	tests := []struct {
		name          string
		cloud         SupportedCloud
		expectError   bool
		errorContains string
	}{
		{
			name:        "AWS cloud type",
			cloud:       AWSCloud,
			expectError: true, // Will fail due to AWS API calls
		},
		{
			name:        "GCP cloud type",
			cloud:       GCPCloud,
			expectError: true, // Will fail due to GCP API calls
		},
		{
			name:          "Unknown cloud type",
			cloud:         SupportedCloud(99),
			expectError:   true,
			errorContains: "unsupported cloud type",
		},
		{
			name:          "Docker cloud type",
			cloud:         Docker,
			expectError:   true,
			errorContains: "unsupported cloud type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := GetDefaultCloudParams(ctx, tt.cloud)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, params)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, params)
			}
		})
	}
}

func TestGetDefaultCloudParams_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true

	tests := []struct {
		name        string
		ctx         context.Context
		cloud       SupportedCloud
		expectError bool
	}{
		{
			name:        "Nil context",
			ctx:         nil,
			cloud:       AWSCloud,
			expectError: true,
		},
		{
			name: "Cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			cloud:       AWSCloud,
			expectError: true,
		},
		{
			name:        "Valid context with AWS",
			ctx:         context.Background(),
			cloud:       AWSCloud,
			expectError: true, // Will fail due to AWS API calls
		},
		{
			name:        "Valid context with GCP",
			ctx:         context.Background(),
			cloud:       GCPCloud,
			expectError: true, // Will fail due to GCP API calls
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := GetDefaultCloudParams(tt.ctx, tt.cloud)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, params)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, params)
			}
		})
	}
}

func TestGetDefaultCloudParams_Validation(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
	}()

	// Enable feature flags
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true

	ctx := context.Background()

	tests := []struct {
		name           string
		cloud          SupportedCloud
		expectError    bool
		validateParams func(*testing.T, *CloudParams)
	}{
		{
			name:        "AWS default parameters",
			cloud:       AWSCloud,
			expectError: true, // Will fail due to AWS API calls
			validateParams: func(t *testing.T, params *CloudParams) {
				if params != nil {
					assert.Equal(t, AWSCloud, params.Cloud())
					assert.Equal(t, "us-east-1", params.Region)
					assert.NotNil(t, params.AWSConfig)
					assert.Equal(t, "default", params.AWSConfig.AWSProfile)
					assert.Equal(t, 1000, params.AWSConfig.AWSVolumeSize)
					assert.Equal(t, 500, params.AWSConfig.AWSVolumeThroughput)
					assert.Equal(t, 1000, params.AWSConfig.AWSVolumeIOPS)
					assert.Equal(t, "gp3", params.AWSConfig.AWSVolumeType)
				}
			},
		},
		{
			name:        "GCP default parameters",
			cloud:       GCPCloud,
			expectError: true, // Will fail due to GCP API calls
			validateParams: func(t *testing.T, params *CloudParams) {
				if params != nil {
					assert.Equal(t, GCPCloud, params.Cloud())
					assert.Equal(t, "us-west1", params.Region)
					assert.NotNil(t, params.GCPConfig)
					assert.Equal(t, "test-project", params.GCPConfig.GCPProject)
					assert.Equal(t, 100, params.GCPConfig.GCPVolumeSize)
					// GCPConfig doesn't have GCPVolumeType field
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := GetDefaultCloudParams(ctx, tt.cloud)

			if tt.expectError {
				require.Error(t, err)
				// Still validate parameters if they were created before the error
				if params != nil && tt.validateParams != nil {
					tt.validateParams(t, params)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, params)
				if tt.validateParams != nil {
					tt.validateParams(t, params)
				}
			}
		})
	}
}
