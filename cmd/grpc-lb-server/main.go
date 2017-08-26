package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	lbpb "github.com/bsm/grpclb/grpclb_backend_v1"
	"github.com/bsm/grpclb/load"
	"github.com/hashicorp/consul/api"
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

var flags struct {
	addr, service string
}

func init() {
	flag.StringVar(&flags.addr, "a", "127.0.0.1:50051", "Backend address. Default: 127.0.0.1:50051")
	flag.StringVar(&flags.service, "s", "service", "Service name. Default: service")
}

func main() {
	flag.Parse()

	deregister, err := register(flags.addr, flags.service)
	if err != nil {
		log.Fatalf("failed to register to consul: %v", err)
		return
	}
	defer deregister()

	if err := run(flags.addr); err != nil {
		log.Fatal("FATAL", err.Error())
	}
}

// register registers the current backend to Consul
func register(address, serviceName string) (deregister func(), err error) {
	host, strPort, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		return nil, err
	}

	var (
		ttl                = 3 * time.Second
		updateTTLInterval  = 1 * time.Second
		deregisterCritical = 10 * time.Second
		id                 = fmt.Sprintf("%v@%v:%v", serviceName, host, port)
	)

	client, err := api.NewClient(&api.Config{})
	if err != nil {
		return nil, err
	}

	s := &api.AgentServiceRegistration{
		ID:      id,
		Name:    serviceName,
		Tags:    []string{},
		Address: host,
		Port:    port,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: deregisterCritical.String(),
			TTL: ttl.String(),
		},
	}
	agent := client.Agent()
	err = agent.ServiceRegister(s)
	if err != nil {
		return nil, err
	}

	closing := make(chan struct{})
	closed := make(chan struct{})

	// Updates TTL asynchronously
	go func() {
		t := time.NewTicker(updateTTLInterval)
		defer t.Stop()

		for {
			select {
			case <-closing:
				close(closed)
				return
			case <-t.C:
				if err := agent.UpdateTTL("service:"+id, "", api.HealthPassing); err != nil {
					log.Printf("Error: %v\n", err)
				}
			}
		}
	}()

	return func() {
		// Deregisters the current backend from Consul
		err := agent.ServiceDeregister(id)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		close(closing)
		<-closed
	}, nil
}

func run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	r := load.NewRateReporter(time.Minute)
	pb.RegisterGreeterServer(s, &GreeterServer{reporter: r})
	lbpb.RegisterLoadReportServer(s, r)
	return s.Serve(lis)
}
