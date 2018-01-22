package vaultrole

import "time"

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

type SearchConfig struct {
	ID            string
	Organizations []string
}

type UpdateConfig struct {
	AllowBareDomains bool
	AllowSubdomains  bool
	AltNames         []string
	ID               string
	Organizations    []string
	TTL              string
}

type Interface interface {
	Create(config CreateConfig) error
	Exists(config ExistsConfig) (bool, error)
	Search(config SearchConfig) (Role, error)
	Update(config UpdateConfig) error
}

type Role struct {
	AllowBareDomains bool
	AllowSubdomains  bool
	AltNames         []string
	ID               string
	Organizations    []string
	TTL              time.Duration
}
