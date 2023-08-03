package project

var (
	description string = "The cert-operator handles certificates for Kubernetes clusters running on Giantnetes."
	gitSHA             = "n/a"
	name        string = "cert-operator"
	source      string = "https://github.com/giantswarm/cert-operator"
	version            = "3.2.1"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}

// ManagementClusterAppVersion is always 0.0.0 for management cluster app CRs. These CRs
// are processed by app-operator-unique which always runs the latest version.
func ManagementClusterAppVersion() string {
	return "0.0.0"
}
