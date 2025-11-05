// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"testing"

	"github.com/DioneProtocol/subnet-evm/precompile/allowlist"
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
