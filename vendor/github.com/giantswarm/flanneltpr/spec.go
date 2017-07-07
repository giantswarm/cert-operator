package flanneltpr

import (
	"github.com/giantswarm/flanneltpr/host"
	"github.com/giantswarm/flanneltpr/network"
)

type Spec struct {
	Host host.Host `json:"host" yaml:"host"`
	// Namespace is the namespace of the guest cluster to watch for. E.g. Flannel
	// bridges should not be removed until there is any workload running within
	// the guest cluster's namespace.
	Namespace string          `json:"namespace" yaml:"namespace"`
	Network   network.Network `json:"network" yaml:"network"`
}
