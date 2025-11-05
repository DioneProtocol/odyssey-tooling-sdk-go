// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple tests that verify the destroy function behavior without complex mocking

func TestNode_Destroy_CloudFunctionalityRemoved(t *testing.T) {
	// Setup
	node := &Node{
		NodeID: "test-node-123",
		IP:     "192.168.1.1",
	}

	// Execute
	ctx := context.Background()
	err := node.Destroy(ctx)

	// Assert - should return error indicating cloud functionality removed
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cloud functionality has been removed")
}
