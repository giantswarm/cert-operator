package label

import (
	"k8s.io/apimachinery/pkg/labels"

	"github.com/giantswarm/cert-operator/v3/pkg/project"
)

const (
	Cluster         = "giantswarm.io/cluster"
	OperatorVersion = "cert-operator.giantswarm.io/version"
)

func AppVersionSelector() labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		OperatorVersion: project.Version(),
	})
}

// KubeconfigSelector selects all certconfigs that use the special version `0.0.0`.
func KubeconfigSelector() labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		OperatorVersion: project.ManagementClusterAppVersion(),
	})
}
