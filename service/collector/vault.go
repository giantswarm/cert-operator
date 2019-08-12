package collector

import (
	"context"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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
		prometheus.BuildFQName("cert_operator", "vault", "token_expire_time_seconds"),
		"A metric of the expire time of Vault tokens as unix seconds.",
		nil,
		nil,
	)
)

type VaultConfig struct {
	Logger      micrologger.Logger
	VaultClient *vault.Client
}

type Vault struct {
	logger      micrologger.Logger
	vaultClient *vault.Client
}

func NewVault(config VaultConfig) (*Vault, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VaultClient must not be empty", config)
	}

	v := &Vault{
		logger:      config.Logger,
		vaultClient: config.VaultClient,
	}

	return v, nil
}

func (v *Vault) Collect(ch chan<- prometheus.Metric) error {
	ctx := context.Background()

	secret, err := v.vaultClient.Auth().Token().LookupSelf()
	if IsVaultAccess(err) {
		v.logger.LogCtx(ctx, "level", "debug", "message", "vault not reachable")
		v.logger.LogCtx(ctx, "level", "debug", "message", "canceling collection")
		return nil

	} else if err != nil {
		return microerror.Mask(err)
	}

	expiration, err := expirationFromSecret(secret)
	if err != nil {
		return microerror.Mask(err)
	}

	ch <- prometheus.MustNewConstMetric(
		tokenExpireTimeDesc,
		prometheus.GaugeValue,
		float64(expiration.Unix()),
	)

	return nil
}

func (v *Vault) Describe(ch chan<- *prometheus.Desc) error {
	ch <- tokenExpireTimeDesc
	return nil
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
		return time.Time{}, microerror.Maskf(executionFailedError, "%#q must be string in order to collect metrics for the Vault token expiration", value)
	}

	split := strings.Split(e, ".")
	if len(split) == 0 {
		return time.Time{}, microerror.Maskf(executionFailedError, "%#q must have at least one item in order to collect metrics for the Vault token expiration", e)
	}

	expiration, err := time.Parse(ExpireTimeLayout, split[0])
	if err != nil {
		return time.Time{}, microerror.Mask(err)
	}

	return expiration, nil
}
