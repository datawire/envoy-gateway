// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

// APIGroup defines a K8s API group (e.g gateway.envoyproxy.io)
type APIGroup string

// ProviderType defines the types of providers supported by Envoy Gateway.
//
// +kubebuilder:validation:Enum=Kubernetes
type ProviderType string

const (
	// ProviderTypeKubernetes defines the "Kubernetes" provider.
	ProviderTypeKubernetes ProviderType = "Kubernetes"

	// ProviderTypeFile defines the "File" provider.
	ProviderTypeFile ProviderType = "File"
)

// TLSType defines the types where TLS certificates can be loaded.
//
// +kubebuilder:validation:Enum=Secret
type TLSType string

const (
	// TLSTypeSecret defines the "Secret" TLS type.
	TLSTypeSecret TLSType = "Secret"

	// TLSTypeFile defines the "File" TLS type.
	TLSTypeFile TLSType = "File"
)

type TLSSecret struct {
	// Name is the secret name to load the TLS certificate from
	Name string `json:"name"`

	// Namespace is the namespace where the secret is located.
	//
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type TLSFile struct {
	// TODO: Add config
}

// KubernetesDeploymentSpec defines the desired state of the Kubernetes deployment resource.
type KubernetesDeploymentSpec struct {
	// Replicas is the number of desired pods. Defaults to 1.
	//
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// TODO: Expose config as use cases are better understood, e.g. labels.
}
