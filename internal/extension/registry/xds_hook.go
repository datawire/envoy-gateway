package registry

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/status"

	clusters "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	resourceTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	"github.com/envoyproxy/gateway/api/config/v1alpha1"
	"github.com/envoyproxy/gateway/internal/extension/types"
	"github.com/envoyproxy/gateway/internal/metrics"
	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

var _ types.XDSHookClient = (*XDSHook)(nil)

const (
	postHTTPListenerTranslation = "PostHTTPListenerTranslation"
)

type XDSHook struct {
	extensionId v1alpha1.ExtensionId
	grpcClient  extension.EnvoyGatewayExtensionClient
}

// PostHTTPListenerTranslation calls out to the extension hook to modify XDS after processing
// the HTTP listener. The method takes two arguments:
//
//   - httpListenerExtCtx: Contains the context necessary for the extension to be able to modify xDS resources
//
//   - xdsResourceVerTbl: the XDS Resource Table that the xDS server will server back to envoy. This table
//     will be modified by the method using the response from the extension service for the updated xDS resources.
func (h *XDSHook) PostHTTPListenerTranslation(
	httpListenerExtCtx *extension.HTTPListenerExtensionContext, xdsResourceVerTbl *xdsTypes.ResourceVersionTable,
) error {
	if httpListenerExtCtx == nil && xdsResourceVerTbl == nil {
		return fmt.Errorf("Both httpListenerExtCtx & xdsResourceVerTbl are nil")
	} else if httpListenerExtCtx == nil {
		return fmt.Errorf("httpListenerExtCtx is nil")
	} else if xdsResourceVerTbl == nil {
		return fmt.Errorf("xdsResourceVerTbl is nil")
	}

	xdsResources := xdsResourceVerTbl.GetXdsResources()

	listeners := xdsResources[resource.ListenerType]
	if listeners == nil {
		return nil
	}

	listn := findListenerByName(xdsResourceVerTbl, httpListenerExtCtx.ListenerName)
	if listn == nil {
		return fmt.Errorf("Cannot find listener with name %s in the xDS resource table", httpListenerExtCtx.ListenerName)
	}

	var routeCfg *route.RouteConfiguration
	routeCfgName := findXdsHTTPRouteConfigName(listn)
	if routeCfgName != "" {
		routeCfg = findXdsRouteConfig(xdsResourceVerTbl, routeCfgName)
		if routeCfg == nil {
			return fmt.Errorf("Listener with name %s has a filter chain that contains an http connection manager filter"+
				"specifying routes in RDS but no RDS route configuration was found in the xDS resource table",
				httpListenerExtCtx.ListenerName)
		}
	}

	ctx := context.Background()
	resp, err := h.callExtensionHook(ctx, httpListenerExtCtx, listn, routeCfg)
	if err != nil {
		return err
	}

	updateXdsTable(xdsResourceVerTbl, resp)

	return nil
}

