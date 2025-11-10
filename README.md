# Odyssey Tooling Go SDK

The official Odyssey Tooling Go SDK library.

*** Please note that this SDK is in experimental mode, major changes to the SDK are to be expected
in between releases ***

Current version (v1.0.0) currently supports: 
- Create Subnet and Create Blockchain in a Subnet in Testnet / Mainnet. 
- Enable Odyssey nodes to validate Primary Network
- Adding Validators to a Subnet
- Ledger SDK
- Working with existing Odyssey Nodes (Validator / API / Monitoring / Load Test Node)

Currently, both stored keys and Ledger hardware wallets are supported for transaction building and signing.


## Getting Started

### Installing
Use `go get` to retrieve the SDK to add it to your project's Go module dependencies.

	go get github.com/DioneProtocol/odyssey-tooling-sdk-go

To update the SDK use `go get -u` to retrieve the latest version of the SDK.

	go get -u github.com/DioneProtocol/odyssey-tooling-sdk-go

### Local Development and Testing

#### Setting up a Local Testnet Node

For local development and testing, you can run a local odyssey testnet node using Docker. This eliminates 429 rate limiting errors and provides faster test execution.

1. **Clone and set up the odysseygo-installer repository:**
   ```bash
   git clone https://github.com/DioneProtocol/odysseygo-installer.git
   cd odysseygo-installer/docker
   ```
   
   Follow the instructions in the repository to run an archival node for testnet.

2. **Verify the node is running:**
   ```bash
   curl -X POST http://127.0.0.1:9650/ext/info \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"info.getNodeID"}'
   ```

#### Running Tests with Local Node

The SDK automatically detects and uses the local node when the `LOCAL_NODE` environment variable is set to `true`.

**Run all tests with local node:**
```bash
LOCAL_NODE=true go test ./... -timeout 10m
```

**Run specific package tests:**
```bash
LOCAL_NODE=true go test ./wallet -v
LOCAL_NODE=true go test ./evm -v
LOCAL_NODE=true go test ./subnet -v
```

**Get test coverage:**
```bash
LOCAL_NODE=true go test ./... -coverprofile=coverage.out -timeout 10m
go tool cover -func=coverage.out
```

**Benefits of using local node:**
- ✅ **No 429 rate limiting errors**
- ✅ **Faster test execution** (local vs remote)
- ✅ **More reliable testing** (no network dependencies)
- ✅ **Consistent test environment**

**Note:** When `LOCAL_NODE=true` is not set, tests will use the official testnet endpoints and may experience rate limiting.

## Quick Examples

### Subnet SDK Example

This example shows how to create a Subnet Genesis, deploy the Subnet into Testnet Network and create
a blockchain in the Subnet. 

This examples also shows how to create a key pair to pay for transactions, how to create a Wallet
object that will be used to build and sign CreateSubnetTx and CreateChainTx and how to commit these 
transactions on chain.

More examples can be found at examples directory.

