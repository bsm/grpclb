package balancer

import (
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
)

type backends struct {
	target      string
	maxFailures int

	set map[string]*backend
	mu  sync.RWMutex

	closing, closed chan struct{}
}

func newBackends(target string, queryInterval time.Duration, maxFailures int) *backends {
	b := &backends{
		target:      target,
		maxFailures: maxFailures,

		set: make(map[string]*backend),

		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	go b.loop(queryInterval)
	return b
}

func (b *backends) Servers() []*balancerpb.Server {
	b.mu.RLock()
	defer b.mu.RUnlock()

	servers := make([]*balancerpb.Server, 0, len(b.set))
	for _, b := range b.set {
		servers = append(servers, b.Server())
	}
	return servers
}

func (b *backends) Update(addrs []string) (err error) {
	var removed []*backend
	var added []string

	sort.Strings(addrs)

	b.mu.Lock()
	for addr, backend := range b.set {
		if pos := sort.SearchStrings(addrs, addr); !(pos < len(addrs) && addrs[pos] == addr) {
			removed = append(removed, backend)
			delete(b.set, addr)
		}
	}

	for _, addr := range addrs {
		if _, ok := b.set[addr]; !ok {
			added = append(added, addr)
		}
	}
	b.mu.Unlock()

	// Close removed backends
	for _, b := range removed {
		_ = b.Close()
	}

	// Connect to added backends, in parallel
	if len(added) != 0 {
		err = b.connectAll(addrs)
	}
	return
}

func (b *backends) Close() error {
	close(b.closing)
	<-b.closed
	return b.Update(nil)
}

func (b *backends) connectAll(addrs []string) (err error) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := range addrs {
		wg.Add(1)
		go func(addr string) {
			if e := b.connect(addr); e != nil {
				mu.Lock()
				err = e
				mu.Unlock()
			}
			wg.Done()
		}(addrs[i])
	}
	wg.Wait()
	return
}

func (b *backends) connect(addr string) error {
	backend, err := newBackend(b.target, addr, b.maxFailures)
	if err != nil {
		return err
	}

	b.mu.Lock()
	b.set[addr] = backend
	b.mu.Unlock()
	return nil
}

func (b *backends) loop(queryInterval time.Duration) {
	t := time.NewTicker(queryInterval)
	defer t.Stop()

	for {
		select {
		case <-b.closing:
			close(b.closed)
			return
		case <-t.C:
			if err := b.updateBackendScores(); err != nil {
				grpclog.Printf("failed to update backend load scores: %s", err)
			}
		}
	}
}

func (b *backends) updateBackendScores() error {
	b.mu.RLock()
	set := b.set
	b.mu.RUnlock()

	succeeded := make([]string, 0, len(set))
	for addr, backend := range set {
		if err := backend.UpdateScore(); err != nil {
			continue
		}
		succeeded = append(succeeded, addr)
	}
	return b.Update(succeeded)
}
