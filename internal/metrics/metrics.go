package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	controllerMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	TranslationCount         *prometheus.CounterVec
	TranslationTime          *prometheus.HistogramVec
	TranslationSuccess       *prometheus.CounterVec
	TranslationError         *prometheus.CounterVec
	ExtensionHookRequests    *prometheus.CounterVec
	ExtensionHookResponse    *prometheus.CounterVec
	ExtensionResponseLatency *prometheus.HistogramVec
)

// RegisterMetrics registers prometheus metrics for the Control Plane.
//
// TODO: For now we're going to register the metrics with the controller-runtime since we only have the
// Kubernetes provider. At some point when we have more providers, we might want to create a standalone
// metrics server.
func RegisterMetrics() {
	TranslationCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "translation_ops_total",
			Help:      "Total number of translation operations",
		},
		[]string{"translator"},
	)

	TranslationSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "translation_success_total",
			Help:      "Number of times a translation operation succeeded",
		},
		[]string{"translator"},
	)

	TranslationError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "translation_error_total",
			Help:      "Number of times a translation operation resulted in an error",
		},
		[]string{"translator"},
	)

	TranslationTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "translation_time",
			Help:      "Histogram of translation operation time (seconds)",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"translator"},
	)

	ExtensionHookRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "extension_hook_request_total",
			Help:      "Total number of requests made to an extension service",
		},
		[]string{"extension_id", "hook_type", "method"},
	)

	ExtensionHookResponse = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "extension_hook_response_total",
			Help:      "Total number of RPCs returned by an extension service, regardless of success or failure",
		},
		[]string{"extension_id", "hook_type", "method", "code"},
	)

	ExtensionResponseLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "envoy_gateway_control_plane",
			Name:      "extension_hook_response_latency_seconds",
			Help:      "Histogram of response latency (seconds) of requests made to an extension service",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"extension_id", "hook_type", "method"},
	)

	controllerMetrics.Registry.MustRegister(
		collectors.NewProcessCollector(
			collectors.ProcessCollectorOpts{Namespace: "envoy_gateway_control_plane"},
		),
		TranslationCount,
		TranslationSuccess,
		TranslationError,
		TranslationTime,
		ExtensionHookRequests,
		ExtensionHookResponse,
		ExtensionResponseLatency,
	)
}
