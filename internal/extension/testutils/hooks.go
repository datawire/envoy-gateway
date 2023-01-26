package testutils

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

type XDSHookClient struct{}

// PostHTTPListenerTranslation modifies xdsResourceVerTbl. It adds extra response headers using the context inputs.
func (c *XDSHookClient) PostHTTPListenerTranslation(
	httpListenerExtCtx *extension.HTTPListenerExtensionContext, xdsResourceVerTbl *xdsTypes.ResourceVersionTable,
) error {
	listenerName := httpListenerExtCtx.ListenerName
	routeCfg := getRouteCfgByName(xdsResourceVerTbl, listenerName)
	if routeCfg == nil {
		// NOTE: this could be a programming error
		return fmt.Errorf("Could not find route configuration for listener with name %s", listenerName)
	}

	for _, routeExtRefCtx := range httpListenerExtCtx.RouteFilterExtensionCtxs {
		routeName := routeExtRefCtx.RouteName
		routeNamespace := routeExtRefCtx.RouteNamespace
		route := getRouteByName(routeCfg, routeName)
		if route == nil {
			continue
		}
		for _, extensionRef := range routeExtRefCtx.RouteFilterExtensionRefs {
			route.ResponseHeadersToAdd = append(route.ResponseHeadersToAdd,
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-listener-name",
						Value: listenerName,
					},
				},
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-route-name",
						Value: routeName,
					},
				},
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-route-namespace",
						Value: routeNamespace,
					},
				},
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-extensionRef-group",
						Value: extensionRef.ApiGroup,
					},
				},
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-extensionRef-kind",
						Value: extensionRef.Kind,
					},
				},
				&core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   "mock-extension-was-here-extensionRef-name",
						Value: extensionRef.Name,
					},
				},
			)
		}
	}
	return nil
}

// getRouteCfgByName finds the RouteConfigration with the name and returns nil if there is no match
func getRouteCfgByName(xdsResourceVerTbl *xdsTypes.ResourceVersionTable, name string) *route.RouteConfiguration {
	if xdsResourceVerTbl == nil ||
		xdsResourceVerTbl.XdsResources == nil ||
		xdsResourceVerTbl.XdsResources[resource.RouteType] == nil {
		return nil
	}

	for _, r := range xdsResourceVerTbl.XdsResources[resource.RouteType] {
		routeCfgResource := r.(*route.RouteConfiguration)
		if routeCfgResource.Name == name {
			return routeCfgResource
		}
	}

	return nil
}

// getRouteByName finds the Route with the name and returns nil if there is no match
func getRouteByName(routeCfg *route.RouteConfiguration, name string) *route.Route {
	if routeCfg == nil {
		return nil
	}

	for _, vh := range routeCfg.VirtualHosts {
		for _, route := range vh.Routes {
			if route.Name == name {
				return route
			}
		}
	}
	return nil
}
