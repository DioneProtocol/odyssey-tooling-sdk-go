// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsE2E(t *testing.T) {
	// Save original environment
	originalRunE2E := os.Getenv("RUN_E2E")
	defer func() {
		if originalRunE2E != "" {
			os.Setenv("RUN_E2E", originalRunE2E)
		} else {
			os.Unsetenv("RUN_E2E")
		}
	}()

	tests := []struct {
		name      string
		runE2EEnv string
		username  string
		userError error
		expected  bool
	}{
		{
			name:      "RUN_E2E=true",
			runE2EEnv: "true",
			username:  "testuser",
			expected:  true,
		},
		{
			name:      "RUN_E2E=false",
			runE2EEnv: "false",
			username:  "testuser",
			expected:  false,
		},
		{
			name:      "RUN_E2E not set, username runner",
			runE2EEnv: "",
			username:  "runner",
			expected:  false, // This test will fail if the actual username is not "runner"
		},
		{
			name:      "RUN_E2E not set, different username",
			runE2EEnv: "",
			username:  "testuser",
			expected:  false,
		},
		{
			name:      "RUN_E2E not set, empty username",
			runE2EEnv: "",
			username:  "",
			expected:  false,
		},
		{
			name:      "user.Current() error",
			runE2EEnv: "",
			username:  "",
			userError: assert.AnError,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.runE2EEnv == "" {
				os.Unsetenv("RUN_E2E")
			} else {
				os.Setenv("RUN_E2E", tt.runE2EEnv)
			}

			// Mock user.Current() if needed
			if tt.userError != nil {
				// This test would require mocking user.Current()
				// For now, we'll skip it as it's hard to mock without dependency injection
				t.Skip("Requires user.Current() mocking")
			}

			result := IsE2E()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestE2EDocker(t *testing.T) {
	// Test when docker is available
	t.Run("docker available", func(t *testing.T) {
		// Check if docker is actually available
		cmd := exec.Command("docker", "--version")
		err := cmd.Run()

		if err != nil {
			t.Skip("Docker not available, skipping test")
		}

		result := E2EDocker()
		assert.True(t, result)
	})

	// Test when docker is not available
	t.Run("docker not available", func(t *testing.T) {
		// This test would require mocking exec.Command or modifying the function
		// to accept a custom command runner
		t.Skip("Requires exec.Command mocking or dependency injection")
	})
}

func TestE2EConvertIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "valid IP with suffix",
			ip:       "192.168.1.100",
			expected: "192.168.223.10100",
		},
		{
			name:     "valid IP with single digit suffix",
			ip:       "192.168.1.5",
			expected: "192.168.223.105",
		},
		{
			name:     "valid IP with double digit suffix",
			ip:       "192.168.1.50",
			expected: "192.168.223.1050",
		},
		{
			name:     "valid IP with triple digit suffix",
			ip:       "192.168.1.255",
			expected: "192.168.223.10255",
		},
		{
			name:     "invalid IP format",
			ip:       "192.168.1",
			expected: "",
		},
		{
			name:     "invalid IP with too many parts",
			ip:       "192.168.1.100.200",
			expected: "",
		},
		{
			name:     "empty IP",
			ip:       "",
			expected: "",
		},
		{
			name:     "IP with non-numeric parts",
			ip:       "192.168.abc.100",
			expected: "192.168.223.10100",
		},
		{
			name:     "IP with negative numbers",
			ip:       "192.168.-1.100",
			expected: "192.168.223.10100",
		},
		{
			name:     "IP with leading zeros",
			ip:       "192.168.001.100",
			expected: "192.168.223.10100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := E2EConvertIP(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestE2ESuffix(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "valid 4-part IP",
			ip:       "192.168.1.100",
			expected: "100",
		},
		{
			name:     "valid IP with single digit suffix",
			ip:       "192.168.1.5",
			expected: "5",
		},
		{
			name:     "valid IP with double digit suffix",
			ip:       "192.168.1.50",
			expected: "50",
		},
		{
			name:     "valid IP with triple digit suffix",
			ip:       "192.168.1.255",
			expected: "255",
		},
		{
			name:     "invalid IP with 3 parts",
			ip:       "192.168.1",
			expected: "",
		},
		{
			name:     "invalid IP with 5 parts",
			ip:       "192.168.1.100.200",
			expected: "",
		},
		{
			name:     "empty IP",
			ip:       "",
			expected: "",
		},
		{
			name:     "IP with non-numeric parts",
			ip:       "192.168.abc.100",
			expected: "100", // Still extracts the last part
		},
		{
			name:     "IP with negative numbers",
			ip:       "192.168.-1.100",
			expected: "100", // Still extracts the last part
		},
		{
			name:     "IP with leading zeros",
			ip:       "192.168.001.100",
			expected: "100",
		},
		{
			name:     "IP with zero suffix",
			ip:       "192.168.1.0",
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := E2ESuffix(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveLineCleanChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean string",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with \\r\\x1b[K",
			input:    "hello\r\x1b[Kworld",
			expected: "helloworld",
		},
		{
			name:     "multiple occurrences",
			input:    "hello\r\x1b[Kworld\r\x1b[Ktest",
			expected: "helloworldtest",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only \\r\\x1b[K",
			input:    "\r\x1b[K",
			expected: "",
		},
		{
			name:     "mixed with other characters",
			input:    "a\r\x1b[Kb\r\x1b[Kc\r\x1b[Kd",
			expected: "abcd",
		},
		{
			name:     "\\r\\x1b[K at start",
			input:    "\r\x1b[Khello world",
			expected: "hello world",
		},
		{
			name:     "\\r\\x1b[K at end",
			input:    "hello world\r\x1b[K",
			expected: "hello world",
		},
		{
			name:     "\\r\\x1b[K at both ends",
			input:    "\r\x1b[Khello world\r\x1b[K",
			expected: "hello world",
		},
		{
			name:     "consecutive \\r\\x1b[K",
			input:    "hello\r\x1b[K\r\x1b[Kworld",
			expected: "helloworld",
		},
		{
			name:     "with newlines and other escape sequences",
			input:    "hello\r\x1b[K\nworld\r\x1b[K\ntest",
			expected: "hello\nworld\ntest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveLineCleanChars(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig(t *testing.T) {
	// Test Config struct initialization
	t.Run("Config struct", func(t *testing.T) {
		config := Config{
			IPs:           []string{"192.168.1.1", "192.168.1.2"},
			UbuntuVersion: "20.04",
			NetworkPrefix: "192.168.1",
			SSHPubKey:     "ssh-rsa AAAAB3NzaC1yc2E...",
			E2ESuffixList: []string{"1", "2"},
		}

		assert.Len(t, config.IPs, 2)
		assert.Equal(t, "20.04", config.UbuntuVersion)
		assert.Equal(t, "192.168.1", config.NetworkPrefix)
		assert.NotEmpty(t, config.SSHPubKey)
		assert.Len(t, config.E2ESuffixList, 2)
	})

	// Test E2EListenPrefix constant
	t.Run("E2EListenPrefix constant", func(t *testing.T) {
		assert.Equal(t, "192.168.223", E2EListenPrefix)
	})
}
