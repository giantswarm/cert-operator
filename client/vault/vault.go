package vault

import (
	"net/http"
	"net/url"

	microerror "github.com/giantswarm/microkit/error"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"

	"github.com/giantswarm/cert-operator/flag"
)

type Config struct {
	// Settings.
	Flag  *flag.Flag
	Viper *viper.Viper
}

func NewClient(config Config) (*vaultapi.Client, error) {
	address := config.Viper.GetString(config.Flag.Service.Vault.Config.Address)
	token := config.Viper.GetString(config.Flag.Service.Vault.Config.Token)

	if address == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault address must not be empty")
	}

	// Check Vault address is valid.
	_, err := url.ParseRequestURI(address)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	if token == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "vault address must not be empty")
	}

	newClientConfig := vaultapi.DefaultConfig()
	newClientConfig.Address = address
	newClientConfig.HttpClient = http.DefaultClient

	newVaultClient, err := vaultapi.NewClient(newClientConfig)
	if err != nil {
		return nil, err
	}
	newVaultClient.SetToken(token)

	return newVaultClient, nil
}
