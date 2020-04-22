package service

import (
	"github.com/giantswarm/versionbundle"

	v2 "github.com/giantswarm/cert-operator/service/controller/v2"
)

func NewVersionBundles() []versionbundle.Bundle {
	var versionBundles []versionbundle.Bundle

	versionBundles = append(versionBundles, v2.VersionBundle())

	return versionBundles
}
