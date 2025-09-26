// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/stretchr/testify/assert"
)

// TestFeatureFlags_AWSIntegration tests AWS integration feature flag
func TestFeatureFlags_AWSIntegration(t *testing.T) {
	// Save original value
	originalValue := constants.AWSIntegrationEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalValue
	}()

	// Test with AWS integration disabled
	constants.AWSIntegrationEnabled = false

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			AWSConfig: &AWSConfig{
				AWSProfile:           "test-profile",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-sg",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
			},
			Region:       "us-east-1",
			ImageID:      "ami-12345678",
			InstanceType: "c5.2xlarge",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should fail with AWS integration disabled
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AWS integration functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.AWSIntegrationEnabled = true to enable")

	// Test with AWS integration enabled
	constants.AWSIntegrationEnabled = true
	// Note: This will still fail due to other missing dependencies, but should pass the AWS flag check
	_, err = CreateNodes(ctx, nodeParams)
	// The error should not be about AWS integration being disabled
	assert.NotContains(t, err.Error(), "AWS integration functionality is disabled")
}

// TestFeatureFlags_GCPIntegration tests GCP integration feature flag
func TestFeatureFlags_GCPIntegration(t *testing.T) {
	// Save original value
	originalValue := constants.GCPIntegrationEnabled
	defer func() {
		constants.GCPIntegrationEnabled = originalValue
	}()

	// Test with GCP integration disabled
	constants.GCPIntegrationEnabled = false

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			GCPConfig: &GCPConfig{
				GCPProject:     "test-project",
				GCPCredentials: "/tmp/test-creds.json",
				GCPZone:        "us-central1-a",
				GCPNetwork:     "default",
				GCPSSHKey:      "test-ssh-key",
				GCPVolumeSize:  100,
			},
			Region:       "us-central1",
			ImageID:      "projects/test-project/global/images/test-image",
			InstanceType: "e2-standard-8",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should fail with GCP integration disabled
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GCP integration functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.GCPIntegrationEnabled = true to enable")

	// Test with GCP integration enabled
	constants.GCPIntegrationEnabled = true
	// Note: This will still fail due to other missing dependencies, but should pass the GCP flag check
	_, err = CreateNodes(ctx, nodeParams)
	// The error should not be about GCP integration being disabled
	assert.NotContains(t, err.Error(), "GCP integration functionality is disabled")
}

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

// TestFeatureFlags_InstanceManagement tests instance management feature flag
func TestFeatureFlags_InstanceManagement(t *testing.T) {
	// Save original value
	originalValue := constants.InstanceManagementEnabled
	defer func() {
		constants.InstanceManagementEnabled = originalValue
	}()

	// Test with instance management disabled
	constants.InstanceManagementEnabled = false

	// Enable AWS to bypass AWS flag check
	constants.AWSIntegrationEnabled = true
	defer func() {
		constants.AWSIntegrationEnabled = false
	}()

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			AWSConfig: &AWSConfig{
				AWSProfile:           "test-profile",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-sg",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
			},
			Region:       "us-east-1",
			ImageID:      "ami-12345678",
			InstanceType: "c5.2xlarge",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should fail with instance management disabled
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instance management functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.InstanceManagementEnabled = true to enable")

	// Test with instance management enabled
	constants.InstanceManagementEnabled = true
	// Note: This will still fail due to other missing dependencies, but should pass the instance management flag check
	_, err = CreateNodes(ctx, nodeParams)
	// The error should not be about instance management being disabled
	assert.NotContains(t, err.Error(), "instance management functionality is disabled")
}

