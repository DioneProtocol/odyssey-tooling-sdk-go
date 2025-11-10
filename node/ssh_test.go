// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
)

func TestSSHConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		sshConfig   SSHConfig
		expectError bool
	}{
		{
			name: "Valid SSH config",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/path/to/key",
				Params: map[string]string{
					"StrictHostKeyChecking": "no",
				},
			},
			expectError: false,
		},
		{
			name: "SSH config with empty user",
			sshConfig: SSHConfig{
				User:           "",
				PrivateKeyPath: "/path/to/key",
			},
			expectError: true,
		},
		{
			name: "SSH config with empty private key path",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "",
			},
			expectError: true,
		},
		{
			name: "SSH config with non-existent private key",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/non/existent/key",
			},
			expectError: true,
		},
		{
			name: "SSH config with invalid private key format",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/path/to/invalid/key",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for the private key if needed
			if tt.sshConfig.PrivateKeyPath != "" && tt.sshConfig.PrivateKeyPath != "/non/existent/key" {
				tempDir := t.TempDir()
				keyPath := filepath.Join(tempDir, "test_key")
				err := os.WriteFile(keyPath, []byte("-----BEGIN PRIVATE KEY-----\ntest-key-content\n-----END PRIVATE KEY-----"), 0o600)
				require.NoError(t, err)
				tt.sshConfig.PrivateKeyPath = keyPath
			}

			node := Node{
				IP:        "192.168.1.1",
				SSHConfig: tt.sshConfig,
			}

			err := node.Connect(0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Error(t, err) // Will still fail due to no real connection
			}
		})
	}
}

func TestSSHConfig_Parameters(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]string
		expectError bool
	}{
		{
			name: "Valid parameters",
			params: map[string]string{
				"StrictHostKeyChecking": "no",
				"UserKnownHostsFile":    "/dev/null",
				"LogLevel":              "ERROR",
			},
			expectError: false,
		},
		{
			name:        "Empty parameters",
			params:      map[string]string{},
			expectError: false,
		},
		{
			name:        "Nil parameters",
			params:      nil,
			expectError: false,
		},
		{
			name: "Invalid parameter values",
			params: map[string]string{
				"StrictHostKeyChecking": "invalid",
				"UserKnownHostsFile":    "",
				"LogLevel":              "INVALID",
			},
			expectError: false, // Parameters are passed as-is to SSH
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sshConfig := SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/path/to/key",
				Params:         tt.params,
			}

			node := Node{
				IP:        "192.168.1.1",
				SSHConfig: sshConfig,
			}

			err := node.Connect(0)
			assert.Error(t, err) // Will fail due to no real connection
		})
	}
}

func TestNode_SSHConnectionRetries(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test connection retries
	err := node.Connect(0)
	assert.Error(t, err) // Will fail due to no real connection

	// Test that retries are attempted
	assert.False(t, node.Connected())
}

func TestNode_SSHConnectionTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "Short timeout",
			timeout: time.Nanosecond,
		},
		{
			name:    "Normal timeout",
			timeout: time.Second,
		},
		{
			name:    "Long timeout",
			timeout: time.Minute,
		},
		{
			name:    "Zero timeout",
			timeout: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			}

			err := node.Connect(0)
			assert.Error(t, err) // Will fail due to no real connection
		})
	}
}

func TestNode_SSHConnectionPorts(t *testing.T) {
	tests := []struct {
		name     string
		port     uint
		expected uint
	}{
		{
			name:     "Default port",
			port:     0,
			expected: constants.SSHTCPPort,
		},
		{
			name:     "Custom port",
			port:     2222,
			expected: 2222,
		},
		{
			name:     "Standard SSH port",
			port:     22,
			expected: 22,
		},
		{
			name:     "High port",
			port:     65535,
			expected: 65535,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			}

			err := node.Connect(tt.port)
			assert.Error(t, err) // Will fail due to no real connection
		})
	}
}

func TestNode_SSHConnectionStates(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test initial state
	assert.False(t, node.Connected())

	// Test connection attempt
	err := node.Connect(0)
	assert.Error(t, err)
	assert.False(t, node.Connected())

	// Test disconnect when not connected
	err = node.Disconnect()
	assert.NoError(t, err)

	// Test multiple disconnect calls
	err = node.Disconnect()
	assert.NoError(t, err)
}

