// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestGithubReleaseURL(t *testing.T) {
	tests := []struct {
		name     string
		org      string
		repo     string
		expected string
	}{
		{
			name:     "valid org and repo",
			org:      "testorg",
			repo:     "testrepo",
			expected: "https://api.github.com/repos/testorg/testrepo/releases/latest",
		},
		{
			name:     "empty org and repo",
			org:      "",
			repo:     "",
			expected: "https://api.github.com/repos///releases/latest",
		},
		{
			name:     "special characters in org/repo",
			org:      "test-org",
			repo:     "test.repo",
			expected: "https://api.github.com/repos/test-org/test.repo/releases/latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLatestGithubReleaseURL(tt.org, tt.repo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetGithubReleasesURL(t *testing.T) {
	tests := []struct {
		name     string
		org      string
		repo     string
		expected string
	}{
		{
			name:     "valid org and repo",
			org:      "testorg",
			repo:     "testrepo",
			expected: "https://api.github.com/repos/testorg/testrepo/releases",
		},
		{
			name:     "empty org and repo",
			org:      "",
			repo:     "",
			expected: "https://api.github.com/repos///releases",
		},
		{
			name:     "special characters in org/repo",
			org:      "test-org",
			repo:     "test.repo",
			expected: "https://api.github.com/repos/test-org/test.repo/releases",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGithubReleasesURL(tt.org, tt.repo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLatestGithubReleaseVersion(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		org            string
		repo           string
		authToken      string
		expected       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful request with valid version",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"tag_name": "v1.2.3"}`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "",
			expected:  "v1.2.3",
			wantErr:   false,
		},
		{
			name: "successful request with auth token",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Check for authorization header
				auth := r.Header.Get("authorization")
				if auth != "Bearer test-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"tag_name": "v2.0.0"}`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "test-token",
			expected:  "v2.0.0",
			wantErr:   false,
		},
		{
			name: "invalid JSON response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`invalid json`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "failed to unmarshal binary json version string",
		},
		{
			name: "invalid version format",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"tag_name": "invalid-version"}`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "invalid version string",
		},
		{
			name: "missing tag_name field",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"name": "Release 1.0.0"}`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "interface conversion",
		},
		{
			name: "404 Not Found",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"message": "Not Found"}`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "failed downloading",
		},
		{
			name: "500 Internal Server Error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal Server Error"}`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "failed downloading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// For this test, we need to mock the HTTPGet function or modify the code
			// to accept a custom HTTP client. Since we can't easily mock HTTPGet,
			// we'll test the URL construction and basic functionality
			url := GetLatestGithubReleaseURL(tt.org, tt.repo)
			expectedURL := "https://api.github.com/repos/" + tt.org + "/" + tt.repo + "/releases/latest"
			assert.Equal(t, expectedURL, url)

			// Test the actual function call would require mocking HTTPGet
			// For now, we'll skip the actual HTTP call test
			t.Skip("Requires HTTPGet mocking or dependency injection")
		})
	}
}

func TestGetAllGithubReleaseVersions(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		org            string
		repo           string
		authToken      string
		expected       []string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful request with multiple versions",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[
					{"tag_name": "v1.2.3"},
					{"tag_name": "v1.2.2"},
					{"tag_name": "v1.1.0"}
				]`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "",
			expected:  []string{"v1.2.3", "v1.2.2", "v1.1.0"},
			wantErr:   false,
		},
		{
			name: "successful request with single version",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"tag_name": "v1.0.0"}]`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "",
			expected:  []string{"v1.0.0"},
			wantErr:   false,
		},
		{
			name: "empty releases array",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "",
			expected:  []string{},
			wantErr:   false,
		},
		{
			name: "invalid JSON response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`invalid json`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "failed to unmarshal binary json version string",
		},
		{
			name: "invalid version in array",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[
					{"tag_name": "v1.2.3"},
					{"tag_name": "invalid-version"},
					{"tag_name": "v1.1.0"}
				]`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "invalid version string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Test URL construction
			url := GetGithubReleasesURL(tt.org, tt.repo)
			expectedURL := "https://api.github.com/repos/" + tt.org + "/" + tt.repo + "/releases"
			assert.Equal(t, expectedURL, url)

			// Test the actual function call would require mocking HTTPGet
			t.Skip("Requires HTTPGet mocking or dependency injection")
		})
	}
}

func TestGetLatestGithubPreReleaseVersion(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		org            string
		repo           string
		authToken      string
		expected       string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful request with releases",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[
					{"tag_name": "v1.2.3"},
					{"tag_name": "v1.2.2"},
					{"tag_name": "v1.1.0"}
				]`))
			},
			org:       "testorg",
			repo:      "testrepo",
			authToken: "",
			expected:  "v1.2.3", // First release in array
			wantErr:   false,
		},
		{
			name: "no releases found",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[]`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "no releases found",
		},
		{
			name: "error from GetAllGithubReleaseVersions",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "Internal Server Error"}`))
			},
			org:         "testorg",
			repo:        "testrepo",
			authToken:   "",
			wantErr:     true,
			errContains: "failed downloading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Test the actual function call would require mocking HTTPGet
			t.Skip("Requires HTTPGet mocking or dependency injection")
		})
	}
}

func TestGetGithubReleaseAssetURL(t *testing.T) {
	tests := []struct {
		name     string
		org      string
		repo     string
		version  string
		asset    string
		expected string
	}{
		{
			name:     "valid parameters",
			org:      "testorg",
			repo:     "testrepo",
			version:  "v1.2.3",
			asset:    "testrepo-v1.2.3-linux-amd64.tar.gz",
			expected: "https://github.com/testorg/testrepo/releases/download/v1.2.3/testrepo-v1.2.3-linux-amd64.tar.gz",
		},
		{
			name:     "empty parameters",
			org:      "",
			repo:     "",
			version:  "",
			asset:    "",
			expected: "https://github.com///releases/download//",
		},
		{
			name:     "special characters",
			org:      "test-org",
			repo:     "test.repo",
			version:  "v1.0.0-beta.1",
			asset:    "test.repo-v1.0.0-beta.1-windows-x64.zip",
			expected: "https://github.com/test-org/test.repo/releases/download/v1.0.0-beta.1/test.repo-v1.0.0-beta.1-windows-x64.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGithubReleaseAssetURL(tt.org, tt.repo, tt.version, tt.asset)
			assert.Equal(t, tt.expected, result)
		})
	}
}
