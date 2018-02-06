package vaultrole

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/vaultrole"
	"github.com/giantswarm/vaultrole/vaultroletest"
)

func Test_Resource_VaultRole_newCreateChange(t *testing.T) {
	testCases := []struct {
		Obj          interface{}
		CurrentState interface{}
		DesiredState interface{}
		ExpectedRole *vaultrole.Role
	}{
		// Case 0 ensures zero value input results in zero value output.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState: nil,
			DesiredState: nil,
			ExpectedRole: nil,
		},

		// Case 1 ensures a given current state results in a nil output when there
		// is a nil desired state.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState: &vaultrole.Role{
				AllowBareDomains: true,
				AllowSubdomains:  true,
				AltNames: []string{
					"kubernetes",
					"kubernetes.default",
				},
				ID: "al9qy",
				Organizations: []string{
					"api",
					"system:masters",
				},
				TTL: 24 * time.Hour,
			},
			DesiredState: &vaultrole.Role{},
			ExpectedRole: nil,
		},

		// Case 2 that the desired state defines the output in case the current
		// state is nil.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						ClusterID: "foobar",
					},
				},
			},
			CurrentState: nil,
			DesiredState: &vaultrole.Role{
				AllowBareDomains: true,
				AllowSubdomains:  true,
				AltNames: []string{
					"kubernetes",
					"kubernetes.default",
				},
				ID: "al9qy",
				Organizations: []string{
					"api",
					"system:masters",
				},
				TTL: 24 * time.Hour,
			},
			ExpectedRole: &vaultrole.Role{
				AllowBareDomains: true,
				AllowSubdomains:  true,
				AltNames: []string{
					"kubernetes",
					"kubernetes.default",
				},
				ID: "al9qy",
				Organizations: []string{
					"api",
					"system:masters",
				},
				TTL: 24 * time.Hour,
			},
		},
	}

	var err error
	var newResource *Resource
	{
		c := DefaultConfig()

		c.Logger = microloggertest.New()
		c.VaultRole = vaultroletest.New()

		newResource, err = New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	for i, tc := range testCases {
		result, err := newResource.newCreateChange(context.TODO(), tc.Obj, tc.CurrentState, tc.DesiredState)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		role := result.(*vaultrole.Role)
		if !reflect.DeepEqual(tc.ExpectedRole, role) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedRole, role)
		}
	}
}
