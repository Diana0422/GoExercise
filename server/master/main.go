package main

import (
	"GoExercise/server/worker"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

const (
	network      = "tcp"
	addressLocal = "localhost:1234"
	address      = "localhost:5678"
	service2     = "Worker.Grep"
	numWorkers   = 5
)

var port string

type GrepRequest struct {
	ArgFile File
	Regex   string
}

type File struct {
	Name    string
	Content []byte
}

type MasterServer struct {
	filepath string
	regex    string
}

type MasterClient struct {
	Workers     []worker.Worker
	currWorkers int
}

// Grep /*---------- REMOTE PROCEDURE - CLIENT SIDE ---------------------------------------*/
func (m *MasterServer) Grep(payload []byte, reply *File) error {
	log.Printf("Received: %s", payload)
	var inArgs GrepRequest

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err)
	log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s", inArgs.ArgFile.Name, inArgs.ArgFile.Content, inArgs.Regex)

	//master := new(MasterClient)
	//master.Grep(inArgs.File, inArgs.Regex)

	return nil
}

// Grep /*---------- REMOTE PROCEDURE - WORKER SIDE ---------------------------------------*/
/*func (mc *MasterClient) Grep(srcFile File, regex string) (*pb.GrepRow, error) {
	// chunk the file using getChunks function
	var chunks []File
	chunks = getChunks(srcFile)
	log.Println(chunks) // just for now

	//prepare results
	grepResp := [numWorkers]worker.GrepResp{}

	//SEND CHUNKS TO WORKERS
	for i,chunk := range chunks {
		//create a TCP connection to localhost on port 5678
		cli, err := rpc.DialHTTP(network, address)
		errorHandler(err)

		mArgs := prepareArguments(chunk, regex)

		//spawn worker connections
		worker := new(worker.Worker)
		grep := cli.Go(service2, mArgs, worker, nil)

		//wait for response
		grepResp[i] := <-grep.Done
	}
}*/

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	// Generate a random port for the client
	rand.Seed(time.Now().UTC().UnixNano())
	max := 50005
	min := 50000

	portNum := rand.Intn(max-min) + min
	port = strconv.Itoa(portNum)

	go serveClients()

	master := new(MasterServer)
	// Publish the receiver methods
	err := rpc.Register(master)
	errorHandler(err)

	for {
	}
}

func serveClients() {
	addr, err := net.ResolveTCPAddr(network, "0.0.0.0:"+port)
	errorHandler(err)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 1234
	listen, err := net.ListenTCP(network, addr)
	errorHandler(err)
	log.Printf("Serving RPC server on address %s , port %s", addr, port)

	for {
		// serve the new client
		rpc.Accept(listen)
		log.Printf("Serving the client.")
		errorHandler(err)
	}
}

/*------------------ LOCAL FUNCTIONS -------------------------------------------------------*/
func getChunks(srcFile File) []File {
	//retrieve number of lines per worker
	numLines := lineCounter(srcFile)
	linesPerWorker := numLines / numWorkers

	//create chunk buffer
	chunks := make([]File, numWorkers)
	scanner := bufio.NewScanner(bytes.NewReader(srcFile.Content))

	for i := 0; i < numWorkers; i++ {
		//create and open with append mode a new chunk
		chunk := new(File)
		chunk.Name = "chunk" + string(i) + ".txt"

		//write 'linesPerWorker' lines from src to chunk
		count := 0
		for scanner.Scan() {
			chunk.Content = append(chunk.Content, scanner.Bytes()...)
			chunk.Content = append(chunk.Content, "\n"...)
			count++

			//chunk complete: append
			if count == linesPerWorker {
				chunks = append(chunks, *chunk)
				break
			}

			//scanner error
			err := scanner.Err()
			errorHandler(err)
		}
	}

	return chunks
}

//count number of lines in a file
func lineCounter(file File) int {
	scanner := bufio.NewScanner(bytes.NewReader(file.Content))

	count := 0
	for scanner.Scan() {
		count++
	}

	err := scanner.Err()
	errorHandler(err)

	return count
}

func prepareArguments(chunk File, regex string) interface{} {
	// Arguments
	grepArgs := new(GrepRequest)
	grepArgs.Regex = regex
	grepArgs.ArgFile = chunk

	// Marshaling
	mArgs, err := json.Marshal(&grepArgs)
	errorHandler(err)
	log.Printf("Marshaled Data: %s", mArgs)

	return mArgs
}

func readFileContent(filename string) []byte {
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
	return content
}

//error handling
func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}

// Get a file name
func getFileName(file File) string {
	return file.Name
}
