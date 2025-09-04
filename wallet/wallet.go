// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package wallet

import (
	"context"
	"errors"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary/common"
)

var ErrNotReadyToCommit = errors.New("tx is not fully signed so can't be committed")

type Wallet struct {
	primary.Wallet
	Keychain keychain.Keychain
	options  []common.Option
	config   *primary.WalletConfig
}

func New(ctx context.Context, config *primary.WalletConfig) (Wallet, error) {
	if config == nil {
		return Wallet{}, errors.New("wallet config cannot be nil")
	}

	wallet, err := primary.MakeWallet(
		ctx,
		config,
	)

	// Determine network from URI and create keychain with proper network
	network := odyssey.NetworkFromURI(config.URI)
	kc := keychain.NewKeychainFromExisting(config.DIONEKeychain, network)

	return Wallet{
		Wallet:   wallet,
		Keychain: kc,
		config:   config,
	}, err
}

// SecureWalletIsChangeOwner ensures that a fee paying address (wallet's keychain) will receive
// the change UTXO and not a randomly selected auth key that may not be paying fees
func (w *Wallet) SecureWalletIsChangeOwner() {
	addrs := w.Addresses()
	changeAddr := addrs[0]
	// sets change to go to wallet addr (instead of any other subnet auth key)
	changeOwner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{changeAddr},
	}
	w.options = append(w.options, common.WithChangeOwner(changeOwner))
	w.Wallet = primary.NewWalletWithOptions(w.Wallet, w.options...)
}

// SetAuthKeys sets auth keys that will be used when signing txs, besides the wallet's Keychain fee
// paying ones
func (w *Wallet) SetAuthKeys(authKeys []ids.ShortID) {
	addrs := w.Addresses()
	addrsSet := set.Set[ids.ShortID]{}
	addrsSet.Add(addrs...)
	addrsSet.Add(authKeys...)
	w.options = append(w.options, common.WithCustomAddresses(addrsSet))
	w.Wallet = primary.NewWalletWithOptions(w.Wallet, w.options...)
}

func (w *Wallet) SetSubnetAuthMultisig(authKeys []ids.ShortID) {
	w.SecureWalletIsChangeOwner()
	w.SetAuthKeys(authKeys)
}

func (w *Wallet) Addresses() []ids.ShortID {
	return w.Keychain.Addresses().List()
}
