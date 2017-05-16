package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	vaultclient "github.com/hashicorp/vault/api"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type inspectFlags struct {
	// Vault
	VaultAddress string
	VaultToken   string
	VaultTLS     *vaultclient.TLSConfig

	// Cluster
	ClusterID string
}

var (
	inspectCmd = &cobra.Command{
		Use:   "inspect",
		Short: "Inspect a Vault PKI backend including all necessary requirements.",
		Run:   inspectRun,
	}

	newInspectFlags = &inspectFlags{
		VaultTLS: &vaultclient.TLSConfig{},
	}
)

func init() {
	CLICmd.AddCommand(inspectCmd)

	inspectCmd.Flags().StringVar(&newInspectFlags.VaultAddress, "vault-addr", fromEnvToString(EnvVaultAddress, "http://127.0.0.1:8200"), "Address used to connect to Vault.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultToken, "vault-token", fromEnvToString(EnvVaultToken, ""), "Token used to authenticate against Vault.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultTLS.CACert, "vault-cacert", fromEnvToString(EnvVaultCACert, ""), "The path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultTLS.CAPath, "vault-capath", fromEnvToString(EnvVaultCAPath, ""), "The path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultTLS.ClientCert, "vault-client-cert", fromEnvToString(EnvVaultClientCert, ""), "The path to the certificate for Vault communication.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultTLS.ClientKey, "vault-client-key", fromEnvToString(EnvVaultClientKey, ""), "The path to the private key for Vault communication.")
	inspectCmd.Flags().StringVar(&newInspectFlags.VaultTLS.TLSServerName, "vault-tls-server-name", fromEnvToString(EnvVaultTLSServerName, ""), "If set, is used to set the SNI host when connecting via TLS.")
	inspectCmd.Flags().BoolVar(&newInspectFlags.VaultTLS.Insecure, "vault-tls-skip-verify", fromEnvBool(EnvVaultInsecure, false), "Do not verify TLS certificate.")

	inspectCmd.Flags().StringVar(&newInspectFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new root CA for.")
}

func inspectValidate(newInspectFlags *inspectFlags) error {
	if newInspectFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newInspectFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}

	return nil
}

func inspectRun(cmd *cobra.Command, args []string) {
	err := inspectValidate(newInspectFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.Address = newInspectFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newInspectFlags.VaultToken
	newVaultFactoryConfig.TLS = newInspectFlags.VaultTLS
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

	// Create a PKI controller to check for PKI backend specific operations.
	var pkiService pki.Service
	{
		pkiConfig := pki.DefaultServiceConfig()
		pkiConfig.VaultClient = newVaultClient
		pkiService, err = pki.NewService(pkiConfig)
		if err != nil {
			log.Fatalf("%#v\n", maskAny(err))
		}
	}

	// Create a token generator to check for token specific operations.
	var tokenService token.Service
	{
		tokenConfig := token.DefaultServiceConfig()
		tokenConfig.VaultClient = newVaultClient
		tokenService, err = token.NewService(tokenConfig)
		if err != nil {
			log.Fatalf("%#v\n", maskAny(err))
		}
	}

	mounted, err := pkiService.IsMounted(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	generated, err := pkiService.IsCAGenerated(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	roleCreated, err := pkiService.IsRoleCreated(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	policyCreated, err := tokenService.IsPolicyCreated(newInspectFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("Inspecting cluster for ID '%s':\n", newInspectFlags.ClusterID)
	fmt.Printf("\n")
	fmt.Printf("    PKI backend mounted: %t\n", mounted)
	fmt.Printf("    Root CA generated:   %t\n", generated)
	fmt.Printf("    PKI role created:    %t\n", roleCreated)
	fmt.Printf("    PKI policy created:  %t\n", policyCreated)
	fmt.Printf("\n")
	fmt.Printf("Tokens may have been generated for this cluster. Created tokens\n")
	fmt.Printf("cannot be shown as they are secret. Information about these\n")
	fmt.Printf("secrets needs to be looked up directly from the location of the\n")
	fmt.Printf("cluster's installation.\n")
}
