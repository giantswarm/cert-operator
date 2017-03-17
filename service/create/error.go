package create

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var keyPairNotFoundError = errgo.New("key pair not found")

// IsKeyPairNotFound asserts keyPairNotFoundError.
func IsKeyPairNotFound(err error) bool {
	return errgo.Cause(err) == keyPairNotFoundError
}
