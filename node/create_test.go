// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNodes_CloudFunctionalityRemoved(t *testing.T) {
	// This test verifies that CreateNodes returns appropriate error
	// when cloud functionality has been removed
	ctx := context.Background()
	nodes, err := CreateNodes(ctx, &NodeParams{
		Count: 1,
		Roles: []SupportedRole{Validator},
	})

	// Should return error and no nodes
	assert.Error(t, err)
	assert.Nil(t, nodes)
	assert.Contains(t, err.Error(), "cloud functionality has been removed")
}
