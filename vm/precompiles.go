// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/DioneProtocol/subnet-evm/precompile/allowlist"
	"github.com/DioneProtocol/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/DioneProtocol/subnet-evm/precompile/contracts/txallowlist"

	// "github.com/DioneProtocol/subnet-evm/precompile/contracts/warp"
	"github.com/ethereum/go-ethereum/common"
)

// func ConfigureWarp(timestamp *uint64) warp.Config {
// 	return ConfigureWarpWithConfig(timestamp, DefaultVMConfig())
// }

// func ConfigureWarpWithConfig(timestamp *uint64, config *VMConfig) warp.Config {
// 	if !config.EnableWarp {
// 		// Return empty config when warp is disabled
// 		return warp.Config{}
// 	}

// 	// Original warp configuration when enabled
// 	warpConfig := warp.Config{
// 		QuorumNumerator: warp.WarpDefaultQuorumNumerator,
// 	}
// 	// Note: Type mismatch between DioneProtocol and DioneProtocol packages
// 	// This will need to be resolved when warp is actually enabled
// 	// For now, return empty config to avoid compilation errors
// 	return warpConfig
// }

// AddTeleporterAddressesToAllowLists adds teleporter-related addresses (main funded key, messenger
// deploy key, relayer key) to the allow list of relevant enabled precompiles
func AddTeleporterAddressesToAllowLists(
	config params.ChainConfig,
	teleporterAddress string,
	teleporterMessengerDeployerAddress string,
	relayerAddress string,
) params.ChainConfig {
	// tx allow list:
	// teleporterAddress funds the other two and also deploys the registry
	// teleporterMessengerDeployerAddress deploys the messenger
	// relayerAddress is used by the relayer to send txs to the target chain
	for _, address := range []string{teleporterAddress, teleporterMessengerDeployerAddress, relayerAddress} {
		precompileConfig := config.GenesisPrecompiles[txallowlist.ConfigKey]
		if precompileConfig != nil {
			txAllowListConfig := precompileConfig.(*txallowlist.Config)
			txAllowListConfig.AllowListConfig = addAddressToAllowed(
				txAllowListConfig.AllowListConfig,
				address,
			)
		}
	}
	// contract deploy allow list:
	// teleporterAddress deploys the registry
	// teleporterMessengerDeployerAddress deploys the messenger
	for _, address := range []string{teleporterAddress, teleporterMessengerDeployerAddress} {
		precompileConfig := config.GenesisPrecompiles[deployerallowlist.ConfigKey]
		if precompileConfig != nil {
			txAllowListConfig := precompileConfig.(*deployerallowlist.Config)
			txAllowListConfig.AllowListConfig = addAddressToAllowed(
				txAllowListConfig.AllowListConfig,
				address,
			)
		}
	}
	return config
}

// adds an address to the given allowlist, as an Allowed address,
// if it is not yet Admin, Manager or Allowed
func addAddressToAllowed(
	allowListConfig allowlist.AllowListConfig,
	addressStr string,
) allowlist.AllowListConfig {
	address := common.HexToAddress(addressStr)
	allowed := false
	if utils.Belongs(
		allowListConfig.AdminAddresses,
		address,
	) {
		allowed = true
	}
	if utils.Belongs(
		allowListConfig.ManagerAddresses,
		address,
	) {
		allowed = true
	}
	if utils.Belongs(
		allowListConfig.EnabledAddresses,
		address,
	) {
		allowed = true
	}
	if !allowed {
		allowListConfig.EnabledAddresses = append(
			allowListConfig.EnabledAddresses,
			address,
		)
	}
	return allowListConfig
}
