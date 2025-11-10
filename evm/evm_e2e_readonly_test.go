// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Env variables:
// - ODX_RPC_URL: RPC endpoint (e.g., http://127.0.0.1:9650/ext/bc/C/rpc). Defaults to https://testnode.dioneprotocol.com/ext/bc/D/rpc
// - ODX_CONTRACT_ADDR (optional): Deployed contract address for bytecode/deployed checks
// - ODX_QUERY_ADDR (optional): Address to query balance/nonce
// - ODX_ERC20_ADDR (optional): ERC20 contract address for balanceOf
// - ODX_NFT_ADDR (optional): ERC721 contract address for metadata/supports and balanceOf
// Note: ERC721 tests use hardcoded addresses (0xA2D8ccB2415a5E88f5dBDd9BdcC827bac1A2224A and 0xEa4A0aA8aD7418f22373567BcA5728cB030f27Af)

func TestE2E_ReadOnly_NodeParams(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	client, err := GetClient(rpcURL)
	require.NoError(t, err)
	defer client.Close()

	chainID, err := GetChainID(client)
	require.NoError(t, err)
	require.NotNil(t, chainID)
	require.True(t, chainID.Sign() >= 0)

	baseFee, err := EstimateBaseFee(client)
	require.NoError(t, err)
	require.NotNil(t, baseFee)
	require.True(t, baseFee.Sign() >= 0)

	tip, err := SuggestGasTipCap(client)
	require.NoError(t, err)
	require.NotNil(t, tip)
	require.True(t, tip.Sign() >= 0)
}

func TestE2E_ReadOnly_GetAddressBalance(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	addr := os.Getenv("ODX_QUERY_ADDR")
	if addr == "" {
		t.Skip("ODX_QUERY_ADDR must be set")
	}
	client, err := GetClient(rpcURL)
	require.NoError(t, err)
	defer client.Close()

	bal, err := GetAddressBalance(client, addr)
	require.NoError(t, err)
	require.NotNil(t, bal)
}

func TestE2E_ReadOnly_ContractBytecodeAndDeployed(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	contract := os.Getenv("ODX_CONTRACT_ADDR")
	if contract == "" {
		t.Skip("ODX_CONTRACT_ADDR must be set")
	}
	client, err := GetClient(rpcURL)
	require.NoError(t, err)
	defer client.Close()

	code, err := GetContractBytecode(client, contract)
	require.NoError(t, err)
	require.NotNil(t, code)
	require.Greater(t, len(code), 0)

	deployed, err := ContractAlreadyDeployed(client, contract)
	require.NoError(t, err)
	require.True(t, deployed)
}

func TestE2E_ReadOnly_CallToMethod_BalanceOf(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	erc20 := os.Getenv("ODX_ERC20_ADDR")
	queryAddr := os.Getenv("ODX_QUERY_ADDR")
	if erc20 == "" || queryAddr == "" {
		t.Skip("ODX_ERC20_ADDR and ODX_QUERY_ADDR must be set")
	}
	contractAddr := common.HexToAddress(erc20)
	out, err := CallToMethod(
		rpcURL,
		contractAddr,
		"balanceOf(address)->(uint256)",
		common.HexToAddress(queryAddr),
	)
	require.NoError(t, err)
	require.Len(t, out, 1)
	bal, ok := out[0].(*big.Int)
	require.True(t, ok)
	require.NotNil(t, bal)
}

func TestE2E_ReadOnly_ERC721_MetadataAndSupports(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	// Hardcoded NFT address: 0xA2D8ccB2415a5E88f5dBDd9BdcC827bac1A2224A
	addr := common.HexToAddress("0xA2D8ccB2415a5E88f5dBDd9BdcC827bac1A2224A")
	// name()
	outName, err := CallToMethod(rpcURL, addr, "name()->(string)")
	require.NoError(t, err)
	require.Len(t, outName, 1)
	_, ok := outName[0].(string)
	require.True(t, ok)
	// symbol()
	outSym, err := CallToMethod(rpcURL, addr, "symbol()->(string)")
	require.NoError(t, err)
	require.Len(t, outSym, 1)
	_, ok = outSym[0].(string)
	require.True(t, ok)
	// supportsInterface(bytes4)
	var iid = [4]byte{0x80, 0xac, 0x58, 0xcd} // ERC-721 interface id
	outSupp, err := CallToMethod(
		rpcURL,
		addr,
		"supportsInterface(bytes4)->(bool)",
		iid,
	)
	require.NoError(t, err)
	require.Len(t, outSupp, 1)
	_, ok = outSupp[0].(bool)
	require.True(t, ok)
}

func TestE2E_ReadOnly_ERC721_BalanceOf(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	rpcURL := os.Getenv("ODX_RPC_URL")
	if rpcURL == "" {
		rpcURL = testnetDefaultRPC
	}
	// Hardcoded NFT address: 0xA2D8ccB2415a5E88f5dBDd9BdcC827bac1A2224A
	// Hardcoded owner address: 0xEa4A0aA8aD7418f22373567BcA5728cB030f27Af
	addr := common.HexToAddress("0xA2D8ccB2415a5E88f5dBDd9BdcC827bac1A2224A")
	owner := "0xEa4A0aA8aD7418f22373567BcA5728cB030f27Af"
	out, err := CallToMethod(
		rpcURL,
		addr,
		"balanceOf(address)->(uint256)",
		common.HexToAddress(owner),
	)
	require.NoError(t, err)
	require.Len(t, out, 1)
	bal, ok := out[0].(*big.Int)
	require.True(t, ok)
	require.NotNil(t, bal)
}
