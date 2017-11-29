package vaultpkiv1

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultpki/vaultpkitest"
	vaultapi "github.com/hashicorp/vault/api"
)

func Test_Resource_VaultPKI_GetDesiredState(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		ExpectedState VaultPKIState
	}{
		// test 0 ensures the desired state is always the same placeholder state.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			ExpectedState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
		},

		// test 1 is the same as 0 but with a different custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "al9qy",
				},
			},
			ExpectedState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
		},
	}

	var err error
	var newResource *Resource
	{
		resourceConfig := DefaultConfig()

		resourceConfig.Logger = microloggertest.New()
		resourceConfig.VaultPKI = vaultpkitest.New()

		newResource, err = New(resourceConfig)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedState, r)
		}
	}
}
