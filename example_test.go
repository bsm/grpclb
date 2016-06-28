package grpclb

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func ExampleNewResolver() {
	const target = "helloworld"

	// Create a round-robin load-balancer
	balancer := grpc.RoundRobin(NewResolver(&Options{
		Address: "127.0.0.1:8383",
	}))

	// Set up a load-balanced connection to the server.
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func ExamplePickFirst() {
	const target = "helloworld"

	// Create a pick-first load-balancer
	balancer := PickFirst(&Options{
		Address: "127.0.0.1:8383",
	})

	// Set up a load-balanced connection to the server.
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancer(balancer))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
