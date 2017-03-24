package vault

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/giantswarm/cert-operator/flag"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		token         string
		expectedError bool
	}{
		{
			name:          "Specify vault address and token. It should return a vault client.",
			address:       "http://localhost:8200",
			token:         "auth-token",
			expectedError: false,
		},
		{
			name:          "Specify vault address but no token. It should return an error.",
			address:       "http://localhost:8200",
			token:         "",
			expectedError: true,
		},
		{
			name:          "Specify a vault token but no address. It should return an error.",
			address:       "",
			token:         "auth-token",
			expectedError: true,
		},
		{
			name:          "Specify an invalid vault address. It should return an error.",
			address:       "http//invalid-address",
			token:         "auth-token",
			expectedError: true,
		},
	}

	for _, tc := range tests {
		f := flag.New()
		v := viper.New()

		v.Set(f.Service.Vault.Config.Address, tc.address)
		v.Set(f.Service.Vault.Config.Token, tc.token)

		config := Config{
			Flag:  f,
			Viper: v,
		}

		vaultClient, err := NewClient(config)
		if tc.expectedError {
			assert.Error(t, err, fmt.Sprintf("[%s] An error was expected", tc.name))
			continue
		} else {
			assert.NotNil(t, vaultClient, "Vault client is nil but should not be")
			assert.Nil(t, err, "Unexpected error creating vault client")
		}
	}
}
