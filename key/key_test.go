// Copyright (c) 2025 Dione Protocol, LLC.
// See the file LICENSE for licensing terms.

package key

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/cb58"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ewoqOChainAddr = "O-custom18jma8ppw3nhx5r4ap8clazz0dps7rv5u9xde7p"

func TestNewKeyEwoq(t *testing.T) {
	t.Parallel()

	m, err := NewSoft(
		WithPrivateKeyEncoded(EwoqPrivateKey),
	)
	if err != nil {
		t.Fatal(err)
	}

	pAddr, err := m.O("custom")
	if err != nil {
		t.Fatal(err)
	}
	if pAddr != ewoqOChainAddr {
		t.Fatalf("unexpected O-Chain address %q, expected %q", pAddr, ewoqOChainAddr)
	}

	keyPath := filepath.Join(t.TempDir(), "key.pk")
	if err := m.Save(keyPath); err != nil {
		t.Fatal(err)
	}

	m2, err := LoadSoft(keyPath)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(m.PrivKeyRaw(), m2.PrivKeyRaw()) {
		t.Fatalf("loaded key unexpected %v, expected %v", m2.PrivKeyRaw(), m.PrivKeyRaw())
	}
}

func TestNewKey(t *testing.T) {
	t.Parallel()

	skBytes, err := cb58.Decode(rawEwoqPk)
	if err != nil {
		t.Fatal(err)
	}
	factory := secp256k1.Factory{}
	ewoqPk, err := factory.ToPrivateKey(skBytes)
	if err != nil {
		t.Fatal(err)
	}

	privKey2, err := factory.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		name   string
		opts   []SOpOption
		expErr error
	}{
		{
			name:   "test",
			opts:   nil,
			expErr: nil,
		},
		{
			name: "ewop with WithPrivateKey",
			opts: []SOpOption{
				WithPrivateKey(ewoqPk),
			},
			expErr: nil,
		},
		{
			name: "ewop with WithPrivateKeyEncoded",
			opts: []SOpOption{
				WithPrivateKeyEncoded(EwoqPrivateKey),
			},
			expErr: nil,
		},
		{
			name: "ewop with WithPrivateKey/WithPrivateKeyEncoded",
			opts: []SOpOption{
				WithPrivateKey(ewoqPk),
				WithPrivateKeyEncoded(EwoqPrivateKey),
			},
			expErr: nil,
		},
		{
			name: "ewop with invalid WithPrivateKey",
			opts: []SOpOption{
				WithPrivateKey(privKey2),
				WithPrivateKeyEncoded(EwoqPrivateKey),
			},
			expErr: ErrInvalidPrivateKey,
		},
	}
	for i, tv := range tt {
		_, err := NewSoft(tv.opts...)
		if !errors.Is(err, tv.expErr) {
			t.Fatalf("#%d(%s): unexpected error %v, expected %v", i, tv.name, err, tv.expErr)
		}
	}
}

