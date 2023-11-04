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

func UselessMiddleware(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctxVal := ctx.Value("message")
	if ctxVal == nil {
		log.Println("Client middleware context value is nil")
		ctxVal = 1
	}
	msg := req.(*pb.CountRequest)
	log.Printf("Client middleware received: %v, with context value: %v", msg.GetMessage(), ctxVal)
	ctxVal = ctxVal.(int) + 1
	ctx = context.WithValue(ctx, "message", ctxVal)
	msg.Message = "ClientMiddleware"
	return invoker(ctx, method, msg, reply, cc, opts...)
}

func SendMessage(message string) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		"localhost:30301", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UselessMiddleware),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCounterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx = context.WithValue(ctx, "message", 1)
	defer cancel()
	r, err := c.Count(ctx, &pb.CountRequest{Message: message})
	if err != nil {
		log.Fatalf("could not count: %v", err)
	}
	log.Printf("Size: %d", r.GetMessageSize())
}