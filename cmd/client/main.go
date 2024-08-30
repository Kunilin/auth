package main

import (
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	desc "gitlab.com/konfka/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
)

const address = "localhost:50052"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %e", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Failed to close connection: %e", err)
		}
	}()

	client := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Create(ctx, &desc.CreateRequest{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Role:     desc.Role_USER,
		Password: gofakeit.Password(true, true, true, false, false, 10),
	})
	if err != nil {
		log.Fatalf("Failed to create user: %e", err)
	}

	log.Printf("User info: %+v", r)
}
