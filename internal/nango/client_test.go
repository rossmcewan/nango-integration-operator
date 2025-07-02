/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nango

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		token    string
		expected string
	}{
		{
			name:     "with custom base URL",
			baseURL:  "https://custom.nango.dev",
			token:    "test-token",
			expected: "https://custom.nango.dev",
		},
		{
			name:     "with empty base URL",
			baseURL:  "",
			token:    "test-token",
			expected: "https://api.nango.dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.token)
			if client.baseURL != tt.expected {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.expected)
			}
			if client.token != tt.token {
				t.Errorf("NewClient() token = %v, want %v", client.token, tt.token)
			}
		})
	}
}

func TestCreateIntegrationRequest(t *testing.T) {
	req := CreateIntegrationRequest{
		UniqueKey:   "test-integration",
		Provider:    "slack",
		DisplayName: "Test Integration",
		Credentials: NangoCredentials{
			Type:         "OAUTH1",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Scopes:       "chat:write",
		},
	}

	if req.UniqueKey != "test-integration" {
		t.Errorf("Expected UniqueKey to be 'test-integration', got %s", req.UniqueKey)
	}

	if req.Provider != "slack" {
		t.Errorf("Expected Provider to be 'slack', got %s", req.Provider)
	}

	if req.Credentials.Type != "OAUTH1" {
		t.Errorf("Expected Credentials.Type to be 'OAUTH1', got %s", req.Credentials.Type)
	}
}
