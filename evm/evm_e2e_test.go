// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Env variables used by these tests:
// - ODX_RPC_URL: RPC endpoint (e.g., http://127.0.0.1:9650/ext/bc/C/rpc). Defaults to https://testnode.dioneprotocol.com
// - ODX_PRIV_KEY: Sender private key (hex, without 0x) with funds
// - ODX_TO_ADDR: Recipient address (0x...)
// - ODX_AMOUNT_WEI (optional): Transfer amount in wei (default 1)
// - ODX_ENABLE_TRACE (optional): If set to 1, enables trace test (requires debug api)
// - ODX_TX_HASH (optional): Tx hash to trace when ODX_ENABLE_TRACE=1

func TestE2ENativeTransfer_RealNode(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	privKey := os.Getenv("ODX_PRIV_KEY")
	toAddr := os.Getenv("ODX_TO_ADDR")
	if privKey == "" || toAddr == "" {
		t.Skip("ODX_PRIV_KEY and ODX_TO_ADDR must be set to run E2E transfer test")
	}
	amount := big.NewInt(1)
	if s := os.Getenv("ODX_AMOUNT_WEI"); s != "" {
		if v, ok := new(big.Int).SetString(s, 10); ok {
			amount = v
		}
	}

	client, err := GetClient(rpcURL)
	require.NoError(t, err)
	defer client.Close()

	_, err = GetChainID(client)
	require.NoError(t, err)

	before, err := GetAddressBalance(client, toAddr)
	require.NoError(t, err)

	require.NoError(t, Transfer(client, privKey, toAddr, amount))

	after, err := GetAddressBalance(client, toAddr)
	require.NoError(t, err)
	expected := new(big.Int).Add(before, amount)
	require.Equal(t, 0, expected.Cmp(after))
}

func TestE2EDebugTrace_RealNode(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	if os.Getenv("ODX_ENABLE_TRACE") != "1" {
		t.Skip("ODX_ENABLE_TRACE!=1; skipping")
	}
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	txHash := os.Getenv("ODX_TX_HASH")
	if txHash == "" {
		t.Skip("ODX_TX_HASH must be set to run trace test")
	}
	// txHash may be provided with or without 0x prefix; ensure hex format
	if !common.IsHexAddress(txHash) && len(txHash) != 66 {
		t.Skip("ODX_TX_HASH must be a 0x-prefixed 32-byte hash")
	}
	trace, err := GetTrace(rpcURL, txHash)
	require.NoError(t, err)
	require.NotNil(t, trace)
	require.NotEmpty(t, trace)
}
