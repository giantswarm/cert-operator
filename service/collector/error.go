package collector

import (
	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
)

var alreadyRegisteredError = &microerror.Error{
	Kind: "alreadyRegisteredError",
}

// IsAlreadyRegisteredError asserts alreadyRegisteredError.
func IsAlreadyRegisteredError(err error) bool {
	c := microerror.Cause(err)
	_, ok := c.(prometheus.AlreadyRegisteredError)
	if ok {
		return true
	}
	if c == alreadyRegisteredError {
		return true
	}

	return false
}

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

var wrongTypeError = &microerror.Error{
	Kind: "wrongTypeError",
}

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