// TestFeatureFlags_SSHKeyManagement tests SSH key management feature flag
func TestFeatureFlags_SSHKeyManagement(t *testing.T) {
	// Save original value
	originalValue := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalValue
	}()

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	// Enable AWS and other required flags to bypass their checks
	constants.AWSIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.DockerSupportEnabled = true
	defer func() {
		constants.AWSIntegrationEnabled = false
		constants.InstanceManagementEnabled = false
		constants.DockerSupportEnabled = false
	}()

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			AWSConfig: &AWSConfig{
				AWSProfile:           "test-profile",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-sg",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
			},
			Region:       "us-east-1",
			ImageID:      "ami-12345678",
			InstanceType: "c5.2xlarge",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should fail with SSH key management disabled
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
	assert.Contains(t, err.Error(), "Set constants.SSHKeyManagementEnabled = true to enable")

	// Test with SSH key management enabled
	constants.SSHKeyManagementEnabled = true
	// Note: This will still fail due to other missing dependencies, but should pass the SSH flag check
	_, err = CreateNodes(ctx, nodeParams)
	// The error should not be about SSH key management being disabled
	assert.NotContains(t, err.Error(), "SSH key management functionality is disabled")
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

	// Test ComposeSSHSetupAWMRelayer (still blocked)
	err = node.ComposeSSHSetupAWMRelayer()
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

// TestFeatureFlags_MultipleFlags tests multiple flags being disabled
func TestFeatureFlags_MultipleFlags(t *testing.T) {
	// Save original values
	originalAWS := constants.AWSIntegrationEnabled
	originalDocker := constants.DockerSupportEnabled
	originalInstance := constants.InstanceManagementEnabled
	originalSSH := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWS
		constants.DockerSupportEnabled = originalDocker
		constants.InstanceManagementEnabled = originalInstance
		constants.SSHKeyManagementEnabled = originalSSH
	}()

	// Disable multiple flags
	constants.AWSIntegrationEnabled = false
	constants.DockerSupportEnabled = false
	constants.InstanceManagementEnabled = false
	constants.SSHKeyManagementEnabled = false

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			AWSConfig: &AWSConfig{
				AWSProfile:           "test-profile",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-sg",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
			},
			Region:       "us-east-1",
			ImageID:      "ami-12345678",
			InstanceType: "c5.2xlarge",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should fail with AWS integration disabled (first check)
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AWS integration functionality is disabled")

	// Enable AWS but disable others
	constants.AWSIntegrationEnabled = true
	_, err = CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "instance management functionality is disabled")

	// Enable instance management but disable others
	constants.InstanceManagementEnabled = true
	_, err = CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")

	// Enable Docker but disable SSH
	constants.DockerSupportEnabled = true
	_, err = CreateNodes(ctx, nodeParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

// TestFeatureFlags_AllEnabled tests that all flags work when enabled
func TestFeatureFlags_AllEnabled(t *testing.T) {
	// Save original values
	originalAWS := constants.AWSIntegrationEnabled
	originalDocker := constants.DockerSupportEnabled
	originalInstance := constants.InstanceManagementEnabled
	originalSSH := constants.SSHKeyManagementEnabled

	defer func() {
		constants.AWSIntegrationEnabled = originalAWS
		constants.DockerSupportEnabled = originalDocker
		constants.InstanceManagementEnabled = originalInstance
		constants.SSHKeyManagementEnabled = originalSSH
	}()

	// Enable all flags
	constants.AWSIntegrationEnabled = true
	constants.DockerSupportEnabled = true
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()
	nodeParams := &NodeParams{
		CloudParams: &CloudParams{
			AWSConfig: &AWSConfig{
				AWSProfile:           "test-profile",
				AWSKeyPair:           "test-keypair",
				AWSSecurityGroupID:   "sg-12345678",
				AWSSecurityGroupName: "test-sg",
				AWSVolumeSize:        100,
				AWSVolumeType:        "gp3",
			},
			Region:       "us-east-1",
			ImageID:      "ami-12345678",
			InstanceType: "c5.2xlarge",
		},
		Count:             1,
		Roles:             []SupportedRole{Validator},
		Network:           odyssey.TestnetNetwork(),
		SSHPrivateKeyPath: "/tmp/test-key",
		OdysseyGoVersion:  "v1.0.0",
		UseStaticIP:       false,
	}

	// Should not fail due to feature flags (will fail due to other missing dependencies)
	_, err := CreateNodes(ctx, nodeParams)
	assert.Error(t, err) // Will fail due to missing AWS credentials, etc.
	// But should not fail due to feature flags
	assert.NotContains(t, err.Error(), "functionality is disabled")
}
