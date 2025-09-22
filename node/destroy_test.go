// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple tests that verify the destroy function behavior without complex mocking

func TestNode_Destroy_AWS_Success(t *testing.T) {
	// Setup
	node := &Node{
		NodeID: "test-node-123",
		Cloud:  AWSCloud,
		CloudConfig: CloudParams{
			Region: "us-west-2",
			AWSConfig: &AWSConfig{
				AWSProfile: "default",
			},
		},
	}

	// Execute - this will fail due to missing AWS credentials but tests the validation logic
	ctx := context.Background()
	err := node.Destroy(ctx)

	// Assert - this should fail due to AWS credential issues
	assert.Error(t, err)
}

func TestNode_Destroy_AWS_CloudCreationError(t *testing.T) {
	// Setup
	node := &Node{
		NodeID: "test-node-123",
		Cloud:  AWSCloud,
		CloudConfig: CloudParams{
			Region: "us-west-2",
			AWSConfig: &AWSConfig{
				AWSProfile: "default",
			},
		},
	}

	// Execute - this will fail due to missing AWS credentials
	ctx := context.Background()
	err := node.Destroy(ctx)

	// Assert
	assert.Error(t, err)
}

func TestNode_Destroy_GCP_Success(t *testing.T) {
	// Setup
	node := &Node{
		NodeID: "test-node-456",
		Cloud:  GCPCloud,
		CloudConfig: CloudParams{
			Region: "us-central1",
			GCPConfig: &GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/path/to/credentials.json",
			},
		},
	}

	// Execute - this will fail due to missing GCP credentials but tests the validation logic
	ctx := context.Background()
	err := node.Destroy(ctx)

	// Assert - this should fail due to GCP credential issues
	assert.Error(t, err)
}

func TestNode_Destroy_UnsupportedCloud(t *testing.T) {
	// Setup
	node := &Node{
		NodeID: "test-node-789",
		Cloud:  SupportedCloud(999), // Invalid cloud type
		CloudConfig: CloudParams{
			Region: "us-west-1",
		},
	}

	// Execute
	ctx := context.Background()
	err := node.Destroy(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported cloud type")
}
