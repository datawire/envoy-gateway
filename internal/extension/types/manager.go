package types

import (
	"sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/envoyproxy/gateway/api/config/v1alpha1"
)

// Extension represents an extension registered by the Manager.
type Extension struct {
	v1alpha1.Extension

	SupportedHookTypes []HookType
}

// Manager handles and maintains registered extensions and returns clients for
// different Hook types.
type Manager interface {
	// HasExtension checks to see whether a given Group and Kind has an
	// associated extension registered for it.
	//
	// If a Group and Kind is registered with an extension, then it should
	// return true and the extemnsion ID, otherwise return (false, nil).
	HasExtension(g v1beta1.Group, k v1beta1.Kind) (bool *v1alpha1.ExtensionId)

	// GetXDSHookClient returns an hook client of type hookType for a given extension ID.
	//
	// If the extension does not support this hook, then it should return
	// (nil, error)
	GetXDSHookClient(extID v1alpha1.ExtensionId) (*XDSHookClient, error)
}