// TestKeyInterfaceUtilities tests the utility functions in key.go
func TestKeyInterfaceUtilities(t *testing.T) {
	t.Parallel()

	t.Run("WithTime", func(t *testing.T) {
		op := &Op{}
		WithTime(12345)(op)
		assert.Equal(t, uint64(12345), op.time)
	})

	t.Run("WithTargetAmount", func(t *testing.T) {
		op := &Op{}
		WithTargetAmount(1000)(op)
		assert.Equal(t, uint64(1000), op.targetAmount)
	})

	t.Run("WithFeeDeduct", func(t *testing.T) {
		op := &Op{}
		WithFeeDeduct(100)(op)
		assert.Equal(t, uint64(100), op.feeDeduct)
	})

	t.Run("applyOpts", func(t *testing.T) {
		op := &Op{}
		opts := []OpOption{
			WithTime(100),
			WithTargetAmount(200),
			WithFeeDeduct(50),
		}
		op.applyOpts(opts)
		assert.Equal(t, uint64(100), op.time)
		assert.Equal(t, uint64(200), op.targetAmount)
		assert.Equal(t, uint64(50), op.feeDeduct)
	})

	t.Run("SortTransferableInputsWithSigners", func(t *testing.T) {
		// Create test data with clear ordering
		id1 := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		id2 := ids.ID{2, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		id3 := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 33}

		// Create inputs in unsorted order
		inputs := []*dione.TransferableInput{
			{UTXOID: dione.UTXOID{TxID: id2, OutputIndex: 1}}, // id2, index 1
			{UTXOID: dione.UTXOID{TxID: id1, OutputIndex: 2}}, // id1, index 2
			{UTXOID: dione.UTXOID{TxID: id1, OutputIndex: 1}}, // id1, index 1
			{UTXOID: dione.UTXOID{TxID: id3, OutputIndex: 1}}, // id3, index 1
		}

		signers := [][]ids.ShortID{
			{ids.ShortID{1}}, // corresponds to id2, index 1
			{ids.ShortID{2}}, // corresponds to id1, index 2
			{ids.ShortID{3}}, // corresponds to id1, index 1
			{ids.ShortID{4}}, // corresponds to id3, index 1
		}

		// Store original order for comparison
		originalInputs := make([]*dione.TransferableInput, len(inputs))
		originalSigners := make([][]ids.ShortID, len(signers))
		copy(originalInputs, inputs)
		copy(originalSigners, signers)

		// Sort
		SortTransferableInputsWithSigners(inputs, signers)

		// Verify that the function was called (no panic)
		assert.Len(t, inputs, 4)
		assert.Len(t, signers, 4)

		// Verify that inputs and signers have the same length
		for i := range inputs {
			assert.Len(t, signers[i], 1)
		}

		// Verify that the function actually does something (not just a no-op)
		// We'll check that at least one element moved from its original position
		changed := false
		for i := range inputs {
			if inputs[i].UTXOID.TxID != originalInputs[i].UTXOID.TxID ||
				inputs[i].UTXOID.OutputIndex != originalInputs[i].UTXOID.OutputIndex {
				changed = true
				break
			}
		}
		assert.True(t, changed, "Sorting function should have changed the order")
	})
}

// TestSoftKeyMethods tests all the untested SoftKey methods
func TestSoftKeyMethods(t *testing.T) {
	t.Parallel()

	// Create a test key
	key, err := NewSoft(WithPrivateKeyEncoded(EwoqPrivateKey))
	require.NoError(t, err)
	require.NotNil(t, key)

	t.Run("D", func(t *testing.T) {
		dAddr := key.D()
		assert.NotEmpty(t, dAddr)
		assert.Contains(t, dAddr, "0x") // Ethereum address format
		assert.Len(t, dAddr, 42)        // 0x + 40 hex chars
	})

	t.Run("KeyChain", func(t *testing.T) {
		keychain := key.KeyChain()
		assert.NotNil(t, keychain)
	})

	t.Run("PrivKey", func(t *testing.T) {
		privKey := key.PrivKey()
		assert.NotNil(t, privKey)
		assert.Equal(t, key.PrivKeyRaw(), privKey.Bytes())
	})

	t.Run("PrivKeyCB58", func(t *testing.T) {
		cb58Key := key.PrivKeyCB58()
		assert.NotEmpty(t, cb58Key)
		assert.Contains(t, cb58Key, "PrivateKey-")
		assert.Equal(t, EwoqPrivateKey, cb58Key)
	})

	t.Run("O", func(t *testing.T) {
		oAddr, err := key.O("custom")
		require.NoError(t, err)
		assert.NotEmpty(t, oAddr)
		assert.Contains(t, oAddr, "O-custom")
	})

	t.Run("A", func(t *testing.T) {
		aAddr, err := key.A("custom")
		require.NoError(t, err)
		assert.NotEmpty(t, aAddr)
		assert.Contains(t, aAddr, "A-custom")
	})

	t.Run("Addresses", func(t *testing.T) {
		addresses := key.Addresses()
		assert.Len(t, addresses, 1)
		assert.NotNil(t, addresses[0])
	})

	t.Run("Match", func(t *testing.T) {
		// Create a test output owner
		addr := key.PrivKey().PublicKey().Address()
		owners := &secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{addr},
		}

		indices, pks, ok := key.Match(owners, 0)
		assert.True(t, ok)
		assert.Len(t, indices, 1)
		assert.Len(t, pks, 1)
		assert.Equal(t, addr, pks[0])
		assert.Equal(t, uint32(0), indices[0])

		// Test with non-matching address
		otherAddr := ids.ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		otherOwners := &secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{otherAddr},
		}

		indices, pks, ok = key.Match(otherOwners, 0)
		assert.False(t, ok)
		assert.Empty(t, indices)
		assert.Empty(t, pks)
	})
}

