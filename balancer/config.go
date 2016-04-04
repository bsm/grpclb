package balancer

import "time"

// Config options
type Config struct {
	// Balancer chooses can decide whether to provide a complete list,
	// a subset, or a specific list of "picked" servers in a particular order.
	// Default: LeastBusyBalancer
	Balancer Balancer

	Discovery struct {
		// Interval between service discovery checks
		// Default: 5s
		Interval time.Duration
	}
	LoadReport struct {
		// Interval between service load pings
		// Default: 5s
		Interval time.Duration
	}
}

func (c *Config) norm() *Config {
	if c.Balancer == nil {
		c.Balancer = NewLeastBusyBalancer()
	}
	if c.Discovery.Interval == 0 {
		c.Discovery.Interval = 5 * time.Second
	}
	if c.LoadReport.Interval == 0 {
		c.LoadReport.Interval = 5 * time.Second
	}
	return c
}
