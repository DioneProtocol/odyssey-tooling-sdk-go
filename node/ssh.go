// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	remoteconfig "github.com/DioneProtocol/odyssey-tooling-sdk-go/node/config"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/node/monitoring"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
)

type scriptInputs struct {
	OdysseyGoVersion     string
	SubnetExportFileName string
	SubnetName           string
	ClusterName          string
	GoVersion            string
	IsDevNet             bool
	NetworkFlag          string
	SubnetVMBinaryPath   string
	SubnetEVMReleaseURL  string
	SubnetEVMArchive     string
	LoadTestRepoDir      string
	LoadTestRepo         string
	LoadTestPath         string
	LoadTestCommand      string
	LoadTestBranch       string
	LoadTestGitCommit    string
	CheckoutCommit       bool
	LoadTestResultFile   string
	GrafanaPkg           string
}

//go:embed shell/*.sh
var script embed.FS

// RunOverSSH runs provided script path over ssh.
// This script can be template as it will be rendered using scriptInputs vars
func (h *Node) RunOverSSH(
	scriptDesc string,
	timeout time.Duration,
	scriptPath string,
	templateVars scriptInputs,
) error {
	startTime := time.Now()
	shellScript, err := script.ReadFile(scriptPath)
	if err != nil {
		return err
	}
	var script bytes.Buffer
	t, err := template.New(scriptDesc).Parse(string(shellScript))
	if err != nil {
		return err
	}
	err = t.Execute(&script, templateVars)
	if err != nil {
		return err
	}

	if output, err := h.Command(nil, timeout, script.String()); err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}
	executionTime := time.Since(startTime)
	h.Logger.Infof("RunOverSSH[%s]%s took %s with err: %v", h.NodeID, scriptDesc, executionTime, err)
	return nil
}

// RunSSHSetupNode runs script to setup sdk dependencies on a remote host over SSH.
func (h *Node) RunSSHSetupNode() error {
	if err := h.RunOverSSH(
		"Setup Node",
		constants.SSHLongRunningScriptTimeout,
		"shell/setupNode.sh",
		scriptInputs{},
	); err != nil {
		return err
	}
	return nil
}

// RunSSHSetupDockerService runs script to setup docker compose service for CLI
func (h *Node) RunSSHSetupDockerService() error {
	if !constants.DockerSupportEnabled {
		return fmt.Errorf("Docker support functionality is disabled. Set constants.DockerSupportEnabled = true to enable")
	}
	if h.HasSystemDAvailable() {
		return h.RunOverSSH(
			"Setup Docker Service",
			constants.SSHLongRunningScriptTimeout,
			"shell/setupDockerService.sh",
			scriptInputs{},
		)
	} else {
		// no need to setup docker service
		return nil
	}
}

// RunSSHRestartOdysseygo runs script to restart odysseygo
func (h *Node) RunSSHRestartOdysseygo() error {
	remoteComposeFile := utils.GetRemoteComposeFile()
	return h.RestartDockerComposeService(remoteComposeFile, constants.ServiceOdysseygo, constants.SSHLongRunningScriptTimeout)
}

// RunSSHUpgradeOdysseygo runs script to upgrade odysseygo
func (h *Node) RunSSHUpgradeOdysseygo(odysseyGoVersion string) error {
	withMonitoring, err := h.WasNodeSetupWithMonitoring()
	if err != nil {
		return err
	}

	if err := h.ComposeOverSSH("Compose Node",
		constants.SSHScriptTimeout,
		"templates/odysseygo.docker-compose.yml",
		dockerComposeInputs{
			OdysseygoVersion: odysseyGoVersion,
			WithMonitoring:   withMonitoring,
			WithOdysseygo:    true,
			E2E:              utils.IsE2E(),
			E2EIP:            utils.E2EConvertIP(h.IP),
			E2ESuffix:        utils.E2ESuffix(h.IP),
		}); err != nil {
		return err
	}
	return h.RestartDockerCompose(constants.SSHLongRunningScriptTimeout)
}

// RunSSHStartOdysseygo runs script to start odysseygo
func (h *Node) RunSSHStartOdysseygo() error {
	return h.StartDockerComposeService(utils.GetRemoteComposeFile(), constants.ServiceOdysseygo, constants.SSHLongRunningScriptTimeout)
}

