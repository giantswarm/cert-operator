package key

import (
	"reflect"
	"sort"
	"testing"

	"github.com/giantswarm/certificatetpr"
)

func Test_Organization_sort(t *testing.T) {
	testCases := []struct {
		Organizations         []string
		ClusterComponent      string
		ExpectedOrganizations []string
	}{
		{
			Organizations:         []string{},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"api"},
		},
		{
			Organizations:         []string{"system:master"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"api", "system:master"},
		},
		{
			Organizations:         []string{"system:master", "giantswarm"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"api", "system:master", "giantswarm"},
		},
		{
			Organizations:         []string{"giantswarm", "system:master"},
			ClusterComponent:      "api",
			ExpectedOrganizations: []string{"api", "giantswarm", "system:master"},
		},
	}

	for i, tc := range testCases {
		customObject := certificatetpr.CustomObject{
			Spec: certificatetpr.Spec{
				ClusterComponent: tc.ClusterComponent,
				Organizations:    tc.Organizations,
			},
		}

		validate := func(result []string) {
			if !reflect.DeepEqual(tc.ExpectedOrganizations, result) {
				t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedOrganizations, result)
			}
		}
		result := Organizations(customObject)
		validate(result)

		result = Organizations(customObject)
		validate(result)
		a := Organizations(customObject)
		sort.Strings(a)
		result = Organizations(customObject)
		validate(result)
	}
}
