package crt

import (
	"github.com/cenk/backoff"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cert-operator/flag"
	"github.com/giantswarm/cert-operator/service/ca"
)

// Config represents the configuration used to create a Crt service.
type Config struct {
	// Dependencies.
	BackOff     backoff.BackOff
	CAService   *ca.Service
	Logger      micrologger.Logger
	K8sClient   kubernetes.Interface
	VaultClient *vaultapi.Client

	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

// certificateSecret stores a cert issued by Vault that will be stored as a k8s secret.
type certificateSecret struct {
	Certificate   certificatetpr.Spec
	IssueResponse spec.IssueResponse
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		BackOff:     nil,
		CAService:   nil,
		K8sClient:   nil,
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		Flag:  nil,
		Viper: nil,
	}
}

// Service implements the Crt service interface.
type Service struct {
	Config
}

// New creates a new configured Crt service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.CAService == nil {
		return nil, microerror.Maskf(invalidConfigError, "ca service must not be empty")
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "kubernetes client must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "vault client must not be empty")
	}

	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "viper must not be empty")
	}

	newService := &Service{
		Config: config,
	}

	return newService, nil
}
