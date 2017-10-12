package spec

import "github.com/giantswarm/flanneltpr/spec/bridge"

// Bridge holds information used for the execution of the Docker entrypoint of
// https://github.com/giantswarm/k8s-network-bridge.
type Bridge struct {
	// Docker describes the docker image running k8s-network-bridge.
	Docker bridge.Docker `json:"docker" yaml:"docker"`
	// Spec contains network bridge configuration.
	Spec bridge.Spec `json:"spec" yaml:"spec"`
}
