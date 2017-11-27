package integration

import "github.com/giantswarm/microerror"

var waitTimeoutError = microerror.New("waitTimeout")

var tooManyResultsError = microerror.New("too many results")

var unexpectedStatusPhase = microerror.New("unexpected status phase")
