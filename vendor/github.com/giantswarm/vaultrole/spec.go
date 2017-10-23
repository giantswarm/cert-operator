package vaultrole

type CreateConfig struct {
	AllowBareDomains bool
	AllowSubdomains  bool
	AltNames         []string
	ID               string
	Organizations    []string
	TTL              string
}

type ExistsConfig struct {
	ID            string
	Organizations []string
}

type Interface interface {
	Create(config CreateConfig) error
	Exists(config ExistsConfig) (bool, error)
}
