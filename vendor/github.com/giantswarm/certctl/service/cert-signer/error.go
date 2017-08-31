package certsigner

import (
	"strings"

	"github.com/giantswarm/microerror"

	"github.com/juju/errgo"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var keyPairNotFoundError = microerror.New("key pair not found")

// IsKeyPairNotFound asserts keyPairNotFoundError.
func IsKeyPairNotFound(err error) bool {
	return microerror.Cause(err) == keyPairNotFoundError
}

// IsNoVaultHandlerDefined asserts a dirty string matching against the error
// message provided by err. This is necessary due to the poor error handling
// design of the Vault library we are using.
func IsNoVaultHandlerDefined(err error) bool {
	cause := errgo.Cause(err)

	if cause != nil && strings.Contains(cause.Error(), "no handler for route") {
		return true
	}

	return false
}
