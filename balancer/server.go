package balancer

import (
	"sync"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Discovery describes a service discovery interface
type Discovery interface {
	// Resolve accepts a target string and returns a list of addresses
	Resolve(target string) ([]string, error)
}

// Server is a gRPC server
type Server struct {
	discovery Discovery
	config    *Config

	cache map[string]*service
	mu    sync.RWMutex
}

// New creates a new Server instance with a given resolver
func New(discovery Discovery, config *Config) *Server {
	if config == nil {
		config = new(Config)
	}
	return &Server{
		discovery: discovery,
		config:    config.norm(),
		cache:     make(map[string]*service),
	}
}

// Servers implements RPC server
func (b *Server) Servers(ctx context.Context, req *balancerpb.ServersRequest) (*balancerpb.ServersResponse, error) {
	if req.Target == "" {
		return &balancerpb.ServersResponse{}, nil
	}

	servers, err := b.GetServers(req.Target)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &balancerpb.ServersResponse{Servers: servers}, nil
}

// GetServers retrieves all known servers for a target service
func (b *Server) GetServers(target string) ([]*balancerpb.Server, error) {
	// Retrieve from cache, if available
	b.mu.RLock()
	svc, cached := b.cache[target]
	b.mu.RUnlock()
	if cached {
		return svc.Servers(), nil
	}

	// Create a new service connection
	newSvc, err := newService(target, b.discovery, b.config.Discovery.Interval, b.config.LoadReport.Interval, b.config.LoadReport.MaxFailures)
	if err != nil {
		return nil, err
	}

	// Apply write lock
	b.mu.Lock()
	svc, cached = b.cache[target]
	if !cached {
		svc = newSvc
		b.cache[target] = svc
	}
	b.mu.Unlock()

	// Close new connection again, if svc was cached
	if cached {
		newSvc.Close()
	}

	// Retrieve and filter/sort servers
	servers := b.config.Balancer.Balance(svc.Servers())
	return servers, nil
}

// Reset resets services cache
func (b *Server) Reset() {
	b.mu.Lock()
	services := make([]*service, 0, len(b.cache))
	for target, svc := range b.cache {
		services = append(services, svc)
		delete(b.cache, target)
	}
	b.mu.Unlock()

	for _, svc := range services {
		svc.Close()
	}
}
