// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/DioneProtocol/subnet-evm/accounts/abi/bind"
	"github.com/DioneProtocol/subnet-evm/core/types"
	"github.com/DioneProtocol/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestParseMethodSignature(t *testing.T) {
	t.Run("constructor_nonpayable", func(t *testing.T) {
		name, abiJSON, err := ParseMethodSignature("(address,uint256)", Constructor, nil, NonPayable, struct {
			Owner common.Address
			Cap   *big.Int
		}{})
		require.NoError(t, err)
		require.Equal(t, "", name)

		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(abiJSON), &out))
		require.Len(t, out, 1)
		require.Equal(t, "constructor", out[0]["type"])           //nolint:forcetypeassert
		require.Equal(t, "nonpayable", out[0]["stateMutability"]) //nolint:forcetypeassert
	})

	t.Run("method_view_with_outputs", func(t *testing.T) {
		name, abiJSON, err := ParseMethodSignature("balanceOf(address)->(uint256)", Method, nil, View, struct{ Account common.Address }{})
		require.NoError(t, err)
		require.Equal(t, "balanceOf", name)

		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(abiJSON), &out))
		require.Equal(t, "function", out[0]["type"])        //nolint:forcetypeassert
		require.Equal(t, "view", out[0]["stateMutability"]) //nolint:forcetypeassert
		outputs := out[0]["outputs"].([]interface{})
		require.Len(t, outputs, 1)
		m := outputs[0].(map[string]interface{})
		require.Equal(t, "uint256", m["type"]) //nolint:forcetypeassert
	})

	t.Run("method_payable", func(t *testing.T) {
		name, abiJSON, err := ParseMethodSignature("deposit()->()", Method, nil, Payable)
		require.NoError(t, err)
		require.Equal(t, "deposit", name)
		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(abiJSON), &out))
		require.Equal(t, "payable", out[0]["stateMutability"]) //nolint:forcetypeassert
	})

	t.Run("method_nonpayable", func(t *testing.T) {
		_, abiJSON, err := ParseMethodSignature("transfer(address,uint256)->(bool)", Method, nil, NonPayable, struct {
			To    common.Address
			Value *big.Int
		}{})
		require.NoError(t, err)
		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(abiJSON), &out))
		require.Equal(t, "nonpayable", out[0]["stateMutability"]) //nolint:forcetypeassert
	})

	t.Run("event_with_indexed_fields", func(t *testing.T) {
		_, abiJSON, err := ParseMethodSignature("Transfer(address,address,uint256)", Event, []int{0, 1}, NonPayable, struct {
			From  common.Address
			To    common.Address
			Value *big.Int
		}{})
		require.NoError(t, err)
		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(abiJSON), &out))
		inputs := out[0]["inputs"].([]interface{})
		require.Equal(t, true, inputs[0].(map[string]interface{})["indexed"]) //nolint:forcetypeassert
		require.Equal(t, true, inputs[1].(map[string]interface{})["indexed"]) //nolint:forcetypeassert
		require.Nil(t, inputs[2].(map[string]interface{})["indexed"])         // non-indexed
	})

	t.Run("unsupported_payment_kind", func(t *testing.T) {
		_, _, err := ParseMethodSignature("foo()", Method, nil, PaymentKind(42))
		require.Error(t, err)
	})

	t.Run("signature_without_parentheses_returns_name_only", func(t *testing.T) {
		name, abiJSON, err := ParseMethodSignature("foo", Method, nil, View)
		require.NoError(t, err)
		require.Equal(t, "foo", name)
		require.Equal(t, "", abiJSON)
	})

	t.Run("malformed_parentheses_errors", func(t *testing.T) {
		_, _, err := ParseMethodSignature("balanceOf(address", Method, nil, View)
		require.Error(t, err)
	})
}

