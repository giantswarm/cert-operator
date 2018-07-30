package vaultcrt

import (
	"github.com/giantswarm/microerror"
)

var missingAnnotationError = &microerror.Error{
	Docs: "https://github.com/giantswarm/ops-recipes",
	Kind: "missingAnnotationError",
}

// IsMissingAnnotation asserts missingAnnotationError.
func IsMissingAnnotation(err error) bool {
	return microerror.Cause(err) == missingAnnotationError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
