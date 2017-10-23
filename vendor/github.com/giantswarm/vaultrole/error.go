package vaultrole

import (
	"strings"

	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

// IsNoVaultHandlerDefined asserts a dirty string matching against the error
// message provided by err. This is necessary due to the poor error handling
// design of the Vault library we are using.
func IsNoVaultHandlerDefined(err error) bool {
	cause := microerror.Cause(err)

	if cause != nil && strings.Contains(cause.Error(), "no handler for route") {
		return true
	}

	return false
}
