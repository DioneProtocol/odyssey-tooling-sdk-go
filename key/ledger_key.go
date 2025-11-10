// Copyright (c) 2025 Dione Limited.
// See the file LICENSE for licensing terms.

package key

import (
	"fmt"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

var _ Key = &LedgerKey{}

type LedgerKey struct {
	index uint32
}

// ledger device should be connected
func NewLedger(index uint32) LedgerKey {
	return LedgerKey{
		index: index,
	}
}

// LoadLedger loads the ledger key info from disk and creates the corresponding LedgerKey.
func LoadLedger(_ string) (*LedgerKey, error) {
	return nil, fmt.Errorf("not implemented")
}

// LoadLedgerFromBytes loads the ledger key info from bytes and creates the corresponding LedgerKey.
func LoadLedgerFromBytes(_ []byte) (*SoftKey, error) {
	return nil, fmt.Errorf("not implemented")
}

func (*LedgerKey) D() string {
	return ""
}

// Returns the KeyChain
func (*LedgerKey) KeyChain() *secp256k1fx.Keychain {
	return nil
}

// Saves the key info to disk
func (*LedgerKey) Save(_ string) error {
	return fmt.Errorf("not implemented")
}

func (*LedgerKey) O(_ string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (*LedgerKey) A(_ string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (*LedgerKey) Spends(_ []*dione.UTXO, _ ...OpOption) (
	totalBalanceToSpend uint64,
	inputs []*dione.TransferableInput,
	signers [][]ids.ShortID,
) {
	return 0, nil, nil
}

func (*LedgerKey) Addresses() []ids.ShortID {
	return nil
}

func (*LedgerKey) Sign(_ *txs.Tx, _ [][]ids.ShortID) error {
	return fmt.Errorf("not implemented")
}

func (*LedgerKey) Match(_ *secp256k1fx.OutputOwners, _ uint64) ([]uint32, []ids.ShortID, bool) {
	return nil, nil, false
}
