package vaultcrtv1

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	"github.com/giantswarm/vaultrole/vaultroletest"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func Test_Resource_VaultCrt_newDeleteChange(t *testing.T) {
	testCases := []struct {
		Obj            interface{}
		CurrentState   interface{}
		DesiredState   interface{}
		ExpectedSecret *apiv1.Secret
	}{
		// Test 0 ensures that zero value input results in zero value output.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState:   nil,
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: nil,
		},

		// Test 1 is the same as 0 but with initialized empty pointer values.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState:   &apiv1.Secret{},
			DesiredState:   &apiv1.Secret{},
			ExpectedSecret: &apiv1.Secret{},
		},

		// Test 2 ensures that the delete state is defined by the current state
		// since we want to remove the current state in case a delete event happens.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
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
			DesiredState: &apiv1.Secret{},
			ExpectedSecret: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
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
		result, err := newResource.newDeleteChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		secret := result.(*apiv1.Secret)
		if !reflect.DeepEqual(tc.ExpectedSecret, secret) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedSecret, secret)
		}
	}
}
