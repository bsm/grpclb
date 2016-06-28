package grpclb

import (
	"time"

	"google.golang.org/grpc"
)

// Options is the picker configuration
type Options struct {
	// Address is the address of the load balancer.
	// Default: 127.0.0.1:8383
	Address string
	// DialOptions define custom dial options.
	// Default: [WithInsecure].
	DialOptions []grpc.DialOption
	// UpdateInterval is the query interval for
	// updating known servers.
	// Default: 2s
	UpdateInterval time.Duration
}

func (c *Options) norm() *Options {
	if c.Address == "" {
		c.Address = "127.0.0.1:8383"
	}
	if c.DialOptions == nil {
		c.DialOptions = []grpc.DialOption{grpc.WithInsecure()}
	}
	if c.UpdateInterval == 0 {
		c.UpdateInterval = 2 * time.Second
	}
	return c
}

func contains(vv []string, v string) bool {
	for _, x := range vv {
		if x == v {
			return true
		}
	}
	return false
}
