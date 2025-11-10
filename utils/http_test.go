// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPGet(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		url            string
		authToken      string
		expectedBody   string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful GET request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "success"}`))
			},
			url:          "/test",
			authToken:    "",
			expectedBody: `{"message": "success"}`,
			wantErr:      false,
		},
		{
			name: "successful GET request with auth token",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Check for authorization header
				auth := r.Header.Get("authorization")
				if auth != "Bearer test-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "authenticated"}`))
			},
			url:          "/test",
			authToken:    "test-token",
			expectedBody: `{"message": "authenticated"}`,
			wantErr:      false,
		},
		{
			name: "GET request without auth token",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Check that no authorization header is present
				auth := r.Header.Get("authorization")
				if auth != "" {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error": "unexpected auth header"}`))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "no auth"}`))
			},
			url:          "/test",
			authToken:    "",
			expectedBody: `{"message": "no auth"}`,
			wantErr:      false,
		},
		{
			name: "404 Not Found",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "not found"}`))
			},
			url:         "/notfound",
			authToken:   "",
			wantErr:     true,
			errContains: "unexpected http status code: 404",
		},
		{
			name: "500 Internal Server Error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "internal server error"}`))
			},
			url:         "/error",
			authToken:   "",
			wantErr:     true,
			errContains: "unexpected http status code: 500",
		},
		{
			name: "empty response body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// No body written
			},
			url:          "/empty",
			authToken:    "",
			expectedBody: "",
			wantErr:      false,
		},
		{
			name: "large response body",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				// Write a large response
				largeData := make([]byte, 1024*1024) // 1MB
				for i := range largeData {
					largeData[i] = byte(i % 256)
				}
				w.Write(largeData)
			},
			url:       "/large",
			authToken: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Make the request
			url := server.URL + tt.url
			result, err := HTTPGet(url, tt.authToken)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, string(result))
			} else {
				// For tests where we don't specify expected body, just check it's not nil
				assert.NotNil(t, result)
			}
		})
	}

	// Test with invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		result, err := HTTPGet("invalid://url", "")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed downloading")
	})

	// Test with network timeout (this would require a more complex setup)
	t.Run("network timeout", func(t *testing.T) {
		// This test would require setting up a server that doesn't respond
		// or modifying the HTTP client timeout, which isn't easily testable
		// without modifying the function to accept a custom HTTP client
		t.Skip("Requires HTTP client configuration or dependency injection")
	})
}
