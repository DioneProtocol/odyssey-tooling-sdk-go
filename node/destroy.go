// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"fmt"
)

// Destroy destroys a node.
// Cloud functionality has been removed from this SDK.
func (h *Node) Destroy(ctx context.Context) error {
	return fmt.Errorf("cloud functionality has been removed from this SDK. Please use local node management instead")
}
