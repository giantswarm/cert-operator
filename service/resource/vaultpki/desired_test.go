package vaultpki

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultpki/vaultpkitest"
)

func Test_Resource_VaultPKI_GetDesiredState(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		Deleted       bool
		ExpectedState VaultPKIState
	}{
		// test 0 ensures the desired state for created/updated object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Deleted: false,
			ExpectedState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
		},

		// test 1 is the same as 0 but with a different custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "al9qy",
				},
			},
			Deleted: false,
			ExpectedState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
		},

		// test 2 ensures the desired state is always the same for deleted custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "foobar",
				},
			},
			Deleted: true,
			ExpectedState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
		},

		// test 3 is the same as 2 but with a different custom object.
		{
			Obj: &certificatetpr.CustomObject{
				Spec: certificatetpr.Spec{
					ClusterID: "al9qy",
				},
			},
			Deleted: true,
			ExpectedState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
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
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj, tc.Deleted)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedState, r)
		}
	}
}
