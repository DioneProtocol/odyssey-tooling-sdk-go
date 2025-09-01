// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package services

import (
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
)

func LokiFoldersToCreate() []string {
	return []string{utils.GetRemoteComposeServicePath(constants.ServiceLoki, "data")}
}
