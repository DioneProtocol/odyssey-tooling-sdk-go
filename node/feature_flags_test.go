// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

// TestFeatureFlags_DockerSupport tests Docker support feature flag
func TestFeatureFlags_DockerSupport(t *testing.T) {
	// Save original value
	originalValue := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalValue
	}()

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	// Create a test node to test Docker functions directly
	testNode := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/tmp/test-key",
		},
	}

	// Test RunSSHSetupDockerService
	err := testNode.RunSSHSetupDockerService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.DockerSupportEnabled = true to enable")

	// Test with Docker support enabled
	constants.DockerSupportEnabled = true
	// This will still fail due to SSH connection, but should pass the Docker flag check
	err = testNode.RunSSHSetupDockerService()
	// The error should not be about Docker support being disabled
	if err != nil {
		assert.NotContains(t, err.Error(), "Docker support functionality is disabled")
	}
}

// TestFeatureFlags_DockerFunctions tests individual Docker function gating
func TestFeatureFlags_DockerFunctions(t *testing.T) {
	// Save original value
	originalValue := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalValue
	}()

	// Create a test node
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/tmp/test-key",
		},
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	// Test RunSSHSetupDockerService (still blocked)
	err := node.RunSSHSetupDockerService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")

	// Test ComposeSSHSetupLoadTest (still blocked)
	err = node.ComposeSSHSetupLoadTest()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")

	// Test ComposeSSHSetupMonitoring (still blocked)
	err = node.ComposeSSHSetupMonitoring()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")

	// Note: ComposeSSHSetupNode, StartDockerCompose, StopDockerCompose, and RestartDockerCompose
	// are no longer blocked by DockerSupportEnabled flag as they are used for core functionality
	// like monitoring and subnet management that should work independently of cloud infrastructure.
}

// TestFeatureFlags_SSHConnect tests SSH connection gating
func TestFeatureFlags_SSHConnect(t *testing.T) {
	// Save original value
	originalValue := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalValue
	}()

	// Create a test node
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/tmp/test-key",
		},
	}

	// Test SSH connection - should fail due to actual SSH issues, not feature flags
	err := node.Connect(22)
	assert.Error(t, err)
	// Should fail due to SSH connection issues, not feature flag issues
	assert.Contains(t, err.Error(), "failed to connect to node")

	// Note: SSH connection is no longer gated by SSHKeyManagementEnabled flag
	// as it's a core functionality that should work independently of cloud infrastructure.
}
