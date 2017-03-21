package create

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClusterCA(t *testing.T) {
	svc, _ := New(DefaultConfig())

	tests := []struct {
		name           string
		certCommonName string
		expectedCA     string
	}{
		{
			name:           "Specify valid api cert",
			certCommonName: "api.test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
			expectedCA:     "test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
		},
		{
			name:           "Specify valid etcd cert",
			certCommonName: "etcd.test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
			expectedCA:     "test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
		},
		{
			name:           "Specify different url structure",
			certCommonName: "api.test-cluster.private.giantswarm.io",
			expectedCA:     "test-cluster.private.giantswarm.io",
		},
	}

	for _, tc := range tests {
		cert := CertificateSpec{
			CommonName: tc.certCommonName,
		}

		ca := svc.getClusterCA(cert)
		assert.Equal(t, tc.expectedCA, ca, fmt.Sprintf("[%s] CA common name should be equal", tc.name))
	}
}

func TestGetAllowedDomainsForCA(t *testing.T) {
	svc, _ := New(DefaultConfig())

	tests := []struct {
		name                string
		certCommonName      string
		altNames            []string
		expectedDomainCount int
	}{
		{
			name:           "Specify valid api cert",
			certCommonName: "api.test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
			altNames: []string{
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				"kubernetes.default.svc.cluster.local",
			},
			expectedDomainCount: 5,
		},
		{
			name:                "Specify etcd cert without allowed names",
			certCommonName:      "etcd.test-cluster.g8s.eu-west-1.aws.test.private.giantswarm.io",
			altNames:            []string{},
			expectedDomainCount: 1,
		},
	}

	for _, tc := range tests {
		cert := CertificateSpec{
			CommonName: tc.certCommonName,
			AltNames:   tc.altNames,
		}

		allowedDomains := svc.getAllowedDomainsForCA(cert)
		domainCount := len(strings.Split(allowedDomains, ","))

		assert.Equal(t, tc.expectedDomainCount, domainCount, fmt.Sprintf("[%s] Allowed domains should be equal", tc.name))
	}
}
