package role

// CreateParams represent the parameters for creating a role.
type CreateParams struct {
	AllowBareDomains bool   `json:"allow_bare_domains"`
	AllowSubdomains  bool   `json:"allow_sub_domains"`
	AllowedDomains   string `json:"allowed_domains"`
	Name             string `json:"name"`
	Organizations    string `json:"organizations"`
	TTL              string `json:"ttl"`
}

// Service manages the setup of Vault's PKI backends and all other required
// steps necessary to be done.
type Service interface {
	// Role management.

	// Create creates a role.
	Create(params CreateParams) error

	// IsRoleCreated checks whether a given role exists.
	IsRoleCreated(roleName string) (bool, error)
}
