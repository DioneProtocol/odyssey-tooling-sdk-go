// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestNode_PullDockerImage(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		image       string
		expectError bool
	}{
		{
			name:        "Valid image name",
			image:       "nginx:latest",
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty image name",
			image:       "",
			expectError: true,
		},
		{
			name:        "Invalid image format",
			image:       "invalid:",
			expectError: true,
		},
		{
			name:        "Image with tag",
			image:       "redis:alpine",
			expectError: true,
		},
		{
			name:        "Image without tag",
			image:       "postgres",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.PullDockerImage(tt.image)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_PullDockerImage_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	err := node.PullDockerImage("nginx:latest")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
}

func TestNode_DockerLocalImageExists(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		image       string
		expectError bool
	}{
		{
			name:        "Valid image name",
			image:       "nginx:latest",
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty image name",
			image:       "",
			expectError: true,
		},
		{
			name:        "Image with special characters",
			image:       "my-registry.com/namespace/image:tag",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			exists, err := node.DockerLocalImageExists(tt.image)

			if tt.expectError {
				assert.Error(t, err)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_DockerLocalImageExists_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	exists, err := node.DockerLocalImageExists("nginx:latest")
	assert.Error(t, err)
	assert.False(t, exists)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
}

func TestParseDockerImageListOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []string
	}{
		{
			name:     "Multiple images",
			input:    []byte("nginx:latest\nredis:alpine\npostgres:13"),
			expected: []string{"nginx:latest", "redis:alpine", "postgres:13"},
		},
		{
			name:     "Single image",
			input:    []byte("nginx:latest"),
			expected: []string{"nginx:latest"},
		},
		{
			name:     "Empty input",
			input:    []byte(""),
			expected: []string{""},
		},
		{
			name:     "With empty lines",
			input:    []byte("nginx:latest\n\nredis:alpine\n"),
			expected: []string{"nginx:latest", "", "redis:alpine", ""},
		},
		{
			name:     "With trailing newline",
			input:    []byte("nginx:latest\nredis:alpine\n"),
			expected: []string{"nginx:latest", "redis:alpine", ""},
		},
		{
			name:     "Images with different formats",
			input:    []byte("nginx:latest\nmy-registry.com/namespace/image:tag\nredis"),
			expected: []string{"nginx:latest", "my-registry.com/namespace/image:tag", "redis"},
		},
		{
			name:     "Only newlines",
			input:    []byte("\n\n\n"),
			expected: []string{"", "", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDockerImageListOutput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_BuildDockerImage(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		image       string
		path        string
		dockerfile  string
		expectError bool
	}{
		{
			name:        "Valid parameters",
			image:       "test-image:latest",
			path:        "/tmp/build",
			dockerfile:  "Dockerfile",
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty image name",
			image:       "",
			path:        "/tmp/build",
			dockerfile:  "Dockerfile",
			expectError: true,
		},
		{
			name:        "Empty path",
			image:       "test-image:latest",
			path:        "",
			dockerfile:  "Dockerfile",
			expectError: true,
		},
		{
			name:        "Empty dockerfile",
			image:       "test-image:latest",
			path:        "/tmp/build",
			dockerfile:  "",
			expectError: true,
		},
		{
			name:        "Custom dockerfile name",
			image:       "test-image:latest",
			path:        "/tmp/build",
			dockerfile:  "CustomDockerfile",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.BuildDockerImage(tt.image, tt.path, tt.dockerfile)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_BuildDockerImage_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	err := node.BuildDockerImage("test-image:latest", "/tmp/build", "Dockerfile")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
}

func TestNode_BuildDockerImageFromGitRepo(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		image       string
		gitRepo     string
		commit      string
		expectError bool
	}{
		{
			name:        "Valid parameters with commit",
			image:       "test-image:latest",
			gitRepo:     "https://github.com/test/repo.git",
			commit:      "abc123",
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Valid parameters without commit (should use HEAD)",
			image:       "test-image:latest",
			gitRepo:     "https://github.com/test/repo.git",
			commit:      "",
			expectError: true,
		},
		{
			name:        "Empty image name",
			image:       "",
			gitRepo:     "https://github.com/test/repo.git",
			commit:      "abc123",
			expectError: true,
		},
		{
			name:        "Empty git repo",
			image:       "test-image:latest",
			gitRepo:     "",
			commit:      "abc123",
			expectError: true,
		},
		{
			name:        "Invalid git repo URL",
			image:       "test-image:latest",
			gitRepo:     "not-a-valid-url",
			commit:      "abc123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.BuildDockerImageFromGitRepo(tt.image, tt.gitRepo, tt.commit)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_BuildDockerImageFromGitRepo_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	err := node.BuildDockerImageFromGitRepo("test-image:latest", "https://github.com/test/repo.git", "abc123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
}

func TestNode_PrepareDockerImageWithRepo(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		image       string
		gitRepo     string
		commit      string
		expectError bool
	}{
		{
			name:        "Valid parameters",
			image:       "test-image:latest",
			gitRepo:     "https://github.com/test/repo.git",
			commit:      "abc123",
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty image name",
			image:       "",
			gitRepo:     "https://github.com/test/repo.git",
			commit:      "abc123",
			expectError: true,
		},
		{
			name:        "Empty git repo",
			image:       "test-image:latest",
			gitRepo:     "",
			commit:      "abc123",
			expectError: true,
		},
		{
			name:        "Image with tag",
			image:       "nginx:1.21",
			gitRepo:     "https://github.com/nginx/nginx.git",
			commit:      "main",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.PrepareDockerImageWithRepo(tt.image, tt.gitRepo, tt.commit)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_PrepareDockerImageWithRepo_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with Docker support disabled
	constants.DockerSupportEnabled = false

	err := node.PrepareDockerImageWithRepo("test-image:latest", "https://github.com/test/repo.git", "abc123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Docker support functionality is disabled")
}

func TestNode_DockerImageOperations_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "PullDockerImage with special characters",
			operation: func() error {
				node := &Node{NodeID: "test-node", IP: "192.168.1.1"}
				return node.PullDockerImage("my-registry.com/namespace/image:tag-with-dashes")
			},
			expectError: true,
		},
		{
			name: "BuildDockerImage with relative path",
			operation: func() error {
				node := &Node{NodeID: "test-node", IP: "192.168.1.1"}
				return node.BuildDockerImage("test:latest", "./relative/path", "Dockerfile")
			},
			expectError: true,
		},
		{
			name: "BuildDockerImageFromGitRepo with branch name",
			operation: func() error {
				node := &Node{NodeID: "test-node", IP: "192.168.1.1"}
				return node.BuildDockerImageFromGitRepo("test:latest", "https://github.com/test/repo.git", "feature-branch")
			},
			expectError: true,
		},
		{
			name: "PrepareDockerImageWithRepo with tag name",
			operation: func() error {
				node := &Node{NodeID: "test-node", IP: "192.168.1.1"}
				return node.PrepareDockerImageWithRepo("test:latest", "https://github.com/test/repo.git", "v1.0.0")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_DockerImageOperations_EmptyNode(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	emptyNode := &Node{}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "PullDockerImage on empty node",
			operation: func() error {
				return emptyNode.PullDockerImage("nginx:latest")
			},
			expectError: true,
		},
		{
			name: "DockerLocalImageExists on empty node",
			operation: func() error {
				_, err := emptyNode.DockerLocalImageExists("nginx:latest")
				return err
			},
			expectError: true,
		},
		{
			name: "BuildDockerImage on empty node",
			operation: func() error {
				return emptyNode.BuildDockerImage("test:latest", "/tmp", "Dockerfile")
			},
			expectError: true,
		},
		{
			name: "BuildDockerImageFromGitRepo on empty node",
			operation: func() error {
				return emptyNode.BuildDockerImageFromGitRepo("test:latest", "https://github.com/test/repo.git", "main")
			},
			expectError: true,
		},
		{
			name: "PrepareDockerImageWithRepo on empty node",
			operation: func() error {
				return emptyNode.PrepareDockerImageWithRepo("test:latest", "https://github.com/test/repo.git", "main")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
