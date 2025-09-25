// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"testing"

	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/DioneProtocol/subnet-evm/precompile/allowlist"
	"github.com/DioneProtocol/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/DioneProtocol/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ethereum/go-ethereum/common"
)

func TestAddAddressToAllowed(t *testing.T) {
	t.Run("add new address to empty allowlist", func(t *testing.T) {
		allowListConfig := allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{},
			ManagerAddresses: []common.Address{},
			EnabledAddresses: []common.Address{},
		}

		addressStr := "0x1234567890123456789012345678901234567890"
		result := addAddressToAllowed(allowListConfig, addressStr)

		expectedAddress := common.HexToAddress(addressStr)
		if len(result.EnabledAddresses) != 1 {
			t.Errorf("Expected 1 enabled address, got %d", len(result.EnabledAddresses))
		}
		if result.EnabledAddresses[0] != expectedAddress {
			t.Errorf("Expected address %s, got %s", expectedAddress.Hex(), result.EnabledAddresses[0].Hex())
		}
	})

	t.Run("do not add address if already admin", func(t *testing.T) {
		adminAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
		allowListConfig := allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{adminAddress},
			ManagerAddresses: []common.Address{},
			EnabledAddresses: []common.Address{},
		}

		result := addAddressToAllowed(allowListConfig, adminAddress.Hex())

		if len(result.EnabledAddresses) != 0 {
			t.Errorf("Expected 0 enabled addresses, got %d", len(result.EnabledAddresses))
		}
		if len(result.AdminAddresses) != 1 {
			t.Errorf("Expected 1 admin address, got %d", len(result.AdminAddresses))
		}
	})

	t.Run("do not add address if already manager", func(t *testing.T) {
		managerAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
		allowListConfig := allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{},
			ManagerAddresses: []common.Address{managerAddress},
			EnabledAddresses: []common.Address{},
		}

		result := addAddressToAllowed(allowListConfig, managerAddress.Hex())

		if len(result.EnabledAddresses) != 0 {
			t.Errorf("Expected 0 enabled addresses, got %d", len(result.EnabledAddresses))
		}
		if len(result.ManagerAddresses) != 1 {
			t.Errorf("Expected 1 manager address, got %d", len(result.ManagerAddresses))
		}
	})

	t.Run("do not add address if already enabled", func(t *testing.T) {
		enabledAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
		allowListConfig := allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{},
			ManagerAddresses: []common.Address{},
			EnabledAddresses: []common.Address{enabledAddress},
		}

		result := addAddressToAllowed(allowListConfig, enabledAddress.Hex())

		if len(result.EnabledAddresses) != 1 {
			t.Errorf("Expected 1 enabled address, got %d", len(result.EnabledAddresses))
		}
		if result.EnabledAddresses[0] != enabledAddress {
			t.Errorf("Expected address %s, got %s", enabledAddress.Hex(), result.EnabledAddresses[0].Hex())
		}
	})

	t.Run("add address to existing enabled addresses", func(t *testing.T) {
		existingAddress := common.HexToAddress("0x1111111111111111111111111111111111111111")
		newAddress := common.HexToAddress("0x2222222222222222222222222222222222222222")

		allowListConfig := allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{},
			ManagerAddresses: []common.Address{},
			EnabledAddresses: []common.Address{existingAddress},
		}

		result := addAddressToAllowed(allowListConfig, newAddress.Hex())

		if len(result.EnabledAddresses) != 2 {
			t.Errorf("Expected 2 enabled addresses, got %d", len(result.EnabledAddresses))
		}

		// Check that both addresses are present
		foundExisting := false
		foundNew := false
		for _, addr := range result.EnabledAddresses {
			if addr == existingAddress {
				foundExisting = true
			}
			if addr == newAddress {
				foundNew = true
			}
		}

		if !foundExisting {
			t.Error("Expected existing address to still be present")
		}
		if !foundNew {
			t.Error("Expected new address to be added")
		}
	})
}

