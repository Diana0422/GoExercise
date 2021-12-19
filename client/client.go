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
	mArgs := prepareArguments()
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

func prepareArguments() []byte {
	// retrieve file to grep TODO better: choose your file
	file := new(File)
	file.Name = fileToGrep()
	regex := getRegex()
	file.Content = readFileContent("client/files/" + file.Name)
	log.Printf("File Content: %s", file.Content)

	grepRequest := new(GrepRequest)
	grepRequest.File = *file
	grepRequest.Regex = regex

	// Marshaling
	s, err := json.Marshal(&grepRequest)
	errorHandler(err, 84)
	log.Printf("Marshaled Data: %s", s)

	return s
}

/*------------------ OTHER FUNCTIONS -------------------------------------------------------*/
func fileToGrep() string {
	var fileNum int
	fileMap := make(map[int]string)
	fmt.Println("\n\nChoose a file to grep:")

	// Read files directory
	file, err := ioutil.ReadDir("client/files")
	errorHandler(err, 94)

	for i := 0; i < len(file); i++ {
		fmt.Printf("-> (%d) %s\n", i+1, file[i].Name())
		fileMap[i+1] = file[i].Name()
	}

	// Input the file chosen
	fmt.Print("Select a number: ")
	fmt.Scanf("%d\n", &fileNum)
	return fileMap[fileNum]
}

func getRegex() string {
	var regex string
	// Input the regex
	fmt.Print("Select a regex: ")
	fmt.Scanf("%s\n", &regex)
	return regex
}

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
	errorHandler(err, 103)
	return string(content)
}

func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}
