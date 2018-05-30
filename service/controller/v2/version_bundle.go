package v2

import (
	"github.com/giantswarm/versionbundle"
)

func VersionBundle() versionbundle.Bundle {
	return versionbundle.Bundle{
		Changelogs: []versionbundle.Changelog{
			{
				Component:   "vault",
				Description: "Vault version updated.",
				Kind:        versionbundle.KindChanged,
			},
		},
		Components: []versionbundle.Component{
			{
				Name:    "vault",
				Version: "0.7.3",
			},
		},
		Name:    "cert-operator",
		Version: "0.1.0",
	}
}
