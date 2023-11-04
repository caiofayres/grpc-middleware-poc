package server

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "tcc/proto"
)

type server struct {
	pb.UnimplementedCounterServer
}

func (s *server) Count(ctx context.Context, in *pb.CountRequest) (*pb.CountReply, error) {
	ctxVal := ctx.Value("message")
	if ctxVal == nil {
		log.Println("Handle server context value is nil")
	}
	log.Printf("Server handler received: %v, with context value: %v", in.GetMessage(), ctxVal)
	resp := pb.CountReply{MessageSize: int32(len(in.GetMessage()))}
	return &resp, nil
}

func UselessMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctxVal := ctx.Value("message")
	if ctxVal == nil {
		log.Println("Middleware server context value is nil")
		ctxVal = 1
	}
	ctx = context.WithValue(ctx, "message", ctxVal)
	msg := req.(*pb.CountRequest)
	log.Printf("Server middleware received: %v, with context value: %v", msg.GetMessage(), ctxVal)
	ctxVal = ctxVal.(int) + 1
	msg.Message = "ServerMiddleware"
	return handler(ctx, msg)
}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":30301")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(UselessMiddleware),
	)
	pb.RegisterCounterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}