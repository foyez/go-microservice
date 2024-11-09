package main

import (
	"context"
	"log"
	"net"

	"github.com/foyez/go-microservice/user-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

// CreateUser creates a new user and returns the user details
func (server *Server) CreateUser(ctx context.Context, req *pb.NewUserRequest) (*pb.UserResponse, error) {
	if _, exists := server.users[req.Id]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "user with ID %s already exists", req.Id)
	}

	// Create a new user and save to the in-memory map
	user := &pb.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}
	server.users[req.Id] = user
	log.Printf("User created: %v", user)

	rsp := &pb.UserResponse{
		User: user,
	}
	return rsp, nil
}

// GetUser retrieves an existing user by ID
func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	// users := map[string]*pb.User{
	// 	"1": {Id: "1", Name: "Alice", Email: "alice@example.com"},
	// }

	user, exists := server.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	rsp := &pb.UserResponse{
		User: user,
	}
	return rsp, nil
}

func main() {
	server := &Server{
		users: make(map[string]*pb.User),
	}
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
