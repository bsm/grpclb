package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	balancerpb "github.com/bsm/grpclb/grpclb_balancer_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var flags struct {
	addr, target string
	version      bool
}

var version string

func init() {
	flag.StringVar(&flags.addr, "a", "127.0.0.1:8383", "Server address. Default: 127.0.0.1:8383")
	flag.StringVar(&flags.target, "t", "service", "Service name/target. Default: service")
	flag.BoolVar(&flags.version, "version", false, "Print version and exit")
}

func main() {
	flag.Parse()

	if flags.version {
		fmt.Println(filepath.Base(os.Args[0]), version)
		return
	}

	if err := run(); err != nil {
		log.Fatal("FATAL", err.Error())
	}
}

func run() error {
	cc, err := grpc.Dial(flags.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	bc := balancerpb.NewLoadBalancerClient(cc)
	resp, err := bc.Servers(context.Background(), &balancerpb.ServersRequest{
		Target: flags.target,
	})
	if err != nil {
		return err
	}

	for _, srv := range resp.Servers {
		fmt.Printf("%d\t%s\n", srv.Score, srv.Address)
	}
	return nil
}
