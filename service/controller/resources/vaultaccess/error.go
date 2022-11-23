package vaultaccess

import (
	"strings"

	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var vaultAccessError = &microerror.Error{
	Kind: "vaultAccessError",
}

// IsVaultAccess asserts vaultAccessError. The matcher also asserts errors
// caused by situations in which Vault is updated strategically and thus
// temporarily replies with HTTP responses. In such cases we intend to cancel
// reconciliation and wait until Vault is fully operational again.
//
//	Get https://vault.g8s.amag.ch:8200/v1/sys/mounts: http: server gave HTTP response to HTTPS client
func IsVaultAccess(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if strings.Contains(c.Error(), "server gave HTTP response to HTTPS client") {
		return true
	}

	if c == vaultAccessError {
		return true
	}

	return false
}
