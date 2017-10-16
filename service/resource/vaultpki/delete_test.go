package vaultpki

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultpki/vaultpkitest"
	vaultapi "github.com/hashicorp/vault/api"
)

func Test_Resource_VaultPKI_GetDeleteState(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		CurrentState  interface{}
		DesiredState  interface{}
		ExpectedState VaultPKIState
	}{
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState:  VaultPKIState{},
			DesiredState:  VaultPKIState{},
			ExpectedState: VaultPKIState{},
		},

		// NOTE deletion is not allowed so delete state will always be false.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
			DesiredState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
			ExpectedState: VaultPKIState{
				Backend:       nil,
				CACertificate: "",
			},
		},

		// NOTE deletion is not allowed so delete state will always be false.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{
				Backend:       nil,
				CACertificate: "",
			},
			DesiredState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
			ExpectedState: VaultPKIState{
				Backend:       nil,
				CACertificate: "",
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
		result, err := newResource.GetDeleteState(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i+1, tc.ExpectedState, r)
		}
	}
}
