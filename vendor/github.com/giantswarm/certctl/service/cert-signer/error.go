package certsigner

import (
	"github.com/giantswarm/microerror"
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
