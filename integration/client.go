// +build k8srequired

package integration

import (
	giantclientset "github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/e2e-harness/pkg/harness"
)

type clients struct {
	K8sCs kubernetes.Interface
	GsCs  *giantclientset.Clientset
}

func newClients() (*clients, error) {
	config, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	gsCs, err := giantclientset.NewForConfig(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	clients := &clients{
		K8sCs: cs,
		GsCs:  gsCs,
	}
	return clients, nil
}
