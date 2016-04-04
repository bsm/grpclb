package grpclb

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/transport"
)

// picker is a load-balanced unicast picker
type picker struct {
	target string
	config *PickerConfig

	conn  *grpc.Conn
	addrs []string
	mu    sync.Mutex

	lbconn *grpc.ClientConn
	lbclnt pb.LoadBalancerClient

	closing, closed chan struct{}
}

// PickerConfig is the picker configuration
type PickerConfig struct {
	// Address is the address of the load balancer.
	// Default: 127.0.0.1:8383
	Address string
	// DialOptions define custom dial options.
	// Default: [WithInsecure].
	DialOptions []grpc.DialOption
	// UpdateInterval is the query interval for
	// updating known servers.
	// Default: 10s
	UpdateInterval time.Duration
}

func (c *PickerConfig) norm() *PickerConfig {
	if c.Address == "" {
		c.Address = "127.0.0.1:8383"
	}
	if c.DialOptions == nil {
		c.DialOptions = []grpc.DialOption{grpc.WithInsecure()}
	}
	if c.UpdateInterval == 0 {
		c.UpdateInterval = 10 * time.Second
	}
	return c
}

// NewPicker creates a Picker to pick addresses from a load-balancer
// to connect.
func NewPicker(target string, config *PickerConfig) grpc.Picker {
	if config == nil {
		config = new(PickerConfig)
	}

	return &picker{
		target: target,
		config: config.norm(),
	}
}

func (p *picker) Init(cc *grpc.ClientConn) error {
	lbconn, err := grpc.Dial(p.config.Address, p.config.DialOptions...)
	if err != nil {
		return err
	}

	p.lbconn = lbconn
	p.lbclnt = pb.NewLoadBalancerClient(lbconn)

	if err := p.update(); err != nil {
		_ = p.lbconn.Close()
		return err
	}

	conn, err := grpc.NewConn(cc)
	if err != nil {
		_ = p.lbconn.Close()
		return err
	}
	p.conn = conn

	p.closing = make(chan struct{})
	p.closed = make(chan struct{})
	go p.loop()
	return nil
}

func (p *picker) Pick(ctx context.Context) (transport.ClientTransport, error) {
	return p.conn.Wait(ctx)
}

func (p *picker) PickAddr() (string, error) {
	p.mu.Lock()
	addrs := p.addrs
	p.mu.Unlock()

	if len(addrs) == 0 {
		return "", fmt.Errorf("there is no address available to pick")
	}

	addr := addrs[0]
	grpclog.Printf("picked %s for %s", addr, p.target)
	return addr, nil
}

func (p *picker) State() (grpc.ConnectivityState, error) {
	return 0, fmt.Errorf("State() is not supported for grpclb.Picker")
}

func (p *picker) WaitForStateChange(_ context.Context, _ grpc.ConnectivityState) (grpc.ConnectivityState, error) {
	return 0, fmt.Errorf("WaitForStateChange is not supported for grpclb.Picker")
}

func (p *picker) Close() (err error) {
	if p.closing != nil {
		close(p.closing)
		<-p.closed
	}

	if p.lbconn != nil {
		_ = p.lbconn.Close()
	}
	if p.conn != nil {
		err = p.conn.Close()
	}
	return
}

func (p *picker) loop() {
	ticker := time.NewTicker(p.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.closing:
			close(p.closed)
			return
		case <-ticker.C:
			if err := p.update(); err != nil {
				grpclog.Printf("Failed load-balancer lookup: %s", err.Error())
			}
		}
	}
}

func (p *picker) update() error {
	resp, err := p.lbclnt.Servers(context.Background(), &pb.ServersRequest{
		Target: p.target,
	})
	if err != nil {
		return err
	}

	addrs := make([]string, len(resp.Servers))
	for i, s := range resp.Servers {
		addrs[i] = s.Address
	}
	p.mu.Lock()
	p.addrs = addrs
	p.mu.Unlock()
	return nil
}
