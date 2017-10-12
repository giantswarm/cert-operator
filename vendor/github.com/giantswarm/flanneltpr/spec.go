package flanneltpr

import "github.com/giantswarm/flanneltpr/spec"

type Spec struct {
	Cluster spec.Cluster `json:"cluster" yaml:"cluster"`

	Bridge  spec.Bridge  `json:"bridge" yaml:"bridge"`
	Flannel spec.Flannel `json:"flannel" yaml:"flannel"`
}
