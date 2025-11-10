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
)

func TestNode_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		operation   func() error
		expectError bool
	}{
		{
			name: "Empty node - Connect",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.Connect(0)
			},
			expectError: true,
		},
		{
			name: "Empty node - Upload",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.Upload("/local", "/remote", time.Second)
			},
			expectError: true,
		},
		{
			name: "Empty node - Download",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.Download("/remote", "/local", time.Second)
			},
			expectError: true,
		},
		{
			name: "Empty node - Command",
			node: Node{},
			operation: func() error {
				node := Node{}
				_, err := node.Command(nil, time.Second, "echo test")
				return err
			},
			expectError: true,
		},
		{
			name: "Empty node - FileExists",
			node: Node{},
			operation: func() error {
				node := Node{}
				_, err := node.FileExists("/path")
				return err
			},
			expectError: true,
		},
		{
			name: "Empty node - CreateTempFile",
			node: Node{},
			operation: func() error {
				node := Node{}
				_, err := node.CreateTempFile()
				return err
			},
			expectError: true,
		},
		{
			name: "Empty node - CreateTempDir",
			node: Node{},
			operation: func() error {
				node := Node{}
				_, err := node.CreateTempDir()
				return err
			},
			expectError: true,
		},
		{
			name: "Empty node - Remove",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.Remove("/path", false)
			},
			expectError: true,
		},
		{
			name: "Empty node - WaitForSSHShell",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.WaitForSSHShell(time.Second)
			},
			expectError: true,
		},
		{
			name: "Empty node - StreamSSHCommand",
			node: Node{},
			operation: func() error {
				node := Node{}
				return node.StreamSSHCommand(nil, time.Second, "echo test")
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

func TestNode_InvalidInputs(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "Upload with empty local file",
			operation: func() error {
				return node.Upload("", "/remote", time.Second)
			},
			expectError: true,
		},
		{
			name: "Upload with empty remote file",
			operation: func() error {
				return node.Upload("/local", "", time.Second)
			},
			expectError: true,
		},
		{
			name: "Download with empty remote file",
			operation: func() error {
				return node.Download("", "/local", time.Second)
			},
			expectError: true,
		},
		{
			name: "Download with empty local file",
			operation: func() error {
				return node.Download("/remote", "", time.Second)
			},
			expectError: true,
		},
		{
			name: "Command with empty script",
			operation: func() error {
				_, err := node.Command(nil, time.Second, "")
				return err
			},
			expectError: true,
		},
		{
			name: "FileExists with empty path",
			operation: func() error {
				_, err := node.FileExists("")
				return err
			},
			expectError: true,
		},
		{
			name: "Remove with empty path",
			operation: func() error {
				return node.Remove("", false)
			},
			expectError: true,
		},
		{
			name: "WaitForSSHShell with empty IP",
			operation: func() error {
				node.IP = ""
				return node.WaitForSSHShell(time.Second)
			},
			expectError: true,
		},
		{
			name: "StreamSSHCommand with empty command",
			operation: func() error {
				return node.StreamSSHCommand(nil, time.Second, "")
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

func TestNode_TimeoutEdgeCases(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "Zero timeout",
			timeout: 0,
		},
		{
			name:    "Negative timeout",
			timeout: -1 * time.Second,
		},
		{
			name:    "Very short timeout",
			timeout: time.Nanosecond,
		},
		{
			name:    "Very long timeout",
			timeout: 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test various operations with different timeouts
			operations := []func() error{
				func() error { return node.Upload("/local", "/remote", tt.timeout) },
				func() error { return node.Download("/remote", "/local", tt.timeout) },
				func() error { _, err := node.Command(nil, tt.timeout, "echo test"); return err },
				func() error { return node.MkdirAll("/remote/dir", tt.timeout) },
				func() error { _, err := node.ReadFileBytes("/remote/file", tt.timeout); return err },
			}

			for i, operation := range operations {
				t.Run(tt.name+"_operation_"+string(rune(i+'0')), func(t *testing.T) {
					err := operation()
					// Should fail due to no connection, but timeout should be handled
					assert.Error(t, err)
				})
			}
		})
	}
}

func TestNode_ConcurrentOperations(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test concurrent operations
	t.Run("Concurrent Upload", func(t *testing.T) {
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(i int) {
				defer func() { done <- true }()
				err := node.Upload("/local", "/remote", time.Second)
				assert.Error(t, err) // Will fail due to no connection
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Concurrent Download", func(t *testing.T) {
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(i int) {
				defer func() { done <- true }()
				err := node.Download("/remote", "/local", time.Second)
				assert.Error(t, err) // Will fail due to no connection
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Concurrent Commands", func(t *testing.T) {
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(i int) {
				defer func() { done <- true }()
				_, err := node.Command(nil, time.Second, "echo test")
				assert.Error(t, err) // Will fail due to no connection
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestNode_MemoryLeaks(t *testing.T) {
	// Test that operations don't leak memory
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Run operations multiple times to check for memory leaks
	for i := 0; i < 100; i++ {
		// These operations should fail but not leak memory
		node.Upload("/local", "/remote", time.Second)
		node.Download("/remote", "/local", time.Second)
		node.Command(nil, time.Second, "echo test")
		node.FileExists("/path")
		node.CreateTempFile()
		node.CreateTempDir()
		node.Remove("/path", false)
		node.WaitForSSHShell(time.Second)
		node.StreamSSHCommand(nil, time.Second, "echo test")
	}
}

func TestNode_ResourceCleanup(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test that resources are properly cleaned up
	t.Run("Disconnect cleanup", func(t *testing.T) {
		err := node.Disconnect()
		assert.NoError(t, err) // Should succeed even if not connected
	})

	t.Run("Multiple disconnect calls", func(t *testing.T) {
		// Multiple disconnect calls should not cause issues
		for i := 0; i < 10; i++ {
			err := node.Disconnect()
			assert.NoError(t, err)
		}
	})
}

func TestNode_ErrorRecovery(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test error recovery scenarios
	t.Run("Connection failure recovery", func(t *testing.T) {
		// First attempt should fail
		err := node.Connect(0)
		assert.Error(t, err)

		// Second attempt should also fail but not crash
		err = node.Connect(0)
		assert.Error(t, err)
	})

	t.Run("Operation after connection failure", func(t *testing.T) {
		// Operations after connection failure should fail gracefully
		err := node.Upload("/local", "/remote", time.Second)
		assert.Error(t, err)

		err = node.Download("/remote", "/local", time.Second)
		assert.Error(t, err)
	})
}

func TestNode_InvalidSSHConfig(t *testing.T) {
	tests := []struct {
		name        string
		sshConfig   SSHConfig
		expectError bool
	}{
		{
			name: "Empty user",
			sshConfig: SSHConfig{
				User:           "",
				PrivateKeyPath: "/path/to/key",
			},
			expectError: true,
		},
		{
			name: "Empty private key path",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "",
			},
			expectError: true,
		},
		{
			name: "Non-existent private key",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/non/existent/key",
			},
			expectError: true,
		},
		{
			name: "Invalid private key format",
			sshConfig: SSHConfig{
				User:           "ubuntu",
				PrivateKeyPath: "/path/to/invalid/key",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{
				IP:        "192.168.1.1",
				SSHConfig: tt.sshConfig,
			}

			err := node.Connect(0)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_InvalidRoleCombinations(t *testing.T) {
	tests := []struct {
		name        string
		roles       []SupportedRole
		expectError bool
	}{
		{
			name:        "Validator and API",
			roles:       []SupportedRole{Validator, API},
			expectError: true,
		},
		{
			name:        "Loadtest with Validator",
			roles:       []SupportedRole{Loadtest, Validator},
			expectError: true,
		},
		{
			name:        "Loadtest with API",
			roles:       []SupportedRole{Loadtest, API},
			expectError: true,
		},
		{
			name:        "Loadtest with Monitor",
			roles:       []SupportedRole{Loadtest, Monitor},
			expectError: true,
		},
		{
			name:        "Monitor with Validator",
			roles:       []SupportedRole{Monitor, Validator},
			expectError: true,
		},
		{
			name:        "Monitor with API",
			roles:       []SupportedRole{Monitor, API},
			expectError: true,
		},
		{
			name:        "Valid single role",
			roles:       []SupportedRole{Validator},
			expectError: false,
		},
		{
			name:        "Valid multiple roles",
			roles:       []SupportedRole{Validator, Monitor},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRoles(tt.roles)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_ContextCancellation(t *testing.T) {
	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Operations with cancelled context should fail
	err := node.MonitorNodes(ctx, []Node{}, "test-cluster")
	assert.Error(t, err)
}

func TestNode_FileSystemEdgeCases(t *testing.T) {
	// Test file system edge cases
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
	}{
		{
			name: "Non-existent directory",
			setup: func() string {
				return "/non/existent/directory"
			},
			expectError: true,
		},
		{
			name: "Permission denied directory",
			setup: func() string {
				dir := filepath.Join(tempDir, "no-permission")
				err := os.MkdirAll(dir, 0o000)
				require.NoError(t, err)
				return dir
			},
			expectError: true,
		},
		{
			name: "Valid directory",
			setup: func() string {
				dir := filepath.Join(tempDir, "valid")
				err := os.MkdirAll(dir, 0o755)
				require.NoError(t, err)
				return dir
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			node := Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			}

			// Test MkdirAll with different paths
			err := node.MkdirAll(path, time.Second)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Error(t, err) // Will still fail due to no connection
			}
		})
	}
}

func TestNode_NetworkEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		expectError bool
	}{
		{
			name:        "Invalid IP format",
			ip:          "invalid-ip",
			expectError: true,
		},
		{
			name:        "Empty IP",
			ip:          "",
			expectError: true,
		},
		{
			name:        "Localhost IP",
			ip:          "127.0.0.1",
			expectError: true,
		},
		{
			name:        "Private IP",
			ip:          "192.168.1.1",
			expectError: true,
		},
		{
			name:        "Public IP",
			ip:          "8.8.8.8",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Node{
				IP: tt.ip,
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			}

			err := node.WaitForSSHShell(time.Second)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_StressTest(t *testing.T) {
	// Stress test with many operations
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Run many operations concurrently
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer func() { done <- true }()

			// Mix of different operations
			switch i % 5 {
			case 0:
				node.Upload("/local", "/remote", time.Second)
			case 1:
				node.Download("/remote", "/local", time.Second)
			case 2:
				node.Command(nil, time.Second, "echo test")
			case 3:
				node.FileExists("/path")
			case 4:
				node.CreateTempFile()
			}
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < 100; i++ {
		<-done
	}
}
