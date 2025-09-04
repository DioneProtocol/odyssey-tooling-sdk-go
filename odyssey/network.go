// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package odyssey

import (
	"fmt"
	"strings"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
	"github.com/DioneProtocol/odysseygo/genesis"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/constants"
	"github.com/DioneProtocol/odysseygo/vms/omegavm"
)

type NetworkKind int64

const (
	Undefined NetworkKind = iota
	Mainnet
	Testnet
	Devnet
)

const (
	TestnetAPIEndpoint = "https://testnode.dioneprotocol.com"
	MainnetAPIEndpoint = "https://node.dioneprotocol.com"
)

func (nk NetworkKind) String() string {
	switch nk {
	case Mainnet:
		return "mainnet"
	case Testnet:
		return "testnet"
	case Devnet:
		return "local"
	}
	return "invalid network"
}

type Network struct {
	Kind     NetworkKind
	ID       uint32
	Endpoint string
}

var UndefinedNetwork = Network{}

func (n Network) HRP() string {
	switch n.ID {
	case constants.TestnetID:
		return constants.TestnetHRP // Returns "testnet"
	case constants.MainnetID:
		return constants.MainnetHRP // Returns "dione"
	default:
		return constants.FallbackHRP // Returns "custom"
	}
}

func NetworkFromNetworkID(networkID uint32) Network {
	switch networkID {
	case constants.MainnetID:
		return MainnetNetwork()
	case constants.TestnetID:
		return TestnetNetwork()
	}
	return UndefinedNetwork
}

func NewNetwork(kind NetworkKind, id uint32, endpoint string) Network {
	return Network{
		Kind:     kind,
		ID:       id,
		Endpoint: endpoint,
	}
}

func TestnetNetwork() Network {
	return NewNetwork(Testnet, constants.TestnetID, TestnetAPIEndpoint)
}

func MainnetNetwork() Network {
	return NewNetwork(Mainnet, constants.MainnetID, MainnetAPIEndpoint)
}

func (n Network) GenesisParams() *genesis.Params {
	switch n.Kind {
	case Devnet:
		return &genesis.LocalParams
	case Testnet:
		return &genesis.TestnetParams
	case Mainnet:
		return &genesis.MainnetParams
	}
	return nil
}

func (n Network) BlockchainEndpoint(blockchainID string) string {
	return fmt.Sprintf("%s/ext/bc/%s/rpc", n.Endpoint, blockchainID)
}

func (n Network) BlockchainWSEndpoint(blockchainID string) string {
	trimmedURI := n.Endpoint
	trimmedURI = strings.TrimPrefix(trimmedURI, "http://")
	trimmedURI = strings.TrimPrefix(trimmedURI, "https://")
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", trimmedURI, blockchainID)
}

func (n Network) GetMinStakingAmount() (uint64, error) {
	pClient := omegavm.NewClient(n.Endpoint)
	ctx, cancel := utils.GetAPIContext()
	defer cancel()
	minValStake, _, err := pClient.GetMinStake(ctx, ids.Empty)
	if err != nil {
		return 0, err
	}
	return minValStake, nil
}

// NetworkFromURI determines the network type from a URI endpoint
func NetworkFromURI(uri string) Network {
	switch uri {
	case TestnetAPIEndpoint:
		return TestnetNetwork()
	case MainnetAPIEndpoint:
		return MainnetNetwork()
	default:
		// For unknown endpoints, return undefined network
		return UndefinedNetwork
	}
}
