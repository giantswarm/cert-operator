package vaultrole

import (
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/vaultrole/key"
	vaultclient "github.com/hashicorp/vault/api"
)

type Config struct {
	Logger      micrologger.Logger
	VaultClient *vaultclient.Client

	CommonNameFormat string
}

func DefaultConfig() Config {
	config := Config{
		Logger:      nil,
		VaultClient: nil,

		CommonNameFormat: "",
	}

	return config
}

type VaultRole struct {
	logger      micrologger.Logger
	vaultClient *vaultclient.Client

	commonNameFormat string
}

func New(config Config) (*VaultRole, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.VaultClient must not be empty")
	}

	if config.CommonNameFormat == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.CommonNameFormat must not be empty")
	}

	r := &VaultRole{
		logger:      config.Logger,
		vaultClient: config.VaultClient,

		commonNameFormat: config.CommonNameFormat,
	}

	return r, nil
}

type writeConfig struct {
	AllowBareDomains bool
	AllowSubdomains  bool
	AltNames         []string
	ID               string
	Organizations    []string
	TTL              string
}

func (r *VaultRole) write(config writeConfig) error {
	k := key.WriteRolePath(config.ID, config.Organizations)
	v := map[string]interface{}{
		"allow_bare_domains": config.AllowBareDomains,
		"allow_subdomains":   config.AllowSubdomains,
		"allowed_domains":    key.AllowedDomains(config.ID, r.commonNameFormat, config.AltNames),
		"organization":       strings.Join(config.Organizations, ","),
		"ttl":                config.TTL,
	}

	_, err := r.vaultClient.Logical().Write(k, v)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
