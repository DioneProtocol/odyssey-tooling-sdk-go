// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

// VMConfig holds configuration for VM features
type VMConfig struct {
	// EnableWarp controls whether warp functionality is enabled
	EnableWarp bool
}

// DefaultVMConfig returns the default VM configuration
func DefaultVMConfig() *VMConfig {
	return &VMConfig{
		EnableWarp: false, // Disabled by default
	}
}

// EnableWarpConfig returns a config with warp enabled
func EnableWarpConfig() *VMConfig {
	return &VMConfig{
		EnableWarp: true,
	}
}

// DisableWarpConfig returns a config with warp disabled
func DisableWarpConfig() *VMConfig {
	return &VMConfig{
		EnableWarp: false,
	}
}
