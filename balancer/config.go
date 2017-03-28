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
		// MaxFailures allows up to this number of failures for backend
		// before removing it from set of backends for service.
		// Negative or zero value ignores failures completely.
		// Default: 3
		MaxFailures int
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
	if c.LoadReport.MaxFailures == 0 {
		c.LoadReport.MaxFailures = 3
	}
	return c
}
