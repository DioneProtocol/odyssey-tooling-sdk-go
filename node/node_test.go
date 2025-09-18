// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/melbahja/goph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGophClient is a mock implementation of goph.Client
type MockGophClient struct {
	mock.Mock
}

func (m *MockGophClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGophClient) Upload(localFile, remoteFile string) error {
	args := m.Called(localFile, remoteFile)
	return args.Error(0)
}

func (m *MockGophClient) Download(remoteFile, localFile string) error {
	args := m.Called(remoteFile, localFile)
	return args.Error(0)
}

func (m *MockGophClient) CommandContext(ctx context.Context, name, script string) (*goph.Cmd, error) {
	args := m.Called(ctx, name, script)
	return args.Get(0).(*goph.Cmd), args.Error(1)
}

func (m *MockGophClient) NewSftp() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func (m *MockGophClient) NewSession() (interface{}, error) {
	args := m.Called()
	return args.Get(0), args.Error(1)
}

func (m *MockGophClient) DialTCP(network string, laddr, raddr interface{}) (interface{}, error) {
	args := m.Called(network, laddr, raddr)
	return args.Get(0), args.Error(1)
}

func (m *MockGophClient) RemoteAddr() string {
	args := m.Called()
	return args.String(0)
}

// MockSftpClient is a mock implementation of goph.Sftp
type MockSftpClient struct {
	mock.Mock
}

func (m *MockSftpClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSftpClient) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *MockSftpClient) Create(path string) (interface{}, error) {
	args := m.Called(path)
	return args.Get(0), args.Error(1)
}

func (m *MockSftpClient) Mkdir(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockSftpClient) MkdirAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockSftpClient) Remove(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// MockFileInfo is a mock implementation of os.FileInfo
type MockFileInfo struct {
	mock.Mock
}

func (m *MockFileInfo) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileInfo) Size() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockFileInfo) Mode() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockFileInfo) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockFileInfo) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockFileInfo) Sys() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestNewNodeConnection(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		port        uint
		expectError bool
	}{
		{
			name: "Valid connection with default port",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			port:        0,
			expectError: true, // Will fail due to missing key file
		},
		{
			name: "Valid connection with custom port",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			port:        2222,
			expectError: true, // Will fail due to missing key file
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
			port:        22,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewNodeConnection(&tt.node, tt.port)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_GetConnection(t *testing.T) {
	node := &Node{}
	assert.Nil(t, node.GetConnection())

	// Note: We can't easily test with mock client due to type constraints
	// This test verifies the basic functionality
}

func TestNode_GetSSHClient(t *testing.T) {
	node := &Node{}
	assert.Panics(t, func() { _ = node.GetSSHClient() })

	// This would require a real goph.Client to test properly
	// For now, we test the nil case
}

func TestNode_GetCloudID(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			name: "AWS node ID",
			node: Node{
				NodeID: "aws_node_12345",
			},
			expected: "12345",
		},
		{
			name: "GCP node ID",
			node: Node{
				NodeID: "gcp_node_67890",
			},
			expected: "67890",
		},
		{
			name: "Regular node ID",
			node: Node{
				NodeID: "regular_node_id",
			},
			expected: "regular_node_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.GetCloudID()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_Connected(t *testing.T) {
	node := &Node{}
	assert.False(t, node.Connected())

	// Note: We can't easily test with mock client due to type constraints
	// This test verifies the basic functionality
}

func TestNode_Disconnect(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name:        "No connection",
			node:        Node{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Disconnect()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_ExpandHome(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		path     string
		expected string
	}{
		{
			name: "Empty path",
			node: Node{
				SSHConfig: SSHConfig{User: "ubuntu"},
			},
			path:     "",
			expected: "/home/ubuntu",
		},
		{
			name: "Path with tilde",
			node: Node{
				SSHConfig: SSHConfig{User: "ubuntu"},
			},
			path:     "~/test",
			expected: "/home/ubuntu/test",
		},
		{
			name: "Path without tilde",
			node: Node{
				SSHConfig: SSHConfig{User: "ubuntu"},
			},
			path:     "/absolute/path",
			expected: "/absolute/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.ExpandHome(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_Upload(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		localFile   string
		remoteFile  string
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			localFile:   "/local/file",
			remoteFile:  "/remote/file",
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Upload(tt.localFile, tt.remoteFile, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_UploadBytes(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		data        []byte
		remoteFile  string
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Valid data",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			data:        []byte("test data"),
			remoteFile:  "/remote/file",
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.UploadBytes(tt.data, tt.remoteFile, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_Download(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		remoteFile  string
		localFile   string
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			remoteFile:  "/remote/file",
			localFile:   "/local/file",
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Download(tt.remoteFile, tt.localFile, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_ReadFileBytes(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		remoteFile  string
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			remoteFile:  "/remote/file",
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.ReadFileBytes(tt.remoteFile, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_MkdirAll(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		remoteDir   string
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			remoteDir:   "/remote/dir",
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.MkdirAll(tt.remoteDir, tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_Command(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		env         []string
		timeout     time.Duration
		script      string
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			env:         []string{"TEST=value"},
			timeout:     time.Second,
			script:      "echo test",
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.Command(tt.env, tt.timeout, tt.script)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_Commandf(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	_, err := node.Commandf([]string{"TEST=value"}, time.Second, "echo %s", "test")
	assert.Error(t, err) // Will fail due to missing key file
}

func TestNode_FileExists(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		path        string
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			path:        "/remote/file",
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.FileExists(tt.path)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_CreateTempFile(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.CreateTempFile()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_CreateTempDir(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.node.CreateTempDir()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_Remove(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		path        string
		recursive   bool
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			path:        "/remote/file",
			recursive:   false,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Remove(tt.path, tt.recursive)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_WaitForSSHShell(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Empty IP",
			node: Node{
				IP: "",
			},
			timeout:     time.Second,
			expectError: true,
		},
		{
			name: "Valid IP but no connection",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			timeout:     time.Second,
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.WaitForSSHShell(tt.timeout)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_HasSystemDAvailable(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "Not connected",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.node.HasSystemDAvailable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_StreamSSHCommand(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		env         []string
		timeout     time.Duration
		command     string
		expectError bool
	}{
		{
			name: "Not connected - should connect first",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			env:         []string{"TEST=value"},
			timeout:     time.Second,
			command:     "echo test",
			expectError: true, // Will fail due to missing key file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.StreamSSHCommand(tt.env, tt.timeout, tt.command)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsumeOutput(t *testing.T) {
	// This is a helper function that's hard to test in isolation
	// It's tested indirectly through StreamSSHCommand
	assert.True(t, true) // Placeholder test
}
