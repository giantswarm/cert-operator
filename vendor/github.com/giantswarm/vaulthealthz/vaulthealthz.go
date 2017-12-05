package vaulthealthz

import (
	"context"
	"fmt"
	"strings"
	"time"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/microendpoint/service/healthz"
)

const (
	// Description describes which functionality this health check implements.
	Description = "Ensure Vault API availability."
	// Name is the identifier of the health check. This can be used for emitting
	// metrics.
	Name = "vault"
	// SuccessMessage is the message returned in case the health check did not
	// fail.
	SuccessMessage = "all good"
	// Timeout is the time being waited until timing out health check, which
	// renders its result unsuccessful.
	Timeout = 5 * time.Second
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

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	Logger      micrologger.Logger
	VaultClient *vaultapi.Client

	// Settings.
	Timeout time.Duration
}

// DefaultConfig provides a default configuration to create a new healthz
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:      nil,
		VaultClient: nil,

		// Settings.
		Timeout: Timeout,
	}
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	logger      micrologger.Logger
	vaultClient *vaultapi.Client

	// Settings.
	timeout time.Duration
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "vault client must not be empty")
	}

	// Settings.
	if config.Timeout.Seconds() == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.Timout must not be empty")
	}

	newService := &Service{
		// Dependencies.
		logger:      config.Logger,
		vaultClient: config.VaultClient,

		// Settings.
		timeout: config.Timeout,
	}

	return newService, nil
}

// GetHealthz implements the health check for Vault. It does this by calling
// the Vault /sys/health endpoint. This checks that the we can connect to the
// Vault API and that the Vault token is valid.
func (s *Service) GetHealthz(ctx context.Context) (healthz.Response, error) {
	failed := false
	message := SuccessMessage
	{
		ch := make(chan string, 1)

		go func() {
			_, err := s.vaultClient.Sys().Health()
			if err != nil {
				ch <- err.Error()
				return
			}

			err = s.updateTokenTTLMetric()
			if err != nil {
				ch <- err.Error()
				return
			}

			ch <- ""
		}()

		select {
		case m := <-ch:
			if m != "" {
				failed = true
				message = m
			}
		case <-time.After(s.timeout):
			failed = true
			message = fmt.Sprintf("timed out after %s", s.timeout)
		}
	}

	response := healthz.Response{
		Description: Description,
		Failed:      failed,
		Message:     message,
		Name:        Name,
	}

	return response, nil
}

func (s *Service) updateTokenTTLMetric() error {
	secret, err := s.vaultClient.Auth().Token().LookupSelf()
	if err != nil {
		return microerror.Mask(err)
	}

	key, ok := secret.Data[ExpireTimeKey]
	if !ok {
		return microerror.Maskf(executionFailedError, "value of '%s' must exist in order to collect metrics for the Vault token expiration", ExpireTimeKey)
	}
	e, ok := key.(string)
	if !ok {
		return microerror.Maskf(executionFailedError, "'%#v' must be string in order to collect metrics for the Vault token expiration", e)
	}
	split := strings.Split(e, ".")
	if len(split) == 0 {
		return microerror.Maskf(executionFailedError, "'%#v' must have at least one item in order to collect metrics for the Vault token expiration", split)
	}
	expireTime := split[0]

	t, err := time.Parse(ExpireTimeLayout, expireTime)
	if err != nil {
		return microerror.Mask(err)
	}

	tokenExpireTimeGauge.Set(float64(t.Unix()))

	return nil
}