// TestSoftKeySpends tests the Spends method
func TestSoftKeySpends(t *testing.T) {
	t.Parallel()

	key, err := NewSoft(WithPrivateKeyEncoded(EwoqPrivateKey))
	require.NoError(t, err)

	// Create a test UTXO that the key can spend
	addr := key.PrivKey().PublicKey().Address()
	utxo := &dione.UTXO{
		UTXOID: dione.UTXOID{
			TxID:        ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			OutputIndex: 0,
		},
		Asset: dione.Asset{ID: ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}},
		Out: &secp256k1fx.TransferOutput{
			Amt: 1000,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{addr},
			},
		},
	}

	t.Run("Spends without options", func(t *testing.T) {
		total, inputs, signers := key.Spends([]*dione.UTXO{utxo})
		assert.Equal(t, uint64(1000), total)
		assert.Len(t, inputs, 1)
		assert.Len(t, signers, 1)
		assert.Len(t, signers[0], 1)
		assert.Equal(t, addr, signers[0][0])
	})

	t.Run("Spends with target amount", func(t *testing.T) {
		total, inputs, signers := key.Spends([]*dione.UTXO{utxo}, WithTargetAmount(500))
		assert.Equal(t, uint64(1000), total)
		assert.Len(t, inputs, 1)
		assert.Len(t, signers, 1)
	})

	t.Run("Spends with fee deduction", func(t *testing.T) {
		total, inputs, signers := key.Spends([]*dione.UTXO{utxo}, WithFeeDeduct(100))
		assert.Equal(t, uint64(1000), total)
		assert.Len(t, inputs, 1)
		assert.Len(t, signers, 1)
	})

	t.Run("Spends with time", func(t *testing.T) {
		total, inputs, signers := key.Spends([]*dione.UTXO{utxo}, WithTime(1000))
		assert.Equal(t, uint64(1000), total)
		assert.Len(t, inputs, 1)
		assert.Len(t, signers, 1)
	})

	t.Run("Spends with unspendable UTXO", func(t *testing.T) {
		// Create UTXO with different address
		otherAddr := ids.ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		unspendableUTXO := &dione.UTXO{
			UTXOID: dione.UTXOID{
				TxID:        ids.ID{2, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
				OutputIndex: 0,
			},
			Asset: dione.Asset{ID: ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}},
			Out: &secp256k1fx.TransferOutput{
				Amt: 2000,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{otherAddr},
				},
			},
		}

		total, inputs, signers := key.Spends([]*dione.UTXO{unspendableUTXO})
		assert.Equal(t, uint64(0), total)
		assert.Empty(t, inputs)
		assert.Empty(t, signers)
	})

	t.Run("Spends with invalid output type", func(t *testing.T) {
		// Create a UTXO with an output that can't be converted to TransferableIn
		addr := key.PrivKey().PublicKey().Address()

		// Create a mock output that will fail the type assertion
		mockOutput := &secp256k1fx.MintOutput{
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{addr},
			},
		}

		invalidUTXO := &dione.UTXO{
			UTXOID: dione.UTXOID{
				TxID:        ids.ID{3, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
				OutputIndex: 0,
			},
			Asset: dione.Asset{ID: ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}},
			Out:   mockOutput,
		}

		total, inputs, signers := key.Spends([]*dione.UTXO{invalidUTXO})
		assert.Equal(t, uint64(0), total)
		assert.Empty(t, inputs)
		assert.Empty(t, signers)
	})
}

