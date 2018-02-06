package vaultpki

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultpki/vaultpkitest"
	vaultapi "github.com/hashicorp/vault/api"
)

func Test_Resource_VaultPKI_NewCreateChange(t *testing.T) {
	testCases := []struct {
		Obj           interface{}
		CurrentState  interface{}
		DesiredState  interface{}
		ExpectedState VaultPKIState
	}{
		// Test 0 ensures that zero value input results in zero value output.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState:  VaultPKIState{},
			DesiredState:  VaultPKIState{},
			ExpectedState: VaultPKIState{},
		},

		// Test 1 ensures that the current state is reversed using the desired
		// state. In case the backend state is nil and the CA certificate state is
		// not empty within the current state, the create state should contain the
		// backend state from the desired state and the CA certificate state should
		// be empty.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState: VaultPKIState{
				Backend:       nil,
				CACertificate: "placeholder",
			},
			DesiredState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "placeholder",
			},
			ExpectedState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
				CACertificate: "",
			},
		},

		// Test 2 ensures that the current state is reversed using the desired
		// state. In case the backend state is not nil and the CA certificate state
		// is empty within the current state, the create state should contain a nil
		// backend state and the CA certificate state should be defined by the
		// desired state.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState: VaultPKIState{
				Backend: &vaultapi.MountOutput{
					Type: "pki",
				},
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
				CACertificate: "placeholder",
			},
		},

		// Test 3 ensures that a complete current state results in a completely
		// empty create state.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
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
		result, err := newResource.newCreateChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		r := result.(VaultPKIState)
		if !reflect.DeepEqual(r, tc.ExpectedState) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedState, r)
		}
	}
}
