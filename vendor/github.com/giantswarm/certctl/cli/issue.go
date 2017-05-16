package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/cert-signer"
	"github.com/giantswarm/certctl/service/spec"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type issueFlags struct {
	VaultAddress string
	VaultToken   string
	VaultTLS     *vaultclient.TLSConfig

	// Cluster
	ClusterID string

	// Certificate
	CommonName string
	IPSANs     string
	AltNames   string
	TTL        string

	// Path
	CrtFilePath string
	KeyFilePath string
	CAFilePath  string
}

var (
	issueCmd = &cobra.Command{
		Use:   "issue",
		Short: "Generate signed certificates for a specific cluster.",
		Run:   issueRun,
	}

	newIssueFlags = &issueFlags{
		VaultTLS: &vaultclient.TLSConfig{},
	}
)

func init() {
	CLICmd.AddCommand(issueCmd)

	issueCmd.Flags().StringVar(&newIssueFlags.VaultAddress, "vault-addr", fromEnvToString(EnvVaultAddress, "http://127.0.0.1:8200"), "Address used to connect to Vault.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultToken, "vault-token", fromEnvToString(EnvVaultToken, ""), "Token used to authenticate against Vault.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultTLS.CACert, "vault-cacert", fromEnvToString(EnvVaultCACert, ""), "The path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultTLS.CAPath, "vault-capath", fromEnvToString(EnvVaultCAPath, ""), "The path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultTLS.ClientCert, "vault-client-cert", fromEnvToString(EnvVaultClientCert, ""), "The path to the certificate for Vault communication.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultTLS.ClientKey, "vault-client-key", fromEnvToString(EnvVaultClientKey, ""), "The path to the private key for Vault communication.")
	issueCmd.Flags().StringVar(&newIssueFlags.VaultTLS.TLSServerName, "vault-tls-server-name", fromEnvToString(EnvVaultTLSServerName, ""), "If set, is used to set the SNI host when connecting via TLS.")
	issueCmd.Flags().BoolVar(&newIssueFlags.VaultTLS.Insecure, "vault-tls-skip-verify", fromEnvBool(EnvVaultInsecure, false), "Do not verify TLS certificate.")

	issueCmd.Flags().StringVar(&newIssueFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new signed certificate for.")

	issueCmd.Flags().StringVar(&newIssueFlags.CommonName, "common-name", "", "Common name used to generate a new signed certificate for.")
	issueCmd.Flags().StringVar(&newIssueFlags.IPSANs, "ip-sans", "", "IPSANs used to generate a new signed certificate for.")
	issueCmd.Flags().StringVar(&newIssueFlags.AltNames, "alt-names", "", "Alternative names used to generate a new signed certificate for.")
	issueCmd.Flags().StringVar(&newIssueFlags.TTL, "ttl", "8640h", "TTL used to generate a new signed certificate for.") // 1 year

	issueCmd.Flags().StringVar(&newIssueFlags.CrtFilePath, "crt-file", "", "File path used to write the generated public key to.")
	issueCmd.Flags().StringVar(&newIssueFlags.KeyFilePath, "key-file", "", "File path used to write the generated private key to.")
	issueCmd.Flags().StringVar(&newIssueFlags.CAFilePath, "ca-file", "", "File path used to write the issuing root CA to.")
}

func issueValidate(newIssueFlags *issueFlags) error {
	if newIssueFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newIssueFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}
	if newIssueFlags.CommonName == "" {
		return maskAnyf(invalidConfigError, "--common-name must not be empty")
	}
	if newIssueFlags.CrtFilePath == "" {
		return maskAnyf(invalidConfigError, "--crt-file name must not be empty")
	}
	if newIssueFlags.KeyFilePath == "" {
		return maskAnyf(invalidConfigError, "--key-file name must not be empty")
	}
	if newIssueFlags.CAFilePath == "" {
		return maskAnyf(invalidConfigError, "--ca-file name must not be empty")
	}

	return nil
}

func issueRun(cmd *cobra.Command, args []string) {
	err := issueValidate(newIssueFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.Address = newIssueFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newIssueFlags.VaultToken
	newVaultFactoryConfig.TLS = newIssueFlags.VaultTLS
	newVaultFactory, err := vaultfactory.New(newVaultFactoryConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client and configure it with the provided admin token
	// through the factory.
	newVaultClient, err := newVaultFactory.NewClient()
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a certificate signer to generate a new signed certificate.
	newCertSignerConfig := certsigner.DefaultConfig()
	newCertSignerConfig.VaultClient = newVaultClient
	newCertSigner, err := certsigner.New(newCertSignerConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Generate a new signed certificate.
	newIssueConfig := spec.IssueConfig{
		ClusterID:  newIssueFlags.ClusterID,
		CommonName: newIssueFlags.CommonName,
		IPSANs:     newIssueFlags.IPSANs,
		AltNames:   newIssueFlags.AltNames,
		TTL:        newIssueFlags.TTL,
	}
	newIssueResponse, err := newCertSigner.Issue(newIssueConfig)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	err = os.MkdirAll(filepath.Dir(newIssueFlags.CrtFilePath), os.FileMode(0744))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = ioutil.WriteFile(newIssueFlags.CrtFilePath, []byte(newIssueResponse.Certificate), os.FileMode(0644))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = os.MkdirAll(filepath.Dir(newIssueFlags.KeyFilePath), os.FileMode(0744))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = ioutil.WriteFile(newIssueFlags.KeyFilePath, []byte(newIssueResponse.PrivateKey), os.FileMode(0644))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = os.MkdirAll(filepath.Dir(newIssueFlags.CAFilePath), os.FileMode(0744))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = ioutil.WriteFile(newIssueFlags.CAFilePath, []byte(newIssueResponse.IssuingCA), os.FileMode(0644))
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("Issued new signed certificate with the following serial number.\n")
	fmt.Printf("\n")
	fmt.Printf("    %s\n", newIssueResponse.SerialNumber)
	fmt.Printf("\n")
	fmt.Printf("Public key written to '%s'.\n", newIssueFlags.CrtFilePath)
	fmt.Printf("Private key written to '%s'.\n", newIssueFlags.KeyFilePath)
	fmt.Printf("Root CA written to '%s'.\n", newIssueFlags.CAFilePath)
}