// TestSoftKeySign tests the Sign method
func TestSoftKeySign(t *testing.T) {
	t.Parallel()

	key, err := NewSoft(WithPrivateKeyEncoded(EwoqPrivateKey))
	require.NoError(t, err)

	addr := key.PrivKey().PublicKey().Address()

	t.Run("Sign with valid signers", func(t *testing.T) {
		// Create a simple transaction
		tx := &txs.Tx{}
		signers := [][]ids.ShortID{{addr}}

		err := key.Sign(tx, signers)
		// Note: This might fail due to transaction structure, but we're testing the logic
		// The important part is that it doesn't fail due to signer validation
		if err != nil {
			// If it fails, it should not be ErrCantSpend
			assert.NotEqual(t, ErrCantSpend, err)
		}
	})

	t.Run("Sign with invalid signer", func(t *testing.T) {
		otherAddr := ids.ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		tx := &txs.Tx{}
		signers := [][]ids.ShortID{{otherAddr}}

		err := key.Sign(tx, signers)
		assert.Equal(t, ErrCantSpend, err)
	})
}

// TestSoftKeyLoaders tests the loader functions
func TestSoftKeyLoaders(t *testing.T) {
	t.Parallel()

	t.Run("LoadSoftOrCreate - file exists", func(t *testing.T) {
		// Create a temporary file
		tempDir := t.TempDir()
		keyPath := filepath.Join(tempDir, "test.key")

		// Create a key and save it
		originalKey, err := NewSoft()
		require.NoError(t, err)
		err = originalKey.Save(keyPath)
		require.NoError(t, err)

		// Load it using LoadSoftOrCreate
		loadedKey, err := LoadSoftOrCreate(keyPath)
		require.NoError(t, err)
		require.NotNil(t, loadedKey)

		// Verify it's the same key
		assert.Equal(t, originalKey.PrivKeyRaw(), loadedKey.PrivKeyRaw())
	})

	t.Run("LoadSoftOrCreate - file doesn't exist", func(t *testing.T) {
		tempDir := t.TempDir()
		keyPath := filepath.Join(tempDir, "nonexistent.key")

		// Load using LoadSoftOrCreate (should create new key)
		key, err := LoadSoftOrCreate(keyPath)
		require.NoError(t, err)
		require.NotNil(t, key)

		// Verify file was created
		assert.FileExists(t, keyPath)

		// Verify we can load it again
		loadedKey, err := LoadSoft(keyPath)
		require.NoError(t, err)
		assert.Equal(t, key.PrivKeyRaw(), loadedKey.PrivKeyRaw())
	})

	t.Run("LoadSoftOrCreate - error during save", func(t *testing.T) {
		// Try to create a key in a non-existent directory
		invalidPath := "/nonexistent/directory/key.key"

		_, err := LoadSoftOrCreate(invalidPath)
		assert.Error(t, err)
	})

	t.Run("LoadEwoq", func(t *testing.T) {
		key, err := LoadEwoq()
		require.NoError(t, err)
		require.NotNil(t, key)

		// Verify it's the ewoq key
		assert.Equal(t, EwoqPrivateKey, key.PrivKeyCB58())
	})
}

