package certificatetpr

type Spec struct {
	AllowBareDomains bool     `json:"allowBareDomains"`
	AltNames         []string `json:"altNames"`
	ClusterComponent string   `json:"clusterComponent"`
	ClusterID        string   `json:"clusterID"`
	CommonName       string   `json:"commonName"`
	IPSANs           []string `json:"ipSans"`
	TTL              string   `json:"ttl"`
}
