package types

import (
	"github.com/envoyproxy/gateway/api/config/v1alpha1"

	"github.com/envoyproxy/gateway/proto/extension"
)

// HTTPListenerExtensionRefs holds all HTTPRouteFilter extensionRefs along with
// its corresponding context from the associated HTTPRoute.
type HTTPListenerExtensionRefs map[v1alpha1.ExtensionId][]extension.HTTPListenerExtensionContext

// TODO: Generate DeepCopy methods?
//
// type HTTPListenerExtensionRefTable struct {
//	HTTPListenerExtensionRefs
// }
