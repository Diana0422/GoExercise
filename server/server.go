package server

import (
	pb "../go_exercise"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type masterServer struct {
	pb.UnimplementedGrepMapReduceServer
}

func (m *masterServer) Grep(ctx context.Context, file *pb.File) (*pb.File, error) {
	log.Printf("Received: %v", file.GetFileName())
	return &pb.File{FileName: "test.txt", FileContent: "test_content"}, nil
}

func main() {
	// Add server listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: #{err}")
	}

	// Start and register the server
	m := grpc.NewServer()
	pb.RegisterGrepMapReduceServer(m, &masterServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := m.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
