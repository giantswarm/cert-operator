package key

import (
	"reflect"
	"sort"
	"testing"

	"github.com/giantswarm/certificatetpr"
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
		customObject := certificatetpr.CustomObject{
			Spec: certificatetpr.Spec{
				ClusterComponent: tc.ClusterComponent,
				Organizations:    tc.Organizations,
			},
		}

		for j := 0; j < 10; j++ {
			if !reflect.DeepEqual(tc.ExpectedOrganizations, customObject.Spec.Organizations) {
				t.Fatalf("case %d iteration %d expected %#v got %#v", i, j, tc.ExpectedOrganizations, customObject.Spec.Organizations)
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
