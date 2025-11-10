// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSCPTargetPath(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		path     string
		expected string
	}{
		{
			name:     "empty IP",
			ip:       "",
			path:     "/home/user/file.txt",
			expected: "/home/user/file.txt",
		},
		{
			name:     "valid IP",
			ip:       "192.168.1.100",
			path:     "/home/user/file.txt",
			expected: "ubuntu@192.168.1.100:/home/user/file.txt",
		},
		{
			name:     "IP with different path",
			ip:       "10.0.0.1",
			path:     "/var/log/app.log",
			expected: "ubuntu@10.0.0.1:/var/log/app.log",
		},
		{
			name:     "empty path",
			ip:       "192.168.1.100",
			path:     "",
			expected: "ubuntu@192.168.1.100:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSCPTargetPath(tt.ip, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitSCPPath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedNode string
		expectedPath string
	}{
		{
			name:         "path with colon",
			path:         "ubuntu@192.168.1.100:/home/user/file.txt",
			expectedNode: "ubuntu@192.168.1.100",
			expectedPath: "/home/user/file.txt",
		},
		{
			name:         "path without colon",
			path:         "/home/user/file.txt",
			expectedNode: "",
			expectedPath: "/home/user/file.txt",
		},
		{
			name:         "empty path",
			path:         "",
			expectedNode: "",
			expectedPath: "",
		},
		{
			name:         "multiple colons",
			path:         "user@host:path:with:colons",
			expectedNode: "user@host",
			expectedPath: "path",
		},
		{
			name:         "only colon",
			path:         ":",
			expectedNode: "",
			expectedPath: "",
		},
		{
			name:         "colon at start",
			path:         ":path",
			expectedNode: "",
			expectedPath: "path",
		},
		{
			name:         "colon at end",
			path:         "node:",
			expectedNode: "node",
			expectedPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, path := SplitSCPPath(tt.path)
			assert.Equal(t, tt.expectedNode, node)
			assert.Equal(t, tt.expectedPath, path)
		})
	}
}

