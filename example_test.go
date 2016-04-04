package grpclb_test

import (
	"log"

	"github.com/bsm/grpclb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func ExampleNewPicker() {
	// Create a load-balanced address picker
	picker := grpclb.NewPicker("helloworld", &grpclb.PickerConfig{
		Address: "127.0.0.1:8383",
	})

	// Set up a load-balanced connection to the server.
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithPicker(picker))
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