func TestAddTeleporterAddressesToAllowLists(t *testing.T) {
	t.Run("add addresses to tx allowlist when precompile exists", func(t *testing.T) {
		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{},
		}

		// Add tx allowlist precompile
		txConfig := &txallowlist.Config{
			AllowListConfig: allowlist.AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{},
			},
		}
		config.GenesisPrecompiles[txallowlist.ConfigKey] = txConfig

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Check that tx allowlist was updated
		txResultConfig := result.GenesisPrecompiles[txallowlist.ConfigKey].(*txallowlist.Config)
		if len(txResultConfig.AllowListConfig.EnabledAddresses) != 3 {
			t.Errorf("Expected 3 enabled addresses in tx allowlist, got %d", len(txResultConfig.AllowListConfig.EnabledAddresses))
		}
	})

	t.Run("add addresses to deployer allowlist when precompile exists", func(t *testing.T) {
		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{},
		}

		// Add deployer allowlist precompile
		deployerConfig := &deployerallowlist.Config{
			AllowListConfig: allowlist.AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{},
			},
		}
		config.GenesisPrecompiles[deployerallowlist.ConfigKey] = deployerConfig

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Check that deployer allowlist was updated
		deployerResultConfig := result.GenesisPrecompiles[deployerallowlist.ConfigKey].(*deployerallowlist.Config)
		if len(deployerResultConfig.AllowListConfig.EnabledAddresses) != 2 {
			t.Errorf("Expected 2 enabled addresses in deployer allowlist, got %d", len(deployerResultConfig.AllowListConfig.EnabledAddresses))
		}
	})

	t.Run("handle missing precompiles gracefully", func(t *testing.T) {
		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{},
		}

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		// Should not panic when precompiles don't exist
		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Result should be the same as input when no precompiles exist
		if len(result.GenesisPrecompiles) != 0 {
			t.Errorf("Expected 0 precompiles, got %d", len(result.GenesisPrecompiles))
		}
	})

	t.Run("handle nil precompile config gracefully", func(t *testing.T) {
		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{
				txallowlist.ConfigKey: nil,
			},
		}

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		// Should not panic when precompile config is nil
		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Result should be the same as input
		if len(result.GenesisPrecompiles) != 1 {
			t.Errorf("Expected 1 precompile, got %d", len(result.GenesisPrecompiles))
		}
	})

	t.Run("add both tx and deployer allowlists", func(t *testing.T) {
		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{},
		}

		// Add both precompiles
		txConfig := &txallowlist.Config{
			AllowListConfig: allowlist.AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{},
			},
		}
		config.GenesisPrecompiles[txallowlist.ConfigKey] = txConfig

		deployerConfig := &deployerallowlist.Config{
			AllowListConfig: allowlist.AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{},
			},
		}
		config.GenesisPrecompiles[deployerallowlist.ConfigKey] = deployerConfig

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Check tx allowlist (should have 3 addresses)
		txResultConfig := result.GenesisPrecompiles[txallowlist.ConfigKey].(*txallowlist.Config)
		if len(txResultConfig.AllowListConfig.EnabledAddresses) != 3 {
			t.Errorf("Expected 3 enabled addresses in tx allowlist, got %d", len(txResultConfig.AllowListConfig.EnabledAddresses))
		}

		// Check deployer allowlist (should have 2 addresses)
		deployerResultConfig := result.GenesisPrecompiles[deployerallowlist.ConfigKey].(*deployerallowlist.Config)
		if len(deployerResultConfig.AllowListConfig.EnabledAddresses) != 2 {
			t.Errorf("Expected 2 enabled addresses in deployer allowlist, got %d", len(deployerResultConfig.AllowListConfig.EnabledAddresses))
		}
	})

	t.Run("preserve existing addresses in allowlists", func(t *testing.T) {
		existingAddress := common.HexToAddress("0x9999999999999999999999999999999999999999")

		config := params.ChainConfig{
			GenesisPrecompiles: params.Precompiles{},
		}

		// Add tx allowlist with existing address
		txConfig := &txallowlist.Config{
			AllowListConfig: allowlist.AllowListConfig{
				AdminAddresses:   []common.Address{},
				ManagerAddresses: []common.Address{},
				EnabledAddresses: []common.Address{existingAddress},
			},
		}
		config.GenesisPrecompiles[txallowlist.ConfigKey] = txConfig

		teleporterAddress := "0x1234567890123456789012345678901234567890"
		teleporterMessengerDeployerAddress := "0x2345678901234567890123456789012345678901"
		relayerAddress := "0x3456789012345678901234567890123456789012"

		result := AddTeleporterAddressesToAllowLists(
			config,
			teleporterAddress,
			teleporterMessengerDeployerAddress,
			relayerAddress,
		)

		// Check that existing address is preserved
		txResultConfig := result.GenesisPrecompiles[txallowlist.ConfigKey].(*txallowlist.Config)
		if len(txResultConfig.AllowListConfig.EnabledAddresses) != 4 {
			t.Errorf("Expected 4 enabled addresses in tx allowlist, got %d", len(txResultConfig.AllowListConfig.EnabledAddresses))
		}

		// Verify existing address is still there
		foundExisting := false
		for _, addr := range txResultConfig.AllowListConfig.EnabledAddresses {
			if addr == existingAddress {
				foundExisting = true
				break
			}
		}
		if !foundExisting {
			t.Error("Expected existing address to be preserved")
		}
	})
}
