package vaultcrt

import (
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/vaultcrt/key"
)

func (c *VaultCrt) Create(config CreateConfig) (CreateResult, error) {
	k := key.IssuePath(config.ID, config.Organizations)
	v := map[string]interface{}{
		"alt_names":   strings.Join(config.AltNames, ","),
		"common_name": config.CommonName,
		"ip_sans":     strings.Join(config.IPSANs, ","),
		"ttl":         config.TTL,
	}

	secret, err := c.vaultClient.Logical().Write(k, v)
	if err != nil {
		return CreateResult{}, microerror.Mask(err)
	}

	var CA string
	{
		v, ok := secret.Data["issuing_ca"]
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "issuing CA missing")
		}
		CA, ok = v.(string)
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "issuing CA must be string")
		}
	}

	var crt string
	{
		v, ok := secret.Data["certificate"]
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "certificate missing")
		}
		crt, ok = v.(string)
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "certificate must be string")
		}
	}

	var key string
	{
		v, ok := secret.Data["private_key"]
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "private key missing")
		}
		key, ok = v.(string)
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "private key must be string")
		}
	}

	var serialNumber string
	{
		v, ok := secret.Data["serial_number"]
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "serial number missing")
		}
		serialNumber, ok = v.(string)
		if !ok {
			return CreateResult{}, microerror.Maskf(executionFailedError, "serial number must be string")
		}
	}

	createResult := CreateResult{
		CA:           CA,
		Crt:          crt,
		Key:          key,
		SerialNumber: serialNumber,
	}

	return createResult, nil
}
