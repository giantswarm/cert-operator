package controller

import (
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var noKindError = &microerror.Error{
	Kind: "noKindError",
}

// IsNoKind asserts noKindError.
func IsNoKind(err error) bool {
	c := microerror.Cause(err)

	_, ok := c.(*meta.NoKindMatchError)
	if ok {
		return true
	}

	if c == noKindError {
		return true
	}

	return false
}
