package worker

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	pb "GoExercise/server/utils"
	"google.golang.org/grpc"
)

const (
	port = ":50052" //for master side -> server behaviour
)

type workerServer struct {
	id uint
	rpc.UnimplementedMasterServer
}

func (ws *workerServer) Grep(ctx context.Context, in *pb.FileChunk) (*pb.GrepRow, error) {
	//TODO
}

func findPattern(file *os.File, pattern string) string {
	fileReader := bufio.NewReader(file)

	for {
		line, err := fileReader.ReadString('\n')
		if err != nil && err == io.EOF {
			//SEND DONE TO MASTER
		}

		if strings.Contains(line, pattern) {
			//SEND STRING TO MASTER
		}
	}
}

//error handling
func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}

func main() {
	// Add server listener
	lis, err := net.Listen("tcp", port)
	errorHandler(err)

	// Start and register the server
	w := grpc.NewServer()
	pb.RegisterGrepMapReduceServer(w, &workerServer{})
	log.Printf("server listening at %v", lis.Addr())

	err = w.Serve(lis)
	errorHandler(err)
}
