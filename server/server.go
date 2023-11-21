package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	cc "tcc/cache"
	db "tcc/data"
	pb "tcc/proto"

	"google.golang.org/grpc"
)

type WriteServer struct {
	CacheService cc.CacheService
	DataService db.PersonDataService
	pb.UnimplementedPersonWriteServer
}
type ReadServer struct {
	CacheService cc.CacheService
	DataService db.PersonDataService
	pb.UnimplementedPersonReadServer
}
// service PersonWrite {
// 	rpc NewPerson (Person) returns (Empty) {}
//   }
  
//   service PersonRead {
// 	rpc GetPerson (PersonId) returns (Person) {}
//   }

func (s *WriteServer) NewPerson(ctx context.Context, in *pb.Person) (*pb.Empty, error) {
	p := db.NewPerson(in.GetId().GetId(), in.GetName(), in.GetSurname())
	s.DataService.Upsert(p)
	return &pb.Empty{}, nil
}

func (s *ReadServer) GetPerson(ctx context.Context, in *pb.PersonId) (*pb.Person, error) {
	p, err := s.DataService.Get(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.Person{
		Id: &pb.PersonId{Id: p.Id},
		Name: p.Name,
		Surname: p.Surname,
	}, nil
}

// func (s *server) Count(ctx context.Context, in *pb.CountRequest) (*pb.CountReply, error) {
// 	ctxVal := ctx.Value("message")
// 	if ctxVal == nil {
// 		log.Println("Handle server context value is nil")
// 	}
// 	log.Printf("Server handler received: %v, with context value: %v", in.GetMessage(), ctxVal)
// 	resp := pb.CountReply{MessageSize: int32(len(in.GetMessage()))}
// 	return &resp, nil
// }

// func UselessMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
// 	ctxVal := ctx.Value("message")
// 	if ctxVal == nil {
// 		log.Println("Middleware server context value is nil")
// 		ctxVal = 1
// 	}
// 	ctx = context.WithValue(ctx, "message", ctxVal)
// 	msg := req.(*pb.CountRequest)
// 	log.Printf("Server middleware received: %v, with context value: %v", msg.GetMessage(), ctxVal)
// 	ctxVal = ctxVal.(int) + 1
// 	msg.Message = "ServerMiddleware"
// 	return handler(ctx, msg)
// }

func (s *ReadServer) CacheReadMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod != "/proto.PersonRead/GetPerson" {
		p, ok := req.(*pb.Person)
		if !ok {
			return nil, errors.New("bad request")
		}
		s.CacheService.Invalidate(p.GetId().Id)
		return handler(ctx, req)
	}
	pid, ok := req.(*pb.PersonId)
	if !ok {
		return nil, errors.New("request is not PersonId")
	}
	p, err := s.CacheService.Get(pid.GetId())
	if err == cc.ErrDataNotInCache {
		resp, err := handler(ctx, pid)
		if err == db.ErrPersonNotFound {
			return nil, err
		}
		s.CacheService.Set(pid.GetId(), resp)
		return resp, nil
	}
	fmt.Println("cache hit!", p)
	return p, nil
}

func StartServer() {
	cs := cc.NewCacheService()
	db := db.NewLocalPersonData()


	ws := WriteServer{CacheService: cs, DataService: db}
	rs := ReadServer{CacheService: cs, DataService: db}
	flag.Parse()
	lis, err := net.Listen("tcp", ":30301")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// rawServer := grpc.NewServer()
	cacheInterceptor := grpc.NewServer(
		grpc.UnaryInterceptor(rs.CacheReadMiddleware),
	)

	pb.RegisterPersonReadServer(cacheInterceptor, &rs)
	pb.RegisterPersonWriteServer(cacheInterceptor, &ws)
	log.Printf("server listening at %v", lis.Addr())
	if err := cacheInterceptor.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}