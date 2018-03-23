package vaultrole

import (
	"github.com/giantswarm/microerror"
	"github.com/hashicorp/vault/api"
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

	role, err := vaultSecretToRole(secret)
	if err != nil {
		return Role{}, microerror.Mask(err)
	}

	role.ID = config.ID
	return role, nil
}

// vaultSecretToRole makes required type casts / type checks and parsing to
// extract role information from Vault api.Secret.
func vaultSecretToRole(secret *api.Secret) (Role, error) {
	var role Role

	{
		v, exists := secret.Data["allow_bare_domains"]
		if !exists {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "allow_bare_domains missing from Vault api.Secret.Data")
		}

		allowBareDomains, ok := v.(bool)
		if !ok {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"allow_bare_domains\"] type is %T, expected %T", secret.Data["allow_bare_domains"], allowBareDomains)
		}

		role.AllowBareDomains = allowBareDomains
	}

	{
		v, exists := secret.Data["allow_subdomains"]
		if !exists {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "allow_subdomains missing from Vault api.Secret.Data")
		}

		allowSubdomains, ok := v.(bool)
		if !ok {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"allow_subdomains\"] type is %T, expected %T", secret.Data["allow_subdomains"], allowSubdomains)
		}

		role.AllowSubdomains = allowSubdomains
	}

	// List types in Vault were earlier joined with comma to single
	// concatenated string. Now they are slice of interfaces which are strings
	// underneath.
	{
		allowedDomains, exists := secret.Data["allowed_domains"]
		if !exists {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "allowed_domains missing from Vault api.Secret.Data")
		}

		switch v := allowedDomains.(type) {
		case string:
			role.AltNames = key.ToAltNames(v)
		case []string:
			role.AltNames = v[1:]
		case []interface{}:
			allowedDomains := make([]string, 0, len(v))
			for i, val := range v {
				if s, ok := val.(string); ok {
					allowedDomains = append(allowedDomains, s)
				} else {
					return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"allowed_domains\"][%d] has unexpected type '%T'. It's not string nor []string.", i, val)
				}
			}

			// TODO: Why first one is dropped (this was in key.ToAltNames()?
			role.AltNames = allowedDomains[1:]
		default:
			return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"allowed_domains\"] type is '%T'. It's not string, []string nor []interface{} (masking strings).", secret.Data["allowed_domains"])
		}
	}

	{
		organizations, exists := secret.Data["organization"]
		if !exists {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "organization missing from Vault api.Secret.Data")
		}

		switch v := organizations.(type) {
		case string:
			role.Organizations = key.ToOrganizations(v)
		case []string:
			role.Organizations = v
		case []interface{}:
			organizations := make([]string, 0, len(v))
			for i, val := range v {
				if s, ok := val.(string); ok {
					organizations = append(organizations, s)
				} else {
					return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"organization\"][%d] has unexpected type '%T'. It's not string nor []string.", i, val)
				}
			}

			role.Organizations = organizations
		default:
			return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"organization\"] type is '%T'. It's not string, []string nor []interface{} (masking strings).", secret.Data["organization"])
		}
	}

	{
		v, exists := secret.Data["ttl"]
		if !exists {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "ttl missing from Vault api.Secret.Data")
		}

		ttlStr, ok := v.(string)
		if !ok {
			return Role{}, microerror.Maskf(invalidVaultResponseError, "Vault secret.Data[\"ttl\"] type is %T, expected %T", secret.Data["ttl"], ttlStr)
		}

		ttl, err := parseutil.ParseDurationSecond(ttlStr)
		if err != nil {
			return Role{}, microerror.Mask(err)
		}

		role.TTL = ttl
	}

	return role, nil
}
