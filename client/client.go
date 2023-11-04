package client

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "tcc/proto"
)

func SendMessage(message string) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:30301", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCounterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Count(ctx, &pb.CountRequest{Message: message})
	if err != nil {
		log.Fatalf("could not count: %v", err)
	}
	log.Printf("Size: %d", r.GetMessageSize())
}