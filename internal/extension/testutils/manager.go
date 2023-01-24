package testutils

import (
	"sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/envoyproxy/gateway/api/config/v1alpha1"
	"github.com/envoyproxy/gateway/internal/extension/types"
)

var _ types.Manager = (*Manager)(nil)

type Manager struct {
	extensions []types.Extension
}

func NewManager(ext ...types.Extension) types.Manager {
	return &Manager{
		extensions: ext,
	}
}

func (m *Manager) HasExtension(g v1beta1.Group, k v1beta1.Kind) (bool, *v1alpha1.ExtensionId) {
	for _, ext := range m.extensions {
		for _, apiGroup := range ext.APIGroups {
			if g == apiGroup {
				extId := ext.Name
				return true, &extId
			}
		}
	}
	return false, nil
}

func (m *Manager) GetXDSHookClient(extID v1alpha1.ExtensionId) (types.XDSHookClient, error) {
	return &XDSHookClient{}, nil
}
