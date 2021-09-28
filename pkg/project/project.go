package project

var (
	description string = "The cert-operator handles certificates for Kubernetes clusters running on Giantnetes."
	gitSHA             = "n/a"
	name        string = "cert-operator"
	source      string = "https://github.com/giantswarm/cert-operator"
	version            = "1.1.0"
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