// RunSSHStopOdysseygo runs script to stop odysseygo
func (h *Node) RunSSHStopOdysseygo() error {
	return h.StopDockerComposeService(utils.GetRemoteComposeFile(), constants.ServiceOdysseygo, constants.SSHLongRunningScriptTimeout)
}

// RunSSHUpgradeSubnetEVM runs script to upgrade subnet evm
func (h *Node) RunSSHUpgradeSubnetEVM(subnetEVMBinaryPath string) error {
	if _, err := h.Commandf(nil, constants.SSHScriptTimeout, "cp -f subnet-evm %s", subnetEVMBinaryPath); err != nil {
		return err
	}
	return nil
}

func (h *Node) RunSSHSetupPrometheusConfig(odysseyGoPorts, machinePorts, loadTestPorts []string) error {
	for _, folder := range remoteconfig.PrometheusFoldersToCreate() {
		if err := h.MkdirAll(folder, constants.SSHFileOpsTimeout); err != nil {
			return err
		}
	}
	nodePrometheusConfigTemp := utils.GetRemoteComposeServicePath(constants.ServicePrometheus, "prometheus.yml")
	promConfig, err := os.CreateTemp("", constants.ServicePrometheus)
	if err != nil {
		return err
	}
	defer os.Remove(promConfig.Name())
	if err := monitoring.WritePrometheusConfig(promConfig.Name(), odysseyGoPorts, machinePorts, loadTestPorts); err != nil {
		return err
	}

	return h.Upload(
		promConfig.Name(),
		nodePrometheusConfigTemp,
		constants.SSHFileOpsTimeout,
	)
}

func (h *Node) RunSSHSetupLokiConfig(port int) error {
	for _, folder := range remoteconfig.LokiFoldersToCreate() {
		if err := h.MkdirAll(folder, constants.SSHFileOpsTimeout); err != nil {
			return err
		}
	}
	nodeLokiConfigTemp := utils.GetRemoteComposeServicePath(constants.ServiceLoki, "loki.yml")
	lokiConfig, err := os.CreateTemp("", constants.ServiceLoki)
	if err != nil {
		return err
	}
	defer os.Remove(lokiConfig.Name())
	if err := monitoring.WriteLokiConfig(lokiConfig.Name(), strconv.Itoa(port)); err != nil {
		return err
	}
	return h.Upload(
		lokiConfig.Name(),
		nodeLokiConfigTemp,
		constants.SSHFileOpsTimeout,
	)
}

func (h *Node) RunSSHSetupPromtailConfig(lokiIP string, lokiPort int, nodeID string, chainID string) error {
	for _, folder := range remoteconfig.PromtailFoldersToCreate() {
		if err := h.MkdirAll(folder, constants.SSHFileOpsTimeout); err != nil {
			return err
		}
	}
	nodePromtailConfigTemp := utils.GetRemoteComposeServicePath(constants.ServicePromtail, "promtail.yml")
	promtailConfig, err := os.CreateTemp("", constants.ServicePromtail)
	if err != nil {
		return err
	}
	defer os.Remove(promtailConfig.Name())

	if err := monitoring.WritePromtailConfig(promtailConfig.Name(), lokiIP, strconv.Itoa(lokiPort), lokiIP, nodeID, chainID); err != nil {
		return err
	}
	return h.Upload(
		promtailConfig.Name(),
		nodePromtailConfigTemp,
		constants.SSHFileOpsTimeout,
	)
}

// RunSSHGetNewSubnetEVMRelease runs script to download new subnet evm
func (h *Node) RunSSHGetNewSubnetEVMRelease(subnetEVMReleaseURL, subnetEVMArchive string) error {
	return h.RunOverSSH(
		"Get Subnet EVM Release",
		constants.SSHScriptTimeout,
		"shell/getNewSubnetEVMRelease.sh",
		scriptInputs{SubnetEVMReleaseURL: subnetEVMReleaseURL, SubnetEVMArchive: subnetEVMArchive},
	)
}

