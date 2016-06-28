package grpclb

import (
	"testing"

	pb "github.com/bsm/grpclb/grpclb_balancer_v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "grpclb")
}

type nopCloser struct{}

func (nopCloser) Close() error { return nil }

type mockClient struct {
	pb.LoadBalancerClient
	S []string
	E error
}

func (c *mockClient) Servers(ctx context.Context, in *pb.ServersRequest, opts ...grpc.CallOption) (*pb.ServersResponse, error) {
	servers := make([]*pb.Server, len(c.S))
	for i, s := range c.S {
		servers[i] = &pb.Server{Address: s}
	}
	return &pb.ServersResponse{Servers: servers}, c.E
}
