package vaulthealthz

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "microendpoint"
	PrometheusSubsystem = "vaulthealthz"
)

var tokenExpireTimeGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "token_expire_time",
		Help:      "A metric of the expire time of Vault tokens as unix seconds.",
	},
)

func init() {
	prometheus.MustRegister(tokenExpireTimeGauge)
}
