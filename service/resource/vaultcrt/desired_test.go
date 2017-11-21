package vaultcrt

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/certificatetpr/spec"
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
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures the desired state is always the same placeholder state.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID:        "foobar",
					ClusterComponent: "api",
					VersionBundle: spec.VersionBundle{
						Version: "0.1.0",
					},
				},
			},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Annotations: map[string]string{
						UpdateTimestampAnnotation:      (time.Time{}).Format(UpdateTimestampLayout),
						VersionBundleVersionAnnotation: "0.1.0",
					},
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
		},

		// Test 1 is the same as 0 but with a different custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID:        "al9qy",
					ClusterComponent: "worker",
					VersionBundle: spec.VersionBundle{
						Version: "0.2.0",
					},
				},
			},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
					Annotations: map[string]string{
						UpdateTimestampAnnotation:      (time.Time{}).Format(UpdateTimestampLayout),
						VersionBundleVersionAnnotation: "0.2.0",
					},
					Labels: map[string]string{
						"clusterID":        "al9qy",
						"clusterComponent": "worker",
					},
				},
				StringData: map[string]string{
					"ca":  "",
					"crt": "",
					"key": "",
				},
			},
		},
	}

	var err error
	var newResource *Resource
	{
		c := DefaultConfig()

		c.CurrentTimeFactory = func() time.Time { return time.Time{} }
		c.K8sClient = fake.NewSimpleClientset()
		c.Logger = microloggertest.New()
		c.VaultCrt = vaultcrttest.New()
		c.VaultRole = vaultroletest.New()

		c.ExpirationThreshold = 24 * time.Hour
		c.Namespace = "default"

		newResource, err = New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}
