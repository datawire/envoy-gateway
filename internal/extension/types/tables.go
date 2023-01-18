package types

import (
	"github.com/envoyproxy/gateway/api/config/v1alpha1"

	"github.com/envoyproxy/gateway/internal/extension/proto"
)

// HTTPListenerExtensionRefs holds all HTTPRouteFilter extensionRefs along with
// its corresponding context from the associated HTTPRoute.
type HTTPListenerExtensionRefs map[v1alpha1.ExtensionId][]proto.HTTPListenerExtensionContext

// TODO: Generate DeepCopy methods?
//
// type HTTPListenerExtensionRefTable struct {
//	HTTPListenerExtensionRefs
// }
