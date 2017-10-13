package vaultpki

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultpki/vaultpkitest"
)

func Test_Resource_VaultPKI_GetCreateState(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		Cur           interface{}
		Des           interface{}
		ExpectedState VaultPKIState
	}{
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Cur:           VaultPKIState{},
			Des:           VaultPKIState{},
			ExpectedState: VaultPKIState{},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Cur: VaultPKIState{},
			Des: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
			},
			ExpectedState: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
			},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Cur: VaultPKIState{
				BackendExists: false,
				CAExists:      true,
			},
			Des: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
			},
			ExpectedState: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
			},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Cur: VaultPKIState{
				BackendExists: false,
				CAExists:      false,
			},
			Des: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
			},
			ExpectedState: VaultPKIState{
				BackendExists: true,
				CAExists:      true,
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
		result, err := newResource.GetCreateState(context.TODO(), tc.Obj, tc.Cur, tc.Des)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i+1, tc.ExpectedState, r)
		}
	}
}
