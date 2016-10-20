package balancer

import (
	"math/rand"
	"sort"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
)

// Balancer algorithm interface
type Balancer interface {
	Balance([]*balancerpb.Server) []*balancerpb.Server
}

// BalancerFunc is a simple balancer function which implements the Balancer interface
type BalancerFunc func([]*balancerpb.Server) []*balancerpb.Server

// Balance implements Balance
func (f BalancerFunc) Balance(s []*balancerpb.Server) []*balancerpb.Server { return f(s) }

// --------------------------------------------------------------------

// NewRandomBalancer returns a balancer which returns all known servers in random order
func NewRandomBalancer() Balancer {
	return BalancerFunc(func(s []*balancerpb.Server) []*balancerpb.Server {
		n := len(s)
		for i := 0; i < n; i++ {
			r := i + rand.Intn(n-i)
			s[r], s[i] = s[i], s[r]
		}
		return s
	})
}

type busyServers []*balancerpb.Server

func (p busyServers) Len() int           { return len(p) }
func (p busyServers) Less(i, j int) bool { return p[i].Score < p[j].Score }
func (p busyServers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p busyServers) Shuffle() {
	for i := len(p) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
func (p busyServers) Sort() {
	sort.Sort(busyServers(p))

}

// NewLeastBusyBalancer returns a balancer which returns all known servers in priority order, from least to most busy
func NewLeastBusyBalancer() Balancer {
	return BalancerFunc(func(s []*balancerpb.Server) []*balancerpb.Server {
		bs := busyServers(s)
		bs.Shuffle()
		bs.Sort()
		return s
	})
}