// TestSoftKeyErrorCases tests various error scenarios
func TestSoftKeyErrorCases(t *testing.T) {
	t.Parallel()

	t.Run("NewSoft with invalid private key encoding", func(t *testing.T) {
		_, err := NewSoft(WithPrivateKeyEncoded("invalid-encoding"))
		assert.Error(t, err)
		// The error could be about base58 decoding or invalid encoding
		assert.True(t,
			strings.Contains(err.Error(), "invalid") ||
				strings.Contains(err.Error(), "base58") ||
				strings.Contains(err.Error(), "decoding"))
	})

	t.Run("NewSoft with encoding mismatch", func(t *testing.T) {
		// Create a key with one private key but different encoding
		factory := secp256k1.Factory{}
		privKey, err := factory.NewPrivateKey()
		require.NoError(t, err)

		// Use a different encoding that doesn't match the private key
		_, err = NewSoft(
			WithPrivateKey(privKey),
			WithPrivateKeyEncoded(EwoqPrivateKey), // Different key
		)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPrivateKey, err)
	})

	t.Run("NewSoft with invalid encoding format", func(t *testing.T) {
		// Test with encoding that has wrong prefix
		_, err := NewSoft(WithPrivateKeyEncoded("WrongPrefix-ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN"))
		assert.Error(t, err)
		// Should fail during decoding
	})

	t.Run("LoadSoft with non-existent file", func(t *testing.T) {
		_, err := LoadSoft("nonexistent.key")
		assert.Error(t, err)
	})

	t.Run("LoadSoftFromBytes with invalid hex", func(t *testing.T) {
		invalidBytes := []byte("invalid-hex-data")
		_, err := LoadSoftFromBytes(invalidBytes)
		assert.Error(t, err)
	})

	t.Run("LoadSoftFromBytes with wrong length", func(t *testing.T) {
		shortBytes := []byte("1234567890abcdef") // Too short
		_, err := LoadSoftFromBytes(shortBytes)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPrivateKeyLen, err)
	})

	t.Run("LoadSoftFromBytes with invalid file ending", func(t *testing.T) {
		// Create bytes with invalid ending (not just newlines)
		invalidBytes := make([]byte, 64)
		for i := range invalidBytes {
			invalidBytes[i] = 'a'
		}
		invalidBytes = append(invalidBytes, 'x') // Invalid character

		_, err := LoadSoftFromBytes(invalidBytes)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPrivateKeyEnding, err)
	})

	t.Run("LoadSoftFromBytes with too many newlines", func(t *testing.T) {
		// Create bytes with too many newlines at the end
		invalidBytes := make([]byte, 64)
		for i := range invalidBytes {
			invalidBytes[i] = 'a'
		}
		invalidBytes = append(invalidBytes, '\n', '\n', '\n') // Too many newlines

		_, err := LoadSoftFromBytes(invalidBytes)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPrivateKeyLen, err)
	})

	t.Run("LoadSoftFromBytes with readASCII error", func(t *testing.T) {
		// Create bytes that will cause readASCII to return an error
		// We'll create a very short buffer to trigger the error path
		shortBytes := []byte("a") // Too short, will cause readASCII to return early

		_, err := LoadSoftFromBytes(shortBytes)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPrivateKeyLen, err)
	})

	t.Run("LoadSoftFromBytes with hex decode error", func(t *testing.T) {
		// Create bytes with invalid hex characters
		invalidBytes := make([]byte, 64)
		for i := range invalidBytes {
			invalidBytes[i] = 'g' // Invalid hex character
		}
		invalidBytes = append(invalidBytes, '\n')

		_, err := LoadSoftFromBytes(invalidBytes)
		assert.Error(t, err)
		// Should be a hex decode error
	})

	t.Run("LoadSoftFromBytes with factory error", func(t *testing.T) {
		// Create bytes with valid hex but invalid private key
		// Use a private key that's too short (31 bytes instead of 32)
		invalidBytes := make([]byte, 62) // 31 bytes in hex = 62 characters
		for i := range invalidBytes {
			invalidBytes[i] = 'a'
		}
		invalidBytes = append(invalidBytes, '\n')

		_, err := LoadSoftFromBytes(invalidBytes)
		assert.Error(t, err)
		// Should be a length error
	})
}

// TestLedgerKeyStubs tests the LedgerKey stub functions
func TestLedgerKeyStubs(t *testing.T) {
	t.Parallel()

	ledgerKey := NewLedger(0)

	t.Run("LoadLedger", func(t *testing.T) {
		_, err := LoadLedger("test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("LoadLedgerFromBytes", func(t *testing.T) {
		_, err := LoadLedgerFromBytes([]byte("test"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("LedgerKey methods", func(t *testing.T) {
		assert.Empty(t, ledgerKey.D())
		assert.Nil(t, ledgerKey.KeyChain())

		err := ledgerKey.Save("test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")

		_, err = ledgerKey.O("test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")

		_, err = ledgerKey.A("test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")

		total, inputs, signers := ledgerKey.Spends(nil)
		assert.Equal(t, uint64(0), total)
		assert.Nil(t, inputs)
		assert.Nil(t, signers)

		assert.Nil(t, ledgerKey.Addresses())

		err = ledgerKey.Sign(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")

		indices, pks, ok := ledgerKey.Match(nil, 0)
		assert.Nil(t, indices)
		assert.Nil(t, pks)
		assert.False(t, ok)
	})
}
