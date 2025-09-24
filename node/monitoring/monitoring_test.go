// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package monitoring

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test Setup function
	err := Setup(tempDir)
	require.NoError(t, err)

	// Verify that dashboards directory was created
	dashboardDir := filepath.Join(tempDir, constants.DashboardsDir)
	assert.DirExists(t, dashboardDir)

	// Verify that dashboard files were written
	files, err := os.ReadDir(dashboardDir)
	require.NoError(t, err)
	assert.NotEmpty(t, files, "Dashboard files should be created")

	// Check for expected dashboard files
	expectedDashboards := []string{
		"a_chain.json",
		"d_chain.json",
		"database.json",
		"logs.json",
		"machine.json",
		"main.json",
		"network.json",
		"o_chain.json",
		"subnets.json",
	}

	for _, expectedFile := range expectedDashboards {
		found := false
		for _, file := range files {
			if file.Name() == expectedFile {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected dashboard file %s not found", expectedFile)
	}
}

func TestWriteMonitoringJSONFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Test WriteMonitoringJSONFiles function
	err := WriteMonitoringJSONFiles(tempDir)
	require.NoError(t, err)

	// Verify directory structure
	dashboardDir := filepath.Join(tempDir, constants.DashboardsDir)
	assert.DirExists(t, dashboardDir)

	// Verify files exist and are not empty
	files, err := os.ReadDir(dashboardDir)
	require.NoError(t, err)

	for _, file := range files {
		filePath := filepath.Join(dashboardDir, file.Name())
		content, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.NotEmpty(t, content, "Dashboard file %s should not be empty", file.Name())
	}
}

func TestWriteMonitoringJSONFiles_InvalidDirectory(t *testing.T) {
	// Test with invalid directory path
	invalidPath := "/invalid/path/that/does/not/exist"
	err := WriteMonitoringJSONFiles(invalidPath)
	assert.Error(t, err)
}

func TestGenerateConfig(t *testing.T) {
	tests := []struct {
		name         string
		configPath   string
		configDesc   string
		templateVars configInputs
		expectError  bool
	}{
		{
			name:       "Valid prometheus config",
			configPath: "configs/prometheus.yml",
			configDesc: "Prometheus Config",
			templateVars: configInputs{
				OdysseyGoPorts: "'9650','9651'",
				MachinePorts:   "'9100'",
				LoadTestPorts:  "'8082'",
			},
			expectError: false,
		},
		{
			name:       "Valid loki config",
			configPath: "configs/loki.yml",
			configDesc: "Loki Config",
			templateVars: configInputs{
				Port: "23101",
			},
			expectError: false,
		},
		{
			name:       "Valid promtail config",
			configPath: "configs/promtail.yml",
			configDesc: "Promtail Config",
			templateVars: configInputs{
				IP:      "127.0.0.1",
				Port:    "23101",
				Host:    "test-host",
				NodeID:  "test-node-id",
				ChainID: "test-chain-id",
			},
			expectError: false,
		},
		{
			name:         "Invalid config path",
			configPath:   "configs/nonexistent.yml",
			configDesc:   "Invalid Config",
			templateVars: configInputs{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := GenerateConfig(tt.configPath, tt.configDesc, tt.templateVars)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, config)
				// Verify that template variables are substituted
				if tt.templateVars.Port != "" {
					assert.Contains(t, config, tt.templateVars.Port)
				}
				if tt.templateVars.IP != "" {
					assert.Contains(t, config, tt.templateVars.IP)
				}
			}
		})
	}
}

func TestWritePrometheusConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "prometheus.yml")

	odysseyGoPorts := []string{"9650", "9651"}
	machinePorts := []string{"9100"}
	loadTestPorts := []string{"8082"}

	err := WritePrometheusConfig(configPath, odysseyGoPorts, machinePorts, loadTestPorts)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)

	// Verify file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configStr := string(content)

	// Check that ports are included in the config
	assert.Contains(t, configStr, "'9650'")
	assert.Contains(t, configStr, "'9651'")
	assert.Contains(t, configStr, "'9100'")
	assert.Contains(t, configStr, "'8082'")
}

func TestWritePrometheusConfig_EmptyPorts(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "prometheus.yml")

	// Test with empty port slices
	err := WritePrometheusConfig(configPath, []string{}, []string{}, []string{})
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)

	// Verify file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configStr := string(content)
	assert.NotEmpty(t, configStr)
}

func TestWriteLokiConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "loki.yml")

	port := "23101"
	err := WriteLokiConfig(configPath, port)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)

	// Verify file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configStr := string(content)

	// Check that port is included in the config
	assert.Contains(t, configStr, port)
}

func TestWriteLokiConfig_EmptyPort(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "loki.yml")

	// Test with empty port
	err := WriteLokiConfig(configPath, "")
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)
}

func TestWritePromtailConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "promtail.yml")

	lokiIP := "127.0.0.1"
	lokiPort := "23101"
	host := "test-host"
	nodeID := "test-node-id"
	chainID := "test-chain-id"

	err := WritePromtailConfig(configPath, lokiIP, lokiPort, host, nodeID, chainID)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)

	// Verify file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	configStr := string(content)

	// Check that all parameters are included in the config
	assert.Contains(t, configStr, lokiIP)
	assert.Contains(t, configStr, lokiPort)
	assert.Contains(t, configStr, host)
	assert.Contains(t, configStr, nodeID)
	assert.Contains(t, configStr, chainID)
}

