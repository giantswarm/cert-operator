package vaultcrt

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "cert_operator"
	PrometheusSubsystem = "vaultcrt_resource"
	// VersionBundleVersionAnnotation = "giantswarm.io/version-bundle-version"
)

var versionGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "version_total",
		Help:      "A metric labeled by major, minor and patch version of the operator in use.",
	},
	[]string{"major", "minor", "patch"},
)

func init() {
	prometheus.MustRegister(versionGauge)
}
