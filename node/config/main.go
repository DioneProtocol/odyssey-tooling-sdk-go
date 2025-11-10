// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package services

import (
	"embed"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
)

//go:embed templates/*
var templates embed.FS

// RemoteFoldersToCreateMonitoring returns a list of folders that need to be created on the remote Monitoring server
func RemoteFoldersToCreateMonitoring() []string {
	return utils.AppendSlices[string](
		GrafanaFoldersToCreate(),
		LokiFoldersToCreate(),
		PrometheusFoldersToCreate(),
		PromtailFoldersToCreate(),
	)
}

// RemoteFoldersToCreateOdysseygo returns a list of folders that need to be created on the remote Odysseygo server
func RemoteFoldersToCreateOdysseygo() []string {
	return utils.AppendSlices[string](
		OdysseyFolderToCreate(),
		PromtailFoldersToCreate(),
	)
}
