package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/kirychukyurii/wasker/internal/config"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type ServiceDiscovery struct {
	Client *api.Client
}

func New(cfg config.Config, logger log.Logger) ServiceDiscovery {
	consulCfg := &api.Config{
		Address:    cfg.Consul.Addr(),
		Scheme:     cfg.Consul.Scheme,
		Datacenter: cfg.Consul.Datacenter,
	}

	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("consul client")
	}

	return ServiceDiscovery{
		Client: consulClient,
	}
}
