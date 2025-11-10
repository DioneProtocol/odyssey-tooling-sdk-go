// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package services

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
)

type OdysseyConfigInputs struct {
	HTTPHost         string
	APIAdminEnabled  bool
	IndexEnabled     bool
	NetworkID        string
	DBDir            string
	LogDir           string
	PublicIP         string
	StateSyncEnabled bool
	PruningEnabled   bool
	TrackSubnets     string
	BootstrapIDs     string
	BootstrapIPs     string
	GenesisPath      string
}

func PrepareOdysseyConfig(publicIP string, networkID string, subnetsToTrack []string) OdysseyConfigInputs {
	return OdysseyConfigInputs{
		HTTPHost:         "0.0.0.0",
		NetworkID:        networkID,
		DBDir:            "/.odysseygo/db/",
		LogDir:           "/.odysseygo/logs/",
		PublicIP:         publicIP,
		StateSyncEnabled: true,
		PruningEnabled:   false,
		TrackSubnets:     strings.Join(subnetsToTrack, ","),
	}
}

func RenderOdysseyTemplate(templateName string, config OdysseyConfigInputs) ([]byte, error) {
	templateBytes, err := templates.ReadFile(templateName)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New("config").Parse(string(templateBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func RenderOdysseyNodeConfig(config OdysseyConfigInputs) ([]byte, error) {
	if output, err := RenderOdysseyTemplate("templates/odyssey-node.tmpl", config); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}

func RenderOdysseyDChainConfig(config OdysseyConfigInputs) ([]byte, error) {
	if output, err := RenderOdysseyTemplate("templates/odyssey-dchain.tmp", config); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}

func GetRemoteBLSKeyFile() string {
	return filepath.Join("/home/ubuntu/.odysseygo/staking/", constants.BLSKeyFileName)
}

func GetRemoteOdysseyNodeConfig() string {
	return filepath.Join("/home/ubuntu/.odysseygo/configs/", "node.json")
}

func GetRemoteOdysseyDChainConfig() string {
	return filepath.Join("/home/ubuntu/.odysseygo/configs/", "chains", "C", "config.json")
}

func GetRemoteOdysseyGenesis() string {
	return filepath.Join("/home/ubuntu/.odysseygo/configs/", "genesis.json")
}

func OdysseyFolderToCreate() []string {
	return []string{
		"/home/ubuntu/.odysseygo/db",
		"/home/ubuntu/.odysseygo/logs",
		"/home/ubuntu/.odysseygo/configs",
		"/home/ubuntu/.odysseygo/configs/subnets/",
		"/home/ubuntu/.odysseygo/configs/chains/C",
		"/home/ubuntu/.odysseygo/staking",
		"/home/ubuntu/.odysseygo/plugins",
	}
}
