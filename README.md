# grpclb

External Load Balancing Service solution for gRPC written in Go. The approach follows the
[proposal](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md) outlined by the
core gRPC team.

grpclb load-balancer provides a neutral API which can be integrated with various service discovery
frameworks. An example service discovery implementation is provided for [Consul](discovery/consul/).

## Usage

### Load Balancer

Please also see the bootstrap for [Consul backed load-balancers](cmd/grpc-lb-consul/main.go)
as a reference for building load balancers. Either use the command directly or build your very own.

### Server

Servers can optionally report load to the Load Balancer. An example:

```go
import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/bsm/grpclb/grpclb_backend_v1"
)

type MyRPCServer struct {}

// Load implements github.com/bsm/grpclb/grpclb_backend_v1.LoadReportServer
// Returns a numeric proportional indicator - the higher the busier.
func (s *MyRPCServer) Load(_ context.Context, _ *pb.LoadRequest) (*pb.LoadResponse, error) {
	return &pb.LoadResponse{Score: 1234}
}

func main() {
	srv := grpc.NewServer()
	rpc := &MyRPCServer{}
	pb.RegisterLoadReportServer(srv, rpc)

	lis, _ := net.Listen("tcp", ":8080")
	defer lis.Close()

	_ = srv.Serve(lis)
}

Full [Documentation](godoc.org/github.com/bsm/grpclb/grpclb_backend_v1#LoadReportServer)

### Client

Coming soon...

Full [Documentation](godoc.org/github.com/bsm/grpclb/grpclb_balancer_v1#LoadBalancerClient)

## TODO

* Implement server helpers for collecting and exposing load information
* Implement client helpers for connecting with load balancers

## Licence

```
Copyright (c) 2016 Black Square Media

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
