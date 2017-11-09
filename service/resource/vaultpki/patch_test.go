package vaultpki

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/vaultpki/vaultpkitest"
)

func Test_Resource_VaultPKI_NewPatch(t *testing.T) {
	testCases := []struct {
		CurrentState      VaultPKIState
		DesiredState      VaultPKIState
		ExpectedPatchFunc func(*framework.Patch)
	}{
		// Test 0
		{
			CurrentState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			DesiredState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetCreateChange([]ChangeType{
					BackendChange,
					CACertificateChange,
				})
			},
		},

		// Test 1
		{
			CurrentState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			DesiredState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: true,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetCreateChange([]ChangeType{
					CACertificateChange,
				})
			},
		},

		// Test 2
		{
			CurrentState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			DesiredState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: false,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetCreateChange([]ChangeType{
					BackendChange,
				})
			},
		},

		// Test 3
		{
			CurrentState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			DesiredState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {},
		},

		// Test 4
		{
			CurrentState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			DesiredState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {},
		},

		// Test 5
		{
			CurrentState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: false,
			},
			DesiredState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: false,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {},
		},

		// Test 6
		{
			CurrentState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: true,
			},
			DesiredState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: true,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {},
		},

		// Test 7
		{
			CurrentState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			DesiredState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: false,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetDeleteChange([]ChangeType{
					BackendChange,
					CACertificateChange,
				})
			},
		},

		// Test 8
		{
			CurrentState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			DesiredState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: false,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetDeleteChange([]ChangeType{
					CACertificateChange,
				})
			},
		},

		// Test 9
		{
			CurrentState: VaultPKIState{
				BackendExists:       true,
				CACertificateExists: true,
			},
			DesiredState: VaultPKIState{
				BackendExists:       false,
				CACertificateExists: true,
			},
			ExpectedPatchFunc: func(patch *framework.Patch) {
				patch.SetDeleteChange([]ChangeType{
					BackendChange,
				})
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
