package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type GrepRequest struct {
	File  File
	Regex []string
}

type File struct {
	Name    string
	Content string
}

const (
	debug    = false // Set to true to activate debug log
	network  = "tcp"
	address  = "localhost"
	service1 = "MasterServer.Grep"
)

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	var reply []byte
	var cli *rpc.Client
	var err error
	welcome()

	// check for open TCP ports
	for p := 50000; p <= 50005; p++ {
		port := strconv.Itoa(p)
		cli, err = rpc.Dial(network, net.JoinHostPort(address, port))
		if err != nil {
			if debug {
				log.Printf("Connection error: port %v is not active", p)
			}
			log.Printf("Connecting to master...")
			continue
		}
		if cli != nil {
			//create a TCP connection to localhost
			net.JoinHostPort(address, port)
			log.Printf("Connected on port %v", p)

			if debug {
				log.Printf("client conn: %p", cli)
			}
			break
		}
	}

	// call the service
	mArgs := prepareArguments()

	// request to grep file to the server
	if debug {
		log.Printf("service: %v", service1)
		log.Printf("args: %v", string(mArgs))
		log.Printf("reply: %p", &reply)
		log.Printf("client: %p", cli)
	}
	cliCall := cli.Go(service1, mArgs, &reply, nil)
	repCall := <-cliCall.Done
	if debug {
		log.Printf("Done %v", repCall)
	}

	// Unmarshalling of reply
	var result File
	err = json.Unmarshal(reply, &result)
	errorHandler(err, 77)

	// Print grep result on screen
	fmt.Println("")
	fmt.Println("-------------------------- GREP RESULT ------------------------------: ")
	fmt.Println(result.Content)
	err = cli.Close()
	errorHandler(err, 84)
}

func prepareArguments() []byte {
	// retrieve file to grep
	file := new(File)
	file.Name = fileToGrep()
	regex := getRegex()
	file.Content = readFileContent("client/files/" + file.Name)
	if debug {
		log.Printf("File Content: %s", file.Content)
	}

	grepRequest := new(GrepRequest)
	grepRequest.File = *file
	grepRequest.Regex = regex

	// Marshaling
	s, err := json.Marshal(&grepRequest)
	errorHandler(err, 102)
	if debug {
		log.Printf("Marshaled Data: %s", s)
	}

	return s
}

/*------------------ OTHER FUNCTIONS -------------------------------------------------------*/
func fileToGrep() string {
	var fileNum int
	fileMap := make(map[int]string)
	fmt.Println("\n----------------------- WELCOME in GoGREP! --------------------------:")
	fmt.Println("Choose a file to grep:")

	// Read files directory
	file, err := ioutil.ReadDir("client/files")
	errorHandler(err, 119)

	for i := 0; i < len(file); i++ {
		fmt.Printf("-> (%d) %s\n", i+1, file[i].Name())
		fileMap[i+1] = file[i].Name()
	}

	// Input the file chosen
	fmt.Print("Select a number: ")
	_, err = fmt.Scanf("%d\n", &fileNum)
	errorHandler(err, 129)
	return fileMap[fileNum]
}

func getRegex() []string {
	var regex []string

	// Input the regex
	fmt.Print("Insert any regex you want to look for (format: regex1[ regex2...]): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		regex = strings.Split(scanner.Text(), " ")
	}

	return regex
}

func readFileContent(filename string) string {
	//read file content
	if debug {
		log.Printf("Reading file: %s", filename)
	}
	content, err := ioutil.ReadFile(filename)
	errorHandler(err, 158)
	return string(content)
}

func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}

func welcome() {
	// Welcome
	for i := 0; i <= 3; i++ {
		fmt.Println("*")
	}
	fmt.Println("Authors: Diana Pasquali, Livia Simoncini")
	for i := 0; i <= 3; i++ {
		fmt.Println("*")
	}
}
