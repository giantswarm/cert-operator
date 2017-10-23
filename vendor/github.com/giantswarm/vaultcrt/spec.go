package vaultcrt

type CreateConfig struct {
	AltNames      []string
	CommonName    string
	ID            string
	IPSANs        []string
	Organizations []string
	TTL           string
}

type CreateResult struct {
	CA           string
	Crt          string
	Key          string
	SerialNumber string
}

type Interface interface {
	Create(config CreateConfig) (CreateResult, error)
}
