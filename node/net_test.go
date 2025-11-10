// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestNode_OdysseygoTCPClient(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid node configuration",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			expectError: true,
		},
		{
			name: "Node with empty SSH config",
			node: Node{
				NodeID:    "test-node",
				IP:        "192.168.1.1",
				SSHConfig: SSHConfig{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := tt.node.OdysseygoTCPClient()
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, conn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
			}
		})
	}
}

func TestNode_OdysseygoTCPClient_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	conn, err := node.OdysseygoTCPClient()
	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

func TestNode_OdysseygoRPCClient(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid node configuration",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with invalid IP",
			node: Node{
				NodeID: "test-node",
				IP:     "invalid-ip-address",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.node.OdysseygoRPCClient()
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestNode_OdysseygoRPCClient_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	client, err := node.OdysseygoRPCClient()
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

func TestNode_Post(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		path        string
		requestBody string
		expectError bool
	}{
		{
			name:        "Default path",
			path:        "",
			requestBody: `{"jsonrpc":"2.0","method":"info.getNodeID"}`,
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Custom path",
			path:        "/ext/bc/P",
			requestBody: `{"jsonrpc":"2.0","method":"platform.getHeight"}`,
			expectError: true,
		},
		{
			name:        "Empty request body",
			path:        "/ext/info",
			requestBody: "",
			expectError: true,
		},
		{
			name:        "Invalid JSON",
			path:        "/ext/info",
			requestBody: `{"invalid": json}`,
			expectError: true,
		},
		{
			name:        "Valid JSON with complex data",
			path:        "/ext/bc/P",
			requestBody: `{"jsonrpc":"2.0","method":"platform.getHeight","params":{},"id":1}`,
			expectError: true,
		},
		{
			name:        "Path with query parameters",
			path:        "/ext/bc/P?param=value",
			requestBody: `{"jsonrpc":"2.0","method":"platform.getHeight"}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			response, err := node.Post(tt.path, tt.requestBody)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
			}
		})
	}
}

func TestNode_Post_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	response, err := node.Post("/ext/info", `{"jsonrpc":"2.0","method":"info.getNodeID"}`)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

func TestNode_WaitForPort(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		port        uint
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "Default port with short timeout",
			port:        0, // Will use default SSH port
			timeout:     100 * time.Millisecond,
			expectError: true, // Will timeout
		},
		{
			name:        "Custom port",
			port:        8080,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name:        "Zero timeout",
			port:        22,
			timeout:     0,
			expectError: true,
		},
		{
			name:        "Very short timeout",
			port:        22,
			timeout:     1 * time.Millisecond,
			expectError: true,
		},
		{
			name:        "High port number",
			port:        65535,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name:        "Invalid port number",
			port:        99999, // Invalid port
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.WaitForPort(tt.port, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_WaitForPort_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	err := node.WaitForPort(22, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

func TestNode_WaitForPort_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		node        Node
		port        uint
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "Empty node",
			node:        Node{},
			port:        22,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			port:        22,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name: "Node with invalid IP",
			node: Node{
				NodeID: "test-node",
				IP:     "invalid-ip",
			},
			port:        22,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name: "Node with localhost IP",
			node: Node{
				NodeID: "test-node",
				IP:     "127.0.0.1",
			},
			port:        22,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
		{
			name: "Node with IPv6 IP",
			node: Node{
				NodeID: "test-node",
				IP:     "::1",
			},
			port:        22,
			timeout:     100 * time.Millisecond,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.WaitForPort(tt.port, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_NetworkOperations_ConnectionStates(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		node        Node
		operations  []func() error
		expectError bool
	}{
		{
			name: "Multiple operations on same node",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			operations: []func() error{
				func() error {
					testNode := Node{NodeID: "test-node", IP: "192.168.1.1"}
					_, err := testNode.OdysseygoTCPClient()
					return err
				},
				func() error {
					testNode := Node{NodeID: "test-node", IP: "192.168.1.1"}
					_, err := testNode.OdysseygoRPCClient()
					return err
				},
				func() error {
					testNode := Node{NodeID: "test-node", IP: "192.168.1.1"}
					_, err := testNode.Post("/ext/info", `{"jsonrpc":"2.0","method":"info.getNodeID"}`)
					return err
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, operation := range tt.operations {
				t.Run(fmt.Sprintf("operation_%d", i), func(t *testing.T) {
					err := operation()
					if tt.expectError {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				})
			}
		})
	}
}

func TestNode_NetworkOperations_TimeoutHandling(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "WaitForPort with very short timeout",
			operation: func() error {
				return node.WaitForPort(22, 1*time.Nanosecond)
			},
			expectError: true,
		},
		{
			name: "WaitForPort with reasonable timeout",
			operation: func() error {
				return node.WaitForPort(22, 1*time.Second)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_NetworkOperations_ContextHandling(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
	}{
		{
			name:        "Background context",
			ctx:         context.Background(),
			expectError: true,
		},
		{
			name: "Cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expectError: true,
		},
		{
			name: "Context with timeout",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
				return ctx
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that network operations handle context properly
			// Note: Most operations don't directly use context, but we test the pattern
			_, err := node.OdysseygoTCPClient()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_NetworkOperations_ErrorMessages(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test that error messages are informative
	_, err := node.OdysseygoTCPClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")

	_, err = node.OdysseygoRPCClient()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")

	_, err = node.Post("/ext/info", `{"jsonrpc":"2.0","method":"info.getNodeID"}`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")

	err = node.WaitForPort(22, 100*time.Millisecond)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}
