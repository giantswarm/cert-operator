package service

import (
	"time"

	"github.com/giantswarm/versionbundle"
)

func NewVersionBundles() []versionbundle.Bundle {
	return []versionbundle.Bundle{
		{
			Changelogs: []versionbundle.Changelog{
				{
					Component:   "vault",
					Description: "Vault version updated.",
					Kind:        "changed",
				},
			},
			Components: []versionbundle.Component{
				{
					Name:    "vault",
					Version: "0.7.3",
				},
			},
			Dependencies: []versionbundle.Dependency{},
			Deprecated:   false,
			Name:         "cert-operator",
			Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
			Version:      "0.1.0",
			WIP:          true,
		},
	}
}
