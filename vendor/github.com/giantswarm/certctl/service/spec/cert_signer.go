package spec

// IssueConfig is used to configure the process of issuing a certificate key
// pair using the CertSigner.
type IssueConfig struct {
	// ClusterID represents the cluster ID a new certificate should be issued
	// for.
	ClusterID string `json:"cluster_id"`

	// CommonName is the common name used to configure the issued certificate
	// that is being requested.
	CommonName string `json:"common_name"`

	// IPSANs represents a comma separate lists of IPs.
	IPSANs string `json:"ip_sans"`

	// AltNames names represents a comma separate list of alternative names.
	AltNames string `json:"alt_names"`

	// TTL configures the time to live for the requested certificate. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

type IssueResponse struct {
	Certificate  string `json:"certificate"`
	PrivateKey   string `json:"private_key"`
	IssuingCA    string `json:"issuing_ca"`
	SerialNumber string `json:"serial_number"`
}

// CertSigner manages the process of issuing new certificate key pairs
type CertSigner interface {
	// Issue generates a new signed certificate with respect to the given
	// configuration.
	Issue(config IssueConfig) (IssueResponse, error)

	// SignedPath returns the path under which a certificate can be generated.
	// This is very specific to Vault. The path structure is the following. See
	// also https://github.com/hashicorp/vault/blob/6f0f46deb622ba9c7b14b2ec0be24cab3916f3d8/website/source/docs/secrets/pki/index.html.md#pkiissue.
	//
	//     pki-<clusterID>/issue/role-<clusterID>
	//
	SignedPath(clusterID string) string
}
