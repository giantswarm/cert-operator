package pki

// CreateConfig is used to configure the setup of a PKI backend done by the
// Service.
type CreateConfig struct {
	// If set, clients can request certificates matching the value of the actual
	// domains themselves; e.g. if a configured domain set with allowed_domains
	// is example.com, this allows clients to actually request a certificate
	// containing the name example.com as one of the DNS values on the final
	// certificate. In some scenarios, this can be considered a security risk.
	// Defaults to false.
	AllowBareDomains bool `json:"allow_bare_domains"`

	// AllowedDomains represents a comma separate list of valid domain names the
	// generated certificate authority is valid for.
	AllowedDomains string `json:"allowed_domains"`

	// ClusterID represents the cluster ID a PKI backend setup should be done
	// for. This ID is used to restrict access on Vault related operations for a
	// specific cluster. E.g. the Vault PKI backend will be mounted on a path
	// containing this ID. That way Vault policies will restrict access to this
	// specific path.
	ClusterID string `json:"cluster_id"`

	// CommonName is the common name used to configure the root CA associated
	// with the current PKI backend.
	CommonName string `json:"common_name"`

	// TTL configures the time to live for the root CA being set up. This is a
	// golang time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

// Service manages the setup of Vault's PKI backends and all other required
// steps necessary to be done.
type Service interface {
	// PKI management.

	// Create sets up a Vault PKI backend according to the given configuration.
	Create(config CreateConfig) error

	// Delete removes the PKI backend associated wit the given cluster ID.
	Delete(clusterID string) error

	// IsCAGenerated checks whether the root CA associated with the given cluster
	// ID is generated.
	IsCAGenerated(clusterID string) (bool, error)

	// IsMounted checks whether the PKI backend associated with the given
	// cluster ID is mounted.
	IsMounted(clusterID string) (bool, error)

	// IsRoleCreated checks whether the PKI role associated with the given
	// cluster ID is created.
	IsRoleCreated(clusterID string) (bool, error)

	// VerifyPKISetup checks if IsMounted, IsCAGenerated and IsRoleCreated are all true
	// for the given cluster ID.
	VerifyPKISetup(clusterID string) (bool, error)

	// RoleName returns the name used to register the PKI backend's role.
	RoleName(clusterID string) string

	// Path management.

	// MountPKIPath returns the path under which a cluster's PKI backend is
	// mounted. This is very specific to Vault. The path structure is the
	// following.
	//
	//     pki-<clusterID>
	//
	MountPKIPath(clusterID string) string

	// WriteCAPath returns the path under which a cluster's certificate authority
	// can be generated. This is very specific to Vault. The path structure is
	// the following. See also
	// https://github.com/hashicorp/vault/blob/6f0f46deb622ba9c7b14b2ec0be24cab3916f3d8/website/source/docs/secrets/pki/index.html.md#pkirootgenerate.
	//
	//     pki-<clusterID>/root/generate/exported
	//
	WriteCAPath(clusterID string) string

	// WriteRolePath returns the path under which a role is registered. This is
	// very specific to Vault. The path structure is the following. See also
	// https://github.com/hashicorp/vault/blob/6f0f46deb622ba9c7b14b2ec0be24cab3916f3d8/website/source/docs/secrets/pki/index.html.md#pkiroles.
	//
	//     pki-<clusterID>/roles/role-<clusterID>
	//
	WriteRolePath(clusterID string) string
}