// RunSSHUploadStakingFiles uploads staking files to a remote host via SSH.
func (h *Node) RunSSHUploadStakingFiles(keyPath string) error {
	localStakingPath := "/home/ubuntu/.odysseygo/staking/"
	if err := h.MkdirAll(
		localStakingPath,
		constants.SSHFileOpsTimeout,
	); err != nil {
		return err
	}
	if err := h.Upload(
		filepath.Join(keyPath, constants.StakerCertFileName),
		filepath.Join(localStakingPath, constants.StakerCertFileName),
		constants.SSHFileOpsTimeout,
	); err != nil {
		return err
	}
	if err := h.Upload(
		filepath.Join(keyPath, constants.StakerKeyFileName),
		filepath.Join(localStakingPath, constants.StakerKeyFileName),
		constants.SSHFileOpsTimeout,
	); err != nil {
		return err
	}
	return h.Upload(
		filepath.Join(keyPath, constants.BLSKeyFileName),
		filepath.Join(localStakingPath, constants.BLSKeyFileName),
		constants.SSHFileOpsTimeout,
	)
}

// RunSSHSetupMonitoringFolders sets up monitoring folders
func (h *Node) RunSSHSetupMonitoringFolders() error {
	for _, folder := range remoteconfig.RemoteFoldersToCreateMonitoring() {
		if err := h.MkdirAll(folder, constants.SSHFileOpsTimeout); err != nil {
			return err
		}
	}
	return nil
}

// MonitorNodes links all the nodes specified with the monitoring node
// so that the monitoring host can start tracking the validator nodes metrics and collecting their
// logs
// Cloud functionality has been removed from this SDK.
func (h *Node) MonitorNodes(ctx context.Context, targets []Node, chainID string) error {
	return fmt.Errorf("cloud functionality has been removed from this SDK. Please use local monitoring setup instead")
}

// SyncSubnets reconfigures odysseygo to sync subnets
func (h *Node) SyncSubnets(subnetsToTrack []string) error {
	// necessary checks
	if !isOdysseyGoNode(*h) {
		return fmt.Errorf("%s is not a odysseygo node", h.NodeID)
	}
	withMonitoring, err := h.WasNodeSetupWithMonitoring()
	if err != nil {
		return err
	}
	if err := h.WaitForSSHShell(constants.SSHScriptTimeout); err != nil {
		return err
	}
	odysseyGoVersion, err := h.GetDockerImageVersion(constants.OdysseyGoDockerImage, constants.SSHScriptTimeout)
	if err != nil {
		return err
	}
	networkName, err := h.GetOdysseyGoNetworkName()
	if err != nil {
		return err
	}
	if err := h.ComposeSSHSetupNode(networkName, subnetsToTrack, odysseyGoVersion, withMonitoring); err != nil {
		return err
	}
	if err := h.RestartDockerCompose(constants.SSHScriptTimeout); err != nil {
		return err
	}

	return nil
}

func (h *Node) RunSSHCopyMonitoringDashboards(monitoringDashboardPath string) error {
	// TODO: download dashboards from github instead
	remoteDashboardsPath := utils.GetRemoteComposeServicePath("grafana", "dashboards")
	if !utils.DirectoryExists(monitoringDashboardPath) {
		return fmt.Errorf("%s does not exist", monitoringDashboardPath)
	}
	if err := h.MkdirAll(remoteDashboardsPath, constants.SSHFileOpsTimeout); err != nil {
		return err
	}
	monitoringDashboardPath = filepath.Join(monitoringDashboardPath, constants.DashboardsDir)
	dashboards, err := os.ReadDir(monitoringDashboardPath)
	if err != nil {
		return err
	}
	for _, dashboard := range dashboards {
		if err := h.Upload(
			filepath.Join(monitoringDashboardPath, dashboard.Name()),
			filepath.Join(remoteDashboardsPath, dashboard.Name()),
			constants.SSHFileOpsTimeout,
		); err != nil {
			return err
		}
	}
	if composeFileExists(*h) {
		return h.RestartDockerComposeService(utils.GetRemoteComposeFile(), constants.ServiceGrafana, constants.SSHScriptTimeout)
	}
	return nil
}
