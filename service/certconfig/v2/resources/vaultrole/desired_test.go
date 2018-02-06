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

func Test_Resource_VaultRole_GetDesiredState(t *testing.T) {
	testCases := []struct {
		Obj          interface{}
		ExpectedRole *vaultrole.Role
	}{
		// Case 0 ensures the creation of the desired state of the Vault role.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						AllowBareDomains: true,
						AltNames: []string{
							"kubernetes",
							"kubernetes.default",
						},
						ClusterComponent: "api",
						ClusterID:        "al9qy",
						Organizations: []string{
							"system:masters",
						},
						TTL: "24h",
					},
				},
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

		// Case 1 is like 0 but with different inputs.
		{
			Obj: &v1alpha1.CertConfig{
				Spec: v1alpha1.CertConfigSpec{
					Cert: v1alpha1.CertConfigSpecCert{
						AllowBareDomains: false,
						AltNames: []string{
							"kubernetes",
							"kubernetes.default",
						},
						ClusterComponent: "calico",
						ClusterID:        "al9qy",
						Organizations:    nil,
						TTL:              "8h",
					},
				},
			},
			ExpectedRole: &vaultrole.Role{
				AllowBareDomains: false,
				AllowSubdomains:  true,
				AltNames: []string{
					"kubernetes",
					"kubernetes.default",
				},
				ID: "al9qy",
				Organizations: []string{
					"calico",
				},
				TTL: 8 * time.Hour,
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
		result, err := newResource.GetDesiredState(context.TODO(), tc.Obj)
		if err != nil {
			t.Fatal("case", i, "expected", nil, "got", err)
		}
		role := result.(*vaultrole.Role)
		if !reflect.DeepEqual(tc.ExpectedRole, role) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedRole, role)
		}
	}
}
