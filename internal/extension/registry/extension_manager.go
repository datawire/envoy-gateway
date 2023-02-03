package registry

import (
	"context"
	"crypto/x509"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	corev1 "k8s.io/api/core/v1"
	k8smachinery "k8s.io/apimachinery/pkg/types"
	k8scli "sigs.k8s.io/controller-runtime/pkg/client"
	k8sclicfg "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/envoyproxy/gateway/api/config/v1alpha1"
	"github.com/envoyproxy/gateway/internal/envoygateway"
	"github.com/envoyproxy/gateway/internal/envoygateway/config"
	"github.com/envoyproxy/gateway/internal/extension/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

type hookConn struct {
	extensionID v1alpha1.ExtensionId
	hookType    types.HookType
}

var _ types.Manager = (*Manager)(nil)

type Manager struct {
	k8sClient     k8scli.Client
	namespace     string
	extensions    map[v1alpha1.ExtensionId]v1alpha1.Extension
	hookConnCache map[hookConn]*grpc.ClientConn
}

// NewManager returns a new Manager
func NewManager(cfg *config.Server) (types.Manager, error) {
	cli, err := k8scli.New(k8sclicfg.GetConfigOrDie(), k8scli.Options{Scheme: envoygateway.GetScheme()})
	if err != nil {
		return nil, err
	}

	// To keep things simple, we're just going to allow only one extension since it gets messy trying to reconcile
	// xDS from multiple extensions even though the Manager can handle multiple extensions.
	//
	// TODO: have the extensions in EnvoyConfig be a data structure and not a list?
	if len(cfg.EnvoyGateway.Extensions) > 1 {
		return nil, fmt.Errorf("More than 1 extension is registered. For the time being, we allow up to 1 extension to be registered")
	}

	extTable := make(map[v1alpha1.ExtensionId]v1alpha1.Extension)
	for _, extension := range cfg.EnvoyGateway.Extensions {
		extTable[extension.Name] = *extension
	}

	hookConnCache := make(map[hookConn]*grpc.ClientConn)

	return &Manager{
		k8sClient:     cli,
		namespace:     cfg.Namespace,
		extensions:    extTable,
		hookConnCache: hookConnCache,
	}, nil
}

// HasExtension checks to see whether a given Group and Kind has an
// associated extension registered for it.
//
// TODO: Add Kind to the Extension config in the EnvoyGateway config.
// For now, assume that Kind is always valid and we'll just validate for the API Group.
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

// GetXDSHookClient returns an XDS Hook Client for an extension with ID extID
func (m *Manager) GetXDSHookClient(extID v1alpha1.ExtensionId) (types.XDSHookClient, error) {
	ctx := context.Background()

	ext, ok := m.extensions[extID]
	if !ok {
		return nil, fmt.Errorf("extension %s not found", extID)
	}

	if ext.Service == nil {
		return nil, fmt.Errorf("Extension service config is nil")
	}

	key := hookConn{
		extensionID: ext.Name,
		hookType:    types.XDSHook,
	}
	conn, cached := m.hookConnCache[key]
	if !cached {
		serverAddr := fmt.Sprintf("%s:%d", ext.Service.Host, ext.Service.Port)

		opts, err := setupGRPCOpts(ctx, m.k8sClient, &ext, m.namespace)
		if err != nil {
			return nil, err
		}

		conn, err = grpc.Dial(serverAddr, opts...)
		if err != nil {
			return nil, err
		}

		m.hookConnCache[key] = conn
	}

	client := extension.NewEnvoyGatewayExtensionClient(conn)
	xdsHookClient := &XDSHook{
		grpcClient: client,
	}
	return xdsHookClient, nil
}

func (m *Manager) CleanupHookConns() {
	for _, conn := range m.hookConnCache {
		conn.Close()
	}
	m.hookConnCache = make(map[hookConn]*grpc.ClientConn)
}

func parseCA(caSecret *corev1.Secret) (*x509.CertPool, error) {
	caCertPEMBytes, ok := caSecret.Data[corev1.TLSCertKey]
	if !ok {
		return nil, fmt.Errorf("no cert found in CA secret!")
	}
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM(caCertPEMBytes); !ok {
		return nil, fmt.Errorf("failed to append certificates")
	}
	return cp, nil
}

func setupGRPCOpts(ctx context.Context, client k8scli.Client, ext *v1alpha1.Extension, namespace string) ([]grpc.DialOption, error) {
	if ext == nil {
		return nil, fmt.Errorf("Extension config is nil")
	}

	if ext.Service == nil {
		return nil, fmt.Errorf("Extension %s doesn't have a service config", ext.Name)
	}

	var opts []grpc.DialOption
	var creds credentials.TransportCredentials
	if ext.Service.TLS != nil {
		switch ext.Service.TLS.Type {
		case v1alpha1.TLSTypeSecret:
			secret := &corev1.Secret{}
			secretNamespace := namespace
			if ext.Service.TLS.Secret.Namespace != "" {
				secretNamespace = ext.Service.TLS.Secret.Namespace
			}
			key := k8smachinery.NamespacedName{
				Namespace: secretNamespace,
				Name:      ext.Service.TLS.Secret.Name,
			}
			if err := client.Get(ctx, key, secret); err != nil {
				return nil, fmt.Errorf("Cannot find TLS Secret %s in namespace %s", ext.Service.TLS.Secret.Name, secretNamespace)
			}
			cp, err := parseCA(secret)
			if err != nil {
				return nil, fmt.Errorf("Error parsing cert in Secret %s in namespace %s", ext.Service.TLS.Secret.Name, secretNamespace)
			}
			creds = credentials.NewClientTLSFromCert(cp, "")
		default:
			return nil, fmt.Errorf("unsupported TLS Type %s", ext.Service.TLS.Type)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return opts, nil
}
