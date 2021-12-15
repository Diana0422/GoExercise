package master

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "DistGrep/server/utils"
)

const (
	port = ":50051" //for client side -> server behaviour
	address = "localhost:50052" //for worker side -> client behaviour
	numWorkers = 5
)

type masterServer struct {
	Workers      []Worker
	numWorkers   int
	filepath string
	rpc.UnimplementedMasterServer
}



//gRPC - client side - function


//gRPC - worker side - function

func (ms *masterServer) Grep(ctx context.Context, in *pb.FileChunk) (*pb.GrepRow, error) {
	chunks := getChunks(ms.filepath)

	//open chunks
	for i := 0; i < ms.numWorkers; i++ {
		c, err := os.Open(chunks[i])
		errorHandler(err)
	}
}

//local functions
func getChunks(srcName string) []string {
	//retrieve number of lines per worker
	srcFile, err := os.Open(srcName)
	errorHandler(err)

	numLines := lineCounter(srcFile)
	linesPerWorker := numLines / numWorkers

	//create chunk buffer
	chunkNames := make([]string, numWorkers)
	scanner := bufio.NewScanner(srcFile)

	for i := 0; i < numWorkers; i++ {
		//create and open with append mode a new chunk
		name := "chunk" + string(i) + ".txt"
		file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		errorHandler(err)

		//write 'linesPerWorker' lines from src to chunk
		count := 0
		for scanner.Scan() {
			_, err = fmt.Fprintln(file, scanner)
			errorHandler(err)
			count++

			//chunk complete: append
			if count == linesPerWorker{
				file.Close()
				chunkNames = append(chunkNames, name)
				break
			}

			//scanner error
			err:= scanner.Err()
			errorHandler(err)
		}
	}

	return chunkNames
}

//count number of lines in a file
func lineCounter(file *os.File) int {
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	err:= scanner.Err()
	errorHandler(err)

	return count
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
	m := grpc.NewServer()
	pb.RegisterGrepMapReduceServer(m, &masterServer{})
	log.Printf("server listening at %v", lis.Addr())

	err = m.Serve(lis)
	errorHandler(err)
}