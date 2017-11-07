package vaultcrt

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	"github.com/giantswarm/vaultrole/vaultroletest"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func Test_Resource_VaultCrt_GetDesiredState(t *testing.T) {
	testCases := []struct {
		Obj            interface{}
		Deleted        bool
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures the desired is a secret for created/updated
		// custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID:        "foobar",
					ClusterComponent: "api",
				},
			},
			Deleted: false,
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
			},
		},

		// Test 1 is the same as 0 but with a different custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID:        "al9qy",
					ClusterComponent: "worker",
				},
			},
			Deleted: false,
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
					Labels: map[string]string{
						"clusterID":        "al9qy",
						"clusterComponent": "worker",
					},
				},
			},
		},

		// Test 2 ensures desired state is nil for deleted custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID:        "whatever",
					ClusterComponent: "worker",
				},
			},
			Deleted:        true,
			ExpectedSecret: nil,
		},
	}

	var err error
	var newResource *Resource
	{
		resourceConfig := DefaultConfig()

		resourceConfig.K8sClient = fake.NewSimpleClientset()
		resourceConfig.Logger = microloggertest.New()
		resourceConfig.VaultCrt = vaultcrttest.New()
		resourceConfig.VaultRole = vaultroletest.New()

		newResource, err = New(resourceConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj, tc.Deleted)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}
