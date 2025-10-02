// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package constants

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		// Version constants
		{"UbuntuVersionLTS", UbuntuVersionLTS, "20.04"},
		{"BuildEnvGolangVersion", BuildEnvGolangVersion, "1.22.1"},

		// Cloud constants
		{"CloudOperationTimeout", CloudOperationTimeout, 2 * time.Minute},
		{"CloudServerStorageSize", CloudServerStorageSize, 1000},

		// AWS constants
		{"AWSCloudServerRunningState", AWSCloudServerRunningState, "running"},
		{"AWSDefaultInstanceType", AWSDefaultInstanceType, "c5.2xlarge"},
		{"AWSNodeIDPrefix", AWSNodeIDPrefix, "aws_node"},

		// GCP constants
		{"GCPDefaultImageProvider", GCPDefaultImageProvider, "dioneprotocol-experimental"},
		{"GCPDefaultInstanceType", GCPDefaultInstanceType, "e2-standard-8"},
		{"GCPImageFilter", GCPImageFilter, "family=odysseycli-ubuntu-2204 AND architecture=x86_64"},
		{"GCPEnvVar", GCPEnvVar, "GOOGLE_APPLICATION_CREDENTIALS"},
		{"GCPDefaultAuthKeyPath", GCPDefaultAuthKeyPath, "~/.config/gcloud/application_default_credentials.json"},
		{"GCPStaticIPPrefix", GCPStaticIPPrefix, "static-ip"},
		{"GCPErrReleasingStaticIP", GCPErrReleasingStaticIP, "failed to release gcp static ip"},
		{"GCPNodeIDPrefix", GCPNodeIDPrefix, "gcp_node"},

		// Port constants
		{"SSHTCPPort", SSHTCPPort, 22},
		{"OdysseygoAPIPort", OdysseygoAPIPort, 9650},
		{"OdysseygoP2PPort", OdysseygoP2PPort, 9651},
		{"OdysseygoGrafanaPort", OdysseygoGrafanaPort, 3000},
		{"OdysseygoLokiPort", OdysseygoLokiPort, 23101},
		{"OdysseygoMonitoringPort", OdysseygoMonitoringPort, 9090},
		{"OdysseygoMachineMetricsPort", OdysseygoMachineMetricsPort, 9100},
		{"OdysseygoLoadTestPort", OdysseygoLoadTestPort, 8082},

		// HTTP constants
		{"APIRequestTimeout", APIRequestTimeout, 30 * time.Second},
		{"APIRequestLargeTimeout", APIRequestLargeTimeout, 2 * time.Minute},

		// SSH constants
		{"SSHSleepBetweenChecks", SSHSleepBetweenChecks, 1 * time.Second},
		{"SSHLongRunningScriptTimeout", SSHLongRunningScriptTimeout, 10 * time.Minute},
		{"SSHFileOpsTimeout", SSHFileOpsTimeout, 100 * time.Second},
		{"SSHPOSTTimeout", SSHPOSTTimeout, 10 * time.Second},
		{"SSHScriptTimeout", SSHScriptTimeout, 2 * time.Minute},
		{"RemoteHostUser", RemoteHostUser, "ubuntu"},

		// Node constants
		{"CloudNodeCLIConfigBasePath", CloudNodeCLIConfigBasePath, "/home/ubuntu/.odyssey-cli/"},
		{"CloudNodeStakingPath", CloudNodeStakingPath, "/home/ubuntu/.odysseygo/staking/"},
		{"CloudNodeConfigPath", CloudNodeConfigPath, "/home/ubuntu/.odysseygo/configs/"},
		{"ServicesDir", ServicesDir, "services"},
		{"DashboardsDir", DashboardsDir, "dashboards"},

		// Service constants
		{"ServiceOdysseygo", ServiceOdysseygo, "odysseygo"},
		{"ServicePromtail", ServicePromtail, "promtail"},
		{"ServiceGrafana", ServiceGrafana, "grafana"},
		{"ServicePrometheus", ServicePrometheus, "prometheus"},
		{"ServiceLoki", ServiceLoki, "loki"},
		{"ServiceAWMRelayer", ServiceAWMRelayer, "awm-relayer"},

		// Permission constants
		{"DefaultPerms755", DefaultPerms755, 0o755},
		{"WriteReadReadPerms", WriteReadReadPerms, 0o644},
		{"WriteReadUserOnlyPerms", WriteReadUserOnlyPerms, 0o600},
		{"IPAddressSuffix", IPAddressSuffix, "/32"},

		// Odysseygo constants
		{"LocalAPIEndpoint", LocalAPIEndpoint, "http://127.0.0.1:9650"},
		{"OdysseyGoDockerImage", OdysseyGoDockerImage, "dionetech/odysseygo:develop"},
		{"OdysseyGoGitRepo", OdysseyGoGitRepo, "https://github.com/DioneProtocol/odysseygo"},
		{"SubnetEVMRepoName", SubnetEVMRepoName, "subnet-evm"},

		// AWM Relayer constants
		{"AWMRelayerInstallDir", AWMRelayerInstallDir, "awm-relayer"},
		{"AWMRelayerConfigFilename", AWMRelayerConfigFilename, "awm-relayer-config.json"},

		// File name constants
		{"StakerCertFileName", StakerCertFileName, "staker.crt"},
		{"StakerKeyFileName", StakerKeyFileName, "staker.key"},
		{"BLSKeyFileName", BLSKeyFileName, "signer.key"},

		// GitHub constants
		{"DioneProtocolOrg", DioneProtocolOrg, "DioneProtocol"},
		{"ICMRepoName", ICMRepoName, "teleporter"},
		{"RelayerRepoName", RelayerRepoName, "awm-relayer"},
		{"RelayerBinName", RelayerBinName, "awm-relayer"},

		// Feature flag constants
		{"OdysseyCoreEnabled", OdysseyCoreEnabled, true},
		{"OdysseyTestnetEnabled", OdysseyTestnetEnabled, true},
		{"OdysseyMainnetEnabled", OdysseyMainnetEnabled, true},
		{"OdysseyDevnetEnabled", OdysseyDevnetEnabled, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant %s = %v, expected %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestFeatureFlags(t *testing.T) {
	tests := []struct {
		name     string
		flag     *bool
		expected bool
	}{
		{"RelayerEnabled", &RelayerEnabled, false},
		{"TeleporterEnabled", &TeleporterEnabled, false},
		{"AWSIntegrationEnabled", &AWSIntegrationEnabled, true},
		{"GCPIntegrationEnabled", &GCPIntegrationEnabled, true},
		{"DockerSupportEnabled", &DockerSupportEnabled, true},
		{"InstanceManagementEnabled", &InstanceManagementEnabled, true},
		{"SecurityGroupsEnabled", &SecurityGroupsEnabled, true},
		{"SSHKeyManagementEnabled", &SSHKeyManagementEnabled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *tt.flag != tt.expected {
				t.Errorf("Feature flag %s = %v, expected %v", tt.name, *tt.flag, tt.expected)
			}
		})
	}
}

