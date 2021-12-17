package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

type File struct {
	Name    string
	Content string
}

const (
	network  = "tcp"
	address  = "localhost:1234"
	service1 = "MasterServer.Grep"
)

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	var reply File

	//create a TCP connection to localhost on port 1234
	cli, err := rpc.DialHTTP(network, address)
	errorHandler(err)

	// retrieve file to grep TODO better: choose your file
	filename := "client/test.txt"
	file := new(File)
	file.Name = filename
	file.Content = readFileContent(filename)
	log.Printf("File Content: %s", file.Content)

	// Marshaling
	s, err := json.Marshal(&file)
	errorHandler(err)
	log.Printf("Marshaled Data: %s", s)

	// request to grep file to the server
	err = cli.Call(service1, s, &reply)
	errorHandler(err)
	log.Println(reply)
}

/*------------------ OTHER FUNCTIONS -------------------------------------------------------*/
func readFileContent(filename string) string {
	//open file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Stat())

	//read file content
	log.Printf("Reading file: %s", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("failed to read file: #{err}")
	}
	return string(content)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}
