package provider

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/envoyproxy/gateway/api/config/v1alpha1"
	"github.com/envoyproxy/gateway/internal/envoygateway/config"
	"github.com/envoyproxy/gateway/internal/provider/kubernetes"
)

type ResourceTable = kubernetes.ResourceTable

func Start(svr *config.Server, k8sTable *ResourceTable) error {
	log := svr.Logger
	if svr.EnvoyGateway.Provider.Type == v1alpha1.ProviderTypeKubernetes {
		log.Info("Using provider", "type", v1alpha1.ProviderTypeKubernetes)
		cfg, err := ctrl.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get kubeconfig: %w", err)
		}
		provider, err := kubernetes.New(cfg, svr, k8sTable)
		if err != nil {
			return fmt.Errorf("failed to create provider %s", v1alpha1.ProviderTypeKubernetes)
		}
		if err := provider.Start(ctrl.SetupSignalHandler()); err != nil {
			return fmt.Errorf("failed to start provider %s", v1alpha1.ProviderTypeKubernetes)
		}
	}
	// Unsupported provider.
	return fmt.Errorf("unsupported provider type %v", svr.EnvoyGateway.Provider.Type)
}
