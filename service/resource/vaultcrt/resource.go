package vaultcrt

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultrole"
	vaultclient "github.com/hashicorp/vault/api"
	"k8s.io/client-go/kubernetes"
)

const (
	// Name is the identifier of the resource.
	Name = "vaultcrt"
	// VaultAllowSubDomains defines whether to allow the generated root CA of the
	// PKI backend to allow sub domains as common names.
	VaultAllowSubDomains = "true"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	VaultCrt  vaultcrt.Interface
	VaultRole vaultrole.Interface
}

// DefaultConfig provides a default configuration to create a new cloud config
// resource by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient:   nil,
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		CATTL:            "",
		CommonNameFormat: "",
	}
}

// Resource implements the cloud config resource.
type Resource struct {
	// Dependencies.
	k8sClient   kubernetes.Interface
	logger      micrologger.Logger
	vaultClient *vaultclient.Client

	// Settings.
	caTTL            string
	commonNameFormat string
}

// New creates a new configured cloud config resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultClient must not be empty")
	}

	// Settings.
	if config.CATTL == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CATTL must not be empty")
	}
	if config.CommonNameFormat == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CommonNameFormat must not be empty")
	}

	newResource := &Resource{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger: config.Logger.With(
			"resource", Name,
		),
		vaultClient: config.VaultClient,

		// Settings.
		caTTL:            config.CATTL,
		commonNameFormat: config.CommonNameFormat,
	}

	return newResource, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) Underlying() framework.Resource {
	return r
}

func toSecret(v interface{}) (*apiv1.Secret, error) {
	if v == nil {
		return nil, nil
	}

	secret, ok := v.(*apiv1.Secret)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &apiv1.Secret{}, v)
	}

	return secret, nil
}
