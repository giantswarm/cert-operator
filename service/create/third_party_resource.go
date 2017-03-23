package create

import (
	"github.com/giantswarm/certificatetpr"
	microerror "github.com/giantswarm/microkit/error"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// createTPR ensures certificatetpr exists on the cluster.
func (s *Service) createTPR() error {
	tpr := &v1beta1.ThirdPartyResource{
		ObjectMeta: v1.ObjectMeta{
			Name: certificatetpr.Name,
		},
		Versions: []v1beta1.APIVersion{
			{Name: "v1"},
		},
		Description: "Managed certificates on Kubernetes clusters",
	}
	_, err := s.Config.K8sClient.Extensions().ThirdPartyResources().Create(tpr)
	if errors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
