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

func Test_Resource_VaultPKI_NewDeleteChange(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		CurrentState  interface{}
		DesiredState  interface{}
		ExpectedState VaultPKIState
	}{
		// Test 0 ensures that zero value input results in zero value output.
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

		// Test 1 ensures that any input results in zero value output because
		// deletion of PKI backends is not allowed. Thus delete state will always be
		// empty.
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

		// Test 2 is the same as 1 but with different input values.
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
		result, err := newResource.newDeleteChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedState, r)
		}
	}
}
