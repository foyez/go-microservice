package main

import (
	"context"
	"log"
	"net"

	"github.com/foyez/microservice-with-go/user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
}

func (server *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	users := map[string]*pb.UserResponse{
		"1": {Id: "1", Name: "Alice", Email: "alice@example.com"},
		"2": {Id: "2", Name: "Bob", Email: "bob@example.com"},
	}

	user, exists := users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return user, nil
}

func main() {
	server := &Server{}
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
