package grpclb

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

var (
	_ naming.Resolver = (*resolver)(nil)
	_ naming.Watcher  = (*watcher)(nil)
)

var errClosed = errors.New("grpclb: closed")

type resolver struct {
	opt *Options
}

// NewResolver creates a naming resolver
func NewResolver(opt *Options) naming.Resolver {
	if opt == nil {
		opt = new(Options)
	}

	return &resolver{opt: opt.norm()}
}

// Resolve implements naming.Resolver
func (r *resolver) Resolve(target string) (naming.Watcher, error) {
	cc, err := grpc.Dial(r.opt.Address, r.opt.DialOptions...)
	if err != nil {
		return nil, err
	}

	return &watcher{
		target: target,
		cc:     cc,
		lb:     pb.NewLoadBalancerClient(cc),
		check:  r.opt.UpdateInterval,
	}, nil
}

type watcher struct {
	target string

	cc io.Closer
	lb pb.LoadBalancerClient

	check  time.Duration
	addrs  []string
	mutex  sync.Mutex
	closed int32
}

// Next implements naming.Watcher
func (w *watcher) Next() ([]*naming.Update, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	ups, err := w.updates()
	if err != nil || ups != nil {
		return ups, err
	}

	ticker := time.NewTicker(w.check)
	defer ticker.Stop()

	for _ = range ticker.C {
		ups, err := w.updates()
		if err != nil || ups != nil {
			return ups, err
		}
	}
	return nil, nil
}

// Close implements naming.Watcher
func (w *watcher) Close() {
	atomic.StoreInt32(&w.closed, 1)
	_ = w.cc.Close()
}

func (w *watcher) updates() ([]*naming.Update, error) {
	if w.isClosed() {
		return nil, errClosed
	}

	latest, err := w.poll()
	if err != nil {
		return nil, err
	}

	var ups []*naming.Update
	for _, addr := range latest {
		if !contains(w.addrs, addr) {
			ups = append(ups, &naming.Update{Op: naming.Add, Addr: addr})
		}
	}
	for _, addr := range w.addrs {
		if !contains(latest, addr) {
			ups = append(ups, &naming.Update{Op: naming.Delete, Addr: addr})
		}
	}

	w.addrs = latest
	return ups, nil
}

func (w *watcher) isClosed() bool {
	return atomic.LoadInt32(&w.closed) != 0
}

func (w *watcher) poll() ([]string, error) {
	resp, err := w.lb.Servers(context.Background(), &pb.ServersRequest{
		Target: w.target,
	})
	if err != nil {
		return nil, err
	}

	addrs := make([]string, len(resp.Servers))
	for i, s := range resp.Servers {
		addrs[i] = s.Address
	}
	return addrs, nil
}
