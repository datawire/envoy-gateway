package registry

import (
	"context"
	"fmt"

	clusters "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	resourceTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	"github.com/envoyproxy/gateway/internal/extension/types"
	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

var _ types.XDSHookClient = (*XDSHook)(nil)

type XDSHook struct {
	grpcClient extension.EnvoyGatewayExtensionClient
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
	switch {
	case httpListenerExtCtx == nil && xdsResourceVerTbl == nil:
		return fmt.Errorf("Both httpListenerExtCtx & xdsResourceVerTbl are nil")
	case httpListenerExtCtx == nil:
		return fmt.Errorf("httpListenerExtCtx is nil")
	case xdsResourceVerTbl == nil:
		return fmt.Errorf("xdsResourceVerTbl is nil")
	}

	xdsResources := xdsResourceVerTbl.GetXdsResources()

	listeners := xdsResources[resource.ListenerType]
	if listeners == nil {
		return nil
	}

	ctx := context.Background()
	listn := findListenerByName(xdsResourceVerTbl, httpListenerExtCtx.ListenerName)
	resp, err := callExtensionHook(ctx, h.grpcClient, httpListenerExtCtx, listn, xdsResourceVerTbl)
	if err != nil {
		return err
	}

	updateXdsTable(xdsResourceVerTbl, resp)

	return nil
}

func callExtensionHook(
	ctx context.Context,
	client extension.EnvoyGatewayExtensionClient,
	eCtx *extension.HTTPListenerExtensionContext,
	listener *listener.Listener,
	xdsResourceTbl *xdsTypes.ResourceVersionTable) (*extension.PostHTTPListenerTranslationResponse, error) {

	req := &extension.PostHTTPListenerTranslationRequest{
		ExtensionContext: eCtx,
		Listener:         listener,
	}
	routeCfgName := findXdsHTTPRouteConfigName(listener)
	if routeCfgName != "" {
		routeCfg := findXdsRouteConfig(xdsResourceTbl, routeCfgName)
		if routeCfg == nil {
			// If the RouteConfiguration cannot be found then it means
			// it wasn't found when intially processing the HTTPListener and an error
			// should have been thrown then. So if we ever arrive at this point, then it
			// means that this is a programming error.
			panic("Cannot find xds route config")
		}
		req.RouteTable = routeCfg
	}

	// TODO: add retries, metrics, logging, all that good stuff
	resp, err := client.PostHTTPListenerTranslation(ctx, req)
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
		newRouteTale := new.(*route.RouteConfiguration)
		if newRouteTale == nil || existingRouteTable == nil {
			return false
		}
		if existingRouteTable.Name == newRouteTale.Name {
			return true
		}
		return false
	})

	// We're assuming that Cluster names are unique.
	for _, cluster := range extResp.Clusters {
		xdsTable.AddOrReplaceXdsResource(resource.RouteType, cluster, func(existing resourceTypes.Resource, new resourceTypes.Resource) bool {
			existingCluster := existing.(*clusters.Cluster)
			newCluster := new.(*route.RouteConfiguration)
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

// findListenerByName finds a listener with the name and returns nil if there is no match
//
// TODO: Consolidate these helper functions?
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
// TODO: Consolidate these helper functions?
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
// TODO: Consolidate these helper functions?
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