func (h *XDSHook) callExtensionHook(
	ctx context.Context,
	eCtx *extension.HTTPListenerExtensionContext,
	listener *listener.Listener,
	routeCfg *route.RouteConfiguration) (*extension.PostHTTPListenerTranslationResponse, error) {

	req := &extension.PostHTTPListenerTranslationRequest{
		ExtensionContext: eCtx,
		Listener:         listener,
		RouteTable:       routeCfg,
	}
	metrics.ExtensionHookRequests.WithLabelValues(
		string(h.extensionId),
		string(types.XDSHook),
		postHTTPListenerTranslation,
	).Inc()
	start := time.Now()
	resp, err := h.grpcClient.PostHTTPListenerTranslation(ctx, req)
	metrics.ExtensionResponseLatency.WithLabelValues(
		string(h.extensionId),
		string(types.XDSHook),
		postHTTPListenerTranslation,
	).Observe(time.Since(start).Seconds())
	status, _ := status.FromError(err)
	metrics.ExtensionHookResponse.WithLabelValues(
		string(h.extensionId),
		string(types.XDSHook),
		postHTTPListenerTranslation,
		status.Code().String(),
	).Inc()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func updateXdsTable(xdsTable *xdsTypes.ResourceVersionTable, extResp *extension.PostHTTPListenerTranslationResponse) {
	if xdsTable == nil || extResp == nil {
		return
	}

	// We're assuming that Listener names are unique.
	xdsTable.AddOrReplaceXdsResource(resource.ListenerType, extResp.Listener, func(existing resourceTypes.Resource, new resourceTypes.Resource) bool {
		existingListener := existing.(*listener.Listener)
		newListener := new.(*listener.Listener)
		if newListener == nil || existingListener == nil {
			return false
		}
		if existingListener.Name == newListener.Name {
			return true
		}
		return false
	})

	// We're assuming that RouteConfiguration names are unique.
	xdsTable.AddOrReplaceXdsResource(resource.RouteType, extResp.RouteTable, func(existing resourceTypes.Resource, new resourceTypes.Resource) bool {
		existingRouteTable := existing.(*route.RouteConfiguration)
		newRouteTable := new.(*route.RouteConfiguration)
		if newRouteTable == nil || existingRouteTable == nil {
			return false
		}
		if existingRouteTable.Name == newRouteTable.Name {
			return true
		}
		return false
	})

	// We're assuming that Cluster names are unique.
	for _, cluster := range extResp.Clusters {
		xdsTable.AddOrReplaceXdsResource(resource.ClusterType, cluster, func(existing resourceTypes.Resource, new resourceTypes.Resource) bool {
			existingCluster := existing.(*clusters.Cluster)
			newCluster := new.(*clusters.Cluster)
			if newCluster == nil || existingCluster == nil {
				return false
			}
			if existingCluster.Name == newCluster.Name {
				return true
			}
			return false
		})
	}

	for _, secret := range extResp.Secrets {
		xdsTable.AddXdsResource(resource.SecretType, secret)
	}
}

// TODO: Consolidate the below helper functions into a shared package?

// findListenerByName finds a listener with the name and returns nil if there is no match
func findListenerByName(xdsRsrcTbl *xdsTypes.ResourceVersionTable, name string) *listener.Listener {
	if xdsRsrcTbl == nil || xdsRsrcTbl.XdsResources == nil || xdsRsrcTbl.XdsResources[resource.ListenerType] == nil {
		return nil
	}

	for _, l := range xdsRsrcTbl.XdsResources[resource.ListenerType] {
		listn := l.(*listener.Listener)
		if listn.Name == name {
			return listn
		}
	}
	return nil
}

// findXdsHTTPRouteConfigName finds the name of the route config associated with the
// http connection manager within the default filter chain and returns an empty string if
// not found.
//
// NOTE: this is wholesale copied from internal/xds/translator/listener.go:findXdsHTTPRouteConfigName
func findXdsHTTPRouteConfigName(xdsListener *listener.Listener) string {
	if xdsListener == nil || xdsListener.DefaultFilterChain == nil || xdsListener.DefaultFilterChain.Filters == nil {
		return ""
	}

	for _, filter := range xdsListener.DefaultFilterChain.Filters {
		if filter.Name == wellknown.HTTPConnectionManager {
			m := new(hcm.HttpConnectionManager)
			if err := filter.GetTypedConfig().UnmarshalTo(m); err != nil {
				return ""
			}
			rds := m.GetRds()
			if rds == nil {
				return ""
			}
			return rds.GetRouteConfigName()
		}
	}
	return ""
}

// findXdsRouteConfig finds an xds route with the name and returns nil if there is no match.
//
// NOTE: this is wholesale copied from internal/xds/translator/listener.go:findXdsRouteConfig
func findXdsRouteConfig(tCtx *xdsTypes.ResourceVersionTable, name string) *route.RouteConfiguration {
	if tCtx == nil || tCtx.XdsResources == nil || tCtx.XdsResources[resource.RouteType] == nil {
		return nil
	}

	for _, r := range tCtx.XdsResources[resource.RouteType] {
		route := r.(*route.RouteConfiguration)
		if route.Name == name {
			return route
		}
	}

	return nil
}
