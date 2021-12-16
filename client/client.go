package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

type File struct {
	name    string
	content string
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

	// retrieve file to grep TODO better
	filename := "test.txt"
	file := new(File)
	file.name = filename
	file.content = readFileContent(filename)

	// request to grep file to the server
	err = cli.Call(service1, file, &reply)
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
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("failed to read file: #{err}")
	}
	fmt.Println(string(content))
	return string(content)
}

func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}
