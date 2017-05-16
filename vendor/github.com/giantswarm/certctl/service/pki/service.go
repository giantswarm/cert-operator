package pki

import (
	"fmt"

	vaultclient "github.com/hashicorp/vault/api"
)

// ServiceConfig represents the configuration used to create a new PKI controller.
type ServiceConfig struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultServiceConfig provides a default configuration to create a PKI controller.
func DefaultServiceConfig() ServiceConfig {
	newClientConfig := vaultclient.DefaultConfig()
	newClientConfig.Address = "http://127.0.0.1:8200"
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		panic(err)
	}

	newConfig := ServiceConfig{
		// Dependencies.
		VaultClient: newVaultClient,
	}

	return newConfig
}

// NewService creates a new configured PKI controller.
func NewService(config ServiceConfig) (Service, error) {
	// Dependencies.
	if config.VaultClient == nil {
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	newService := &service{
		ServiceConfig: config,
	}

	return newService, nil
}

type service struct {
	ServiceConfig
}

// PKI management.

func (s *service) Delete(clusterID string) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Unmount the PKI backend, if it exists.
	mounted, err := s.IsMounted(clusterID)
	if err != nil {
		return maskAny(err)
	}
	if mounted {
		err = sysBackend.Unmount(s.MountPKIPath(clusterID))
		if err != nil {
			return maskAny(err)
		}
	}

	return nil
}

func (s *service) IsCAGenerated(clusterID string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := s.VaultClient.Logical()

	// Check if a root CA for the given cluster ID exists.
	secret, err := logicalBackend.Read(s.ReadCAPath(clusterID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
	}

	// If the secret is nil, the CA has not been generated.
	if secret == nil {
		return false, nil
	}

	if certificate, ok := secret.Data["certificate"]; ok && certificate == "" {
		return false, nil
	}
	if err, ok := secret.Data["error"]; ok && err != "" {
		return false, nil
	}

	return true, nil
}

func (s *service) IsMounted(clusterID string) (bool, error) {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Check if a PKI for the given cluster ID exists.
	mounts, err := sysBackend.ListMounts()
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
	}
	mountOutput, ok := mounts[s.ListMountsPath(clusterID)+"/"]
	if !ok || mountOutput.Type != "pki" {
		return false, nil
	}

	return true, nil
}

func (s *service) IsRoleCreated(clusterID string) (bool, error) {
	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's PKI backend.
	logicalBackend := s.VaultClient.Logical()

	// Check if a PKI for the given cluster ID exists.
	secret, err := logicalBackend.List(s.ListRolesPath(clusterID))
	if IsNoVaultHandlerDefined(err) {
		return false, nil
	} else if err != nil {
		return false, maskAny(err)
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
				if str, ok := k.(string); ok && str == s.RoleName(clusterID) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (s *service) VerifyPKISetup(clusterID string) (bool, error) {
	mounted, err := s.IsMounted(clusterID)
	if err != nil {
		return false, maskAny(err)
	}
	if !mounted {
		return false, nil
	}

	caGenerated, err := s.IsCAGenerated(clusterID)
	if err != nil {
		return false, maskAny(err)
	}
	if !caGenerated {
		return false, nil
	}

	roleCreated, err := s.IsRoleCreated(clusterID)
	if !roleCreated || err != nil {
		return false, maskAny(err)
	}
	if !roleCreated {
		return false, nil
	}

	// PKI setup is valid.
	return true, nil
}

func (s *service) RoleName(clusterID string) string {
	return fmt.Sprintf("role-%s", clusterID)
}

func (s *service) Create(config CreateConfig) error {
	// Create a client for the system backend configured with the Vault token
	// used for the current cluster's PKI backend.
	sysBackend := s.VaultClient.Sys()

	// Mount a new PKI backend for the cluster, if it does not already exist.
	mounted, err := s.IsMounted(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !mounted {
		newMountConfig := &vaultclient.MountInput{
			Type:        "pki",
			Description: fmt.Sprintf("PKI backend for cluster ID '%s'", config.ClusterID),
			Config: vaultclient.MountConfigInput{
				MaxLeaseTTL: config.TTL,
			},
		}
		err = sysBackend.Mount(s.MountPKIPath(config.ClusterID), newMountConfig)
		if err != nil {
			return maskAny(err)
		}
	}

	// Create a client for the logical backend configured with the Vault token
	// used for the current cluster's root CA and role.
	logicalBackend := s.VaultClient.Logical()

	// Generate a certificate authority for the PKI backend, if it does not
	// already exist.
	generated, err := s.IsCAGenerated(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !generated {
		data := map[string]interface{}{
			"ttl":         config.TTL,
			"common_name": config.CommonName,
		}
		_, err = logicalBackend.Write(s.WriteCAPath(config.ClusterID), data)
		if err != nil {
			return maskAny(err)
		}
	}

	// Create a role for the mounted PKI backend, if it does not already exist.
	created, err := s.IsRoleCreated(config.ClusterID)
	if err != nil {
		return maskAny(err)
	}
	if !created {
		data := map[string]interface{}{
			"allowed_domains":    config.AllowedDomains,
			"allow_subdomains":   "true",
			"ttl":                config.TTL,
			"allow_bare_domains": config.AllowBareDomains,
		}

		_, err = logicalBackend.Write(s.WriteRolePath(config.ClusterID), data)
		if err != nil {
			return maskAny(err)
		}
	}

	return nil
}

// Path management.

func (s *service) ReadCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/cert/ca", clusterID)
}

func (s *service) MountPKIPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (s *service) ListMountsPath(clusterID string) string {
	return fmt.Sprintf("pki-%s", clusterID)
}

func (s *service) ListRolesPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/roles/", clusterID)
}

func (s *service) WriteCAPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/root/generate/internal", clusterID)
}

func (s *service) WriteRolePath(clusterID string) string {
	return fmt.Sprintf("pki-%s/roles/%s", clusterID, s.RoleName(clusterID))
}
