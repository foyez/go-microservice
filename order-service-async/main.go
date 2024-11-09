package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

// getUserDetails returns user details by userID
func getUserDetails(userID string) (*User, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"user_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Set a timeout context for listening to messages
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("No matching user found within timeout.")
			return nil, errors.New("user not found")
		case d := <-msgs:
			var user User
			if err := json.Unmarshal(d.Body, &user); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Check if the user ID matches
			if user.ID == userID {
				return &user, nil
			}
		}
	}
}

// getOrderHandler retrieves a placed order with user details
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
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
