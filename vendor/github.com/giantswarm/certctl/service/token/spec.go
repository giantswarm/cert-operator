package token

// CreateConfig is a data structure used to configure the token creation process
// implemented by Service.Create.
type CreateConfig struct {
	// ClusterID represents the cluster ID a token is requested for. This ID is
	// used to restrict access on Vault related operations for a specific cluster.
	// E.g. the generated token will only be allowed to issue certificates for the
	// Vault PKI backend associated with the given cluster ID.
	ClusterID string `json:"cluster_id"`

	// Num represents the number of tokens the generator should create.
	Num int `json:"num"`

	// TTL configures the time to live for the requested token. This is a golang
	// time string with the allowed units s, m and h.
	TTL string `json:"ttl"`
}

// Service creates new Vault policies to restrict access capabilities
// of e.g. Vault tokens.
type Service interface {
	// Create generates new Vault tokens allowed to be used to issue signed
	// certificates with respect to the given configuration.
	Create(config CreateConfig) ([]string, error)

	// CreatePolicy creates a new policy to restrict access to only being able to
	// issue signed certificates on the Vault PKI backend specific to the given
	// cluster ID. Here the given cluster ID is used to create the policy name and
	// the policy specific rules matching certain paths within the Vault file
	// system like path structure. This policy name can be used to e.g. apply it
	// to some Vault token.
	CreatePolicy(clusterID string) error

	// DeletePolicy removes a policy from Vault using its name.
	DeletePolicy(clusterID string) error

	// IsPolicyCreated checks whether the PKI issue policy already exists.
	IsPolicyCreated(clusterID string) (bool, error)

	// PolicyName returns the name of a policy used to restrict access to Vault
	// for PKI issue requests. This policy is scoped to the given cluster ID.
	PolicyName(clusterID string) string
}
