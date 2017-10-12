package vaultpki

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	vaultclient "github.com/hashicorp/vault/api"
)

const (
	// Name is the identifier of the resource.
	Name = "vaultpki"
	// VaultAllowSubDomains defines whether to allow the generated root CA of the
	// PKI backend to allow sub domains as common names.
	VaultAllowSubDomains = "true"
	// VaultMountType is the mount type used to mount a PKI backend in Vault.
	VaultMountType = "pki"
)

// Config represents the configuration used to create a new cloud config resource.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	VaultClient *vaultclient.Client

	// Settings.
	CATTL            string
	CommonNameFormat string
}

// DefaultConfig provides a default configuration to create a new cloud config
// resource by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
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
	logger      micrologger.Logger
	vaultClient *vaultclient.Client

	// Settings.
	caTTL            string
	commonNameFormat string
}

// New creates a new configured cloud config resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
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

func toVaultPKIState(v interface{}) (VaultPKIState, error) {
	if v == nil {
		return VaultPKIState{}, nil
	}

	vaultPKIState, ok := v.(VaultPKIState)
	if !ok {
		return VaultPKIState{}, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", VaultPKIState{}, v)
	}

	return vaultPKIState, nil
}
