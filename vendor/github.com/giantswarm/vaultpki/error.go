package vaultpki

import (
	"strings"

	"github.com/giantswarm/microerror"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
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
