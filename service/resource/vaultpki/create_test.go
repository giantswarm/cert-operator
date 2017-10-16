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

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{},
			DesiredState: VaultPKIState{},
			ExpectedState: VaultPKIState{
				BackendMissing: false,
				CAMissing:      false,
			},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{
				BackendMissing: false,
				CAMissing:      true,
			},
			DesiredState: VaultPKIState{},
			ExpectedState: VaultPKIState{
				BackendMissing: false,
				CAMissing:      true,
			},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{
				BackendMissing: true,
				CAMissing:      false,
			},
			DesiredState: VaultPKIState{},
			ExpectedState: VaultPKIState{
				BackendMissing: true,
				CAMissing:      false,
			},
		},

		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			CurrentState: VaultPKIState{
				BackendMissing: true,
				CAMissing:      true,
			},
			DesiredState: VaultPKIState{},
			ExpectedState: VaultPKIState{
				BackendMissing: true,
				CAMissing:      true,
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
		result, err := newResource.GetCreateState(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i+1, tc.ExpectedState, r)
		}
	}
}
