package main

import (
	"crypto/rand"
	"log"
	"math"
	"math/big"
	"net"

	"fmt"

	"github.com/brianvoe/gofakeit"
	desc "gitlab.com/konfka/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"context"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	_ = ctx
	log.Printf("Creating user: %+v", req.Name)
	num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{Id: num.Int64()}, nil
}

func (s server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	_ = ctx
	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      desc.Role_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Updating user: %d", req.GetId())
	_ = ctx
	return &emptypb.Empty{}, nil
}

func (s server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_ = ctx
	log.Printf("Deleting user: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalln("Failed to listen server: ", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("Server listening at %+v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %e", err)
	}
}