```go
package main

import (
	"context"
	"fmt"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/node"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/subnet"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/vm"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

// Creates a Subnet Genesis, deploys the Subnet into Testnet Network using CreateSubnetTx
// and creates a blockchain in the Subnet using CreateChainTx
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
	// Number of keys in subnetAuthKeys has to be more than or equal to the threshold value provided during 
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
			DIONEKeychain:     keychain.Keychain,
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

// Add a validator to Subnet
func AddSubnetValidator() {
	// We are using existing Subnet that we have already deployed on Testnet
	subnetParams := subnet.SubnetParams{
		GenesisFilePath: "GENESIS_FILE_PATH",
		Name:            "SUBNET_NAME",
	}

	newSubnet, err := subnet.New(&subnetParams)
	if err != nil {
		panic(err)
	}

	subnetID, err := ids.FromString("SUBNET_ID")
	if err != nil {
		panic(err)
	}

	// Genesis doesn't contain the deployed Subnet's SubnetID, we need to first set the Subnet ID
	newSubnet.SetSubnetID(subnetID)

	// We are using existing host
	node := node.Node{
		// NodeID is Odyssey Node ID of the node
		NodeID: "NODE_ID",
		// IP address of the node
		IP: "NODE_IP_ADDRESS",
		// SSH configuration for the node
		SSHConfig: node.SSHConfig{
			User:           constants.RemoteHostUser,
			PrivateKeyPath: "NODE_KEYPAIR_PRIVATE_KEY_PATH",
		},
		// Role is the role that we expect the host to be (Validator, API, Loadtest or Monitor)
		Roles: []node.SupportedRole{node.Validator},
	}

	// Here we are assuming that the node is currently validating the Primary Network, which is
	// a requirement before the node can start validating a Subnet.
	// To have a node validate the Primary Network, call node.ValidatePrimaryNetwork
	// Now we are calling the node to start tracking the Subnet
	subnetIDsToValidate := []string{newSubnet.SubnetID.String()}
	if err := node.SyncSubnets(subnetIDsToValidate); err != nil {
		panic(err)
	}

	// Node is now tracking the Subnet

	// Key that will be used for paying the transaction fees of Subnet AddValidator Tx
	//
	// In our example, this Key is also the control Key to the Subnet, so we are going to use
	// this key to also sign the Subnet AddValidator tx
	network := odyssey.TestnetNetwork()
	keychain, err := keychain.NewKeychain(network, "PRIVATE_KEY_FILEPATH", nil)
	if err != nil {
		panic(err)
	}

	wallet, err := wallet.New(
		context.Background(),
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:     keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: set.Of(subnetID),
		},
	)
	if err != nil {
		panic(err)
	}

	nodeID, err := ids.NodeIDFromString(node.NodeID)
	if err != nil {
		panic(err)
	}

	validatorParams := validator.SubnetValidatorParams{
		NodeID: nodeID,
		// Validate Subnet for 48 hours
		Duration: 48 * time.Hour,
		// Weight is the validator's weight when sampling validators
		// Weight defaults to 20 if not set
		Weight: 20,
	}

	// We need to set Subnet Auth Keys for this transaction since Subnet AddValidator is
	// a Subnet-changing transaction
	//
	// In this example, the example Subnet was created with only 1 key as control key with a threshold of 1
	// and the control key is the key contained in the keychain object, so we are going to use the
	// key contained in the keychain object as the Subnet Auth Key for Subnet AddValidator tx
	subnetAuthKeys := keychain.Addresses().List()
	newSubnet.SetSubnetAuthKeys(subnetAuthKeys)

	addValidatorTx, err := newSubnet.AddValidator(wallet, validatorParams)
	if err != nil {
		panic(err)
	}

	// Since it has the required signatures, we will now commit the transaction on chain
	txID, err := newSubnet.Commit(*addValidatorTx, wallet, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("obtained tx id %s", txID.String())
}

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
```

### Nodes SDK Example

This example shows how to work with existing Odyssey Nodes and enable them to validate the Primary Network.

**Note:** Cloud functionality for creating nodes has been removed from this SDK. This SDK now focuses on managing and interacting with existing nodes that you have already set up locally or through other means.

More examples can be found at examples directory.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/node"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/units"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

func ValidatePrimaryNetwork() {
	// We are using existing host
	node := node.Node{
		// NodeID is Odyssey Node ID of the node
		NodeID: "NODE_ID",
		// IP address of the node
		IP: "NODE_IP_ADDRESS",
		// SSH configuration for the node
		SSHConfig: node.SSHConfig{
			User:           constants.RemoteHostUser,
			PrivateKeyPath: "NODE_KEYPAIR_PRIVATE_KEY_PATH",
		},
		// Role of the node can be Validator, API, Loadtest, or Monitor
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
		// Stake 2 DIONE
		StakeAmount: 2 * units.Dione,
		// DelegationFee is optional - if not set, the minimum delegation fee for the network will be used
	}

	// Key that will be used for paying the transaction fee of AddValidator Tx
	network := odyssey.TestnetNetwork()
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

	txID, err := node.ValidatePrimaryNetwork(network, validatorParams, wallet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("obtained tx id %s", txID.String())
}

```
