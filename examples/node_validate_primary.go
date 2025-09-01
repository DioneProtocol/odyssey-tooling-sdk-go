// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/avalanche"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/node"
)

func ValidatePrimaryNetwork() {
	// We are using existing host
	node := node.Node{
		// NodeID is Avalanche Node ID of the node
		NodeID: "NODE_ID",
		// IP address of the node
		IP: "NODE_IP_ADDRESS",
		// SSH configuration for the node
		SSHConfig: node.SSHConfig{
			User:           constants.RemoteHostUser,
			PrivateKeyPath: "NODE_KEYPAIR_PRIVATE_KEY_PATH",
		},
		// Role of the node can be 	Validator, API, AWMRelayer, Loadtest, or Monitor
		Roles: []node.SupportedRole{node.Validator},
	}

	nodeID, err := ids.NodeIDFromString(node.NodeID)
	if err != nil {
		panic(err)
	}

	validatorParams := validator.PrimaryNetworkValidatorParams{
		NodeID: nodeID,
		// Validate Primary Network for 48 hours
		Duration: 48 * time.Hour,
		// Stake 2 AVAX
		StakeAmount: 2 * units.Dione,
	}

	// Key that will be used for paying the transaction fee of AddValidator Tx
	network := avalanche.TestnetNetwork()
	keychain, err := keychain.NewKeychain(network, "PRIVATE_KEY_FILEPATH", nil)
	if err != nil {
		panic(err)
	}

	wallet, err := wallet.New(
		context.Background(),
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		},
	)
	if err != nil {
		panic(err)
	}

	txID, err := node.ValidatePrimaryNetwork(avalanche.TestnetNetwork(), validatorParams, wallet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("obtained tx id %s", txID.String())
}
