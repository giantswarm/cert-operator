package bridge

import (
	"github.com/giantswarm/flanneltpr/network/bridge/docker"
)

type Bridge struct {
	Docker docker.Docker `json:"docker" yaml:"docker"`
}
