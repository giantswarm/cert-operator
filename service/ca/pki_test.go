package ca

import (
	"fmt"
	"strings"
	"testing"

	"github.com/giantswarm/certificatetpr"
	"github.com/stretchr/testify/assert"
)

func TestGetAllowedDomainsForCA(t *testing.T) {
	tests := []struct {
		name                string
		caCommonName        string
		cert                certificatetpr.Spec
		expectedDomainCount int
	}{
		{
			name:         "Specify api cert with alt names",
			caCommonName: "cluster-test.g8s.test.giantswarm.io",
			cert: certificatetpr.Spec{
				AltNames: []string{
					"kubernetes",
					"kubernetes.default",
					"kubernetes.default.svc",
					"kubernetes.default.svc.cluster.local",
				},
			},
			expectedDomainCount: 5,
		},
		{
			name:                "Specify etcd cert without alt names",
			caCommonName:        "cluster-test.g8s.test.giantswarm.io",
			cert:                certificatetpr.Spec{},
			expectedDomainCount: 1,
		},
	}

	for _, tc := range tests {
		allowedDomains := getAllowedDomainsForCA(tc.caCommonName, tc.cert)
		domainCount := len(strings.Split(allowedDomains, ","))

		assert.Equal(t, tc.expectedDomainCount, domainCount, fmt.Sprintf("[%s] Allowed domains should be equal", tc.name))
	}
}
