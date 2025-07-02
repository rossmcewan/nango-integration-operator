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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SecretOrStringValue represents a value that can be either a string or a secret reference
type SecretOrStringValue struct {
	// String value (mutually exclusive with secretKeyRef)
	Value string `json:"value,omitempty"`
	// Secret reference (mutually exclusive with value)
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// NangoCredentials defines the OAuth credentials for the integration
type NangoCredentials struct {
	// Type of authentication (e.g., OAUTH1, OAUTH2)
	Type string `json:"type"`
	// Client ID for the OAuth application (either direct value or secret reference)
	ClientID SecretOrStringValue `json:"client_id"`
	// Client Secret for the OAuth application (either direct value or secret reference)
	ClientSecret SecretOrStringValue `json:"client_secret"`
	// Scopes required for the integration
	Scopes string `json:"scopes,omitempty"`
}

// NangoIntegrationSpec defines the desired state of NangoIntegration.
type NangoIntegrationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Unique key for the integration (required)
	UniqueKey string `json:"unique_key"`
	// Provider name (e.g., "slack", "github") (required)
	Provider string `json:"provider"`
	// Display name for the integration (required)
	DisplayName string `json:"display_name"`
	// OAuth credentials for the integration (required)
	Credentials NangoCredentials `json:"credentials"`
	// Nango API token for authentication (either direct value or secret reference)
	NangoToken SecretOrStringValue `json:"nango_token,omitempty"`
	// Nango API base URL (optional, defaults to https://api.nango.dev)
	NangoBaseURL string `json:"nango_base_url,omitempty"`
}

// NangoIntegrationStatus defines the observed state of NangoIntegration.
type NangoIntegrationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Integration ID from Nango API
	IntegrationID string `json:"integration_id,omitempty"`
	// Status of the integration (Created, Failed, etc.)
	Status string `json:"status,omitempty"`
	// Error message if creation failed
	ErrorMessage string `json:"error_message,omitempty"`
	// Last time the status was updated
	LastUpdated *metav1.Time `json:"last_updated,omitempty"`
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NangoIntegration is the Schema for the nangointegrations API.
type NangoIntegration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NangoIntegrationSpec   `json:"spec,omitempty"`
	Status NangoIntegrationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NangoIntegrationList contains a list of NangoIntegration.
type NangoIntegrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NangoIntegration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NangoIntegration{}, &NangoIntegrationList{})
}
