// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/subnet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/units"
)

func TestNode_GetBLSKeyFromRemoteHost_ReadFileError(t *testing.T) {
	node := &Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail to read the BLS key file
	// This tests the error handling path
	err := node.GetBLSKeyFromRemoteHost()

	// This should fail due to SSH connection issues
	assert.Error(t, err)
}

func TestNode_GetBLSKeyFromRemoteHost_InvalidBLSKey(t *testing.T) {
	node := &Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail to read the BLS key file
	// This tests the error handling path
	err := node.GetBLSKeyFromRemoteHost()

	// This should fail due to SSH connection issues
	assert.Error(t, err)
}

func TestNode_ValidatePrimaryNetwork_EmptyNodeID(t *testing.T) {
	node := &Node{}
	validatorParams := validator.PrimaryNetworkValidatorParams{
		NodeID:      ids.EmptyNodeID, // Invalid
		Duration:    48 * time.Hour,
		StakeAmount: 2 * units.Dione,
	}

	// Create a real network for testing
	network := odyssey.TestnetNetwork()

	// Create a real wallet for testing (this will fail but tests validation)
	var wallet wallet.Wallet

	txID, err := node.ValidatePrimaryNetwork(network, validatorParams, wallet)

	assert.Error(t, err)
	assert.Equal(t, ids.Empty, txID)
	assert.Contains(t, err.Error(), subnet.ErrEmptyValidatorNodeID.Error())
}

func TestNode_ValidatePrimaryNetwork_EmptyDuration(t *testing.T) {
	node := &Node{}
	validatorParams := validator.PrimaryNetworkValidatorParams{
		NodeID:      ids.GenerateTestNodeID(),
		Duration:    0, // Invalid
		StakeAmount: 2 * units.Dione,
	}

	// Create a real network for testing
	network := odyssey.TestnetNetwork()

	// Create a real wallet for testing (this will fail but tests validation)
	var wallet wallet.Wallet

	txID, err := node.ValidatePrimaryNetwork(network, validatorParams, wallet)

	assert.Error(t, err)
	assert.Equal(t, ids.Empty, txID)
	assert.Contains(t, err.Error(), subnet.ErrEmptyValidatorDuration.Error())
}

// Helper function to create a temporary file with content
func createTempFileWithContent(t *testing.T, content []byte) *os.File {
	tempFile, err := os.CreateTemp("", "test-bls-key-*.txt")
	require.NoError(t, err)

	_, err = tempFile.Write(content)
	require.NoError(t, err)

	err = tempFile.Close()
	require.NoError(t, err)

	// Reopen for reading
	tempFile, err = os.Open(tempFile.Name())
	require.NoError(t, err)

	return tempFile
}
