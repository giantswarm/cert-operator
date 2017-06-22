package healthz

import (
	"context"
	"sync"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	healthCheckRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "health_check_request_total",
			Help: "Number of health check requests",
		},
		[]string{"success"},
	)
	healthCheckRequestTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_check_request_milliseconds",
		Help: "Time taken to respond to health check, in milliseconds",
	})
)

func init() {
	prometheus.MustRegister(healthCheckRequests)
	prometheus.MustRegister(healthCheckRequestTime)
}

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client
}

// DefaultConfig provides a default configuration to create a new healthz
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:      nil,
		VaultClient: nil,
	}
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault client must not be empty")
	}

	newService := &Service{
		// Dependencies.
		vaultClient: config.VaultClient,
		logger:      config.Logger,

		// Internals.
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	vaultClient *vaultapi.Client
	logger      micrologger.Logger

	// Internals.
	bootOnce sync.Once
}

// Check implements the health check which lists the mounts for the sys backend.
// This checks that the operator can connect to the Vault API and the token is
// valid.
func (s *Service) Check(ctx context.Context, request Request) (*Response, error) {
	start := time.Now()
	defer func() {
		healthCheckRequestTime.Set(float64(time.Since(start) / time.Millisecond))
	}()

	sysBackend := s.vaultClient.Sys()

	_, err := sysBackend.ListMounts()
	if err != nil {
		healthCheckRequests.WithLabelValues("failed").Inc()
		return nil, microerror.MaskAny(err)
	}

	healthCheckRequests.WithLabelValues("successfull").Inc()

	response := DefaultResponse()
	response.Code = microserver.CodeSuccess
	response.Message = "Everything OK."

	return response, nil
}
