package create

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllowedDomainsForCA(t *testing.T) {
	tests := []struct {
		name                string
		caCommonName        string
		cert                CertificateSpec
		expectedDomainCount int
	}{
		{
			name:         "Specify api cert with alt names",
			caCommonName: "cluster-test.g8s.test.giantswarm.io",
			cert: CertificateSpec{
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
			cert:                CertificateSpec{},
			expectedDomainCount: 1,
		},
	}

	for _, tc := range tests {
		allowedDomains := getAllowedDomainsForCA(tc.caCommonName, tc.cert)
		domainCount := len(strings.Split(allowedDomains, ","))

		assert.Equal(t, tc.expectedDomainCount, domainCount, fmt.Sprintf("[%s] Allowed domains should be equal", tc.name))
	}
}
