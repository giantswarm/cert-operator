package create

import (
	"fmt"
	"sync"

	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/flag"
)

// Config represents the configuration used to create a create service.
type Config struct {
	// Dependencies.
	K8sClient   kubernetes.Interface
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient:   nil,
		Logger:      nil,
		VaultClient: nil,

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
	if config.VaultClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault client must not be empty")
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
		if err := s.createTPR(); err != nil {
			panic(fmt.Sprintf("could not create cluster resource: %#v", err))
		}
		s.Config.Logger.Log("info", "successfully created third-party resource")

		cert := certificatetpr.Spec{
			ClusterID:  "cert-test",
			CommonName: "api.cert-test.g8s.eu-west-1.aws.test.private.giantswarm.io",
			IPSANs:     []string{"10.0.0.4", "10.0.0.5"},
			AltNames: []string{
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				"kubernetes.default.svc.cluster.local",
			},
			AllowBareDomains: true,
			TTL:              "720h",
		}
	s.Config.Logger.Log("info", fmt.Sprintf("creating certificate '%s'", cert.Spec.CommonName))

	if err := s.setupPKIBackend(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki backend '%#v'", err))
		return
	}
	if err := s.setupPKIPolicy(cert.Spec); err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not setup pki backend '%#v'", err))
		return
	}

	_, err := s.Issue(cert.Spec)
	if err != nil {
		s.Config.Logger.Log("error", fmt.Sprintf("could not issue cert '%#v'", err))
		return
	}
	s.Config.Logger.Log("info", fmt.Sprintf("certificate issued %s", cert.Spec.CommonName))

			}

}
