package certsigner

import (
	"fmt"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/spec"
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
		return nil, maskAnyf(invalidConfigError, "Vault client must not be empty")
	}

	return newCertSigner, nil
}

type certSigner struct {
	Config
}

func (cs *certSigner) Issue(config spec.IssueConfig) (spec.IssueResponse, error) {
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

	secret, err := logicalStore.Write(cs.SignedPath(config.ClusterID), data)
	if err != nil {
		return spec.IssueResponse{}, maskAny(err)
	}

	// Collect the certificate data from the secret response.
	vCrt, ok := secret.Data["certificate"]
	if !ok {
		return spec.IssueResponse{}, maskAnyf(keyPairNotFoundError, "public key missing")
	}
	crt := vCrt.(string)
	vKey, ok := secret.Data["private_key"]
	if !ok {
		return spec.IssueResponse{}, maskAnyf(keyPairNotFoundError, "private key missing")
	}
	key := vKey.(string)
	vCA, ok := secret.Data["issuing_ca"]
	if !ok {
		return spec.IssueResponse{}, maskAnyf(keyPairNotFoundError, "root CA missing")
	}
	ca := vCA.(string)
	vSerial, ok := secret.Data["serial_number"]
	if !ok {
		return spec.IssueResponse{}, maskAnyf(keyPairNotFoundError, "root CA missing")
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

func (cs *certSigner) SignedPath(clusterID string) string {
	return fmt.Sprintf("pki-%s/issue/role-%s", clusterID, clusterID)
}
