// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package odyssey

import (
	"os"
	"testing"

	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/stretchr/testify/assert"
)

func TestNetworkKind_String(t *testing.T) {
	tests := []struct {
		name     string
		kind     NetworkKind
		expected string
	}{
		{
			name:     "Mainnet",
			kind:     Mainnet,
			expected: "mainnet",
		},
		{
			name:     "Testnet",
			kind:     Testnet,
			expected: "testnet",
		},
		{
			name:     "Devnet",
			kind:     Devnet,
			expected: "local",
		},
		{
			name:     "Undefined",
			kind:     Undefined,
			expected: "invalid network",
		},
		{
			name:     "Invalid value",
			kind:     NetworkKind(999),
			expected: "invalid network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.kind.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNetwork_HRP(t *testing.T) {
	tests := []struct {
		name     string
		network  Network
		expected string
	}{
		{
			name: "Testnet network",
			network: Network{
				ID: constants.TestnetID,
			},
			expected: constants.TestnetHRP,
		},
		{
			name: "Mainnet network",
			network: Network{
				ID: constants.MainnetID,
			},
			expected: constants.MainnetHRP,
		},
		{
			name: "Custom network",
			network: Network{
				ID: 12345,
			},
			expected: constants.FallbackHRP,
		},
		{
			name: "Zero ID network",
			network: Network{
				ID: 0,
			},
			expected: constants.FallbackHRP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.network.HRP()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNetworkFromNetworkID(t *testing.T) {
	tests := []struct {
		name      string
		networkID uint32
		expected  Network
	}{
		{
			name:      "Mainnet ID",
			networkID: constants.MainnetID,
			expected:  MainnetNetwork(),
		},
		{
			name:      "Testnet ID",
			networkID: constants.TestnetID,
			expected:  TestnetNetwork(),
		},
		{
			name:      "Unknown ID",
			networkID: 99999,
			expected:  UndefinedNetwork,
		},
		{
			name:      "Zero ID",
			networkID: 0,
			expected:  UndefinedNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NetworkFromNetworkID(tt.networkID)
			assert.Equal(t, tt.expected.Kind, result.Kind)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Endpoint, result.Endpoint)
		})
	}
}

func TestNewNetwork(t *testing.T) {
	tests := []struct {
		name     string
		kind     NetworkKind
		id       uint32
		endpoint string
		expected Network
	}{
		{
			name:     "Mainnet network",
			kind:     Mainnet,
			id:       constants.MainnetID,
			endpoint: "https://node.dioneprotocol.com",
			expected: Network{
				Kind:     Mainnet,
				ID:       constants.MainnetID,
				Endpoint: "https://node.dioneprotocol.com",
			},
		},
		{
			name:     "Testnet network",
			kind:     Testnet,
			id:       constants.TestnetID,
			endpoint: "https://testnode.dioneprotocol.com",
			expected: Network{
				Kind:     Testnet,
				ID:       constants.TestnetID,
				Endpoint: "https://testnode.dioneprotocol.com",
			},
		},
		{
			name:     "Devnet network",
			kind:     Devnet,
			id:       0,
			endpoint: "http://127.0.0.1:9650",
			expected: Network{
				Kind:     Devnet,
				ID:       0,
				Endpoint: "http://127.0.0.1:9650",
			},
		},
		{
			name:     "Empty endpoint",
			kind:     Mainnet,
			id:       123,
			endpoint: "",
			expected: Network{
				Kind:     Mainnet,
				ID:       123,
				Endpoint: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewNetwork(tt.kind, tt.id, tt.endpoint)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTestnetNetwork(t *testing.T) {
	// Test with LOCAL_NODE environment variable not set
	t.Run("default testnet endpoint", func(t *testing.T) {
		// Clear any existing LOCAL_NODE environment variable
		os.Unsetenv("LOCAL_NODE")

		network := TestnetNetwork()
		assert.Equal(t, Testnet, network.Kind)
		assert.Equal(t, constants.TestnetID, network.ID)
		assert.Equal(t, TestnetAPIEndpoint, network.Endpoint)
	})

	// Test with LOCAL_NODE environment variable set to true
	t.Run("local node endpoint", func(t *testing.T) {
		os.Setenv("LOCAL_NODE", "true")
		defer os.Unsetenv("LOCAL_NODE")

		network := TestnetNetwork()
		assert.Equal(t, Testnet, network.Kind)
		assert.Equal(t, constants.TestnetID, network.ID)
		assert.Equal(t, "http://127.0.0.1:9650", network.Endpoint)
	})

	// Test with LOCAL_NODE environment variable set to false
	t.Run("local node false", func(t *testing.T) {
		os.Setenv("LOCAL_NODE", "false")
		defer os.Unsetenv("LOCAL_NODE")

		network := TestnetNetwork()
		assert.Equal(t, Testnet, network.Kind)
		assert.Equal(t, constants.TestnetID, network.ID)
		assert.Equal(t, TestnetAPIEndpoint, network.Endpoint)
	})
}

func TestMainnetNetwork(t *testing.T) {
	network := MainnetNetwork()

	assert.Equal(t, Mainnet, network.Kind)
	assert.Equal(t, constants.MainnetID, network.ID)
	assert.Equal(t, MainnetAPIEndpoint, network.Endpoint)
}

func TestDevnetNetwork(t *testing.T) {
	network := DevnetNetwork()

	assert.Equal(t, Devnet, network.Kind)
	assert.Equal(t, uint32(0), network.ID)
	assert.Equal(t, "http://127.0.0.1:9650", network.Endpoint)
}

func TestNetwork_GenesisParams(t *testing.T) {
	tests := []struct {
		name     string
		network  Network
		expected interface{} // We can't import genesis package, so we use interface{}
	}{
		{
			name: "Devnet network",
			network: Network{
				Kind: Devnet,
			},
			expected: "not nil", // We expect a non-nil value
		},
		{
			name: "Testnet network",
			network: Network{
				Kind: Testnet,
			},
			expected: "not nil", // We expect a non-nil value
		},
		{
			name: "Mainnet network",
			network: Network{
				Kind: Mainnet,
			},
			expected: "not nil", // We expect a non-nil value
		},
		{
			name: "Undefined network",
			network: Network{
				Kind: Undefined,
			},
			expected: nil,
		},
		{
			name: "Invalid network kind",
			network: Network{
				Kind: NetworkKind(999),
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.network.GenesisParams()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestNetwork_BlockchainEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		network      Network
		blockchainID string
		expected     string
	}{
		{
			name: "Testnet with blockchain ID",
			network: Network{
				Endpoint: "https://testnode.dioneprotocol.com",
			},
			blockchainID: "abc123",
			expected:     "https://testnode.dioneprotocol.com/ext/bc/abc123/rpc",
		},
		{
			name: "Mainnet with blockchain ID",
			network: Network{
				Endpoint: "https://node.dioneprotocol.com",
			},
			blockchainID: "def456",
			expected:     "https://node.dioneprotocol.com/ext/bc/def456/rpc",
		},
		{
			name: "Local network with blockchain ID",
			network: Network{
				Endpoint: "http://127.0.0.1:9650",
			},
			blockchainID: "local123",
			expected:     "http://127.0.0.1:9650/ext/bc/local123/rpc",
		},
		{
			name: "Empty blockchain ID",
			network: Network{
				Endpoint: "https://testnode.dioneprotocol.com",
			},
			blockchainID: "",
			expected:     "https://testnode.dioneprotocol.com/ext/bc//rpc",
		},
		{
			name: "Complex blockchain ID",
			network: Network{
				Endpoint: "https://node.dioneprotocol.com",
			},
			blockchainID: "subnet-abc123-def456",
			expected:     "https://node.dioneprotocol.com/ext/bc/subnet-abc123-def456/rpc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.network.BlockchainEndpoint(tt.blockchainID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNetwork_BlockchainWSEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		network      Network
		blockchainID string
		expected     string
	}{
		{
			name: "HTTPS endpoint",
			network: Network{
				Endpoint: "https://testnode.dioneprotocol.com",
			},
			blockchainID: "abc123",
			expected:     "ws://testnode.dioneprotocol.com/ext/bc/abc123/ws",
		},
		{
			name: "HTTP endpoint",
			network: Network{
				Endpoint: "http://127.0.0.1:9650",
			},
			blockchainID: "def456",
			expected:     "ws://127.0.0.1:9650/ext/bc/def456/ws",
		},
		{
			name: "Endpoint without protocol",
			network: Network{
				Endpoint: "testnode.dioneprotocol.com",
			},
			blockchainID: "abc123",
			expected:     "ws://testnode.dioneprotocol.com/ext/bc/abc123/ws",
		},
		{
			name: "Complex blockchain ID",
			network: Network{
				Endpoint: "https://node.dioneprotocol.com",
			},
			blockchainID: "subnet-abc123-def456",
			expected:     "ws://node.dioneprotocol.com/ext/bc/subnet-abc123-def456/ws",
		},
		{
			name: "Empty blockchain ID",
			network: Network{
				Endpoint: "https://testnode.dioneprotocol.com",
			},
			blockchainID: "",
			expected:     "ws://testnode.dioneprotocol.com/ext/bc//ws",
		},
		{
			name: "Endpoint with port",
			network: Network{
				Endpoint: "https://testnode.dioneprotocol.com:8080",
			},
			blockchainID: "abc123",
			expected:     "ws://testnode.dioneprotocol.com:8080/ext/bc/abc123/ws",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.network.BlockchainWSEndpoint(tt.blockchainID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNetworkFromURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected Network
	}{
		{
			name:     "Testnet API endpoint",
			uri:      TestnetAPIEndpoint,
			expected: TestnetNetwork(),
		},
		{
			name:     "Mainnet API endpoint",
			uri:      MainnetAPIEndpoint,
			expected: MainnetNetwork(),
		},
		{
			name:     "Local node endpoint with LOCAL_NODE=true",
			uri:      "http://127.0.0.1:9650",
			expected: TestnetNetwork(), // Will be handled in the test case
		},
		{
			name:     "Local node endpoint with LOCAL_NODE=false",
			uri:      "http://127.0.0.1:9650",
			expected: UndefinedNetwork,
		},
		{
			name:     "Unknown endpoint",
			uri:      "https://unknown.example.com",
			expected: UndefinedNetwork,
		},
		{
			name:     "Empty URI",
			uri:      "",
			expected: UndefinedNetwork,
		},
		{
			name:     "Different local endpoint",
			uri:      "http://localhost:9650",
			expected: UndefinedNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle special case for LOCAL_NODE=true
			if tt.name == "Local node endpoint with LOCAL_NODE=true" {
				os.Setenv("LOCAL_NODE", "true")
				defer os.Unsetenv("LOCAL_NODE")
				expected := TestnetNetwork() // Get expected value with env var set
				result := NetworkFromURI(tt.uri)
				assert.Equal(t, expected.Kind, result.Kind)
				assert.Equal(t, expected.ID, result.ID)
				assert.Equal(t, expected.Endpoint, result.Endpoint)
			} else {
				// Clear LOCAL_NODE environment variable for other tests
				os.Unsetenv("LOCAL_NODE")
				result := NetworkFromURI(tt.uri)
				assert.Equal(t, tt.expected.Kind, result.Kind)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Endpoint, result.Endpoint)
			}
		})
	}
}

func TestNetworkFromURI_LocalNodeEnvironment(t *testing.T) {
	// Test the LOCAL_NODE environment variable behavior specifically
	t.Run("LOCAL_NODE=true", func(t *testing.T) {
		os.Setenv("LOCAL_NODE", "true")
		defer os.Unsetenv("LOCAL_NODE")

		result := NetworkFromURI("http://127.0.0.1:9650")
		expected := TestnetNetwork()

		assert.Equal(t, expected.Kind, result.Kind)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Endpoint, result.Endpoint)
	})

	t.Run("LOCAL_NODE=false", func(t *testing.T) {
		os.Setenv("LOCAL_NODE", "false")
		defer os.Unsetenv("LOCAL_NODE")

		result := NetworkFromURI("http://127.0.0.1:9650")
		expected := UndefinedNetwork

		assert.Equal(t, expected.Kind, result.Kind)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Endpoint, result.Endpoint)
	})

	t.Run("LOCAL_NODE not set", func(t *testing.T) {
		os.Unsetenv("LOCAL_NODE")

		result := NetworkFromURI("http://127.0.0.1:9650")
		expected := UndefinedNetwork

		assert.Equal(t, expected.Kind, result.Kind)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Endpoint, result.Endpoint)
	})
}

func TestUndefinedNetwork(t *testing.T) {
	// Test that UndefinedNetwork is properly initialized
	assert.Equal(t, NetworkKind(0), UndefinedNetwork.Kind)
	assert.Equal(t, uint32(0), UndefinedNetwork.ID)
	assert.Equal(t, "", UndefinedNetwork.Endpoint)
}

func TestNetworkConstants(t *testing.T) {
	// Test that the API endpoint constants are properly defined
	assert.NotEmpty(t, TestnetAPIEndpoint)
	assert.NotEmpty(t, MainnetAPIEndpoint)
	assert.Contains(t, TestnetAPIEndpoint, "testnode")
	assert.Contains(t, MainnetAPIEndpoint, "node")
	assert.Contains(t, TestnetAPIEndpoint, "dioneprotocol.com")
	assert.Contains(t, MainnetAPIEndpoint, "dioneprotocol.com")
}

// TestNetwork_GetMinStakingAmount tests the GetMinStakingAmount method
// Note: This test requires a real network connection, so we'll test the error case
func TestNetwork_GetMinStakingAmount_Error(t *testing.T) {
	// Create a network with an invalid endpoint to test error handling
	network := Network{
		Kind:     Testnet,
		ID:       constants.TestnetID,
		Endpoint: "http://invalid-endpoint:9999",
	}

	amount, err := network.GetMinStakingAmount()

	// We expect an error since the endpoint is invalid
	assert.Error(t, err)
	assert.Equal(t, uint64(0), amount)
}

// TestNetwork_GetMinStakingAmount_Success tests with a valid endpoint
// This test is commented out as it requires a real network connection
// Uncomment and run when you have network access
/*
func TestNetwork_GetMinStakingAmount_Success(t *testing.T) {
	network := TestnetNetwork()

	amount, err := network.GetMinStakingAmount()

	require.NoError(t, err)
	assert.Greater(t, amount, uint64(0))
}
*/
