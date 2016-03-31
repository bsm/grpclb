package grpclb

import (
	"sync/atomic"
	"time"

	backendpb "github.com/bsm/grpclb/grpclb_backend_v1"
	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

type backend struct {
	cc  *grpc.ClientConn
	cln backendpb.LoadReportClient

	target  string
	address string
	score   int64

	closing chan struct{}
	closed  chan error
}

func newBackend(target, address string, queryInterval time.Duration) (*backend, error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	b := &backend{
		cc:  cc,
		cln: backendpb.NewLoadReportClient(cc),

		target:  target,
		address: address,

		closing: make(chan struct{}),
		closed:  make(chan error),
	}

	if err := b.updateScore(); err != nil {
		cc.Close()
		return nil, err
	}

	go b.loop(queryInterval)
	return b, nil
}

func (b *backend) Server() *balancerpb.Server {
	return &balancerpb.Server{
		Address: b.address,
		Score:   b.Score(),
	}
}

func (b *backend) Score() int64 {
	return atomic.LoadInt64(&b.score)
}

func (b *backend) Close() error {
	close(b.closing)
	return <-b.closed
}

func (b *backend) loop(queryInterval time.Duration) {
	t := time.NewTicker(queryInterval)
	defer t.Stop()

	for {
		select {
		case <-b.closing:
			b.closed <- b.cc.Close()
			close(b.closed)
			return
		case <-t.C:
			err := b.updateScore()
			if err != nil {
				grpclog.Printf("error retrieving load score for %s from %s: %s", b.target, b.address, err)
			}
		}
	}
}

func (b *backend) updateScore() error {
	resp, err := b.cln.Load(context.Background(), &backendpb.LoadRequest{})
	if err != nil {
		if grpc.Code(err) == codes.Unimplemented {
			return nil
		}
		return err
	}
	atomic.StoreInt64(&b.score, resp.Score)
	return nil
}
