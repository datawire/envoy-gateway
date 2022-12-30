package translator

import (
	"fmt"
	"time"

	"github.com/envoyproxy/gateway/internal/ir"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	extauthzv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// MockInjectHeaderFilterSpec testing
type MockInjectHeaderFilterSpec struct {
	Name    string
	Headers map[string]string
}

/******************** This is a mock for experimentation *************/
// simulate the Vendor side of the story as if it had recieved the
// Hook over the wire Protocol and is now handling it
//
// How would I on the extension side get the results I want?
type EGPlusSnapshot struct {
	// InjectHeaderFilters is an in-memory Snapshot from our Controller watching any CR's
	// for this CRD. EG doesn't know about these other than having a Reference
	InjectHeaderFilters map[string]MockInjectHeaderFilterSpec

	// MissingCustomResources provides a way to capture "missing resources" that were not
	// available during Hook processing. They might not exist at all because the user didn't
	// create them or due to eventual-consistency (race-condition) between when it was applied
	// and when the hook was called.
	//
	// Once the Vendor controller finds these resources it would need to "re-trigger" EG's
	// pipeline which would lead to calling the hook again.
	//
	// Note: this needs a mutex and might not reside here but to show-case the concept
	MissingCustomResource map[string]struct{}
}

var mockEGPlusSnapshot = EGPlusSnapshot{
	InjectHeaderFilters: map[string]MockInjectHeaderFilterSpec{
		"my-inject-header-filter": {
			Name: "my-inject-header-filter",
			Headers: map[string]string{
				"x-moon": ":rocket:",
			},
		},
	},
}

// RouteExtensionRefContext is part of the Vendor Extension Protocol which provides
// enough context about an HTTPRoute with an ExtensionRef filter. The context provided
// allows vendors to identify which filterchain/s, route and the vendor Custom Resource's
// that were referenced so they can modify
type RouteExtensionRefContext struct {
	// RouteName needed to identify route that has filter applied
	RouteName string `json:"routeName"`
	// HostNames give us the ability to determine which FilterChain the route is associated with
	// an alternative would be to uniquely name FilterChains and pass them here but currently
	// FilterChains have no "Name" field
	HostNames []string `json:"hostNames"`
	// FilterExtensionRef uniquely identifies the ExtensionRef defined in the HTTPRoute which
	// was registered to be handled by this Vendor
	FilterExtensionRefs []ir.ExtensionFilterRef `json:"filterExtensionRefs"`
}

// ExtensionRefContext is part of the Vendor Extension Protocol which provides context on
// the users intent so that the provided XDS can be modified by a Vendor.
type ExtensionRefContext struct {
	// RouteExtensionRefContexts provides a list of routes that have Vendor ExtensionRefs that
	// are handled by a Vendor. The list is in the order it was processed in XDS Translation
	RouteExtensionRefContexts []RouteExtensionRefContext
}

// PostListenerHookRequest is the vendor protocol request that is sent after a listener is processed
// by `xds.Translate` before it is added to the snapShotCache
type PostListenerHookRequest struct {
	Listener            *listenerv3.Listener `json:"listener"`
	ExtensionRefContext ExtensionRefContext  `json:"extensionRefContext"`
}

// PostListenerHookResponse is the vendor extension protocol response that would be sent over the wire
// protocol. The listener provided should be taken as-is.
type PostListenerHookResponse struct {
	// Listener modified by vendor that will replace the EG generated listener in the XDS snapshot cache
	Listener *listenerv3.Listener `json:"listener"`

	// ??? should we change this to a generic resource.Type and

	Clusters []*clusterv3.Cluster

	// Errors is a list of errors that occur during processing.
	// This allows the vendor to provide full scope of potential issues with CustomResource, etc... so that
	// an EG admin can effectively
	// TODO(lance): does this makes sense to be logged back in EG so that errors like this are centralized or
	// better to keep logging here and just pass back the count of errors and have EG mention an error occured
	// see vendor extension logs for details ???, after writing this I'm leaning towards the latter...
	Errors []error `json:"errors"`
}

// getPostHTTPListenerHooks is just mocking the TBD design of the ExtensionRegistry where one would fetch
// all the registed vendors that want to implement it (TBD whether this is an opt-in to "capabilities" type
// thing or not.
//
// Multiple vendors could extend EG (I'm guessing rare but design should factor this in upfront)
// A single "vendor" is registered below which simulates the request received over the wire
// and the vendor processing it and return the proper xDS that it wants EG to use.
//
// Note: at this point it is up to the vendor to get this right. From EG's perspective it
// sent a valid Listener that conforms to Gateway-API conformance.
func getPostHTTPListenerHooks() []func(PostListenerHookRequest) (PostListenerHookResponse, error) {
	return []func(PostListenerHookRequest) (PostListenerHookResponse, error){
		// Ambassadors registered as a vendorf
		func(req PostListenerHookRequest) (PostListenerHookResponse, error) {

			respClusters := make([]*clusterv3.Cluster, 0)
			respListener := req.Listener
			errors := []error{}
			extAuthzClusterAdded := false

			/***************** Process RouteFilterExtensions ****************/
			// Note: a couple of things:
			//	1. the below logic is just quick and dirty and doesn't cover all nuances we will need for our
			// features. However, it proves out that generally speaking we have full control of xDS for
			// a Listener and its config and other resources we may need to add such as clusters, secrets, etc...
			//
			// 2. once extended to support PolicyAttachments then a vendor
			// may decided to process them first or have some more complex logic
			// where the two are weaved (TBD and up to vendor). The example below
			// handles the initial use case of only supporting HTTPRoute ExtensionRef filters
			//
			// 3. starting by looping the extension refs vs looping the listener routes is an
			// implementation detail subject to change when it comes time to do real implementations

			// for each HTTPRoute that has at least one extension ref that is registered to be handled by us
			for _, routeExtRefCtx := range req.ExtensionRefContext.RouteExtensionRefContexts {
				// a single HTTPRoute could have multiple ExtensionRef's in its rules
				for _, filterExtRef := range routeExtRefCtx.FilterExtensionRefs {

					extRefID := filterExtRef.GroupVersionKind.String()

					// Lookup which ExternalRef it wants processed, this example only supports InjectHeaderFilter
					switch extRefID {
					case "eg-plus.getambassador.io/v1alpha1, Kind=InjectHeaderFilter":
						// TODO (lance) - name is not unique enough and we would also need `namespace`
						if _, ok := mockEGPlusSnapshot.InjectHeaderFilters[filterExtRef.Name]; !ok {
							errors = append(errors, fmt.Errorf("custom resource not found: %s", extRefID))
							continue
						}

						// we need to make sure our cluster exists to handle ext_authz
						if !extAuthzClusterAdded {
							cluster := buildInjectHeaderExtAuthzCluster()
							respClusters = append(respClusters, cluster)
							extAuthzClusterAdded = true
						}

						// now we need to make sure the HCM.httpFilter for ext_authz is added to the correct Filter Chains
						for _, hostName := range routeExtRefCtx.HostNames {
							filterChain, err := getFilterChainByHostName(hostName, respListener)
							if err != nil {
								// note: I don't think this ever could happen, if it did then its a programming error either here or on the EG side.
								errors = append(errors, err)
								continue
							}

							connMgr, connMgrIndex, err := getHCM(filterChain)
							if err != nil {
								// note: I don't think this every could happen, if it did then its a programming error either here or on the EG side.
								errors = append(errors, err)
							}

							// check if it already exists else added it
							exists := false
							for _, httpFilter := range connMgr.HttpFilters {
								if httpFilter.Name == wellknown.HTTPExternalAuthorization {
									exists = true
									break
								}
							}

							if exists {
								continue
							}

							// add in filter to conn Manager
							authFilter := extauthzv3.ExtAuthz{
								TransportApiVersion: corev3.ApiVersion_V3,
								ClearRouteCache:     true,
								Services: &extauthzv3.ExtAuthz_GrpcService{
									GrpcService: &corev3.GrpcService{
										Timeout: durationpb.New(5 * time.Second),
										TargetSpecifier: &corev3.GrpcService_EnvoyGrpc_{
											EnvoyGrpc: &corev3.GrpcService_EnvoyGrpc{
												ClusterName: "ambassador-eg-plus_headerinjection_extauthz",
											},
										},
									},
								},
							}

							anyAuthFilter, _ := anypb.New(&authFilter)

							extAuthzHTTPFilter := &hcm.HttpFilter{
								Name: wellknown.HTTPExternalAuthorization,
								ConfigType: &hcm.HttpFilter_TypedConfig{
									TypedConfig: anyAuthFilter,
								},
							}
							// inject before router filter (note: could be more efficient with allocs but good enough for demonstration)
							connMgr.HttpFilters = append([]*hcm.HttpFilter{extAuthzHTTPFilter}, connMgr.HttpFilters...)

							// since the Listener Filters are of Protobuf "Any" and we had to unmarshall the HCM into a new allocated object,
							// so we need to write that new allocated object back to the filter chain
							anyConnMgr, _ := anypb.New(connMgr)
							filterChain.Filters[connMgrIndex].ConfigType = &listenerv3.Filter_TypedConfig{TypedConfig: anyConnMgr}
						}

						// Note: leaving off for brevity but we would want some logic handling disabling routes that were not affected
						// for each route that doesn't have a "inject header" filter we need to override it
						//     typed_per_filter_config:
						// envoy.filters.http.ext_authz:
						// "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute
						// disabled: true

					default:
						errors = append(errors, fmt.Errorf("unsupported extensionRef: %s", extRefID))
					}
				}
			}

			resp := PostListenerHookResponse{
				Listener: respListener,
				Clusters: respClusters,
				Errors:   errors,
			}

			return resp, nil
		},
	}
}

/********************** A bunch of helpers for working xDS protos ****************************************/

func getFilterChainByHostName(hostName string, listener *listenerv3.Listener) (*listenerv3.FilterChain, error) {

	if hostName == "*" {
		if listener.DefaultFilterChain == nil {
			return nil, fmt.Errorf("hostname is wildcard (*) but no default filter chain found for listener %s", listener.Name)
		}
		return listener.DefaultFilterChain, nil
	}

	for _, filterChain := range listener.GetFilterChains() {

		if filterChain.FilterChainMatch != nil {
			// check filter Chain match for sni since it is https
			for _, serverName := range filterChain.FilterChainMatch.ServerNames {
				if serverName == hostName {
					return filterChain, nil
				}
			}
			continue
		}

		// no matchers so we assume sni is not being used

		hcm, _, err := getHCM(filterChain)
		if err != nil {
			// filter chain doesn't have HCM
			continue
		}

		for _, vh := range hcm.GetRouteConfig().VirtualHosts {
			for _, domain := range vh.Domains {
				if domain == hostName {
					return filterChain, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no filter chain found for hostname %s", hostName)
}

func getHCM(filterChain *listenerv3.FilterChain) (*hcm.HttpConnectionManager, int, error) {

	for hcmIndex, filter := range filterChain.Filters {

		if filter.Name == wellknown.HTTPConnectionManager {
			connMgr := new(hcm.HttpConnectionManager)
			if err := filter.GetTypedConfig().UnmarshalTo(connMgr); err != nil {
				return nil, -1, err
			}

			return connMgr, hcmIndex, nil
		}
	}

	return nil, -1, fmt.Errorf("no HTTPConnectionManager found in FilterChain: %s", filterChain.Name)
}

func buildInjectHeaderExtAuthzCluster() *clusterv3.Cluster {
	// explicitConfig := clusterv3.Protol

	typedConfig := &corev3.TypedExtensionConfig{
		Name:        "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions",
		TypedConfig: &anypb.Any{},
	}

	typedOptions, _ := anypb.New(typedConfig)

	return &clusterv3.Cluster{
		Name:                 "ambassador-eg-plus_headerinjection_extauthz",
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_LOGICAL_DNS},
		TypedExtensionProtocolOptions: map[string]*anypb.Any{
			"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": typedOptions,
		},
	}

}
