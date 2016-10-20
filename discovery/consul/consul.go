// Package consul implements grpclb service discovery via Consul (consul.io).
// Service names may be specified as: svcname,tag1,tag2,tag3
package consul

import (
	"fmt"
	"strings"

	"github.com/bsm/grpclb/balancer"
	"github.com/hashicorp/consul/api"
)

type discovery struct {
	*api.Client
}

// New returns an implementation of balancer.Discovery interface.
func New(config *api.Config) (balancer.Discovery, error) {
	if config == nil {
		config = api.DefaultConfig()
	}

	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &discovery{Client: c}, nil
}

// Resolve implements balancer.Discovery
func (d *discovery) Resolve(target string) ([]string, error) {
	service, tags := splitTarget(target)
	return d.query(service, tags)
}

func (d *discovery) query(service string, tags []string) ([]string, error) {
	var tag string
	if len(tags) > 0 {
		tag = tags[0]
	}

	// Query watcher
	entries, _, err := d.Health().Service(service, tag, true, nil)
	if err != nil {
		return nil, err
	}

	// If more than one tag is passed, we need to filter
	if len(tags) > 1 {
		entries = filterEntries(entries, tags[1:])
	}

	return parseEntries(entries), nil
}

// --------------------------------------------------------------------

func splitTarget(target string) (service string, tags []string) {
	parts := strings.Split(target, ",")
	service = parts[0]
	if len(parts) > 1 {
		tags = parts[1:]
	}
	return
}

func parseEntries(entries []*api.ServiceEntry) []string {
	res := make([]string, len(entries))
	for i, entry := range entries {
		addr := entry.Node.Address
		if entry.Service.Address != "" {
			addr = entry.Service.Address
		}
		res[i] = fmt.Sprintf("%s:%d", addr, entry.Service.Port)
	}
	return res
}

func filterEntries(entries []*api.ServiceEntry, requiredTags []string) (res []*api.ServiceEntry) {
EntriesLoop:
	for _, entry := range entries {
		for _, required := range requiredTags {
			var found bool
			for _, tag := range entry.Service.Tags {
				if tag == required {
					found = true
					break
				}
			}
			if !found {
				continue EntriesLoop
			}
		}
		res = append(res, entry)
	}
	return
}
