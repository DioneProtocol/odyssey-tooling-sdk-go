// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"testing"
)

func TestVMConfig(t *testing.T) {
	t.Run("VMConfig struct fields", func(t *testing.T) {
		config := &VMConfig{
			EnableWarp: true,
		}

		if !config.EnableWarp {
			t.Error("Expected EnableWarp to be true")
		}

		config.EnableWarp = false
		if config.EnableWarp {
			t.Error("Expected EnableWarp to be false")
		}
	})
}

func TestDefaultVMConfig(t *testing.T) {
	config := DefaultVMConfig()

	if config == nil {
		t.Fatal("Expected DefaultVMConfig to return non-nil config")
	}

	if config.EnableWarp {
		t.Error("Expected default config to have EnableWarp disabled")
	}
}

func TestEnableWarpConfig(t *testing.T) {
	config := EnableWarpConfig()

	if config == nil {
		t.Fatal("Expected EnableWarpConfig to return non-nil config")
	}

	if !config.EnableWarp {
		t.Error("Expected EnableWarpConfig to have EnableWarp enabled")
	}
}

func TestDisableWarpConfig(t *testing.T) {
	config := DisableWarpConfig()

	if config == nil {
		t.Fatal("Expected DisableWarpConfig to return non-nil config")
	}

	if config.EnableWarp {
		t.Error("Expected DisableWarpConfig to have EnableWarp disabled")
	}
}

func TestConfigFunctionsReturnDifferentInstances(t *testing.T) {
	defaultConfig := DefaultVMConfig()
	enableConfig := EnableWarpConfig()
	disableConfig := DisableWarpConfig()

	// Ensure they are different instances
	if defaultConfig == enableConfig {
		t.Error("Expected DefaultVMConfig and EnableWarpConfig to return different instances")
	}

	if defaultConfig == disableConfig {
		t.Error("Expected DefaultVMConfig and DisableWarpConfig to return different instances")
	}

	if enableConfig == disableConfig {
		t.Error("Expected EnableWarpConfig and DisableWarpConfig to return different instances")
	}
}

func TestConfigModificationDoesNotAffectOthers(t *testing.T) {
	defaultConfig := DefaultVMConfig()
	enableConfig := EnableWarpConfig()

	// Modify one config
	defaultConfig.EnableWarp = true

	// Check that the other config is not affected
	if !enableConfig.EnableWarp {
		t.Error("Expected EnableWarpConfig to still have EnableWarp enabled")
	}

	// Verify the modification worked
	if !defaultConfig.EnableWarp {
		t.Error("Expected modified default config to have EnableWarp enabled")
	}
}
