// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package examples

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/DioneProtocol/odysseygo/utils/formatting/address"
	"github.com/DioneProtocol/odysseygo/utils/set"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/subnet"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/vm"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

func getDefaultSubnetEVMGenesis() subnet.SubnetParams {
	allocation := core.GenesisAlloc{}
	defaultAmount, _ := new(big.Int).SetString(vm.DefaultEvmAirdropAmount, 10)
	allocation[common.HexToAddress("INITIAL_ALLOCATION_ADDRESS")] = core.GenesisAccount{
		Balance: defaultAmount,
	}
	return subnet.SubnetParams{
		SubnetEVM: &subnet.SubnetEVMParams{
			ChainID:     big.NewInt(123456),
			FeeConfig:   vm.StarterFeeConfig,
			Allocation:  allocation,
			Precompiles: params.Precompiles{},
		},
		Name: "TestSubnet",
	}
}

func DeploySubnet() {
	subnetParams := getDefaultSubnetEVMGenesis()
	// Create new Subnet EVM genesis
	newSubnet, _ := subnet.New(&subnetParams)

	network := odyssey.TestnetNetwork()

	// Key that will be used for paying the transaction fees of CreateSubnetTx and CreateChainTx
	// NewKeychain will generate a new key pair in the provided path if no .pk file currently
	// exists in the provided path
	keychain, _ := keychain.NewKeychain(network, "KEY_PATH", nil)

	// In this example, we are using the fee-paying key generated above also as control key
	// and subnet auth key

	// Control keys are a list of keys that are permitted to make changes to a Subnet
	// such as creating a blockchain in the Subnet and adding validators to the Subnet
	controlKeys := keychain.Addresses().List()

	// Subnet auth keys are a subset of control keys that will be used to sign transactions that
	// modify a Subnet (such as creating a blockchain in the Subnet and adding validators to the
	// Subnet)
	//
	// Number of keys in subnetAuthKeys has to be equal to the threshold value provided during
	// CreateSubnetTx.
	//
	// All keys in subnetAuthKeys have to sign the transaction before the transaction
	subnetAuthKeys := keychain.Addresses().List()
	threshold := 1
	newSubnet.SetSubnetControlParams(controlKeys, uint32(threshold))

	wallet, _ := wallet.New(
		context.Background(),
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		},
	)

	// Build and Sign CreateSubnetTx with our fee paying key
	deploySubnetTx, _ := newSubnet.CreateSubnetTx(wallet)
	// Commit our CreateSubnetTx on chain
	subnetID, _ := newSubnet.Commit(*deploySubnetTx, wallet, true)
	fmt.Printf("subnetID %s \n", subnetID.String())

	// we need to wait to allow the transaction to reach other nodes in Testnet
	time.Sleep(2 * time.Second)

	newSubnet.SetSubnetAuthKeys(subnetAuthKeys)
	// Build and Sign CreateChainTx with our fee paying key (which is also our subnet auth key)
	deployChainTx, _ := newSubnet.CreateBlockchainTx(wallet)
	// Commit our CreateChainTx on chain
	// Since we are using the fee paying key as control key too, we can commit the transaction
	// on chain immediately since the number of signatures has been reached
	blockchainID, _ := newSubnet.Commit(*deployChainTx, wallet, true)
	fmt.Printf("blockchainID %s \n", blockchainID.String())
}

func DeploySubnetWithLedger() {
	subnetParams := getDefaultSubnetEVMGenesis()
	newSubnet, _ := subnet.New(&subnetParams)
	network := odyssey.TestnetNetwork()

	// Create keychain with a specific Ledger address. More than 1 address can be used.
	//
	// Alternatively, keychain can also be created from Ledger without specifying any Ledger address
	// by stating the amount of DIONE required to pay for transaction fees. Keychain SDK will
	// then look through all indexes of all addresses in the Ledger until sufficient DIONE balance
	// is reached. For example:
	//
	// fee := network.GenesisParams().CreateBlockchainTxFee + network.GenesisParams().CreateSubnetTxFee
	// ledgerInfo := keychain.LedgerParams{
	//	 RequiredFunds: fee,
	// }
	//
	// To view Ledger addresses and their balances, you can use Odyssey CLI and use the command
	// odyssey key list --ledger [0,1,2,3,4]
	// The example command above will list the first five addresses in your Ledger
	//
	// To transfer funds between addresses in Ledger, refer to https://docs.dione.network/tooling/cli-transfer-funds/how-to-transfer-funds
	ledgerInfo := keychain.LedgerParams{
		LedgerAddresses: []string{"O-testnetxxxxxxxxx"},
	}

	// Here we are creating keychain A which will be used as fee-paying key for CreateSubnetTx
	// and CreateChainTx
	keychainA, _ := keychain.NewKeychain(network, "", &ledgerInfo)
	walletA, _ := wallet.New(
		context.Background(),
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychainA.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		},
	)

	// In this example, we are using a key different from fee-paying key generated above
	// as control key and subnet auth key
	addressesIDs, _ := address.ParseToIDs([]string{"O-testnetyyyyyyyy"})
	controlKeys := addressesIDs
	subnetAuthKeys := addressesIDs
	threshold := 1
	newSubnet.SetSubnetControlParams(controlKeys, uint32(threshold))

	// Pay and Sign CreateSubnet Tx with fee paying key A using Ledger
	deploySubnetTx, _ := newSubnet.CreateSubnetTx(walletA)
	subnetID, _ := newSubnet.Commit(*deploySubnetTx, walletA, true)
	fmt.Printf("subnetID %s \n", subnetID.String())

	// we need to wait to allow the transaction to reach other nodes in Testnet
	time.Sleep(2 * time.Second)

	newSubnet.SetSubnetAuthKeys(subnetAuthKeys)

	// Pay and sign CreateChain Tx with fee paying key A using Ledger
	deployChainTx, _ := newSubnet.CreateBlockchainTx(walletA)

	// We have to first disconnect Ledger to avoid errors when signing with our subnet auth
	// keys later
	_ = keychainA.Ledger.LedgerDevice.Disconnect()

	// Here we are creating keychain B using the Ledger address that was used as the control key and
	// subnet auth key for our subnet.
	ledgerInfoB := keychain.LedgerParams{
		LedgerAddresses: []string{"O-testnetyyyyyyyy"},
	}
	keychainB, _ := keychain.NewKeychain(network, "", &ledgerInfoB)

	// include subnetID in OChainTxsToFetch when creating second wallet
	walletB, _ := wallet.New(
		context.Background(),
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychainB.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: set.Of(subnetID),
		},
	)

	// Sign with our Subnet auth key
	_ = walletB.O().Signer().Sign(context.Background(), deployChainTx.OChainTx)

	// Now that the number of signatures have been met, we can commit our transaction
	// on chain
	blockchainID, _ := newSubnet.Commit(*deployChainTx, walletB, true)
	fmt.Printf("blockchainID %s \n", blockchainID.String())
}
