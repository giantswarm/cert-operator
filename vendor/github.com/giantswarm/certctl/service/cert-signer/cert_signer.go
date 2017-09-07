package certsigner

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/role"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/microerror"
)

// Config represents the configuration used to create a new certificate signer.
type Config struct {
	// Dependencies.
	VaultClient *vaultclient.Client
}

// DefaultConfig provides a default configuration to create a certificate
// signer.
func DefaultConfig() Config {
	newClientConfig := vaultclient.DefaultConfig()
	newClientConfig.Address = "http://127.0.0.1:8200"
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		panic(err)
	}

	newConfig := Config{
		// Dependencies.
		VaultClient: newVaultClient,
	}

	return newConfig
}

// New creates a new configured certificate signer.
func New(config Config) (spec.CertSigner, error) {
	newCertSigner := &certSigner{
		Config: config,
	}

	// Dependencies.
	if newCertSigner.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "Vault client must not be empty")
	}

	return newCertSigner, nil
}

type certSigner struct {
	Config
}

func (cs *certSigner) Issue(config spec.IssueConfig) (spec.IssueResponse, error) {
	var roleService role.Service
	var err error
	{
		roleServiceConfig := role.DefaultConfig()
		roleServiceConfig.VaultClient = cs.VaultClient
		roleServiceConfig.PKIMountpoint = fmt.Sprintf("pki-%s", config.ClusterID)
		roleService, err = role.New(roleServiceConfig)
		if err != nil {
			return spec.IssueResponse{}, microerror.Mask(err)
		}
	}

	// Ensure a role exists exists that can issue a cert with the desired Organizations
	// before trying to issue a cert.
	isRoleCreated, err := roleService.IsRoleCreated(roleName(config.ClusterID, config.Organizations))
	if err != nil {
		return spec.IssueResponse{}, microerror.Mask(err)
	}

	if !isRoleCreated {
		createRoleParams := role.CreateParams{
			AllowBareDomains: config.AllowBareDomains,
			AllowedDomains:   config.AllowedDomains,
			AllowSubdomains:  true,
			TTL:              config.RoleTTL,
			Name:             roleName(config.ClusterID, config.Organizations),
			Organizations:    config.Organizations,
		}

		err = roleService.Create(createRoleParams)
		if err != nil {
			return spec.IssueResponse{}, microerror.Mask(err)
		}
	}

	// Create a client for issuing a new signed certificate.
	logicalStore := cs.VaultClient.Logical()

	// Generate a certificate for the PKI backend signed by the certificate
	// authority associated with the configured cluster ID.
	data := map[string]interface{}{
		"ttl":         config.TTL,
		"common_name": config.CommonName,
		"ip_sans":     config.IPSANs,
		"alt_names":   config.AltNames,
	}

	secret, err := logicalStore.Write(cs.SignedPath(config.ClusterID, config.Organizations), data)
	if err != nil {
		return spec.IssueResponse{}, microerror.Mask(err)
	}

	// Collect the certificate data from the secret response.
	vCrt, ok := secret.Data["certificate"]
	if !ok {
		return spec.IssueResponse{}, microerror.Maskf(keyPairNotFoundError, "public key missing")
	}
	crt := vCrt.(string)
	vKey, ok := secret.Data["private_key"]
	if !ok {
		return spec.IssueResponse{}, microerror.Maskf(keyPairNotFoundError, "private key missing")
	}
	key := vKey.(string)
	vCA, ok := secret.Data["issuing_ca"]
	if !ok {
		return spec.IssueResponse{}, microerror.Maskf(keyPairNotFoundError, "root CA missing")
	}
	ca := vCA.(string)
	vSerial, ok := secret.Data["serial_number"]
	if !ok {
		return spec.IssueResponse{}, microerror.Maskf(keyPairNotFoundError, "root CA missing")
	}
	serial := vSerial.(string)

	newIssueResponse := spec.IssueResponse{
		Certificate:  crt,
		PrivateKey:   key,
		IssuingCA:    ca,
		SerialNumber: serial,
	}

	return newIssueResponse, nil
}

func (cs *certSigner) SignedPath(clusterID string, organizations string) string {
	return fmt.Sprintf("pki-%s/issue/%s", clusterID, roleName(clusterID, organizations))
}

func roleName(clusterID string, organizations string) string {
	if organizations == "" {
		// If organizations isn't set, use the role that was created when the PKI
		// for this cluster was first setup.
		return fmt.Sprintf("role-%s", clusterID)
	}

	// Compute a url-safe hash of the organizations that stays the same regardless
	// of the order of the organizations supplied.
	return fmt.Sprintf("role-org-%s", computeRoleHash(organizations))
}

// computeRoleHash computes a hash for the role that can issue these organizations.
// Since we want to reuse roles when possible, we should try to make sure that
// the same list of organizations returns the same hash (regardless of the order).
// The reason we don't use just the organizations that the user provided is because
// that could potentially be a very long list, or otherwise contain characters
// that are not allowed in URLs.
func computeRoleHash(organizations string) string {
	// Sort organizations alphabetically
	organizationsSlice := strings.Split(organizations, ",")
	sort.Strings(organizationsSlice)
	organizations = strings.Join(organizationsSlice, ",")

	h := sha1.New()
	h.Write([]byte(organizations))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func writeRolePath(clusterID string, organizations string) string {
	return fmt.Sprintf("pki-%s/roles/%s", clusterID, roleName(clusterID, organizations))
}
