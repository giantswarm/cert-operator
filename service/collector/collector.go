package collector

import (
	"context"
	"fmt"
	"sync"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	vault "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Logger      micrologger.Logger
	VaultClient *vault.Client
}

type Collector struct {
	logger      micrologger.Logger
	vaultClient *vault.Client

	bootOnce sync.Once
}

func New(config Config) (*Collector, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.VaultClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VaultClient must not be empty", config)
	}

	c := &Collector{
		logger:      config.Logger,
		vaultClient: config.VaultClient,

		bootOnce: sync.Once{},
	}

	return c, nil
}

func (c *Collector) Boot(ctx context.Context) {
	c.bootOnce.Do(func() {
		c.logger.LogCtx(ctx, "level", "debug", "message", "registering collector")

		err := prometheus.Register(prometheus.Collector(c))
		if IsAlreadyRegisteredError(err) {
			c.logger.LogCtx(ctx, "level", "debug", "message", "collector already registered")
		} else if err != nil {
			c.logger.Log("level", "error", "message", "registering collector failed", "stack", fmt.Sprintf("%#v", err))
		} else {
			c.logger.LogCtx(ctx, "level", "debug", "message", "registered collector")
		}
	})
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Log("level", "debug", "message", "collecting metrics")

	var wg sync.WaitGroup

	collectFuncs := []func(chan<- prometheus.Metric){
		c.collectVaultMetrics,
	}

	for _, f := range collectFuncs {
		wg.Add(1)

		go func(f func(chan<- prometheus.Metric)) {
			defer wg.Done()
			f(ch)
		}(f)
	}

	wg.Wait()

	c.logger.Log("level", "debug", "message", "collected metrics")
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tokenExpireTimeDesc
}
