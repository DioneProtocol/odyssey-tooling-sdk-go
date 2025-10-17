// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
	"github.com/DioneProtocol/subnet-evm/precompile/allowlist"

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
