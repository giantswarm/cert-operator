package key

import (
	"reflect"
	"sort"
	"testing"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
)

func Test_Organization(t *testing.T) {
	testCases := []struct {
		Organizations         []string
		ClusterComponent      string
		ExpectedOrganizations []string
		ExpectedOutput        []string
	}{
		{
			Organizations:         []string{},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{},
			ExpectedOutput:        []string{"api"},
		},
		{
			Organizations:         []string{"system:master"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"system:master"},
			ExpectedOutput:        []string{"api", "system:master"},
		},
		{
			Organizations:         []string{"system:master", "giantswarm"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"system:master", "giantswarm"},
			ExpectedOutput:        []string{"api", "giantswarm", "system:master"},
		},
		{
			Organizations:         []string{"giantswarm", "system:master"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"giantswarm", "system:master"},
			ExpectedOutput:        []string{"api", "giantswarm", "system:master"},
		},
	}

	for i, tc := range testCases {
		customObject := v1alpha1.CertConfig{
			Spec: v1alpha1.CertConfigSpec{
				Cert: v1alpha1.CertConfigSpecCert{
					ClusterComponent: tc.ClusterComponent,
					Organizations:    tc.Organizations,
				},
			},
		}

		for j := 0; j < 10; j++ {
			if !reflect.DeepEqual(tc.ExpectedOrganizations, customObject.Spec.Cert.Organizations) {
				t.Fatalf("case %d iteration %d expected %#v got %#v", i, j, tc.ExpectedOrganizations, customObject.Spec.Cert.Organizations)
			}

			Organizations(customObject)

			result := Organizations(customObject)
			sort.Strings(result)

			if !reflect.DeepEqual(tc.ExpectedOutput, result) {
				t.Fatalf("case %d iteration %d expected %#v got %#v", i, j, tc.ExpectedOutput, result)
			}
		}
	}
}

func TestOrganizationCapacity(t *testing.T) {
	// create a slice of capacity greater than the number of elements
	// that the copy is going to have
	orgs := make([]string, 1, 4)
	orgs[0] = "myorg"

	customObject := v1alpha1.CertConfig{
		Spec: v1alpha1.CertConfigSpec{
			Cert: v1alpha1.CertConfigSpecCert{
				ClusterComponent: "api",
				Organizations:    orgs,
			},
		},
	}

	// here create an extended copy of orgs
	o := Organizations(customObject)

	// call sort on the copy, this will create havok in the original
	sort.Strings(o)

	expected := "myorg"
	actual := customObject.Spec.Cert.Organizations[0]
	if expected != actual {
		t.Errorf("customObject organizations changed by sorting an unrelated slice, expected %s, actual %s", expected, actual)
	}
}
