// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package services

import (
	"github.com/ava-labs/avalanche-tooling-sdk-go/constants"
	"github.com/ava-labs/avalanche-tooling-sdk-go/utils"
)

func PrometheusFoldersToCreate() []string {
	return []string{
		utils.GetRemoteComposeServicePath(constants.ServicePrometheus),
		utils.GetRemoteComposeServicePath(constants.ServicePrometheus, "data"),
	}
}
