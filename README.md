# Odyssey Tooling Go SDK

The official Odyssey Tooling Go SDK library.

*** Please note that this SDK is in experimental mode, major changes to the SDK are to be expected
in between releases ***

Current version (v0.0.1) currently supports: 
- Create Subnet and Create Blockchain in a Subnet in Fuji / Mainnet. 
- Create Odyssey Node (Validator / API / Monitoring / Load Test Node) & install all required
dependencies (odysseygo, gcc, Promtail, Grafana, etc).
- Enable Odyssey nodes to validate Primary Network
- Adding Validators to a Subnet
- Ledger SDK

Currently, only stored keys are supported for transaction building and signing, ledger support is 
coming soon.

Future SDK releases will contain the following features (in order of priority): 
- Additional Nodes SDK features (Devnet)
- Additional Subnet SDK features (Custom Subnets)
- Teleporter SDK

## Getting Started

### Installing
Use `go get` to retrieve the SDK to add it to your project's Go module dependencies.

	go get github.com/DioneProtocol/odyssey-tooling-sdk-go

To update the SDK use `go get -u` to retrieve the latest version of the SDK.

	go get -u github.com/DioneProtocol/odyssey-tooling-sdk-go

## Quick Examples

### Subnet SDK Example

This example shows how to create a Subnet Genesis, deploy the Subnet into Fuji Network and create
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
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/subnet"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/teleporter"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/vm"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

// Creates a Subnet Genesis, deploys the Subnet into Fuji Network using CreateSubnetTx
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

	// we need to wait to allow the transaction to reach other nodes in Fuji
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
	// We are using existing Subnet that we have already deployed on Fuji
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
		// Role is the role that we expect the host to be (Validator, API, AWMRelayer, Loadtest or
		// Monitor)
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
		Weight:   20,
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

This example shows how to create Odyssey Validator Nodes (SDK function for nodes to start 
validating Primary Network / Subnet will be available in the next Odyssey Tooling SDK release).

This examples also shows how to create an Odyssey Monitoring Node, which enables you to have a 
centralized Grafana Dashboard where you can view metrics relevant to any Validator & API nodes that
the monitoring node is linked to as well as a centralized logs for the A/O/D Chain and Subnet logs 
for the Validator & API nodes. An example on how the dashboard and logs look like can be found at https://docs.dione.network/tooling/cli-create-nodes/create-a-validator-aws

More examples can be found at examples directory.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	awsAPI "github.com/DioneProtocol/odyssey-tooling-sdk-go/cloud/aws"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/node"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
)

func CreateNodes() {
	ctx := context.Background()

	// Get the default cloud parameters for AWS
	cp, err := node.GetDefaultCloudParams(ctx, node.AWSCloud)
	if err != nil {
		panic(err)
	}

	securityGroupName := "SECURITY_GROUP_NAME"
	// Create a new security group in AWS if you do not currently have one in the selected
	// AWS region.
	sgID, err := awsAPI.CreateSecurityGroup(ctx, securityGroupName, cp.AWSConfig.AWSProfile, cp.Region)
	if err != nil {
		panic(err)
	}
	// Set the security group we are using when creating our Odyssey Nodes
	cp.AWSConfig.AWSSecurityGroupID = sgID

	keyPairName := "KEY_PAIR_NAME"
	sshPrivateKeyPath := utils.ExpandHome("PRIVATE_KEY_FILEPATH")
	// Create a new AWS SSH key pair if you do not currently have one in your selected AWS region.
	// Note that the created key pair can only be used in the region that it was created in.
	// The private key to the created key pair will be stored in the filepath provided in
	// sshPrivateKeyPath.
	if err := awsAPI.CreateSSHKeyPair(ctx, cp.AWSConfig.AWSProfile, cp.Region, keyPairName, sshPrivateKeyPath); err != nil {
		panic(err)
	}
	// Set the key pair we are using when creating our Odyssey Nodes
	cp.AWSConfig.AWSKeyPair = keyPairName

	// Odyssey-CLI is installed in nodes to enable them to join subnets as validators
	// Odyssey-CLI dependency by Odyssey nodes will be deprecated in the next release
	// of Odyssey Tooling SDK
	const (
		odysseygoVersion  = "v1.11.8"
	)

	// Create two new Odyssey Validator nodes on Fuji Network on AWS without Elastic IPs
	// attached. Once CreateNodes is completed, the validators will begin bootstrapping process
	// to Primary Network in Fuji Network. Nodes need to finish bootstrapping process
	// before they can validate Odyssey Primary Network / Subnet.
	//
	// SDK function for nodes to start validating Primary Network / Subnet will be available
	// in the next Odyssey Tooling SDK release.
	hosts, err := node.CreateNodes(ctx,
		&node.NodeParams{
			CloudParams:         cp,
			Count:               2,
			Roles:               []node.SupportedRole{node.Validator},
			Network:             odyssey.TestnetNetwork(),
			odysseygoVersion:  odysseygoVersion,
			UseStaticIP:         false,
			SSHPrivateKeyPath:   sshPrivateKeyPath,
		})
	if err != nil {
		panic(err)
	}

	const (
		sshTimeout        = 120 * time.Second
		sshCommandTimeout = 10 * time.Second
	)

	// Examples showing how to run ssh commands on the created nodes
	for _, h := range hosts {
		// Wait for the host to be ready (only needs to be done once for newly created nodes)
		fmt.Println("Waiting for SSH shell")
		if err := h.WaitForSSHShell(sshTimeout); err != nil {
			panic(err)
		}
		fmt.Println("SSH shell ready to execute commands")
		// Run a command on the host
		if output, err := h.Commandf(nil, sshCommandTimeout, "echo 'Hello, %s!'", "World"); err != nil {
			panic(err)
		} else {
			fmt.Println(string(output))
		}
		// sleep for 10 seconds allowing odysseygo container to start
		time.Sleep(10 * time.Second)
		// check if odysseygo is running
		if output, err := h.Commandf(nil, sshCommandTimeout, "docker ps"); err != nil {
			panic(err)
		} else {
			fmt.Println(string(output))
		}
	}

	// Create a monitoring node.
	// Monitoring node enables you to have a centralized Grafana Dashboard where you can view
	// metrics relevant to any Validator & API nodes that the monitoring node is linked to as well
	// as a centralized logs for the A/O/D Chain and Subnet logs for the Validator & API nodes.
	// An example on how the dashboard and logs look like can be found at https://docs.dione.network/tooling/cli-create-nodes/create-a-validator-aws
	monitoringHosts, err := node.CreateNodes(ctx,
		&node.NodeParams{
			CloudParams:       cp,
			Count:             1,
			Roles:             []node.SupportedRole{node.Monitor},
			UseStaticIP:       false,
			SSHPrivateKeyPath: sshPrivateKeyPath,
		})
	if err != nil {
		panic(err)
	}

	// Link the 2 validator nodes previously created with the monitoring host so that
	// the monitoring host can start tracking the validator nodes metrics and collecting their logs
	if err := monitoringHosts[0].MonitorNodes(ctx, hosts, ""); err != nil {
		panic(err)
	}
}

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
		// Stake 2 DIONE
		StakeAmount: 2 * units.Dione,
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
			DIONEKeychain:     keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		},
	)
	if err != nil {
		panic(err)
	}

	txID, err := node.ValidatePrimaryNetwork(odyssey.TestnetNetwork(), validatorParams, wallet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("obtained tx id %s", txID.String())
}

```
