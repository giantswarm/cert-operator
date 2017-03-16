package create

import (
	"sync"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/giantswarm/cert-operator/flag"
)

// Config represents the configuration used to create a create service.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient: nil,
		Logger:    nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// New creates a new configured version service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "kubernetes client must not be empty")
	}

	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "viper must not be empty")
	}

	newService := &Service{
		Config: config,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

// Service implements the version service interface.
type Service struct {
	Config

	// Internals.
	bootOnce sync.Once
}

// Boot starts the service
func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		s.Config.Logger.Log("info", "booted cert-operator")

		// TODO Add watch for certificate TPR
	})
}
