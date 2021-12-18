package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type GrepRequest struct {
	ArgFile File
	Regex   string
}

type File struct {
	Name    string
	Content string
}

const (
	network  = "tcp"
	address  = "localhost"
	service1 = "MasterServer.Grep"

	filename = "client/test.txt"
	regex    = "ciao"
)

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	reply := new(File)
	var cli *rpc.Client
	var err error

	// check for open TCP ports
	for p := 50000; p <= 50005; p++ {
		port := strconv.Itoa(p)
		cli, err = rpc.Dial(network, net.JoinHostPort(address, port))
		if err != nil {
			log.Printf("Connection error: port %v is not active", p)
			continue
		}
		if cli != nil {
			//create a TCP connection to localhost
			net.JoinHostPort(address, port)
			log.Printf("Connected on port %v", p)
			log.Printf("client conn: %p", cli)
			break
		}
	}

	// call the service
	mArgs := prepareArguments(filename, regex)
	fmt.Println(mArgs)
	// request to grep file to the server
	log.Printf("service: %v", service1)
	log.Printf("args: %v", mArgs)
	log.Printf("reply: %p", &reply)
	log.Printf("client: %p", cli)
	cliCall := cli.Go(service1, mArgs, &reply, nil)
	repCall := <-cliCall.Done
	log.Printf("Done %v", repCall)

	log.Println(reply)
	cli.Close()
}

func prepareArguments(f string, r string) []byte {
	// retrieve file to grep TODO better: choose your file
	file := new(File)
	file.Name = filename
	file.Content = readFileContent(filename)
	log.Printf("File Content: %s", file.Content)

	grepRequest := new(GrepRequest)
	grepRequest.ArgFile = *file
	grepRequest.Regex = regex

	// Marshaling
	s, err := json.Marshal(&grepRequest)
	errorHandler(err)
	log.Printf("Marshaled Data: %s", s)

	return s
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
