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
	log.Printf("Received: %v", in.GetMessage())
	resp := pb.CountReply{MessageSize: int32(len(in.GetMessage()))}
	return &resp, nil
}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":30301")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCounterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}