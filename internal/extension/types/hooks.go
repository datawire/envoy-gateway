package types

import (
	"github.com/envoyproxy/gateway/proto/extension"
	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
)

type HookType string

const (
	XDSHook HookType = "xds"
)

// XDSHook specifies the entrypoints for extensions to hook into the XDS
// translation pipeline.
type XDSHookClient interface {
	// PostHTTPListenerTranslation is the entrypoint for modifying XDS after processing
	// the HTTP listener. The method takes two arguments:
	//
	// * httpListenerExtCtx: Contains the context necessary for the extension to be able to modify xDS resources
	//
	// * xdsResourceVerTbl: the XDS Resource Table that the xDS server will server back to envoy. This table
	//   will be modified by the method using the response from the extension service for the updated xDS resources.
	//   It should not be sent over the wire.
	PostHTTPListenerTranslation(
		httpListenerExtCtx *extension.HTTPListenerExtensionContext, xdsResourceVerTbl *xdsTypes.ResourceVersionTable,
	) error
}
