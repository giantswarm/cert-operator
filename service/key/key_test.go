package key

import (
	"reflect"
	"testing"

	"github.com/giantswarm/certificatetpr"
)

func Test_Organization_sort(t *testing.T) {
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
			ExpectedOutput:        []string{"api", "system:master", "giantswarm"},
		},
		{
			Organizations:         []string{"giantswarm", "system:master"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"giantswarm", "system:master"},
			ExpectedOutput:        []string{"api", "giantswarm", "system:master"},
		},
	}

	for i, tc := range testCases {
		customObject := certificatetpr.CustomObject{
			Spec: certificatetpr.Spec{
				ClusterComponent: tc.ClusterComponent,
				Organizations:    tc.Organizations,
			},
		}

		result := Organizations(customObject)
		if !reflect.DeepEqual(tc.ExpectedOutput, result) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedOutput, result)
		}

		if !reflect.DeepEqual(tc.ExpectedOrganizations, customObject.Spec.Organizations) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedOrganizations, customObject.Spec.Organizations)
		}
	}
}
