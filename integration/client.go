// +build k8srequired

package integration

import (
	giantclientset "github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	corev1alpha1 "github.com/giantswarm/apiextensions/pkg/clientset/versioned/typed/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/e2e-harness/pkg/harness"
)

type k8sClient interface {
	CoreV1() corev1.CoreV1Interface
	ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface
	CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface
}

type k8sClientImpl struct {
	cs   kubernetes.Interface
	gsCs *giantclientset.Clientset
}

func (k *k8sClientImpl) CoreV1() corev1.CoreV1Interface {
	return k.cs.CoreV1()
}

func (k *k8sClientImpl) ExtensionsV1beta1() extensionsv1beta1.ExtensionsV1beta1Interface {
	return k.cs.ExtensionsV1beta1()
}

func (k *k8sClientImpl) CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface {
	return k.gsCs.CoreV1alpha1()
}

func getK8sClient() (*k8sClientImpl, error) {
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

	client := &k8sClientImpl{
		cs:   cs,
		gsCs: gsCs,
	}
	return client, nil
}
