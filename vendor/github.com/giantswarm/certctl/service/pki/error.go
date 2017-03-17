package pki

import (
	"fmt"
	"strings"

	"github.com/juju/errgo"
)

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

func maskAnyf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	f = fmt.Sprintf("%s: %s", err.Error(), f)
	newErr := errgo.WithCausef(nil, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(1)

	return newErr
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
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
