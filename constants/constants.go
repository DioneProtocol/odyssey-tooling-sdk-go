// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package constants

import "time"

const (
	// versions
	UbuntuVersionLTS      = "20.04"
	BuildEnvGolangVersion = "1.22.1"

	// clouds
	CloudOperationTimeout  = 2 * time.Minute
	CloudServerStorageSize = 1000

	AWSCloudServerRunningState = "running"
	AWSDefaultInstanceType     = "c5.2xlarge"
	AWSNodeIDPrefix            = "aws_node"

	GCPDefaultImageProvider = "dioneprotocol-experimental"
	GCPDefaultInstanceType  = "e2-standard-8"
	GCPImageFilter          = "family=odysseycli-ubuntu-2204 AND architecture=x86_64"
	GCPEnvVar               = "GOOGLE_APPLICATION_CREDENTIALS"
	GCPDefaultAuthKeyPath   = "~/.config/gcloud/application_default_credentials.json"
	GCPStaticIPPrefix       = "static-ip"
	GCPErrReleasingStaticIP = "failed to release gcp static ip"
	GCPNodeIDPrefix         = "gcp_node"

	// ports
	SSHTCPPort                  = 22
	OdysseygoAPIPort            = 9650
	OdysseygoP2PPort            = 9651
	OdysseygoGrafanaPort        = 3000
	OdysseygoLokiPort           = 23101
	OdysseygoMonitoringPort     = 9090
	OdysseygoMachineMetricsPort = 9100
	OdysseygoLoadTestPort       = 8082

	// http
	APIRequestTimeout      = 30 * time.Second
	APIRequestLargeTimeout = 2 * time.Minute

	// ssh
	SSHSleepBetweenChecks       = 1 * time.Second
	SSHLongRunningScriptTimeout = 10 * time.Minute
	SSHFileOpsTimeout           = 100 * time.Second
	SSHPOSTTimeout              = 10 * time.Second
	SSHScriptTimeout            = 2 * time.Minute
	RemoteHostUser              = "ubuntu"

	// node
	CloudNodeCLIConfigBasePath = "/home/ubuntu/.odyssey-cli/"
	CloudNodeStakingPath       = "/home/ubuntu/.odysseygo/staking/"
	CloudNodeConfigPath        = "/home/ubuntu/.odysseygo/configs/"
	ServicesDir                = "services"
	DashboardsDir              = "dashboards"
	// services
	ServiceOdysseygo  = "odysseygo"
	ServicePromtail   = "promtail"
	ServiceGrafana    = "grafana"
	ServicePrometheus = "prometheus"
	ServiceLoki       = "loki"
	ServiceAWMRelayer = "awm-relayer"

	// misc
	DefaultPerms755        = 0o755
	WriteReadReadPerms     = 0o644
	WriteReadUserOnlyPerms = 0o600
	IPAddressSuffix        = "/32"

	// odysseygo
	LocalAPIEndpoint = "http://127.0.0.1:9650"
	// TODO: change to the latest release version
	OdysseyGoDockerImage = "dionetech/odysseygo:develop"
	OdysseyGoGitRepo     = "https://github.com/DioneProtocol/odysseygo"
	SubnetEVMRepoName    = "subnet-evm"

	AWMRelayerInstallDir     = "awm-relayer"
	AWMRelayerConfigFilename = "awm-relayer-config.json"

	StakerCertFileName = "staker.crt"
	StakerKeyFileName  = "staker.key"
	BLSKeyFileName     = "signer.key"

	// github
	DioneProtocolOrg = "DioneProtocol"
	ICMRepoName      = "teleporter"
	RelayerRepoName  = "awm-relayer"
	RelayerBinName   = "awm-relayer"

	// Odyssey Chain specific feature flags
	OdysseyCoreEnabled    = true
	OdysseyTestnetEnabled = true
	OdysseyMainnetEnabled = true
	OdysseyDevnetEnabled  = false
)

// Feature flags - Set to false to disable functionality
var (
	RelayerEnabled    = false
	TeleporterEnabled = false

	// Cloud integration feature flags
	AWSIntegrationEnabled = true
	GCPIntegrationEnabled = true

	// Infrastructure feature flags
	DockerSupportEnabled      = true
	InstanceManagementEnabled = true
	SecurityGroupsEnabled     = true
	SSHKeyManagementEnabled   = true
)
