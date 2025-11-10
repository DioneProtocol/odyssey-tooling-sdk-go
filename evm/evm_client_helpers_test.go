// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package evm

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/DioneProtocol/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetContractBytecodeAndAlreadyDeployed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	mock.EXPECT().CodeAt(gomock.Any(), addr, gomock.Nil()).Return([]byte{0x01, 0x02}, nil).Times(2)

	bs, err := GetContractBytecode(mock, addr.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, bs)

	b, err := ContractAlreadyDeployed(mock, addr.Hex())
	require.NoError(t, err)
	require.True(t, b)
}

func TestGetAddressBalanceNonceTipBaseFee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	addr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	mock.EXPECT().BalanceAt(gomock.Any(), addr, gomock.Nil()).Return(big.NewInt(123), nil)
	mock.EXPECT().NonceAt(gomock.Any(), addr, gomock.Nil()).Return(uint64(7), nil)
	mock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(2_500_000_000), nil)
	mock.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(25_000_000_000), nil)

	bal, err := GetAddressBalance(mock, addr.Hex())
	require.NoError(t, err)
	require.Equal(t, big.NewInt(123), bal)

	nonce, err := NonceAt(mock, addr.Hex())
	require.NoError(t, err)
	require.Equal(t, uint64(7), nonce)

	gasTip, err := SuggestGasTipCap(mock)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(MaxPriorityFeePerGas), gasTip)

	baseFee, err := EstimateBaseFee(mock)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(25_000_000_000), baseFee)
}

func TestCalculateTxParams_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	addr := common.HexToAddress("0x3333333333333333333333333333333333333333")
	mock.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(10), nil)
	mock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(5), nil)
	mock.EXPECT().NonceAt(gomock.Any(), addr, gomock.Nil()).Return(uint64(42), nil)

	feeCap, tip, nonce, err := CalculateTxParams(mock, addr.Hex())
	require.NoError(t, err)

	expected := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(BaseFeeFactor))
	expected.Add(expected, big.NewInt(MaxPriorityFeePerGas))
	require.Equal(t, 0, expected.Cmp(feeCap))
	require.Equal(t, big.NewInt(5), tip)
	require.Equal(t, uint64(42), nonce)
}

func TestSetMinBalance_FundsWhenBelow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	addr := common.HexToAddress("0x4444444444444444444444444444444444444444")
	required := big.NewInt(1000)
	mock.EXPECT().BalanceAt(gomock.Any(), addr, gomock.Nil()).Return(big.NewInt(999), nil)

	mock.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(1), nil)
	mock.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(1), nil)
	mock.EXPECT().NonceAt(gomock.Any(), gomock.Any(), gomock.Nil()).Return(uint64(0), nil)
	mock.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)
	mock.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)
	mock.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptSuccessful, nil).AnyTimes()

	priv := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	require.NoError(t, SetMinBalance(mock, priv, addr.Hex(), required))
}

var receiptSuccessful = types.Receipt{Status: types.ReceiptStatusSuccessful}
var receiptFailed = types.Receipt{Status: types.ReceiptStatusFailed}

func TestSendTransaction_SuccessImmediate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	mock.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)
	require.NoError(t, SendTransaction(mock, nil))
}

func TestWaitForTransaction_SuccessAndFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	tx := types.NewTx(&types.LegacyTx{})
	mock.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptSuccessful, nil)
	rec, ok, err := WaitForTransaction(mock, tx)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, types.ReceiptStatusSuccessful, rec.Status)

	mock2 := NewMockClient(ctrl)
	tx2 := types.NewTx(&types.LegacyTx{})
	mock2.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptFailed, nil)
	rec2, ok2, err2 := WaitForTransaction(mock2, tx2)
	require.NoError(t, err2)
	require.False(t, ok2)
	require.Equal(t, types.ReceiptStatusFailed, rec2.Status)
}

func TestTransfer_SuccessAndFailedReceipt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	priv := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	srcAddrHex := "0x90f8bf6a479f320ead074411a4b0e7944ea8c9c1" // address from the above priv
	tgtAddr := common.HexToAddress("0x5555555555555555555555555555555555555555")
	amount := big.NewInt(1)

	// success path
	m1 := NewMockClient(ctrl)
	m1.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(1), nil)
	m1.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(1), nil)
	m1.EXPECT().NonceAt(gomock.Any(), common.HexToAddress(srcAddrHex), gomock.Nil()).Return(uint64(0), nil)
	m1.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)
	m1.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)
	m1.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptSuccessful, nil).AnyTimes()
	require.NoError(t, Transfer(m1, priv, tgtAddr.Hex(), amount))

	// failed receipt path
	m2 := NewMockClient(ctrl)
	m2.EXPECT().EstimateBaseFee(gomock.Any()).Return(big.NewInt(1), nil)
	m2.EXPECT().SuggestGasTipCap(gomock.Any()).Return(big.NewInt(1), nil)
	m2.EXPECT().NonceAt(gomock.Any(), common.HexToAddress(srcAddrHex), gomock.Nil()).Return(uint64(0), nil)
	m2.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)
	m2.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)
	m2.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptFailed, nil).AnyTimes()
	require.Error(t, Transfer(m2, priv, tgtAddr.Hex(), amount))
}

func TestGetTxOptsWithSigner_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)
	mock.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1), nil)
	priv := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	opts, err := GetTxOptsWithSigner(mock, priv)
	require.NoError(t, err)
	require.NotNil(t, opts)
}

func TestIssueTx_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := NewMockClient(ctrl)

	// build a small legacy tx and encode to hex
	tx := types.NewTx(&types.LegacyTx{Nonce: 0})
	bs, err := tx.MarshalBinary()
	require.NoError(t, err)
	hexStr := "0x" + hex.EncodeToString(bs)

	mock.EXPECT().SendTransaction(gomock.Any(), gomock.Any()).Return(nil)
	mock.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(&receiptSuccessful, nil).AnyTimes()
	require.NoError(t, IssueTx(mock, hexStr))
}
