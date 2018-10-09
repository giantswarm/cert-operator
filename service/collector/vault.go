package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	vault "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// ExpireTimeKey is the data key provided by the secret when looking up the
	// used Vault token. This key is specific to Vault as they define it.
	ExpireTimeKey = "expire_time"
	// ExpireTimeLayout is the layout used for time parsing when inspecting the
	// expiration date of the used Vault token. This layout is specific to Vault
	// as they define it.
	ExpireTimeLayout = "2006-01-02T15:04:05"
)

var (
	tokenExpireTimeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("cert_operator", "vault", "token_expire_time"),
		"A metric of the expire time of Vault tokens as unix seconds.",
		nil,
		nil,
	)
)

func (c *Collector) collectVaultMetrics(ch chan<- prometheus.Metric) {
	c.logger.Log("level", "debug", "message", "start collecting metrics")

	secret, err := c.vaultClient.Auth().Token().LookupSelf()
	if err != nil {
		c.logger.Log("level", "error", "message", "Vault token lookup failed", "stack", fmt.Sprintf("%#v", err))
		return
	}

	expiration, err := expirationFromSecret(secret)
	if err != nil {
		c.logger.Log("level", "error", "message", "parsing token expiration failed", "stack", fmt.Sprintf("%#v", err))
		return
	}

	ch <- prometheus.MustNewConstMetric(
		tokenExpireTimeDesc,
		prometheus.GaugeValue,
		float64(expiration.Unix()),
	)

	c.logger.Log("level", "debug", "message", "finished collecting metrics")
}

func expirationFromSecret(secret *vault.Secret) (time.Time, error) {
	value, ok := secret.Data[ExpireTimeKey]
	if !ok {
		return time.Time{}, microerror.Maskf(executionFailedError, "value of %q must exist in order to collect metrics for the Vault token expiration", ExpireTimeKey)
	}

	if value == nil {
		return time.Time{}, microerror.Maskf(executionFailedError, "Vault token does not expire, skipping metric update")
	}

	e, ok := value.(string)
	if !ok {
		return time.Time{}, microerror.Maskf(executionFailedError, "'%#v' must be string in order to collect metrics for the Vault token expiration", value)
	}

	split := strings.Split(e, ".")
	if len(split) == 0 {
		return time.Time{}, microerror.Maskf(executionFailedError, "'%#v' must have at least one item in order to collect metrics for the Vault token expiration", e)
	}

	expiration, err := time.Parse(ExpireTimeLayout, split[0])
	if err != nil {
		return time.Time{}, microerror.Mask(err)
	}

	return expiration, nil
}
