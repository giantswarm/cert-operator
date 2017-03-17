package create

import (
	"fmt"
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/flag"
)

// TODO Replace with Certificate TPR
type CertificateSpec struct {
	ClusterID        string
	CommonName       string
	IPSANs           []string
	AltNames         []string
	AllowBareDomains bool
	TTL              string
}

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
		s.Config.Logger.Log("info", "Booted cert-operator")

		s.Config.Logger.Log("info", "Test issuing a cert")

		cert := CertificateSpec{
			ClusterID:        "cert-test",
			CommonName:       "api.cert-test.giantswarm.io",
			IPSANs:           []string{"10.0.0.4", "10.0.0.5"},
			AltNames:         []string{"api.k8s.cert-test.giantswarm.io"},
			AllowBareDomains: false,
			TTL:              "720h",
		}

		issueResp, err := s.Issue(cert)
		if err == nil {
			s.Config.Logger.Log("info", "Cert issued")
			s.Config.Logger.Log("info", cert.CommonName)
			s.Config.Logger.Log("info", issueResp.SerialNumber)
		} else {
			s.Config.Logger.Log("error", fmt.Sprintf("Failed to issue cert - %v", err))
		}
	})
}
