// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"math/big"
	"testing"

	"github.com/DioneProtocol/subnet-evm/commontype"
)

func TestConstants(t *testing.T) {
	t.Run("DefaultEvmAirdropAmount", func(t *testing.T) {
		expected := "1000000000000000000000000"
		if DefaultEvmAirdropAmount != expected {
			t.Errorf("Expected DefaultEvmAirdropAmount to be %s, got %s", expected, DefaultEvmAirdropAmount)
		}
	})
}

func TestDifficulty(t *testing.T) {
	t.Run("Difficulty is zero", func(t *testing.T) {
		expected := big.NewInt(0)
		if Difficulty.Cmp(expected) != 0 {
			t.Errorf("Expected Difficulty to be %s, got %s", expected.String(), Difficulty.String())
		}
	})

	t.Run("Difficulty is not nil", func(t *testing.T) {
		if Difficulty == nil {
			t.Error("Expected Difficulty to not be nil")
		}
	})
}

func TestStarterFeeConfig(t *testing.T) {
	t.Run("GasLimit", func(t *testing.T) {
		expected := big.NewInt(8_000_000)
		if StarterFeeConfig.GasLimit.Cmp(expected) != 0 {
			t.Errorf("Expected GasLimit to be %s, got %s", expected.String(), StarterFeeConfig.GasLimit.String())
		}
	})

	t.Run("MinBaseFee", func(t *testing.T) {
		expected := big.NewInt(25_000_000_000)
		if StarterFeeConfig.MinBaseFee.Cmp(expected) != 0 {
			t.Errorf("Expected MinBaseFee to be %s, got %s", expected.String(), StarterFeeConfig.MinBaseFee.String())
		}
	})

	t.Run("TargetGas", func(t *testing.T) {
		expected := big.NewInt(15_000_000)
		if StarterFeeConfig.TargetGas.Cmp(expected) != 0 {
			t.Errorf("Expected TargetGas to be %s, got %s", expected.String(), StarterFeeConfig.TargetGas.String())
		}
	})

	t.Run("BaseFeeChangeDenominator", func(t *testing.T) {
		expected := big.NewInt(36)
		if StarterFeeConfig.BaseFeeChangeDenominator.Cmp(expected) != 0 {
			t.Errorf("Expected BaseFeeChangeDenominator to be %s, got %s", expected.String(), StarterFeeConfig.BaseFeeChangeDenominator.String())
		}
	})

	t.Run("MinBlockGasCost", func(t *testing.T) {
		expected := big.NewInt(0)
		if StarterFeeConfig.MinBlockGasCost.Cmp(expected) != 0 {
			t.Errorf("Expected MinBlockGasCost to be %s, got %s", expected.String(), StarterFeeConfig.MinBlockGasCost.String())
		}
	})

	t.Run("MaxBlockGasCost", func(t *testing.T) {
		expected := big.NewInt(1_000_000)
		if StarterFeeConfig.MaxBlockGasCost.Cmp(expected) != 0 {
			t.Errorf("Expected MaxBlockGasCost to be %s, got %s", expected.String(), StarterFeeConfig.MaxBlockGasCost.String())
		}
	})

	t.Run("TargetBlockRate", func(t *testing.T) {
		expected := uint64(2)
		if StarterFeeConfig.TargetBlockRate != expected {
			t.Errorf("Expected TargetBlockRate to be %d, got %d", expected, StarterFeeConfig.TargetBlockRate)
		}
	})

	t.Run("BlockGasCostStep", func(t *testing.T) {
		expected := big.NewInt(200_000)
		if StarterFeeConfig.BlockGasCostStep.Cmp(expected) != 0 {
			t.Errorf("Expected BlockGasCostStep to be %s, got %s", expected.String(), StarterFeeConfig.BlockGasCostStep.String())
		}
	})
}

func TestStarterFeeConfigType(t *testing.T) {
	t.Run("StarterFeeConfig is of correct type", func(t *testing.T) {
		var _ commontype.FeeConfig = StarterFeeConfig
		// This test will compile only if StarterFeeConfig is of type commontype.FeeConfig
	})
}

func TestStarterFeeConfigValuesAreValid(t *testing.T) {
	t.Run("GasLimit is positive", func(t *testing.T) {
		if StarterFeeConfig.GasLimit.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected GasLimit to be positive")
		}
	})

	t.Run("MinBaseFee is positive", func(t *testing.T) {
		if StarterFeeConfig.MinBaseFee.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected MinBaseFee to be positive")
		}
	})

	t.Run("TargetGas is positive", func(t *testing.T) {
		if StarterFeeConfig.TargetGas.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected TargetGas to be positive")
		}
	})

	t.Run("BaseFeeChangeDenominator is positive", func(t *testing.T) {
		if StarterFeeConfig.BaseFeeChangeDenominator.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected BaseFeeChangeDenominator to be positive")
		}
	})

	t.Run("MinBlockGasCost is non-negative", func(t *testing.T) {
		if StarterFeeConfig.MinBlockGasCost.Cmp(big.NewInt(0)) < 0 {
			t.Error("Expected MinBlockGasCost to be non-negative")
		}
	})

	t.Run("MaxBlockGasCost is positive", func(t *testing.T) {
		if StarterFeeConfig.MaxBlockGasCost.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected MaxBlockGasCost to be positive")
		}
	})

	t.Run("TargetBlockRate is positive", func(t *testing.T) {
		if StarterFeeConfig.TargetBlockRate == 0 {
			t.Error("Expected TargetBlockRate to be positive")
		}
	})

	t.Run("BlockGasCostStep is positive", func(t *testing.T) {
		if StarterFeeConfig.BlockGasCostStep.Cmp(big.NewInt(0)) <= 0 {
			t.Error("Expected BlockGasCostStep to be positive")
		}
	})
}

func TestStarterFeeConfigRelationships(t *testing.T) {
	t.Run("MaxBlockGasCost is greater than MinBlockGasCost", func(t *testing.T) {
		if StarterFeeConfig.MaxBlockGasCost.Cmp(StarterFeeConfig.MinBlockGasCost) <= 0 {
			t.Error("Expected MaxBlockGasCost to be greater than MinBlockGasCost")
		}
	})

	t.Run("GasLimit is less than TargetGas", func(t *testing.T) {
		if StarterFeeConfig.GasLimit.Cmp(StarterFeeConfig.TargetGas) >= 0 {
			t.Error("Expected GasLimit to be less than TargetGas")
		}
	})
}
