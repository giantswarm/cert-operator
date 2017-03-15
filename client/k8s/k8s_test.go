package k8s

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestGetRawClientConfig(t *testing.T) {
	tests := []struct {
		name            string
		config          Config
		expectedError   bool
		expectedHost    string
		expectedCrtFile string
		expectedKeyFile string
		expectedCAFile  string
	}{
		{
			name: "Specify only in-cluster config. It should return it. Use cert auth files.",
			config: Config{
				InCluster: true,
				inClusterConfigProvider: func() (*rest.Config, error) {
					return &rest.Config{
						Host: "http://in-cluster-host",
						TLSClientConfig: rest.TLSClientConfig{
							CertFile: "/var/run/kubernetes/client-admin.crt",
							KeyFile:  "/var/run/kubernetes/client-admin.key",
							CAFile:   "/var/run/kubernetes/server-ca.crt",
						},
					}, nil
				},
			},
			expectedHost:    "http://in-cluster-host",
			expectedCrtFile: "/var/run/kubernetes/client-admin.crt",
			expectedKeyFile: "/var/run/kubernetes/client-admin.key",
			expectedCAFile:  "/var/run/kubernetes/server-ca.crt",
		},
		{
			name: "Do not specify anything while using in-cluster config. It should return an error.",
			config: Config{
				InCluster: true,
				inClusterConfigProvider: func() (*rest.Config, error) {
					return nil, fmt.Errorf("No in-cluster config")
				},
			},
			expectedError: true,
		},
		{
			name: "Specify both in-cluster config and CLI config. It should return the in-cluster config. Use cert auth files.",
			config: Config{
				InCluster: true,
				Address:   "http://host-from-cli",
				TLSClientConfig: TLSClientConfig{
					CertFile: "/var/run/kubernetes/client-cli.crt",
					KeyFile:  "/var/run/kubernetes/client-cli.key",
					CAFile:   "/var/run/kubernetes/server-ca-cli.crt",
				},
				inClusterConfigProvider: func() (*rest.Config, error) {
					return &rest.Config{
						Host: "http://in-cluster-host",
						TLSClientConfig: rest.TLSClientConfig{
							CertFile: "/var/run/kubernetes/client-admin.crt",
							KeyFile:  "/var/run/kubernetes/client-admin.key",
							CAFile:   "/var/run/kubernetes/server-ca.crt",
						},
					}, nil
				},
			},
			expectedHost:    "http://in-cluster-host",
			expectedCrtFile: "/var/run/kubernetes/client-admin.crt",
			expectedKeyFile: "/var/run/kubernetes/client-admin.key",
			expectedCAFile:  "/var/run/kubernetes/server-ca.crt",
		},
		{
			name: "Specify both in-cluster config and CLI config. It should return the CLI config. Use cert auth files.",
			config: Config{
				Address: "http://host-from-cli",
				TLSClientConfig: TLSClientConfig{
					CertFile: "/var/run/kubernetes/client-cli.crt",
					KeyFile:  "/var/run/kubernetes/client-cli.key",
					CAFile:   "/var/run/kubernetes/server-ca-cli.crt",
				},
				inClusterConfigProvider: func() (*rest.Config, error) {
					return &rest.Config{
						Host: "http://in-cluster-host",
						TLSClientConfig: rest.TLSClientConfig{
							CertFile: "/var/run/kubernetes/client-admin.crt",
							KeyFile:  "/var/run/kubernetes/client-admin.key",
							CAFile:   "/var/run/kubernetes/server-ca.crt",
						},
					}, nil
				},
			},
			expectedHost:    "http://host-from-cli",
			expectedCrtFile: "/var/run/kubernetes/client-cli.crt",
			expectedKeyFile: "/var/run/kubernetes/client-cli.key",
			expectedCAFile:  "/var/run/kubernetes/server-ca-cli.crt",
		},
		{
			name: "Specify only CLI config. It should return it. Use cert auth files.",
			config: Config{
				Address: "http://host-from-cli",
				TLSClientConfig: TLSClientConfig{
					CertFile: "/var/run/kubernetes/client-cli.crt",
					KeyFile:  "/var/run/kubernetes/client-cli.key",
					CAFile:   "/var/run/kubernetes/server-ca-cli.key",
				},
				inClusterConfigProvider: func() (*rest.Config, error) {
					return nil, fmt.Errorf("No in-cluster config")
				},
			},
			expectedHost:    "http://host-from-cli",
			expectedCrtFile: "/var/run/kubernetes/client-cli.crt",
			expectedKeyFile: "/var/run/kubernetes/client-cli.key",
			expectedCAFile:  "/var/run/kubernetes/server-ca-cli.key",
		},
	}
	for _, tc := range tests {
		rawClientConfig, err := getRawClientConfig(tc.config)
		if tc.expectedError {
			assert.Error(t, err, fmt.Sprintf("[%s] An error was expected", tc.name))
			continue
		}
		assert.Nil(t, err, fmt.Sprintf("[%s] An error was unexpected", tc.name))
		assert.Equal(t, tc.expectedHost, rawClientConfig.Host, fmt.Sprintf("[%s] Hosts should be equal", tc.name))
		assert.Equal(t, tc.expectedCrtFile, rawClientConfig.TLSClientConfig.CertFile, fmt.Sprintf("[%s] CertFiles should be equal", tc.name))
		assert.Equal(t, tc.expectedKeyFile, rawClientConfig.TLSClientConfig.KeyFile, fmt.Sprintf("[%s] KeyFiles should be equal", tc.name))
		assert.Equal(t, tc.expectedCAFile, rawClientConfig.TLSClientConfig.CAFile, fmt.Sprintf("[%s] CAFiles should be equal", tc.name))
	}
}
