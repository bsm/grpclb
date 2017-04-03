package balancer

import (
	"time"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"google.golang.org/grpc/grpclog"
)

type service struct {
	target    string
	discovery Discovery
	backends  *backends

	closing, closed chan struct{}
}

func newService(target string, discovery Discovery, discoveryInterval, loadReportInterval time.Duration, maxFailures int) (*service, error) {
	s := &service{
		target:    target,
		discovery: discovery,
		backends:  newBackends(target, loadReportInterval, maxFailures),

		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}
	if err := s.updateBackends(); err != nil {
		// close ALL backend connections (some of them could succeed, don't leak these):
		_ = s.backends.Close()
		return nil, err
	}

	go s.loop(discoveryInterval)
	return s, nil
}

func (s *service) Servers() []*balancerpb.Server { return s.backends.Servers() }

func (s *service) Close() {
	close(s.closing)
	<-s.closed
}

func (s *service) loop(discoveryInterval time.Duration) {
	t := time.NewTicker(discoveryInterval)
	defer t.Stop()

	for {
		select {
		case <-s.closing:
			_ = s.backends.Close()
			close(s.closed)
			return
		case <-t.C:
			err := s.updateBackends()
			if err != nil {
				grpclog.Printf("error on service discovery of %s: %s", s.target, err)
			}
		}
	}
}

func (s *service) updateBackends() error {
	addrs, err := s.discovery.Resolve(s.target)
	if err != nil {
		return err
	}

	return s.backends.Update(addrs)
}
