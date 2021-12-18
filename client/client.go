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

type GrepArgs struct {
	File  File
	Regex string
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
	var reply File
	var cli *rpc.Client

	// check for open TCP ports
	for p := 50000; p <= 50005; p++ {
		port := strconv.Itoa(p)
		cli, err := rpc.Dial(network, net.JoinHostPort(address, port))
		if err != nil {
			log.Printf("Connection error: port %v is not active", p)
			continue
		}
		if cli != nil {
			//create a TCP connection to localhost
			defer cli.Close()
			net.JoinHostPort(address, port)
			log.Printf("Connected on port %v", p)
			break
		}
	}

	// call the service
	mArgs := prepareArguments(filename, regex)
	fmt.Println(mArgs)
	// request to grep file to the server
	cliCall := cli.Go(service1, mArgs, &reply, nil)
	repCall := <-cliCall.Done
	if repCall != nil {
		log.Println("Done")
	}

	log.Println(reply)
}

func prepareArguments(f string, r string) interface{} {
	// retrieve file to grep TODO better: choose your file
	file := new(File)
	file.Name = filename
	file.Content = readFileContent(filename)
	log.Printf("File Content: %s", file.Content)

	grepArgs := new(GrepArgs)
	grepArgs.File = *file
	grepArgs.Regex = regex

	// Marshaling
	s, err := json.Marshal(&grepArgs)
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