func TestWritePromtailConfig_InvalidIP(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "promtail.yml")

	// Test with invalid IP
	invalidIP := "invalid-ip"
	err := WritePromtailConfig(configPath, invalidIP, "23101", "host", "node", "chain")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid IP address")
}

func TestWritePromtailConfig_ValidIPs(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "promtail.yml")

	validIPs := []string{
		"127.0.0.1",
		"192.168.1.1",
		"10.0.0.1",
		"::1",
		"2001:db8::1",
	}

	for _, ip := range validIPs {
		t.Run("IP_"+strings.ReplaceAll(ip, ":", "_"), func(t *testing.T) {
			err := WritePromtailConfig(configPath, ip, "23101", "host", "node", "chain")
			assert.NoError(t, err, "IP %s should be valid", ip)
		})
	}
}

func TestGetGrafanaURL(t *testing.T) {
	tests := []struct {
		name           string
		monitoringHost string
		expectedURL    string
	}{
		{
			name:           "Localhost IP",
			monitoringHost: "127.0.0.1",
			expectedURL:    "http://127.0.0.1:3000/dashboards",
		},
		{
			name:           "Remote IP",
			monitoringHost: "192.168.1.100",
			expectedURL:    "http://192.168.1.100:3000/dashboards",
		},
		{
			name:           "Domain name",
			monitoringHost: "monitoring.example.com",
			expectedURL:    "http://monitoring.example.com:3000/dashboards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := GetGrafanaURL(tt.monitoringHost)
			assert.Equal(t, tt.expectedURL, url)
		})
	}
}

func TestConfigInputs_Struct(t *testing.T) {
	// Test that configInputs struct can be instantiated and used
	inputs := configInputs{
		OdysseyGoPorts: "'9650','9651'",
		MachinePorts:   "'9100'",
		LoadTestPorts:  "'8082'",
		IP:             "127.0.0.1",
		Port:           "23101",
		Host:           "test-host",
		NodeID:         "test-node-id",
		ChainID:        "test-chain-id",
	}

	// Verify all fields are set correctly
	assert.Equal(t, "'9650','9651'", inputs.OdysseyGoPorts)
	assert.Equal(t, "'9100'", inputs.MachinePorts)
	assert.Equal(t, "'8082'", inputs.LoadTestPorts)
	assert.Equal(t, "127.0.0.1", inputs.IP)
	assert.Equal(t, "23101", inputs.Port)
	assert.Equal(t, "test-host", inputs.Host)
	assert.Equal(t, "test-node-id", inputs.NodeID)
	assert.Equal(t, "test-chain-id", inputs.ChainID)
}

func TestEmbeddedFiles(t *testing.T) {
	// Test that embedded files are accessible
	dashboardFiles, err := dashboards.ReadDir("dashboards")
	require.NoError(t, err)
	assert.NotEmpty(t, dashboardFiles, "Dashboard files should be embedded")

	configFiles, err := configs.ReadDir("configs")
	require.NoError(t, err)
	assert.NotEmpty(t, configFiles, "Config files should be embedded")

	// Verify specific files exist
	expectedConfigs := []string{"loki.yml", "prometheus.yml", "promtail.yml"}
	for _, expectedConfig := range expectedConfigs {
		found := false
		for _, file := range configFiles {
			if file.Name() == expectedConfig {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected config file %s not found in embedded files", expectedConfig)
	}
}

func TestWritePrometheusConfig_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "prometheus.yml")

	err := WritePrometheusConfig(configPath, []string{"9650"}, []string{"9100"}, []string{"8082"})
	require.NoError(t, err)

	// Check file permissions
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	// Should have read/write permissions for owner and read for group/others (644)
	expectedMode := os.FileMode(constants.WriteReadReadPerms)
	assert.Equal(t, expectedMode, info.Mode())
}

func TestWriteLokiConfig_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "loki.yml")

	err := WriteLokiConfig(configPath, "23101")
	require.NoError(t, err)

	// Check file permissions
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	// Should have read/write permissions for owner and read for group/others (644)
	expectedMode := os.FileMode(constants.WriteReadReadPerms)
	assert.Equal(t, expectedMode, info.Mode())
}

func TestWritePromtailConfig_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "promtail.yml")

	err := WritePromtailConfig(configPath, "127.0.0.1", "23101", "host", "node", "chain")
	require.NoError(t, err)

	// Check file permissions
	info, err := os.Stat(configPath)
	require.NoError(t, err)

	// Should have read/write permissions for owner and read for group/others (644)
	expectedMode := os.FileMode(constants.WriteReadReadPerms)
	assert.Equal(t, expectedMode, info.Mode())
}

func TestWriteMonitoringJSONFiles_DirectoryPermissions(t *testing.T) {
	tempDir := t.TempDir()

	err := WriteMonitoringJSONFiles(tempDir)
	require.NoError(t, err)

	// Check directory permissions
	dashboardDir := filepath.Join(tempDir, constants.DashboardsDir)
	info, err := os.Stat(dashboardDir)
	require.NoError(t, err)

	// Should have read/write/execute permissions for owner and read/execute for group/others (755)
	// Note: On some systems, the mode might include additional bits, so we check the base permissions
	expectedMode := os.FileMode(constants.DefaultPerms755)
	actualMode := info.Mode() & os.ModePerm // Mask out non-permission bits
	assert.Equal(t, expectedMode, actualMode)
}
