package main

import (
	"context"
	"log"
	"net"
	"sync"

	//import the generated protobuf code
	pb "github.com/viseoandrew/shippy-service-consignment/proto/consignment"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignement) (*pb.Consignment, error)
}

//Repository Simulating Datastore
type Repository struct {
	mu            sync.RWMutex
	consignements []*pb.Consignement
}

//Create a new Consignment
func (repo *Repository) Create(consignment *pb.Consignement) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignements = updated
	repo.mu.Unlock()
	return consignment, nil
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Consignment, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository

	//setup grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterShippingServiceServer(s, &service{repo})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Running on Port: ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
