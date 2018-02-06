package key

import (
	"reflect"
	"testing"
)

func Test_AllowedDomains(t *testing.T) {
	testCases := []struct {
		ID               string
		CommonNameFormat string
		AltNames         []string
		ExpectedResult   string
	}{
		{
			ID:               "al9qy",
			CommonNameFormat: "%s.g8s.gigantic.io",
			AltNames: []string{
				"kubernetes",
				"kubernetes.default.svc.cluster.local",
			},
			ExpectedResult: "al9qy.g8s.gigantic.io,kubernetes,kubernetes.default.svc.cluster.local",
		},

		{
			ID:               "al9qy",
			CommonNameFormat: "%s.g8s.gigantic.io",
			AltNames:         []string{},
			ExpectedResult:   "al9qy.g8s.gigantic.io",
		},

		{
			ID:               "al9qy",
			CommonNameFormat: "%s.g8s.gigantic.io",
			AltNames:         nil,
			ExpectedResult:   "al9qy.g8s.gigantic.io",
		},
	}

	for i, tc := range testCases {
		result := AllowedDomains(tc.ID, tc.CommonNameFormat, tc.AltNames)

		if result != tc.ExpectedResult {
			t.Fatalf("case %d expected %#v got %#v", i+1, tc.ExpectedResult, result)
		}
	}
}

func Test_RoleName(t *testing.T) {
	testCases := []struct {
		ID             string
		Organizations  []string
		ExpectedResult string
	}{
		// Case 0: Without orgs, we should just get a role identified by the cluster id.
		{
			ID:             "123",
			Organizations:  nil,
			ExpectedResult: "role-123",
		},

		// Case 1: same as 0. but with initialized empty slice instead of nil.
		{
			ID:             "123",
			Organizations:  []string{},
			ExpectedResult: "role-123",
		},

		// Case 2: With orgs, we should get a role name that has a org hash in it.
		{
			ID: "123",
			Organizations: []string{
				"blue",
				"green",
			},
			ExpectedResult: "role-org-ae04e382ff1b455a454bfde83bdda9dc8d077649",
		},

		// Case 3: The order of the orgs should not impact the hash.
		{
			ID: "123",
			Organizations: []string{
				"green",
				"blue",
			},
			ExpectedResult: "role-org-ae04e382ff1b455a454bfde83bdda9dc8d077649",
		},

		// Case 4: A different orgs list should yield a different hash.
		{
			ID: "123",
			Organizations: []string{
				"green",
				"blue",
				"red",
			},
			ExpectedResult: "role-org-40c7be91742c1d2343d32ea489e169b1121bc674",
		},

		// Case 5: A common case we see n production. The created hash here is of
		// interest to identify roles used for API certs.
		{
			ID: "al9qy",
			Organizations: []string{
				"api",
				"system:masters",
			},
			ExpectedResult: "role-org-f1b776344f5a19dbc38243e915767ff2ef9234df",
		},
	}

	for i, tc := range testCases {
		result := RoleName(tc.ID, tc.Organizations)

		if result != tc.ExpectedResult {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedResult, result)
		}
	}
}

func Test_ToAltNames(t *testing.T) {
	testCases := []struct {
		AllowedDomains   string
		ExpectedAltNames []string
	}{
		{
			AllowedDomains:   "",
			ExpectedAltNames: nil,
		},

		{
			AllowedDomains: "al9qy.g8s.gigantic.io,kubernetes,kubernetes.default.svc.cluster.local",
			ExpectedAltNames: []string{
				"kubernetes",
				"kubernetes.default.svc.cluster.local",
			},
		},

		{
			AllowedDomains: "kubernetes,kubernetes.default.svc.cluster.local",
			ExpectedAltNames: []string{
				"kubernetes.default.svc.cluster.local",
			},
		},
	}

	for i, tc := range testCases {
		result := ToAltNames(tc.AllowedDomains)

		if !reflect.DeepEqual(result, tc.ExpectedAltNames) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedAltNames, result)
		}
	}
}

func Test_ToOrganizations(t *testing.T) {
	testCases := []struct {
		Organizations         string
		ExpectedOrganizations []string
	}{
		{
			Organizations:         "",
			ExpectedOrganizations: nil,
		},

		{
			Organizations: "api,system:masters",
			ExpectedOrganizations: []string{
				"api",
				"system:masters",
			},
		},

		{
			Organizations: "api",
			ExpectedOrganizations: []string{
				"api",
			},
		},
	}

	for i, tc := range testCases {
		result := ToOrganizations(tc.Organizations)

		if !reflect.DeepEqual(result, tc.ExpectedOrganizations) {
			t.Fatalf("case %d expected %#v got %#v", i, tc.ExpectedOrganizations, result)
		}
	}
}
