// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"crypto/rand"
	"encoding/pem"
	"testing"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BLS staking tests for NewBlsSecretKeyBytes and ToNodeID functions

func TestNewBlsSecretKeyBytes(t *testing.T) {
	// Test successful key generation
	keyBytes, err := NewBlsSecretKeyBytes()
	require.NoError(t, err)
	require.NotNil(t, keyBytes)
	require.Greater(t, len(keyBytes), 0, "Generated key should not be empty")

	// Test that multiple calls generate different keys
	keyBytes2, err := NewBlsSecretKeyBytes()
	require.NoError(t, err)
	require.NotEqual(t, keyBytes, keyBytes2, "Multiple calls should generate different keys")

	// Test key format (BLS secret keys are typically 32 bytes)
	// Note: The exact length depends on the BLS implementation
	assert.True(t, len(keyBytes) > 0, "Key should have positive length")
	assert.True(t, len(keyBytes) <= 128, "Key should be reasonably sized") // Upper bound check
}

func TestToNodeID(t *testing.T) {
	// Create a mock certificate for testing
	// This is a simplified test - in practice, you'd need a real certificate
	mockCertData := []byte("mock certificate data")

	// Create a PEM block
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: mockCertData,
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(block)
	require.NotNil(t, certPEM, "PEM encoding should succeed")

	// Test with invalid certificate data
	t.Run("invalid certificate", func(t *testing.T) {
		invalidCert := []byte("invalid certificate data")
		nodeID, err := ToNodeID(invalidCert)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with empty certificate
	t.Run("empty certificate", func(t *testing.T) {
		emptyCert := []byte("")
		nodeID, err := ToNodeID(emptyCert)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with non-PEM data
	t.Run("non-PEM data", func(t *testing.T) {
		nonPEMData := []byte("this is not PEM data")
		nodeID, err := ToNodeID(nonPEMData)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with valid PEM but invalid certificate content
	t.Run("valid PEM invalid certificate", func(t *testing.T) {
		// Create PEM with invalid certificate data
		block := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte("invalid certificate bytes"),
		}
		invalidCertPEM := pem.EncodeToMemory(block)

		nodeID, err := ToNodeID(invalidCertPEM)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with valid PEM structure but wrong type
	t.Run("wrong PEM type", func(t *testing.T) {
		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: []byte("some data"),
		}
		wrongTypePEM := pem.EncodeToMemory(block)

		nodeID, err := ToNodeID(wrongTypePEM)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with nil input
	t.Run("nil input", func(t *testing.T) {
		nodeID, err := ToNodeID(nil)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with random data
	t.Run("random data", func(t *testing.T) {
		randomData := make([]byte, 100)
		_, err := rand.Read(randomData)
		require.NoError(t, err)

		nodeID, err := ToNodeID(randomData)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with malformed PEM
	t.Run("malformed PEM", func(t *testing.T) {
		malformedPEM := []byte("-----BEGIN CERTIFICATE-----\ninvalid base64\n-----END CERTIFICATE-----")
		nodeID, err := ToNodeID(malformedPEM)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with incomplete PEM
	t.Run("incomplete PEM", func(t *testing.T) {
		incompletePEM := []byte("-----BEGIN CERTIFICATE-----\n")
		nodeID, err := ToNodeID(incompletePEM)
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test with multiple PEM blocks
	t.Run("multiple PEM blocks", func(t *testing.T) {
		block1 := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte("first certificate"),
		}
		block2 := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte("second certificate"),
		}

		multiPEM := append(pem.EncodeToMemory(block1), pem.EncodeToMemory(block2)...)

		nodeID, err := ToNodeID(multiPEM)
		// This should fail because the function expects a single certificate
		assert.Error(t, err)
		assert.Equal(t, ids.EmptyNodeID, nodeID)
	})

	// Test error message content
	t.Run("error message validation", func(t *testing.T) {
		invalidCert := []byte("invalid")
		_, err := ToNodeID(invalidCert)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode certificate")
	})
}
