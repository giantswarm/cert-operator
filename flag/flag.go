package flag

import (
	"github.com/giantswarm/microkit/flag"

	"github.com/giantswarm/cert-operator/flag/kubernetes"
)

type Flag struct {
	Kubernetes kubernetes.Kubernetes
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
