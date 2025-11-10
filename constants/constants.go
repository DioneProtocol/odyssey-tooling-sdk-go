// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package constants

import "time"

const (
	// versions
	UbuntuVersionLTS      = "20.04"
	BuildEnvGolangVersion = "1.22.1"

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
	ServicesDir   = "services"
	DashboardsDir = "dashboards"
	// services
	ServiceOdysseygo  = "odysseygo"
	ServicePromtail   = "promtail"
	ServiceGrafana    = "grafana"
	ServicePrometheus = "prometheus"
	ServiceLoki       = "loki"

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

	StakerCertFileName = "staker.crt"
	StakerKeyFileName  = "staker.key"
	BLSKeyFileName     = "signer.key"

	// github
	DioneProtocolOrg = "DioneProtocol"

	// Odyssey Chain specific feature flags
	OdysseyCoreEnabled    = true
	OdysseyTestnetEnabled = true
	OdysseyMainnetEnabled = true
	OdysseyDevnetEnabled  = false
)

// Feature flags - Set to false to disable functionality
var (
	// Infrastructure feature flags
	DockerSupportEnabled      = true
	InstanceManagementEnabled = true
	SecurityGroupsEnabled     = true
	SSHKeyManagementEnabled   = true
)
