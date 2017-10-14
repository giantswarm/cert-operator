package vaultcrt

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultcrt"
	"github.com/giantswarm/vaultrole"
	"k8s.io/client-go/kubernetes"
)

const (
	Name = "vaultcrt"
	// AllowSubDomains defines whether to allow the generated root CA of the PKI
	// backend to allow sub domains as common names.
	AllowSubDomains = true
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	VaultCrt  vaultcrt.Interface
	VaultRole vaultrole.Interface

	Namespace string
}

func DefaultConfig() Config {
	return Config{
		K8sClient: nil,
		Logger:    nil,
		VaultCrt:  nil,
		VaultRole: nil,

		Namespace: "default",
	}
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger
	vaultCrt  vaultcrt.Interface
	vaultRole vaultrole.Interface
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultCrt == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultCrt must not be empty")
	}
	if config.VaultRole == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultRole must not be empty")
	}

	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Namespace must not be empty")
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		logger: config.Logger.With(
			"resource", Name,
		),
		vaultCrt:  config.VaultCrt,
		vaultRole: config.VaultRole,

		namespace: config.Namespace,
	}

	return r, nil
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
