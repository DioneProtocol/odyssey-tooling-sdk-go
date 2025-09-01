// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avalanche

import "github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"

type VMType string

const (
	SubnetEvm = "Subnet-EVM"
)

func (v VMType) RepoName() string {
	switch v {
	case SubnetEvm:
		return constants.SubnetEVMRepoName
	default:
		return "unknown"
	}
}