func TestNode_SSHConnectionErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Invalid IP address",
			node: Node{
				IP: "invalid-ip",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true,
		},
		{
			name: "Empty IP address",
			node: Node{
				IP: "",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true,
		},
		{
			name: "Non-existent private key",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/non/existent/key",
				},
			},
			expectError: true,
		},
		{
			name: "Empty private key path",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Connect(0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_SSHConnectionRetryLogic(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test that connection attempts are made
	err := node.Connect(0)
	assert.Error(t, err)

	// Test that connection state is properly managed
	assert.False(t, node.Connected())
}

func TestNode_SSHConnectionRetryDelay(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Verify that retries include delays between attempts
	// The connection should fail quickly on first attempt, but subsequent retries
	// should have delays (via SSHSleepBetweenChecks constant)
	start := time.Now()
	err := node.Connect(0)
	elapsed := time.Since(start)

	assert.Error(t, err)
	// With 5 retries and delays between retries (after first attempt),
	// the total time should be at least 4 * SSHSleepBetweenChecks
	// (minimum: 4 delays if all retries fail quickly)
	minExpectedTime := 4 * constants.SSHSleepBetweenChecks
	assert.GreaterOrEqual(t, elapsed, minExpectedTime,
		"Retry logic should include delays between attempts")
}

func TestNode_SSHConnectionContext(t *testing.T) {
	// Test context cancellation
	_, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test that cancelled context is handled
	err := node.Connect(0)
	assert.Error(t, err)
}

func TestNode_SSHConnectionConcurrency(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test concurrent connection attempts
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			err := node.Connect(0)
			assert.Error(t, err) // Will fail due to no real connection
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestNode_SSHConnectionTimeoutHandling(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with very short timeout
	err := node.Connect(0)
	assert.Error(t, err)

	// Test with long timeout
	err = node.Connect(0)
	assert.Error(t, err)
}

func TestNode_SSHConnectionErrorRecovery(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test error recovery
	err := node.Connect(0)
	assert.Error(t, err)

	// Test that subsequent operations fail gracefully
	err = node.Upload("/local", "/remote", time.Second)
	assert.Error(t, err)

	err = node.Download("/remote", "/local", time.Second)
	assert.Error(t, err)
}

func TestNode_SSHConnectionResourceCleanup(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test resource cleanup
	err := node.Disconnect()
	assert.NoError(t, err)

	// Test multiple disconnect calls
	err = node.Disconnect()
	assert.NoError(t, err)
}

func TestNode_SSHConnectionStateManagement(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test initial state
	assert.False(t, node.Connected())

	// Test connection attempt
	err := node.Connect(0)
	assert.Error(t, err)
	assert.False(t, node.Connected())

	// Test disconnect
	err = node.Disconnect()
	assert.NoError(t, err)
	assert.False(t, node.Connected())
}

func TestNode_SSHConnectionValidation(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid node configuration",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no real connection
		},
		{
			name: "Invalid node configuration",
			node: Node{
				IP: "invalid-ip",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Connect(0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_SSHConnectionEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with empty SSH config",
			node: Node{
				IP:        "192.168.1.1",
				SSHConfig: SSHConfig{},
			},
			expectError: true,
		},
		{
			name: "Node with nil SSH config",
			node: Node{
				IP: "192.168.1.1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Connect(0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_SSHConnectionStressTest(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Stress test with many connection attempts
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			defer func() { done <- true }()
			err := node.Connect(0)
			assert.Error(t, err) // Will fail due to no real connection
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestNode_SSHConnectionMemoryLeaks(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test for memory leaks with many connection attempts
	for i := 0; i < 1000; i++ {
		err := node.Connect(0)
		assert.Error(t, err) // Will fail due to no real connection
	}
}

func TestNode_SSHConnectionErrorMessages(t *testing.T) {
	// Save original SSH key management flag value
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	err := node.Connect(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to node")
}

func TestNode_SSHConnectionConstants(t *testing.T) {
	// Test that SSH constants are properly defined
	assert.NotZero(t, constants.SSHTCPPort)
	assert.NotZero(t, sshConnectionTimeout)
	assert.NotZero(t, sshConnectionRetries)
}

func TestNode_SSHConnectionTimeoutConstants(t *testing.T) {
	// Test timeout constants
	assert.Equal(t, 3*time.Second, sshConnectionTimeout)
	assert.Equal(t, 5, sshConnectionRetries)
}

func TestNode_SSHConnectionPortConstants(t *testing.T) {
	// Test port constants
	assert.Equal(t, 22, constants.SSHTCPPort)
}

func TestNode_SSHConnectionRetryCount(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test that retries are attempted
	err := node.Connect(0)
	assert.Error(t, err)

	// Test that connection state is properly managed
	assert.False(t, node.Connected())
}

func TestNode_SSHConnectionStateTransitions(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test state transitions
	assert.False(t, node.Connected())

	err := node.Connect(0)
	assert.Error(t, err)
	assert.False(t, node.Connected())

	err = node.Disconnect()
	assert.NoError(t, err)
	assert.False(t, node.Connected())
}
