package spec

import "github.com/giantswarm/flanneltpr/spec/flannel"

// Flannel holds the configuration to run falnneld and create etcd VNI
// configuration for it.
// https://github.com/coreos/flannel/blob/master/Documentation/configuration.md
type Flannel struct {
	// Docker describes the docker image running flanneld.
	Docker flannel.Docker `json:"docker" yaml:"docker"`
	// Spec contains flannel configuration.
	Spec flannel.Spec `json:"spec" yaml:"spec"`
}
