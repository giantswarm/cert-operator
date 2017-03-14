package vault

import (
	"net/http"

	vaultclient "github.com/hashicorp/vault/api"
)

type Config struct {
	// Dependencies.
	HTTPClient *http.Client

	Address string
	Token   string
}

func NewClient(config Config) (*vaultclient.Client, error) {
	newClientConfig := vaultclient.DefaultConfig()

	newClientConfig.Address = config.Address
	newClientConfig.HttpClient = config.HTTPClient
	newVaultClient, err := vaultclient.NewClient(newClientConfig)
	if err != nil {
		return nil, maskAny(err)
	}
	newVaultClient.SetToken(config.Token)

	return newVaultClient, nil
}
