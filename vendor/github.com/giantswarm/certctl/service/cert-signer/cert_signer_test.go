package certsigner

import (
	"testing"
)

func Test_roleName(t *testing.T) {
	testCases := []struct {
		ClusterID      string
		Organizations  string
		ExpectedResult string
	}{
		// Case 1: Without orgs, we should just get a role identified by the cluster id.
		{
			ClusterID:      "123",
			Organizations:  "",
			ExpectedResult: "role-123",
		},
		// Case 2: With orgs, we should get a role name that has a org hash in it.
		{
			ClusterID:      "123",
			Organizations:  "blue,green",
			ExpectedResult: "role-org-ae04e382ff1b455a454bfde83bdda9dc8d077649",
		},
		// Case 3: The order of the orgs should not impact the hash.
		{
			ClusterID:      "123",
			Organizations:  "green,blue",
			ExpectedResult: "role-org-ae04e382ff1b455a454bfde83bdda9dc8d077649",
		},
		// Case 4: A different orgs list should yield a different hash.
		{
			ClusterID:      "123",
			Organizations:  "green,blue,red",
			ExpectedResult: "role-org-40c7be91742c1d2343d32ea489e169b1121bc674",
		},
	}

	for i, testCase := range testCases {
		result := roleName(testCase.ClusterID, testCase.Organizations)

		if result != testCase.ExpectedResult {
			t.Fatalf("case %d expected %#v got %#v", i+1, testCase.ExpectedResult, result)
		}
	}
}
