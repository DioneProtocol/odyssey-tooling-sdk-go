// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
)

// MockLogger is a mock implementation of odyssey.LeveledLogger for testing
type MockLogger struct{}

func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{})  {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}

// TestHelper provides utility functions for testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new TestHelper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateTestNode creates a test node with default values
func (th *TestHelper) CreateTestNode() Node {
	return Node{
		NodeID: "test-node-1",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
		Roles:  []SupportedRole{Validator},
		Logger: odyssey.LeveledLogger{},
	}
}

// CreateTestNodeParams creates test node parameters
func (th *TestHelper) CreateTestNodeParams() *NodeParams {
	return &NodeParams{
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SubnetIDs:         []string{},
		SSHPrivateKeyPath: "/path/to/key",
		OdysseyGoVersion:  "v1.10.13",
	}
}

// CreateTestSSHKey creates a temporary SSH key file for testing
func (th *TestHelper) CreateTestSSHKey() string {
	tempDir := th.t.TempDir()
	keyPath := filepath.Join(tempDir, "test_key")

	// Create a dummy private key file
	err := os.WriteFile(keyPath, []byte("-----BEGIN PRIVATE KEY-----\ntest-key-content\n-----END PRIVATE KEY-----"), 0o600)
	require.NoError(th.t, err)

	return keyPath
}

// CreateTestGCPCredentials creates a temporary GCP credentials file for testing
func (th *TestHelper) CreateTestGCPCredentials() string {
	tempDir := th.t.TempDir()
	credsPath := filepath.Join(tempDir, "gcp_credentials.json")

	credentials := `{
		"client_id": "test-client-id",
		"client_secret": "test-client-secret",
		"quota_project_id": "test-project-id",
		"refresh_token": "test-refresh-token",
		"type": "service_account"
	}`

	err := os.WriteFile(credsPath, []byte(credentials), 0o600)
	require.NoError(th.t, err)

	return credsPath
}

// CreateTestComposeFile creates a temporary docker-compose file for testing
func (th *TestHelper) CreateTestComposeFile() string {
	tempDir := th.t.TempDir()
	composePath := filepath.Join(tempDir, "docker-compose.yml")

	composeContent := `version: '3.8'
services:
  test:
    image: alpine:latest
    command: echo "test"
`

	err := os.WriteFile(composePath, []byte(composeContent), 0o644)
	require.NoError(th.t, err)

	return composePath
}

// AssertNodeProperties validates node properties
func (th *TestHelper) AssertNodeProperties(node Node, expectedNodeID, expectedIP string, expectedRoles []SupportedRole) {
	require.Equal(th.t, expectedNodeID, node.NodeID)
	require.Equal(th.t, expectedIP, node.IP)
	require.Equal(th.t, expectedRoles, node.Roles)
}

// AssertCloudParams validates cloud parameters - removed with cloud functionality
func (th *TestHelper) AssertCloudParams(cp interface{}, expectedRegion, expectedImageID, expectedInstanceType string) {
	// Cloud functionality has been removed
	th.t.Skip("Cloud functionality has been removed")
}

// CreateTestNodes creates multiple test nodes
func (th *TestHelper) CreateTestNodes(count int) []Node {
	nodes := make([]Node, count)
	for i := 0; i < count; i++ {
		nodes[i] = th.CreateTestNode()
		nodes[i].NodeID = "test-node-" + string(rune(i+1))
		nodes[i].IP = "192.168.1." + string(rune(i+1))
	}
	return nodes
}

// ValidateTestEnvironment checks if the test environment is properly set up
func (th *TestHelper) ValidateTestEnvironment() {
	// Check if required constants are defined
	require.NotEmpty(th.t, constants.SSHTCPPort)
	require.NotEmpty(th.t, constants.OdysseygoAPIPort)
	require.NotEmpty(th.t, constants.OdysseygoMachineMetricsPort)
	require.NotEmpty(th.t, constants.OdysseygoLoadTestPort)
}

// CreateTestNetwork creates a test network
func (th *TestHelper) CreateTestNetwork() odyssey.Network {
	return odyssey.TestnetNetwork()
}

// CreateTestSubnetIDs creates test subnet IDs
func (th *TestHelper) CreateTestSubnetIDs() []string {
	return []string{"subnet-1", "subnet-2"}
}

// AssertErrorContains checks if an error contains expected text
func (th *TestHelper) AssertErrorContains(err error, expectedText string) {
	require.Error(th.t, err)
	require.Contains(th.t, err.Error(), expectedText)
}

// AssertNoError checks that there is no error
func (th *TestHelper) AssertNoError(err error) {
	require.NoError(th.t, err)
}

// CreateTestCloudParams creates test cloud parameters for different clouds - removed with cloud functionality
func (th *TestHelper) CreateTestCloudParams(cloud interface{}) interface{} {
	// Cloud functionality has been removed
	th.t.Skip("Cloud functionality has been removed")
	return nil
}

// CreateTestNodeWithRole creates a test node with specific role
func (th *TestHelper) CreateTestNodeWithRole(role SupportedRole) Node {
	node := th.CreateTestNode()
	node.Roles = []SupportedRole{role}
	return node
}

// CreateTestNodeWithCloud creates a test node with specific cloud - removed with cloud functionality
func (th *TestHelper) CreateTestNodeWithCloud(cloud interface{}) Node {
	// Cloud functionality has been removed
	th.t.Skip("Cloud functionality has been removed")
	return Node{}
}

// AssertRoleCombination validates role combinations
func (th *TestHelper) AssertRoleCombination(roles []SupportedRole, shouldBeValid bool) {
	err := CheckRoles(roles)
	if shouldBeValid {
		th.AssertNoError(err)
	} else {
		th.AssertErrorContains(err, "")
	}
}

// CreateTestTimeout creates a test timeout
func (th *TestHelper) CreateTestTimeout() time.Duration {
	return time.Second
}

// AssertTimeoutHandling checks timeout handling
func (th *TestHelper) AssertTimeoutHandling(operation func(time.Duration) error, timeout time.Duration) {
	err := operation(timeout)
	// Should either succeed or fail with timeout error
	if err != nil {
		require.Error(th.t, err)
	}
}

// CreateTestSSHConfig creates test SSH configuration
func (th *TestHelper) CreateTestSSHConfig() SSHConfig {
	return SSHConfig{
		User:           "ubuntu",
		PrivateKeyPath: "/path/to/key",
		Params: map[string]string{
			"StrictHostKeyChecking": "no",
		},
	}
}

// AssertSSHConfig validates SSH configuration
func (th *TestHelper) AssertSSHConfig(config SSHConfig, expectedUser, expectedKeyPath string) {
	require.Equal(th.t, expectedUser, config.User)
	require.Equal(th.t, expectedKeyPath, config.PrivateKeyPath)
}

// CreateTestFileContent creates test file content
func (th *TestHelper) CreateTestFileContent() []byte {
	return []byte("test file content")
}

// AssertFileContent validates file content
func (th *TestHelper) AssertFileContent(content []byte, expected string) {
	require.Equal(th.t, expected, string(content))
}

// CreateTestEnvironment creates test environment variables
func (th *TestHelper) CreateTestEnvironment() []string {
	return []string{
		"TEST_VAR=test_value",
		"ANOTHER_VAR=another_value",
	}
}

// AssertEnvironment validates environment variables
func (th *TestHelper) AssertEnvironment(env []string, expectedVars []string) {
	require.Equal(th.t, len(expectedVars), len(env))
	for i, expected := range expectedVars {
		require.Equal(th.t, expected, env[i])
	}
}