func TestUnpackLog(t *testing.T) {
	t.Run("decode_indexed_and_nonindexed", func(t *testing.T) {
		_, eventABIJSON, err := ParseMethodSignature("Transfer(address,address,uint256)", Event, []int{0, 1}, NonPayable, struct {
			From  common.Address
			To    common.Address
			Value *big.Int
		}{})
		require.NoError(t, err)

		meta := &bind.MetaData{ABI: eventABIJSON}
		parsedABI, err := meta.GetAbi()
		require.NoError(t, err)
		event := parsedABI.Events["Transfer"]

		from := common.HexToAddress("0x1111111111111111111111111111111111111111")
		to := common.HexToAddress("0x2222222222222222222222222222222222222222")
		value := big.NewInt(123456789)

		topics := []common.Hash{event.ID, common.BytesToHash(common.LeftPadBytes(from.Bytes(), 32)), common.BytesToHash(common.LeftPadBytes(to.Bytes(), 32))}
		data, err := event.Inputs.NonIndexed().Pack(value)
		require.NoError(t, err)

		log := types.Log{Topics: topics, Data: data}

		type TransferEvent struct {
			From  common.Address
			To    common.Address
			Value *big.Int
		}
		var out TransferEvent
		require.NoError(t, UnpackLog("Transfer(address,address,uint256)", []int{0, 1}, log, &out))
		require.Equal(t, from, out.From)
		require.Equal(t, to, out.To)
		require.Equal(t, value, out.Value)
	})

	t.Run("invalid_log_should_error", func(t *testing.T) {
		_, eventABIJSON, err := ParseMethodSignature("Transfer(address,address,uint256)", Event, []int{0, 1}, NonPayable, struct {
			From  common.Address
			To    common.Address
			Value *big.Int
		}{})
		require.NoError(t, err)
		meta := &bind.MetaData{ABI: eventABIJSON}
		parsedABI, err := meta.GetAbi()
		require.NoError(t, err)
		event := parsedABI.Events["Transfer"]

		topics := []common.Hash{event.ID, common.BytesToHash(common.LeftPadBytes(common.HexToAddress("0x1").Bytes(), 32))}
		log := types.Log{Topics: topics, Data: []byte{}}

		var out struct {
			From  common.Address
			To    common.Address
			Value *big.Int
		}
		require.Error(t, UnpackLog("Transfer(address,address,uint256)", []int{0, 1}, log, &out))
	})
}

func TestGetEventFromLogs(t *testing.T) {
	logs := []*types.Log{{}, {}}
	count := -1
	parser := func(l types.Log) (string, error) {
		count++
		if count == 0 {
			return "", assertErr("first fail")
		}
		return "ok", nil
	}
	res, err := GetEventFromLogs[string](logs, parser)
	require.NoError(t, err)
	require.Equal(t, "ok", res)

	parserAlwaysErr := func(l types.Log) (string, error) { return "", assertErr("bad log") }
	_, err = GetEventFromLogs[string](logs, parserAlwaysErr)
	require.Error(t, err)
}

func TestIssueTx_InvalidBytes(t *testing.T) {
	var dummyClient ethclient.Client
	err := IssueTx(dummyClient, "0x01")
	require.Error(t, err)
}

func TestGetTxOptsWithSigner_InvalidKey(t *testing.T) {
	var dummyClient ethclient.Client
	_, err := GetTxOptsWithSigner(dummyClient, "deadbeef")
	require.Error(t, err)
}

func TestGetABIMaps_ErrorPaths(t *testing.T) {
	typesSlice := []string{"bool", "int"}
	values := []interface{}{true, 1, 2}
	_, err := getABIMaps(typesSlice, values)
	require.Error(t, err)

	typesSlice = []string{"bool", "int"}
	type OneField struct{ A bool }
	_, err = getABIMaps(typesSlice, OneField{})
	require.Error(t, err)

	typesSlice = []string{"[(bool)]"}
	_, err = getABIMaps(typesSlice, []interface{}{struct{ FieldName bool }{}})
	require.Error(t, err)
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
