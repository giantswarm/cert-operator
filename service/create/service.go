package create

import (
	"sync"

	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
)

// Config represents the configuration used to create a version service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new version service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient: nil,
		Logger:    nil,
	}
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	newService := &Service{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	// Internals.
	bootOnce sync.Once
}

// Boot starts the service
func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		s.logger.Log("info", "booted cert-operator")

		// TODO Add watch for certificate TPR
	})
}
