package k8s

type TLSClientConfig struct {
	// Files containing keys/certificates.
	CertFile string
	KeyFile  string
	CAFile   string
}

type Config struct {
	InCluster bool
	Address   string
	TLSClientConfig
}
