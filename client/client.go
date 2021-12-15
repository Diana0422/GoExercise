package client

import (
	pb "GoExercise/go_exercise"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
)

var (
	serverAddr = flag.String("server_addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	//dial with the server
	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		log.Fatalf("failed to dial: #{err}")
	}
	defer conn.Close()
	cli := pb.NewGrepMapReduceClient(conn)

	// retrieve file to grep TODO better
	filename := "test.txt"
	content := readFileContent(filename)
	// request to grep file to the server
	grep, err := cli.Grep(context.Background(), &pb.File{FileName: filename, FileContent: content})
	if err != nil {
		log.Fatalf("failed to grep file: #{err}")
	}
	log.Println(grep)
}

func readFileContent(filename string) string {
	//open file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Stat())

	//read file content
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("failed to read file: #{err}")
	}
	fmt.Println(string(content))
	return string(content)
}