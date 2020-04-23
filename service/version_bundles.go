package service

import (
	"github.com/giantswarm/versionbundle"

	"github.com/giantswarm/cert-operator/service/controller"
)

func NewVersionBundles() []versionbundle.Bundle {
	var versionBundles []versionbundle.Bundle

	versionBundles = append(versionBundles, controller.VersionBundle())

	return versionBundles
}
