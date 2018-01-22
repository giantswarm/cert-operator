package vaultrole

import (
	"encoding/json"

	"github.com/giantswarm/microerror"
	"github.com/hashicorp/vault/helper/parseutil"

	"github.com/giantswarm/vaultrole/key"
)

func (r *VaultRole) Exists(config ExistsConfig) (bool, error) {
	// Check if a PKI for the given cluster ID exists.
	secret, err := r.vaultClient.Logical().List(key.ListRolesPath(config.ID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, microerror.Mask(err)
	}

	// In case there is not a single role for this PKI backend, secret is nil.
	if secret == nil {
		return false, nil
	}

	// When listing roles a list of role names is returned. Here we iterate over
	// this list and if we find the desired role name, it means the role has
	// already been created.
	if keys, ok := secret.Data["keys"]; ok {
		if list, ok := keys.([]interface{}); ok {
			for _, k := range list {
				if str, ok := k.(string); ok && str == key.RoleName(config.ID, config.Organizations) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

type role struct {
	AllowBareDomains bool   `json:"allow_bare_domains"`
	AllowSubdomains  bool   `json:"allow_subdomains"`
	AllowedDomains   string `json:"allowed_domains"`
	Organizations    string `json:"organization"` // NOTE the singular form here.
	TTL              string `json:"ttl"`
}

func (r *VaultRole) Search(config SearchConfig) (Role, error) {
	// Check if a PKI for the given cluster ID exists.
	secret, err := r.vaultClient.Logical().Read(key.ReadRolePath(config.ID, config.Organizations))
	if IsNoVaultHandlerDefined(err) {
		return Role{}, microerror.Maskf(notFoundError, "no vault handler defined")
	} else if err != nil {
		return Role{}, microerror.Mask(err)
	}

	// In case there is not a single role for this PKI backend, secret is nil.
	if secret == nil {
		return Role{}, microerror.Maskf(notFoundError, "no vault secret at path '%s'", key.RoleName(config.ID, config.Organizations))
	}

	b, err := json.Marshal(secret.Data)
	if err != nil {
		return Role{}, microerror.Mask(err)
	}

	var internalRole role
	err = json.Unmarshal(b, &internalRole)
	if err != nil {
		return Role{}, microerror.Mask(err)
	}

	altNames := key.ToAltNames(internalRole.AllowedDomains)
	organizations := key.ToOrganizations(internalRole.Organizations)
	ttl, err := parseutil.ParseDurationSecond(internalRole.TTL)
	if err != nil {
		return Role{}, microerror.Mask(err)
	}

	newRole := Role{
		AllowBareDomains: internalRole.AllowBareDomains,
		AllowSubdomains:  internalRole.AllowSubdomains,
		AltNames:         altNames,
		ID:               config.ID,
		Organizations:    organizations,
		TTL:              ttl,
	}

	return newRole, nil
}
