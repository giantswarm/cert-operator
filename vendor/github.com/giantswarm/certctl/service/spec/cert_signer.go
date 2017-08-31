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

	// Organizations is a comma seperated list of organizations ("O"'s) for the issued cert's
	// subject line.
	Organizations string `json:"organizations"`

	// IPSANs represents a comma separate lists of IPs.
	IPSANs string `json:"ip_sans"`

	// AltNames names represents a comma separate list of alternative names.
	AltNames string `json:"alt_names"`

	// TTL configures the time to live for the requested certificate. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`

	//// QUESTIONABLE ATTRIBUTES
	///

	// It seem weird to have these attributes here (AllowedDomains, AllowBareDomains, and RoleTTL)
	// and in the PKI setup call, but we need to know them again because issuing a certificate
	// might also create a role in vault on the fly, and these attributes are part of a role
	// definition.

	AllowedDomains   string `json:"allowed_domains"`
	AllowBareDomains bool   `json:"allow_bare_domains"`
	RoleTTL          string `json:"role_ttl"`

	///
	//// END QUESTIONABLE
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
	//		 When organizations is blank:
	//     pki-<clusterID>/issue/role-<clusterID>
	//
	//     When organizations is not blank
	//     pki-<clusterID>/issue/role-org-<organizationsHash>
	//
	//     organizationsHash is a deterministic urlsafe hash that is always the
	//     same regardless of what order you give the organizations in.
	//
	SignedPath(clusterID string, organizations string) string
}
