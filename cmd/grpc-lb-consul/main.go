package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/bsm/grpclb/balancer"
	"github.com/bsm/grpclb/discovery/consul"
	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var flags struct {
	addr    string
	consul  string
	version bool
}

var version string

func init() {
	flag.StringVar(&flags.addr, "addr", ":8383", "Bind address. Default: :8383")
	flag.StringVar(&flags.consul, "consul", "127.0.0.1:8500", "Consul API address. Default: 127.0.0.1:8500")
	flag.BoolVar(&flags.version, "version", false, "Print version and exit")
}

func main() {
	flag.Parse()

	if flags.version {
		fmt.Println(filepath.Base(os.Args[0]), version)
		return
	}

	if err := listenAndServe(); err != nil {
		log.Fatal("FATAL", err.Error())
	}
}

func listenAndServe() error {
	config := api.DefaultConfig()
	config.Address = flags.consul

	discovery, err := consul.New(config)
	if err != nil {
		return err
	}

	lb := balancer.New(discovery, nil)
	defer lb.Reset()

	srv := grpc.NewServer()
	balancerpb.RegisterLoadBalancerServer(srv, lb)

	lis, err := net.Listen("tcp", flags.addr)
	if err != nil {
		return err
	}

	return srv.Serve(lis)
}
