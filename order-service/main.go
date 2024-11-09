package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/foyez/go-microservice/user-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Order struct {
	ID     string `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
	User   User   `json:"user"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUserDetails(userID string) (*User, error) {
	// Set up connection options with insecure credentials (for development/testing purposes)
	conn, err := grpc.NewClient("localhost:4000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create a new UserService client from the connection
	client := pb.NewUserServiceClient(conn)

	// Set up a context with a timeout for the gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser method on the client
	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
	log.Println(userID)
	log.Println(res)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	user := res.User

	return &User{ID: user.Id, Name: user.Name, Email: user.Email}, nil
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	// userID := strings.TrimPrefix(r.URL.Path, "/orders/")

	userID := r.URL.Query().Get("userId")

	if userID == "" {
		http.Error(w, "userId parameter is required", http.StatusBadRequest)
		return
	}

	user, err := getUserDetails(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	order := Order{
		ID:     "101",
		Item:   "Laptop",
		Amount: 1200,
		User:   *user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	http.HandleFunc("/order", getOrderHandler)
	log.Println("Order Service is running on port 6000...")
	log.Fatal(http.ListenAndServe(":6000", nil))
}
