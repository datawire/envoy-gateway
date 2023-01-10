// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KindEnvoyGateway is the name of the EnvoyGateway kind.
	KindEnvoyGateway = "EnvoyGateway"
	// GatewayControllerName is the name of the GatewayClass controller.
	GatewayControllerName = "gateway.envoyproxy.io/gatewayclass-controller"
)

//+kubebuilder:object:root=true

// EnvoyGateway is the Schema for the envoygateways API.
type EnvoyGateway struct {
	metav1.TypeMeta `json:",inline"`

	// EnvoyGatewaySpec defines the desired state of Envoy Gateway.
	EnvoyGatewaySpec `json:",inline"`
}

// EnvoyGatewaySpec defines the desired state of Envoy Gateway.
type EnvoyGatewaySpec struct {
	// Gateway defines desired Gateway API specific configuration. If unset,
	// default configuration parameters will apply.
	//
	// +optional
	Gateway *Gateway `json:"gateway,omitempty"`

	// Provider defines the desired provider and provider-specific configuration.
	// If unspecified, the Kubernetes provider is used with default configuration
	// parameters.
	//
	// +optional
	Provider *Provider `json:"provider,omitempty"`

	// Extensions defines the list of extensions for the Envoy Gateway Control Plane.
	//
	// +optional
	Extensions []*Extension `json:"extensions,omitempty"`
}

// Gateway defines the desired Gateway API configuration of Envoy Gateway.
type Gateway struct {
	// ControllerName defines the name of the Gateway API controller. If unspecified,
	// defaults to "gateway.envoyproxy.io/gatewayclass-controller". See the following
	// for additional details:
	//
	// https://gateway-api.sigs.k8s.io/v1alpha2/references/spec/#gateway.networking.k8s.io/v1alpha2.GatewayClass
	//
	// +optional
	ControllerName string `json:"controllerName,omitempty"`
}

// Provider defines the desired configuration of a provider.
// +union
type Provider struct {
	// Type is the type of provider to use. Supported types are:
	//
	//   * Kubernetes: A provider that provides runtime configuration via the Kubernetes API.
	//
	// +unionDiscriminator
	Type ProviderType `json:"type"`
	// Kubernetes defines the configuration of the Kubernetes provider. Kubernetes
	// provides runtime configuration via the Kubernetes API.
	//
	// +optional
	Kubernetes *KubernetesProvider `json:"kubernetes,omitempty"`

	// File defines the configuration of the File provider. File provides runtime
	// configuration defined by one or more files.
	//
	// +optional
	File *FileProvider `json:"file,omitempty"`
}

// KubernetesProvider defines configuration for the Kubernetes provider.
type KubernetesProvider struct {
	// TODO: Add config as use cases are better understood.
}

// FileProvider defines configuration for the File provider.
type FileProvider struct {
	// TODO: Add config as use cases are better understood.
}

type Extension struct {
	// Name defines the name to register with for the extension.
	Name string

	// APIGroups defines the set of K8s api groups the extension will handle.
	APIGroups []APIGroup

	// Service defines the configuration of the extension service that the Envoy
	// Gateway Control Plane will call through extension hooks.
	Service *ExtensionService
}

type ExtensionService struct {
	// Host define the extension service hostname.
	Host string `json:"host"`

	// Port defines the port the extension service is exposed on.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=80
	Port int32 `json:"port,omitempty"`

	// TLS defines TLS configuration for communication between Envoy Gateway and
	// the extension service.
	//
	// +optional
	TLS *ExtensionTLS `json:"tls,omitempty"`
}

type ExtensionTLS struct {
	// Type is the method for how the TLS certificate is loaded. Supported types are:
	//
	//   * Secret: Load the TLS certificate from a K8s secret.
	//
	// +unionDiscriminator
	Type TLSType `json:"type"`

	// Secret defines which K8s secret to load the TLS certificate from.
	//
	// +optional
	Secret *TLSSecret `json:"secret,omitempty"`

	// File defines the configuration for loading the TLS certificate from the filesystem.
	//
	// +optional
	File *TLSFile `json:"file,omitempty"`
}

func init() {
	SchemeBuilder.Register(&EnvoyGateway{})
}
