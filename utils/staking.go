// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"encoding/pem"
	"fmt"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/staking"
	"github.com/DioneProtocol/odysseygo/utils/crypto/bls"
)

func NewBlsSecretKeyBytes() ([]byte, error) {
	blsSignerKey, err := bls.NewSecretKey()
	if err != nil {
		return nil, err
	}
	return bls.SecretKeyToBytes(blsSignerKey), nil
}

func ToNodeID(certBytes []byte) (ids.NodeID, error) {
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return ids.EmptyNodeID, fmt.Errorf("failed to decode certificate")
	}
	cert, err := staking.ParseCertificate(block.Bytes)
	if err != nil {
		return ids.EmptyNodeID, err
	}
	return ids.NodeIDFromCert(cert), nil
}
