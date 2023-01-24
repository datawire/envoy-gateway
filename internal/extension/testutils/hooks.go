package testutils

import (
	xdsTypes "github.com/envoyproxy/gateway/internal/xds/types"
	"github.com/envoyproxy/gateway/proto/extension"
)

type XDSHookClient struct{}

func (c *XDSHookClient) PostHTTPListenerTranslation(
	httpListenerExtCtx *extension.HTTPListenerExtensionContext, xdsResourceVerTbl *xdsTypes.ResourceVersionTable,
) error {
	return nil
}