func TestCombineSCPPath(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		path     string
		expected string
	}{
		{
			name:     "empty host",
			host:     "",
			path:     "/home/user/file.txt",
			expected: "/home/user/file.txt",
		},
		{
			name:     "valid host and path",
			host:     "ubuntu@192.168.1.100",
			path:     "/home/user/file.txt",
			expected: "ubuntu@192.168.1.100:/home/user/file.txt",
		},
		{
			name:     "empty path",
			host:     "ubuntu@192.168.1.100",
			path:     "",
			expected: "ubuntu@192.168.1.100:",
		},
		{
			name:     "both empty",
			host:     "",
			path:     "",
			expected: "",
		},
		{
			name:     "host with port",
			host:     "ubuntu@192.168.1.100:22",
			path:     "/home/user/file.txt",
			expected: "ubuntu@192.168.1.100:22:/home/user/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CombineSCPPath(tt.host, tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSSHAgentAvailable(t *testing.T) {
	// Save original environment
	originalSSHAuthSock := os.Getenv("SSH_AUTH_SOCK")
	defer func() {
		if originalSSHAuthSock != "" {
			os.Setenv("SSH_AUTH_SOCK", originalSSHAuthSock)
		} else {
			os.Unsetenv("SSH_AUTH_SOCK")
		}
	}()

	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "SSH_AUTH_SOCK set",
			envValue: "/tmp/ssh-agent.sock",
			expected: true,
		},
		{
			name:     "SSH_AUTH_SOCK not set",
			envValue: "",
			expected: false,
		},
		{
			name:     "SSH_AUTH_SOCK empty",
			envValue: " ",
			expected: true, // Non-empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue == "" {
				os.Unsetenv("SSH_AUTH_SOCK")
			} else {
				os.Setenv("SSH_AUTH_SOCK", tt.envValue)
			}

			result := IsSSHAgentAvailable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSSHPubKey(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "Valid RSA key",
			key:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pBe3b2m5zJKvVlWfk7F0uRcxJ5LnA73LJ2+AW+JQLCRg5P5RPnRg3U4aV4n07a/x33UvCff3Dv+5G2E7QKQGxLizHcBKkE1dFpnO5BPNjSFK/4q+TFDdgA2YC47PODqDxXzOdb+et+1db/f4wYfPgqF2n1A1UXkG5pSzxNzMWvMEW6LeqA5Zq8cVnR51fESsWGDoqAptZ0J2B7s/UMGMbhZZqWflP1p6gAV3dpePC3F2Qf/SjCXHh4rpqDvHBLR4IKmI0zRiZ/Vq+H7Z39a6zXNyAT4PN/YCrX2q4VfljE4oJH3MC6+Vvjfg3tzRZkFshIJg0K4tOP1zqDAFj user@hostname",
			expected: true,
		},
		{
			name:     "Valid ed25519 key",
			key:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMiP5BOZuRjjE7V/HjDmGtBf/2YhUoZ1Fn5O8ss+nG3",
			expected: true,
		},
		{
			name:     "Valid ecdsa key",
			key:      "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGe3ZjpyUBwM/43Alp2H7zQd48+4mV+//kjcsdqJcIK1FwR7d+8tjLjrUq8Ow2V1YLZMfdk9sC8Knyyl8Z5+Y=",
			expected: true,
		},
		{
			name:     "Invalid key type",
			key:      "ssh-rsa-invalid AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pBe3b2m5zJKvVlWfk7F0uRcxJ5LnA73LJ2+AW+JQLCRg5P5RPnRg3U4aV4n07a/x33UvCff3Dv+5G2E7QKQGxLizHcBKkE1dFpnO5BPNjSFK/4q+TFDdgA2YC47PODqDxXzOdb+et+1db/f4wYfPgqF2n1A1UXkG5pSzxNzMWvMEW6LeqA5Zq8cVnR51fESsWGDoqAptZ0J2B7s/UMGMbhZZqWflP1p6gAV3dpePC3F2Qf/SjCXHh4rpqDvHBLR4IKmI0zRiZ/Vq+H7Z39a6zXNyAT4PN/YCrX2q4VfljE4oJH3MC6+Vvjfg3tzRZkFshIJg0K4tOP1zqDAFj user@hostname",
			expected: false,
		},
		{
			name:     "Empty key",
			key:      "",
			expected: false,
		},
		{
			name:     "Key with quotes",
			key:      "\"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pBe3b2m5zJKvVlWfk7F0uRcxJ5LnA73LJ2+AW+JQLCRg5P5RPnRg3U4aV4n07a/x33UvCff3Dv+5G2E7QKQGxLizHcBKkE1dFpnO5BPNjSFK/4q+TFDdgA2YC47PODqDxXzOdb+et+1db/f4wYfPgqF2n1A1UXkG5pSzxNzMWvMEW6LeqA5Zq8cVnR51fESsWGDoqAptZ0J2B7s/UMGMbhZZqWflP1p6gAV3dpePC3F2Qf/SjCXHh4rpqDvHBLR4IKmI0zRiZ/Vq+H7Z39a6zXNyAT4PN/YCrX2q4VfljE4oJH3MC6+Vvjfg3tzRZkFshIJg0K4tOP1zqDAFj user@hostname\"",
			expected: true,
		},
		{
			name:     "Key with single quotes",
			key:      "'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pBe3b2m5zJKvVlWfk7F0uRcxJ5LnA73LJ2+AW+JQLCRg5P5RPnRg3U4aV4n07a/x33UvCff3Dv+5G2E7QKQGxLizHcBKkE1dFpnO5BPNjSFK/4q+TFDdgA2YC47PODqDxXzOdb+et+1db/f4wYfPgqF2n1A1UXkG5pSzxNzMWvMEW6LeqA5Zq8cVnR51fESsWGDoqAptZ0J2B7s/UMGMbhZZqWflP1p6gAV3dpePC3F2Qf/SjCXHh4rpqDvHBLR4IKmI0zRiZ/Vq+H7Z39a6zXNyAT4PN/YCrX2q4VfljE4oJH3MC6+Vvjfg3tzRZkFshIJg0K4tOP1zqDAFj user@hostname'",
			expected: true,
		},
		{
			name:     "Invalid base64",
			key:      "ssh-rsa invalid-base64-here user@hostname",
			expected: false,
		},
		{
			name:     "Missing key data",
			key:      "ssh-rsa",
			expected: false,
		},
		{
			name:     "Key with comment",
			key:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pBe3b2m5zJKvVlWfk7F0uRcxJ5LnA73LJ2+AW+JQLCRg5P5RPnRg3U4aV4n07a/x33UvCff3Dv+5G2E7QKQGxLizHcBKkE1dFpnO5BPNjSFK/4q+TFDdgA2YC47PODqDxXzOdb+et+1db/f4wYfPgqF2n1A1UXkG5pSzxNzMWvMEW6LeqA5Zq8cVnR51fESsWGDoqAptZ0J2B7s/UMGMbhZZqWflP1p6gAV3dpePC3F2Qf/SjCXHh4rpqDvHBLR4IKmI0zRiZ/Vq+H7Z39a6zXNyAT4PN/YCrX2q4VfljE4oJH3MC6+Vvjfg3tzRZkFshIJg0K4tOP1zqDAFj user@hostname my-comment",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSSHPubKey(tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Note: The following tests for SSH agent functions would require mocking or actual SSH agent setup
// For production tests, these would need to be integration tests or use dependency injection

func TestListSSHAgentIdentities(t *testing.T) {
	// This test will only work if SSH agent is actually available
	// In a real test environment, you might want to skip this test or mock it
	if !IsSSHAgentAvailable() {
		t.Skip("SSH agent not available, skipping test")
	}

	identities, err := ListSSHAgentIdentities()
	if err != nil {
		t.Logf("SSH agent error (expected in test environment): %v", err)
		t.Skip("SSH agent not accessible, skipping test")
	}

	// If we get here, we should have a valid result
	assert.NotNil(t, identities)
	// Note: We can't assert the exact content as it depends on the test environment
}

func TestIsSSHAgentIdentityValid(t *testing.T) {
	// This test will only work if SSH agent is actually available
	if !IsSSHAgentAvailable() {
		t.Skip("SSH agent not available, skipping test")
	}

	// Test with a non-existent identity
	valid, err := IsSSHAgentIdentityValid("non-existent-identity")
	if err != nil {
		t.Logf("SSH agent error (expected in test environment): %v", err)
		t.Skip("SSH agent not accessible, skipping test")
	}

	assert.False(t, valid)
}

func TestReadSSHAgentIdentityPublicKey(t *testing.T) {
	// This test will only work if SSH agent is actually available
	if !IsSSHAgentAvailable() {
		t.Skip("SSH agent not available, skipping test")
	}

	// Test with a non-existent identity
	key, err := ReadSSHAgentIdentityPublicKey("non-existent-identity")
	if err != nil {
		t.Logf("SSH agent error (expected in test environment): %v", err)
		t.Skip("SSH agent not accessible, skipping test")
	}

	assert.Empty(t, key)
}
