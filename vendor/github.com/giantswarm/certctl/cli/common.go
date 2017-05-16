package cli

import (
	"log"
	"os"
	"strconv"
)

const (
	EnvVaultAddress       = "VAULT_ADDR"
	EnvVaultCACert        = "VAULT_CACERT"
	EnvVaultCAPath        = "VAULT_CAPATH"
	EnvVaultClientCert    = "VAULT_CLIENT_CERT"
	EnvVaultClientKey     = "VAULT_CLIENT_KEY"
	EnvVaultInsecure      = "VAULT_SKIP_VERIFY"
	EnvVaultTLSServerName = "VAULT_TLS_SERVER_NAME"
	EnvVaultToken         = "VAULT_TOKEN"
)

func fromEnvToString(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}

func fromEnvBool(key string, def bool) bool {
	if value := os.Getenv(key); value != "" {
		parsedValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Cannot parse %s: %s", key, err)
		}

		return parsedValue
	}

	return def
}
