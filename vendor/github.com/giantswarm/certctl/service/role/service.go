package role

import (
	"fmt"

	"github.com/giantswarm/microerror"
	vaultclient "github.com/hashicorp/vault/api"
)

// Config defines configurable aspects (such as dependencies) of this service.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client

	// Settings.
	PKIMountpoint string
}

// DefaultConfig returns a default configuration that can be used to create this service.
func DefaultConfig() Config {
	config := Config{}

	return config
}

// New takes a configuration and returns a configured service.
func New(config Config) (Service, error) {
	// Dependencies.
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "Vault client must not be empty")
	}

	if config.PKIMountpoint == "" {
		return nil, microerror.Maskf(invalidConfigError, "PKIMountpoint must not be empty")
	}

	service := &service{
		vaultClient:   config.VaultClient,
		pkiMountpoint: config.PKIMountpoint,
	}

	return service, nil
}

type service struct {
	// Dependencies.
	vaultClient *vaultclient.Client

	// Settings.
	pkiMountpoint string
}

// Create creates a role if it doesn't exist yet. Creating roles is idempotent
// in the vault api, so no need to check if it already exists.
func (s *service) Create(params CreateParams) error {
	logicalStore := s.vaultClient.Logical()

	data := map[string]interface{}{
		"allowed_domains":    params.AllowedDomains,
		"allow_subdomains":   params.AllowSubdomains,
		"ttl":                params.TTL,
		"allow_bare_domains": params.AllowBareDomains,
		"organization":       params.Organizations,
	}

	_, err := logicalStore.Write(fmt.Sprintf("%s/roles/%s", s.pkiMountpoint, params.Name), data)
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (s *service) IsRoleCreated(roleName string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := s.vaultClient.Logical()

	// Check if a PKI for the given cluster ID exists.
	secret, err := logicalBackend.List(s.listRolesPath())
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
				if str, ok := k.(string); ok && str == roleName {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (s *service) listRolesPath() string {
	return fmt.Sprintf("%s/roles/", s.pkiMountpoint)
}
