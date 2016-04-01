package load_test

import (
	"log"
	"net"
	"time"

	lbpb "github.com/bsm/grpclb/grpclb_backend_v1"
	"github.com/bsm/grpclb/load"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// GreeterServer is used to implement helloworld.GreeterServer.
type GreeterServer struct {
	reporter *load.RateReporter
}

// SayHello implements helloworld.GreeterServer
// It increments rate to report load metrics to the load balancer.
func (s *GreeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.reporter.Increment(1)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func Example() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	r := load.NewRateReporter(time.Minute)
	pb.RegisterGreeterServer(s, &GreeterServer{reporter: r})
	lbpb.RegisterLoadReportServer(s, r)
	s.Serve(lis)
}
