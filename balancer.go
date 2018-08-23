package grpclb

import (
	"io"
	"sync"
	"time"

	"golang.org/x/net/context"

	pb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// PickFirst returns a Balancer that always selects the first address returned
// from the Resolver
func PickFirst(opt *Options) grpc.Balancer {
	if opt == nil {
		opt = new(Options)
	}

	return &pickFirst{
		nfy: make(chan []grpc.Address, 3),
		opt: opt.norm(),
	}
}

type pickFirst struct {
	target string
	opt    *Options

	addrs []grpc.Address
	nfy   chan []grpc.Address
	mu    sync.Mutex

	cc io.Closer
	lb pb.LoadBalancerClient

	closing, closed chan struct{}
}

func (p *pickFirst) Start(target string, _ grpc.BalancerConfig) error {
	cc, err := grpc.Dial(p.opt.Address, p.opt.DialOptions...)
	if err != nil {
		return err
	}

	p.target = target
	p.cc = cc
	p.lb = pb.NewLoadBalancerClient(cc)

	p.closing = make(chan struct{})
	p.closed = make(chan struct{})
	go p.loop()
	return nil
}

func (p *pickFirst) Up(_ grpc.Address) func(error) { return nil }
func (p *pickFirst) Notify() <-chan []grpc.Address { return p.nfy }

func (p *pickFirst) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, _ func(), err error) {
	if first, ok := p.first(); ok {
		addr = first
		return
	} else if !opts.BlockingWait || p.closing == nil {
		return
	}

	ticker := time.NewTicker(p.opt.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			err = status.FromContextError(ctx.Err()).Err()
			return
		case <-p.closing:
			err = grpc.ErrClientConnClosing
			return
		case <-ticker.C:
			if first, ok := p.first(); ok {
				addr = first
				return
			}
		}
	}
}

func (p *pickFirst) Close() (err error) {
	if p.closing != nil {
		close(p.closing)
		<-p.closed
		p.closing = nil
	}

	if p.cc != nil {
		err = p.cc.Close()
		p.cc = nil
	}
	return
}

func (p *pickFirst) first() (_ grpc.Address, _ bool) {
	p.mu.Lock()
	addrs := p.addrs
	p.mu.Unlock()

	if len(addrs) != 0 {
		return addrs[0], true
	}
	return
}

func (p *pickFirst) update() error {
	resp, err := p.lb.Servers(context.Background(), &pb.ServersRequest{
		Target: p.target,
	})
	if err != nil {
		return err
	}

	addrs := make([]grpc.Address, len(resp.Servers))
	for i, s := range resp.Servers {
		addrs[i].Addr = s.Address
	}

	p.mu.Lock()
	p.addrs = addrs
	p.mu.Unlock()
	p.nfy <- addrs

	return nil
}

func (p *pickFirst) loop() {
	ticker := time.NewTicker(p.opt.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.closing:
			close(p.closed)
			return
		case <-ticker.C:
			if err := p.update(); err != nil {
				grpclog.Printf("Failed PickFirst balancer lookup: %s", err.Error())
			}
		}
	}
}
