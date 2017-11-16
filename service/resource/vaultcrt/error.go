package vaultcrt

import (
	"github.com/giantswarm/microerror"
)

var missingAnnotationError = microerror.New("missing annotation")

// IsMissingAnnotation asserts missingAnnotationError.
func IsMissingAnnotation(err error) bool {
	return microerror.Cause(err) == missingAnnotationError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var wrongTypeError = microerror.New("wrong type")

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
