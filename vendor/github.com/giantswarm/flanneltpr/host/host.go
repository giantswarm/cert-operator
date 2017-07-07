package host

// Host holds information used for the execution of the Docker entrypoint of
// https://github.com/giantswarm/k8s-network-bridge.
type Host struct {
	// PrivateNetwork contains the value for the environment variable
	// ${HOST_PRIVATE_NETWORK} in the Docker entrypoint of
	// https://github.com/giantswarm/k8s-network-bridge.
	PrivateNetwork string `json:"privateNetwork" yaml:"privateNetwork"`
}
