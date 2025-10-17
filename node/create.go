// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"fmt"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
)

// NodeParams is an input for CreateNodes
type NodeParams struct {
	// Count is how many Odyssey Nodes to be created during CreateNodes
	Count int

	// Roles pertain to whether the created node is going to be a Validator / API / Monitoring
	// node. See CheckRoles to see which combination of roles for a node is supported.
	Roles []SupportedRole

	// Network is whether the Validator / API node is meant to track OdysseyGo Primary Network
	// in Testnet / Mainnet / Devnet
	Network odyssey.Network

	// SubnetIDs is the list of subnet IDs that the created nodes will be tracking
	// For primary network, it should be empty
	SubnetIDs []string

	// SSHPrivateKeyPath is the file path to the private key of the SSH key pair that is used
	// to gain access to the created nodes
	SSHPrivateKeyPath string

	// OdysseyGoVersion is the version of Odyssey Go to install in the created node
	OdysseyGoVersion string
}

// CreateNodes is a placeholder function for node creation.
// Cloud functionality has been removed from this SDK.
// This function now returns an error indicating that cloud functionality is not available.
func CreateNodes(
	ctx context.Context,
	nodeParams *NodeParams,
) ([]Node, error) {
	return nil, fmt.Errorf("cloud functionality has been removed from this SDK. Please use local node setup instead")
}

// provisionHost provisions a host with the given roles.
func provisionHost(node Node, nodeParams *NodeParams) error {
	if err := CheckRoles(nodeParams.Roles); err != nil {
		return err
	}
	if err := node.Connect(constants.SSHTCPPort); err != nil {
		return err
	}
	for _, role := range nodeParams.Roles {
		switch role {
		case Validator:
			if err := provisionOdysseyGoHost(node, nodeParams); err != nil {
				return err
			}
		case API:
			if err := provisionOdysseyGoHost(node, nodeParams); err != nil {
				return err
			}
		case Loadtest:
			if err := provisionLoadTestHost(node); err != nil {
				return err
			}
		case Monitor:
			if err := provisionMonitoringHost(node); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported role %v", role)
		}
	}
	return nil
}

func provisionOdysseyGoHost(node Node, nodeParams *NodeParams) error {
	const withMonitoring = true
	if err := node.RunSSHSetupNode(); err != nil {
		return err
	}
	if err := node.RunSSHSetupDockerService(); err != nil {
		return err
	}
	// provide dummy config for promtail
	if err := node.RunSSHSetupPromtailConfig("127.0.0.1", constants.OdysseygoLokiPort, node.NodeID, ""); err != nil {
		return err
	}
	if err := node.ComposeSSHSetupNode(nodeParams.Network.HRP(), nodeParams.SubnetIDs, nodeParams.OdysseyGoVersion, withMonitoring); err != nil {
		return err
	}
	if err := node.StartDockerCompose(constants.SSHScriptTimeout); err != nil {
		return err
	}
	return nil
}

func provisionLoadTestHost(node Node) error { // stub
	if err := node.ComposeSSHSetupLoadTest(); err != nil {
		return err
	}
	if err := node.RestartDockerCompose(constants.SSHScriptTimeout); err != nil {
		return err
	}
	return nil
}

func provisionMonitoringHost(node Node) error {
	if err := node.RunSSHSetupDockerService(); err != nil {
		return err
	}
	if err := node.RunSSHSetupMonitoringFolders(); err != nil {
		return err
	}
	if err := node.ComposeSSHSetupMonitoring(); err != nil {
		return err
	}
	if err := node.RestartDockerCompose(constants.SSHScriptTimeout); err != nil {
		return err
	}
	return nil
}
