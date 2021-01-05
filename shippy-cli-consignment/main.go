package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/viseoandrew/shippy/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshall(data, &consignment)
	return consignment, err
}

func main() {

	//Setup a Connection to grpc Server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not Connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewShippingServiceClient(conn)

	//Contact Server and print response
	file := defaultFilename
	if len(os.Args) > 1 {
		file := Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not Create Consignment: %v", err)
	}
	log.Printf("Created: %t", r.Created)
}
