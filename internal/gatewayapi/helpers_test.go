// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// This file contains code derived from Contour,
// https://github.com/projectcontour/contour
// and is provided here subject to the following:
// Copyright Project Contour Authors
// SPDX-License-Identifier: Apache-2.0

package gatewayapi

import (
	"testing"

	"github.com/stretchr/testify/require"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	egv1a1cfg "github.com/envoyproxy/gateway/api/config/v1alpha1"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/extension/testutils"
)

func TestValidateAuthenFilterRef(t *testing.T) {
	testCases := []struct {
		name     string
		filter   *gwapiv1b1.HTTPRouteFilter
		expected bool
	}{
		{
			name: "request mirror filter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterRequestMirror,
			},
			expected: true,
		},
		{
			name: "url rewrite filter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterURLRewrite,
			},
			expected: true,
		},
		{
			name: "request header modifier filter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterRequestHeaderModifier,
			},
			expected: true,
		},
		{
			name: "request redirect filter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterRequestRedirect,
			},
			expected: true,
		},
		{
			name: "unsupported extended filter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwapiv1b1.LocalObjectReference{
					Group: "UnsupportedGroup",
					Kind:  "UnsupportedKind",
					Name:  "test",
				},
			},
			expected: false,
		},
		{
			name: "extended filter with missing reference",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
			},
			expected: false,
		},
		{
			name: "invalid authenticationfilter group",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwapiv1b1.LocalObjectReference{
					Group: "UnsupportedGroup",
					Kind:  egv1a1.AuthenticationFilterKind,
					Name:  "test",
				},
			},
			expected: false,
		},
		{
			name: "invalid authenticationfilter kind",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwapiv1b1.LocalObjectReference{
					Group: gwapiv1b1.Group(egv1a1.GroupVersion.Group),
					Kind:  "UnsupportedKind",
					Name:  "test",
				},
			},
			expected: false,
		},
		{
			name: "valid authenticationfilter",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwapiv1b1.LocalObjectReference{
					Group: gwapiv1b1.Group(egv1a1.GroupVersion.Group),
					Kind:  egv1a1.AuthenticationFilterKind,
					Name:  "test",
				},
			},
			expected: true,
		},
		{
			name: "valid extension group",
			filter: &gwapiv1b1.HTTPRouteFilter{
				Type: gwapiv1b1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwapiv1b1.LocalObjectReference{
					Group: gwapiv1b1.Group("foo.example.io"),
					Kind:  gwapiv1b1.Kind("Bar"),
					Name:  "test",
				},
			},
			expected: true,
		},
	}

	ext := egv1a1cfg.Extension{
		APIGroups: []gwapiv1b1.Group{
			"foo.example.io",
		},
	}
	em := testutils.NewManager(ext)
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateHTTPRouteFilter(em, tc.filter)
			if tc.expected {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
