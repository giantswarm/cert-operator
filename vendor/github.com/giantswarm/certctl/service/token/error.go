package token

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var policyAlreadyExistsError = microerror.New("policy already exists")

// IsPolicyAlreadyExists asserts policyAlreadyExistsError.
func IsPolicyAlreadyExists(err error) bool {
	return microerror.Cause(err) == policyAlreadyExistsError
}
