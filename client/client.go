package client

import (
	"context"
	"log"

	pb "tcc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// func UselessMiddleware(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	ctxVal := ctx.Value("message")
// 	if ctxVal == nil {
// 		log.Println("Client middleware context value is nil")
// 		ctxVal = 1
// 	}
// 	msg := req.(*pb.CountRequest)
// 	log.Printf("Client middleware received: %v, with context value: %v", msg.GetMessage(), ctxVal)
// 	ctxVal = ctxVal.(int) + 1
// 	ctx = context.WithValue(ctx, "message", ctxVal)
// 	msg.Message = "ClientMiddleware"
// 	return invoker(ctx, method, msg, reply, cc, opts...)
// }
func GetPerson(id string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		"localhost:30301", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(UselessMiddleware),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	r := pb.NewPersonReadClient(conn)

	p, err := r.GetPerson(context.Background(), &pb.PersonId{Id: id})
		if err != nil {
				e, _ := status.FromError(err)
				log.Fatalf("could not getPerson: %v", e.Message())
				return
		}
		log.Printf("%v", p)
}

func NewPerson(id, name, surname string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		"localhost:30301", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithUnaryInterceptor(UselessMiddleware),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	w := pb.NewPersonWriteClient(conn)
	
	w.NewPerson(context.Background(), &pb.Person{Id: &pb.PersonId{Id: id}, Name: name, Surname: surname})

	// Contact the server and print out its response.
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// ctx = context.WithValue(ctx, "message", 1)
	// defer cancel()
	// r, err := c.Count(ctx, &pb.CountRequest{Message: message})
	// if err != nil {
	// 	log.Fatalf("could not count: %v", err)
	// }
	// log.Printf("Size: %d", r.GetMessageSize())
}