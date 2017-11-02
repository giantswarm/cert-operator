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

func Test_Resource_VaultCrt_newCreateChange(t *testing.T) {
	testCases := []struct {
		Obj            interface{}
		CurrentState   interface{}
		DesiredState   interface{}
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures a non-nil current state results in the create state to be
		// empty.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState:   &apiv1.Secret{},
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: nil,
		},

		// Test 1 is the same 1 but with different content for the current state.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
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
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: nil,
		},

		// Test 2 ensures an empty current state results in a create state that
		// equals the desired state. NOTE that the secret data is extended with
		// actual certificate content, which in this case is some fake content from
		// the fake VaultCrt service.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: nil,
			DesiredState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
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
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "foobar-api",
					Labels: map[string]string{
						"clusterID":        "foobar",
						"clusterComponent": "api",
					},
				},
				StringData: map[string]string{
					"ca":  "test CA",
					"crt": "test crt",
					"key": "test key",
				},
			},
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
		result, err := newResource.newCreateChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}
