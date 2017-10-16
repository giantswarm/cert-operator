package vaultrole

import (
	"fmt"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/vaultrole/key"
)

func (r *VaultRole) Create(config CreateConfig) error {
	fmt.Printf("start VaultRole.Create\n")
	defer fmt.Printf("end VaultRole.Create\n")

	fmt.Printf("key.RoleName: %#v\n", key.RoleName(config.ID, config.Organizations))

	k := key.WriteRolePath(config.ID, config.Organizations)
	v := map[string]interface{}{
		"allow_bare_domains": config.AllowBareDomains,
		"allow_subdomains":   config.AllowSubdomains,
		"allowed_domains":    key.AllowedDomains(config.ID, r.commonNameFormat, config.AltNames),
		"organization":       config.Organizations,
		"ttl":                config.TTL,
	}
	fmt.Printf("k: %#v\n", k)
	fmt.Printf("v: %#v\n", v)

	_, err := r.vaultClient.Logical().Write(k, v)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
