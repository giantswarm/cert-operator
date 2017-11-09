package vaultcrt

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultcrt/vaultcrttest"
	"github.com/giantswarm/vaultrole/vaultroletest"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func Test_Resource_VaultCrt_NewPatch(t *testing.T) {
	testCases := []struct {
		CurrentState      *apiv1.Secret
		DesiredState      *apiv1.Secret
		ExpectedPatchFunc func(*framework.Patch)
	}{
		// Test 0 shows that if current state is nil and desired is
		// not, then desired state should be created.
		{
			CurrentState: nil,
			DesiredState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
					Labels: map[string]string{
						"clusterID":        "al9qy",
						"clusterComponent": "worker",
					},
				},
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				secretToCreate :=
					&apiv1.Secret{
						ObjectMeta: apismetav1.ObjectMeta{
							Name: "al9qy-worker",
							Labels: map[string]string{
								"clusterID":        "al9qy",
								"clusterComponent": "worker",
							},
						},
					}

				patch.SetCreateChange(secretToCreate)
			},
		},

		// Test 1 shows that if current state is not nil and desired is nil
		{
			CurrentState: &apiv1.Secret{
				ObjectMeta: apismetav1.ObjectMeta{
					Name: "al9qy-worker",
					Labels: map[string]string{
						"clusterID":        "al9qy",
						"clusterComponent": "worker",
					},
				},
			},
			DesiredState: nil,
			ExpectedPatchFunc: func(patch *framework.Patch) {
				secretToDelete :=
					&apiv1.Secret{
						ObjectMeta: apismetav1.ObjectMeta{
							Name: "al9qy-worker",
							Labels: map[string]string{
								"clusterID":        "al9qy",
								"clusterComponent": "worker",
							},
						},
					}

				patch.SetDeleteChange(secretToDelete)
			},
		},

		// Test 2
		{
			CurrentState:      nil,
			DesiredState:      nil,
			ExpectedPatchFunc: func(patch *framework.Patch) {},
		},

		// Test 3
		{
			CurrentState:      &apiv1.Secret{},
			DesiredState:      &apiv1.Secret{},
			ExpectedPatchFunc: func(patch *framework.Patch) {},
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
		patch, err := newResource.NewPatch(context.TODO(), nil, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}

		expectedPatch := framework.NewPatch()
		tc.ExpectedPatchFunc(expectedPatch)

		if !reflect.DeepEqual(expectedPatch, patch) {
			t.Fatalf("case %d expected %#v got %#v", i, *expectedPatch, *patch)
		}
	}
}
