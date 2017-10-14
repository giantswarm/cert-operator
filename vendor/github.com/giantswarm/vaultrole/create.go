package vaultrole

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultrole/key"
)

func (r *VaultRole) Create(config CreateConfig) error {
	k := key.WriteRolePath(config.ID, config.Organizations)
	v := map[string]interface{}{
		"allow_bare_domains": config.AllowBareDomains,
		"allow_subdomains":   config.AllowSubdomains,
		"allowed_domains":    key.AllowedDomains(config.ID, r.commonNameFormat, config.AltNames),
		"organization":       config.Organizations,
		"ttl":                config.TTL,
	}

	_, err := r.vaultClient.Logical().Write(k, v)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
