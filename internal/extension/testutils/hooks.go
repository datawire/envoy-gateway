package testutils

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

type MockInjectHeaderFilterSpec struct {
	Name    string
	Headers map[string]string
}

type XDSHookClient struct{}

// PostHTTPListenerTranslation modifies xdsResourceVerTbl. It justs adds an extra response header.
func (c *XDSHookClient) PostHTTPListenerTranslation(
	httpListenerExtCtx *extension.HTTPListenerExtensionContext, xdsResourceVerTbl *xdsTypes.ResourceVersionTable,
) error {
	listenerName := httpListenerExtCtx.ListenerName

	for _, routeExtRefCtx := range httpListenerExtCtx.RouteFilterExtensionCtxs {
		for i := 0; i < len(routeExtRefCtx.RouteFilterExtensionRefs); i++ {
			for _, hostname := range routeExtRefCtx.Hostnames {
				routeCfg, found := getRouteCfgByHostnameAndListenerName(hostname, listenerName, xdsResourceVerTbl)
				if !found {
					// NOTE: this could be a programming error
					return fmt.Errorf("Could not find route configuration for listener with name %s and hostname %s", listenerName, hostname)
				}
				for _, vh := range routeCfg.VirtualHosts {
					for _, xdsRoute := range vh.Routes {
						xdsRoute.ResponseHeadersToAdd = append(xdsRoute.ResponseHeadersToAdd, &core.HeaderValueOption{
							Header: &core.HeaderValue{
								Key:   "mock-extension-was-here",
								Value: "some-value",
							},
						})
					}
				}
			}
		}
	}
	return nil
}

func getRouteCfgByHostnameAndListenerName(
	hostname string,
	listenerName string,
	xdsResourceVerTbl *xdsTypes.ResourceVersionTable) (*route.RouteConfiguration, bool) {

	if xdsResourceVerTbl == nil || xdsResourceVerTbl.XdsResources == nil || xdsResourceVerTbl.XdsResources[resource.RouteType] == nil {
		return nil, false
	}

	for _, r := range xdsResourceVerTbl.XdsResources[resource.RouteType] {
		routeCfgResource := r.(*route.RouteConfiguration)
		if routeCfgResource.Name == listenerName {
			for _, vh := range routeCfgResource.VirtualHosts {
				for _, domain := range vh.Domains {
					if domain == hostname {
						return routeCfgResource, true
					}
				}
			}
		}
	}

	return nil, false
}
