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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents a Nango API client
type Client struct {
	httpClient *resty.Client
	baseURL    string
	token      string
}

// CreateIntegrationRequest represents the request body for creating an integration
type CreateIntegrationRequest struct {
	UniqueKey   string           `json:"unique_key"`
	Provider    string           `json:"provider"`
	DisplayName string           `json:"display_name"`
	Credentials NangoCredentials `json:"credentials"`
}

// NangoCredentials represents the OAuth credentials
type NangoCredentials struct {
	Type         string `json:"type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scopes       string `json:"scopes,omitempty"`
}

// CreateIntegrationResponse represents the response from creating an integration
type CreateIntegrationResponse struct {
	Data IntegrationData `json:"data"`
}

// IntegrationData represents the integration data returned by the API
type IntegrationData struct {
	UniqueKey   string    `json:"unique_key"`
	DisplayName string    `json:"display_name"`
	Provider    string    `json:"provider"`
	Logo        string    `json:"logo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewClient creates a new Nango API client
func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "https://api.nango.dev"
	}

	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetTimeout(30 * time.Second)

	return &Client{
		httpClient: client,
		baseURL:    baseURL,
		token:      token,
	}
}

// CreateIntegration creates a new integration via the Nango API
func (c *Client) CreateIntegration(req CreateIntegrationRequest) (*CreateIntegrationResponse, error) {
	resp, err := c.httpClient.R().
		SetBody(req).
		Post("/integrations")

	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to create integration: status %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var response CreateIntegrationResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetIntegration retrieves an integration by unique key
func (c *Client) GetIntegration(uniqueKey string) (*CreateIntegrationResponse, error) {
	resp, err := c.httpClient.R().
		Get(fmt.Sprintf("/integrations/%s", uniqueKey))

	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("integration not found")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get integration: status %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var response CreateIntegrationResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
