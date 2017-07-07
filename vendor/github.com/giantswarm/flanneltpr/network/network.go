package network

import (
	"github.com/giantswarm/flanneltpr/network/bridge"
)

// Network holds information used for the execution of the Docker entrypoint of
// https://github.com/giantswarm/k8s-network-bridge.
type Network struct {
	Bridge bridge.Bridge `json:"bridge" yaml:"bridge"`
	// BridgeName contains the value for the environment variable
	// ${NETWORK_BRIDGE_NAME} in the Docker entrypoint of
	// https://github.com/giantswarm/k8s-network-bridge.
	BridgeName string `json:"bridgeName" yaml:"bridgeName"`
	// EnvFilePath contains the value for the environment variable
	// ${NETWORK_ENV_FILE_PATH} in the Docker entrypoint of
	// https://github.com/giantswarm/k8s-network-bridge.
	EnvFilePath string `json:"envFilePath" yaml:"envFilePath"`
	// InterfaceName contains the value for the environment variable
	// ${NETWORK_INTERFACE_NAME} in the Docker entrypoint of
	// https://github.com/giantswarm/k8s-network-bridge.
	InterfaceName string `json:"interfaceName" yaml:"interfaceName"`
}
