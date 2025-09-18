// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odysseygo/utils/units"
)

func TestNodesValidatePrimaryNetwork(t *testing.T) {
	require := require.New(t)
	// Note: This test requires real wallet/SSH, so keep it skipped to avoid false failures
	require.NotNil(t) // prevent unused linter on require
	t.Skip("Skipping test that requires wallet implementation")

	_ = Node{ // prevent unused lint if code below is re-enabled later
		NodeID: "NodeID-11111111111111111111111111111111LpoYY",
		IP:     "NODE_IP_ADDRESS",
		SSHConfig: SSHConfig{
			User:           constants.RemoteHostUser,
			PrivateKeyPath: "NODE_KEYPAIR_PRIVATE_KEY_PATH",
		},
		Roles: []SupportedRole{Validator},
	}

	_ = validator.PrimaryNetworkValidatorParams{
		Duration:    48 * time.Hour,
		StakeAmount: 2 * units.Dione,
	}
}
