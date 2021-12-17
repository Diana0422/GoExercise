package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
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
	address  = "localhost:1234"
	service1 = "MasterServer.Grep"

	filename = "client/test.txt"
	regex    = "ciao"
)

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	var reply File

	//create a TCP connection to localhost on port 1234
	cli, err := rpc.DialHTTP(network, address)
	errorHandler(err)

	mArgs := prepareArguments(filename, regex)

	// request to grep file to the server
	err = cli.Call(service1, mArgs, &reply)
	errorHandler(err)

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