func TestFeatureFlagModification(t *testing.T) {
	// Test that feature flags can be modified
	originalRelayerEnabled := RelayerEnabled
	originalTeleporterEnabled := TeleporterEnabled

	// Modify flags
	RelayerEnabled = true
	TeleporterEnabled = true

	// Verify changes
	if !RelayerEnabled {
		t.Error("RelayerEnabled should be true after modification")
	}
	if !TeleporterEnabled {
		t.Error("TeleporterEnabled should be true after modification")
	}

	// Restore original values
	RelayerEnabled = originalRelayerEnabled
	TeleporterEnabled = originalTeleporterEnabled
}

func TestTimeConstants(t *testing.T) {
	// Test time-based constants
	if CloudOperationTimeout != 2*time.Minute {
		t.Errorf("CloudOperationTimeout = %v, expected %v", CloudOperationTimeout, 2*time.Minute)
	}

	if APIRequestTimeout != 30*time.Second {
		t.Errorf("APIRequestTimeout = %v, expected %v", APIRequestTimeout, 30*time.Second)
	}

	if APIRequestLargeTimeout != 2*time.Minute {
		t.Errorf("APIRequestLargeTimeout = %v, expected %v", APIRequestLargeTimeout, 2*time.Minute)
	}

	if SSHSleepBetweenChecks != 1*time.Second {
		t.Errorf("SSHSleepBetweenChecks = %v, expected %v", SSHSleepBetweenChecks, 1*time.Second)
	}

	if SSHLongRunningScriptTimeout != 10*time.Minute {
		t.Errorf("SSHLongRunningScriptTimeout = %v, expected %v", SSHLongRunningScriptTimeout, 10*time.Minute)
	}

	if SSHFileOpsTimeout != 100*time.Second {
		t.Errorf("SSHFileOpsTimeout = %v, expected %v", SSHFileOpsTimeout, 100*time.Second)
	}

	if SSHPOSTTimeout != 10*time.Second {
		t.Errorf("SSHPOSTTimeout = %v, expected %v", SSHPOSTTimeout, 10*time.Second)
	}

	if SSHScriptTimeout != 2*time.Minute {
		t.Errorf("SSHScriptTimeout = %v, expected %v", SSHScriptTimeout, 2*time.Minute)
	}
}

func TestPortConstants(t *testing.T) {
	// Test port constants are within valid range
	ports := []int{
		SSHTCPPort,
		OdysseygoAPIPort,
		OdysseygoP2PPort,
		OdysseygoGrafanaPort,
		OdysseygoLokiPort,
		OdysseygoMonitoringPort,
		OdysseygoMachineMetricsPort,
		OdysseygoLoadTestPort,
	}

	for _, port := range ports {
		if port < 1 || port > 65535 {
			t.Errorf("Port %d is not within valid range (1-65535)", port)
		}
	}
}

func TestPermissionConstants(t *testing.T) {
	// Test permission constants are valid octal values
	permissions := []int{
		DefaultPerms755,
		WriteReadReadPerms,
		WriteReadUserOnlyPerms,
	}

	for _, perm := range permissions {
		if perm < 0 || perm > 0777 {
			t.Errorf("Permission %o is not within valid range (0-0777)", perm)
		}
	}
}
