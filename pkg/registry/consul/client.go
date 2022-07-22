package consul

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/gole-dev/gole/pkg/registry"
)

// Client is consul client config
type Client struct {
	client *api.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates consul client
func NewClient(cli *api.Client) *Client {
	c := &Client{client: cli}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

// Service get services from consul
func (d *Client) Service(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := d.client.Health().Service(service, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}

	services := make([]*registry.ServiceInstance, 0)

	for _, entry := range entries {
		var version string
		for _, tag := range entry.Service.Tags {
			strs := strings.SplitN(tag, "=", 2)
			if len(strs) == 2 && strs[0] == "version" {
				version = strs[1]
			}
		}
		var endpoints []string
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			endpoints = append(endpoints, addr.Address)
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}
	return services, meta.LastIndex, nil
}

// Register register service instance to consul
func (d *Client) Register(ctx context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress)
	var addr string
	var port uint64
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr = raw.Hostname()
		port, _ = strconv.ParseUint(raw.Port(), 10, 16)
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
		Address:         addr,
		Port:            int(port),
		Checks: []*api.AgentServiceCheck{
			{
				CheckID:                        "service:" + svc.ID,
				TTL:                            "30s",
				Status:                         "passing",
				DeregisterCriticalServiceAfter: "90s",
			},
		},
	}
	if enableHealthCheck {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", addr, port),
			Interval:                       "20s",
			Status:                         "passing",
			DeregisterCriticalServiceAfter: "90s",
		})
	}
	err := d.client.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(time.Second * 20)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = d.client.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass")
			case <-d.ctx.Done():
				return
			}
		}
	}()
	return nil
}

// Deregister deregister service by service ID
func (d *Client) Deregister(ctx context.Context, serviceID string) error {
	d.cancel()
	return d.client.Agent().ServiceDeregister(serviceID)
}
