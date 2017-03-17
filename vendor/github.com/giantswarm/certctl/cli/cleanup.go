package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/giantswarm/certctl/service/pki"
	"github.com/giantswarm/certctl/service/token"
	"github.com/giantswarm/certctl/service/vault-factory"
)

type cleanupFlags struct {
	// Vault
	VaultAddress string
	VaultToken   string

	// Cluster
	ClusterID string
}

var (
	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup a Vault PKI backend including all necessary requirements.",
		Run:   cleanupRun,
	}

	newCleanupFlags = &cleanupFlags{}
)

func init() {
	CLICmd.AddCommand(cleanupCmd)

	cleanupCmd.Flags().StringVar(&newCleanupFlags.VaultAddress, "vault-addr", fromEnv("VAULT_ADDR", "http://127.0.0.1:8200"), "Address used to connect to Vault.")
	cleanupCmd.Flags().StringVar(&newCleanupFlags.VaultToken, "vault-token", fromEnv("VAULT_TOKEN", ""), "Token used to authenticate against Vault.")

	cleanupCmd.Flags().StringVar(&newCleanupFlags.ClusterID, "cluster-id", "", "Cluster ID used to generate a new root CA for.")
}

func cleanupValidate(newCleanupFlags *cleanupFlags) error {
	if newCleanupFlags.VaultToken == "" {
		return maskAnyf(invalidConfigError, "Vault token must not be empty")
	}
	if newCleanupFlags.ClusterID == "" {
		return maskAnyf(invalidConfigError, "cluster ID must not be empty")
	}

	return nil
}

func cleanupRun(cmd *cobra.Command, args []string) {
	err := cleanupValidate(newCleanupFlags)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	// Create a Vault client factory.
	newVaultFactoryConfig := vaultfactory.DefaultConfig()
	newVaultFactoryConfig.HTTPClient = &http.Client{}
	newVaultFactoryConfig.Address = newCleanupFlags.VaultAddress
	newVaultFactoryConfig.AdminToken = newCleanupFlags.VaultToken
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

	// Create a PKI controller to cleanup PKI backend specific operations.
	var pkiService pki.Service
	{
		pkiConfig := pki.DefaultServiceConfig()
		pkiConfig.VaultClient = newVaultClient
		pkiService, err = pki.NewService(pkiConfig)
		if err != nil {
			log.Fatalf("%#v\n", maskAny(err))
		}
	}

	// Create a token generator to cleanup token specific operations.
	var tokenService token.Service
	{
		tokenConfig := token.DefaultServiceConfig()
		tokenConfig.VaultClient = newVaultClient
		tokenService, err = token.NewService(tokenConfig)
		if err != nil {
			log.Fatalf("%#v\n", maskAny(err))
		}
	}

	err = pkiService.Delete(newCleanupFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}
	err = tokenService.DeletePolicy(newCleanupFlags.ClusterID)
	if err != nil {
		log.Fatalf("%#v\n", maskAny(err))
	}

	fmt.Printf("Cleaning up cluster for ID '%s':\n", newCleanupFlags.ClusterID)
	fmt.Printf("\n")
	fmt.Printf("    - PKI backend unmounted\n")
	fmt.Printf("    - Root CA deleted\n")
	fmt.Printf("    - PKI role deleted\n")
	fmt.Printf("    - PKI policy deleted\n")
	fmt.Printf("\n")
	fmt.Printf("Tokens may have been generated for this cluster. Created tokens\n")
	fmt.Printf("cannot be revoked here as they are secret. Tokens need to be\n")
	fmt.Printf("revoked manually. In case a cluster with the same ID will be\n")
	fmt.Printf("generated, tokens generated for this cluster will be able to\n")
	fmt.Printf("access this new cluster again. Information about these secrets\n")
	fmt.Printf("needs to be looked up directly from the location of the cluster's\n")
	fmt.Printf("installation.\n")
}
