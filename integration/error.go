package integration

import "github.com/giantswarm/microerror"

var waitTimeoutError = &microerror.Error{
	Kind: "waitTimeoutError",
}

var tooManyResultsError = &microerror.Error{
	Kind: "tooManyResultsError",
}

var unexpectedStatusPhase = &microerror.Error{
	Kind: "unexpectedStatusPhase",
}
