package token

import (
	"bytes"
	"text/template"
)

// pkiIssuePolicyContext is the template context provided to the rendering of
// the pkiIssuePolicyTemplate.
type pkiIssuePolicyContext struct {
	ClusterID string
}

// pkiIssuePolicyTemplate provides a template of Vault policies used to
// restrict access to only being able to issue signed certificates specific to
// a Vault PKI backend of a cluster ID.
var pkiIssuePolicyTemplate = `
	path "pki-{{.ClusterID}}/issue/role-{{.ClusterID}}" {
		policy = "write"
	}
`

func execTemplate(t string, v interface{}) (string, error) {
	var result bytes.Buffer

	tmpl, err := template.New("policy-template").Parse(t)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&result, v)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
