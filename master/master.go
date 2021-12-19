package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"time"
)

const (
	network  = "tcp"
	address  = "localhost:5678"
	service2 = "Worker.Grep"
	maxLoad  = 10 //every worker operates on a maximum of 'maxLoad' lines
)

var port string

type GrepResp struct {
	Key   string
	Value bool
}

type GrepRequest struct {
	File  File
	Regex string
}

type File struct {
	Name    string
	Content string
}

type MasterServer int

type MasterClient struct {
	numWorkers int
}

// Grep /*---------- REMOTE PROCEDURE - CLIENT SIDE ---------------------------------------*/
func (m *MasterServer) Grep(payload []byte, reply *File) error {
	//log.Printf("Received: %v", string(payload))
	var inArgs GrepRequest

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err, 51)

	//log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s",
	//inArgs.File.Name, inArgs.File.Content, inArgs.Regex)

	master := new(MasterClient)
	reply, err = master.Grep(inArgs.File, inArgs.Regex)

	return err
}

// Grep /*---------- REMOTE PROCEDURE - WORKER SIDE ---------------------------------------*/
func (mc *MasterClient) Grep(srcFile File, regex string) (*File, error) {
	// chunk the file using getChunks function
	var chunks []File
	chunks = getChunks(srcFile, mc)
	//log.Println(chunks)

	//prepare results
	grepChan := make([]*rpc.Call, mc.numWorkers)
	grepResp := make([]byte, mc.numWorkers)
	//log.Printf("grepResp: %v", grepResp)

	//SEND CHUNKS TO WORKERS
	for i, chunk := range chunks {
		//create a TCP connection to localhost on port 5678
		cli, err := rpc.DialHTTP(network, address)
		errorHandler(err, 77)

		mArgs := prepareArguments(chunk, regex)

		//spawn worker connections
		grepChan[i] = cli.Go(service2, mArgs, grepResp[i], nil)
		log.Printf("Spawned worker connection #%d", i)
	}

	//wait for response
	for i := 0; i < mc.numWorkers; i++ {
		<-grepChan[i].Done
		log.Printf("Worker #%d DONE", i)
	}

	//merge results
	log.Println("Merging results...")
	log.Printf("grepResp: %v", grepResp)
	reply, err := mergeMapResults(grepResp)
	return reply, err
}

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
	errorHandler(err, 110)

	for {
	}
}

func serveClients() {
	addr, err := net.ResolveTCPAddr(network, "0.0.0.0:"+port)
	errorHandler(err, 117)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 1234
	listen, err := net.ListenTCP(network, addr)
	errorHandler(err, 123)
	log.Printf("Serving RPC server on address %s , port %s", addr, port)

	for {
		// serve the new client
		rpc.Accept(listen)
		log.Printf("Serving the client.")
	}
}

/*------------------ LOCAL FUNCTIONS -------------------------------------------------------*/
func getChunks(srcFile File, mc *MasterClient) []File {

	//Separate the content of the original file in lines
	lines := strings.Split(srcFile.Content, "\n")
	numLines := len(lines)
	//log.Printf("TOTAL LINES %d\n", numLines)

	//Distribute equal amount of lines per chunk
	var linesPerWorker int
	if numLines < maxLoad {
		mc.numWorkers = 1
		linesPerWorker = numLines
	} else {
		mc.numWorkers = numLines / maxLoad
		if numLines%maxLoad != 0 {
			mc.numWorkers++
		}
		linesPerWorker = maxLoad
	}

	//log.Printf("lines per worker %d\nnumber of workers %d\n", linesPerWorker, mc.numWorkers)

	//create and populate chunk buffer
	chunks := make([]File, mc.numWorkers)
	currLine := 0
	for i := 0; i < mc.numWorkers; i++ {
		//create a new chunk
		chunk := new(File)
		chunk.Name = "chunk" + strconv.Itoa(i) + ".txt"

		//write 'linesPerWorker' lines from src to chunk
		if i == mc.numWorkers-1 && i != 0 && numLines%maxLoad != 0 {
			linesPerWorker = numLines % maxLoad
		}

		for j := 0; j < linesPerWorker; j++ {
			chunk.Content += lines[currLine] + "\n"
			//log.Printf("Worker %d, Line %d\nContent: %s\n", i, j+1, chunk.Content)
			currLine++
		}
		chunks[i] = *chunk
	}

	return chunks
}

func prepareArguments(chunk File, regex string) interface{} {
	// Arguments
	grepArgs := new(GrepRequest)
	grepArgs.Regex = regex
	grepArgs.File = chunk

	// Marshaling
	mArgs, err := json.Marshal(&grepArgs)
	errorHandler(err, 185)
	log.Printf("Marshaled Data: %s", mArgs)

	return mArgs
}

func mergeMapResults(resp []byte) (*File, error) {
	file := new(File)
	file.Name = "result.txt"

	log.Printf("Received: %v", string(resp))
	var outArgs []GrepResp

	// Unmarshalling
	err := json.Unmarshal(resp, &outArgs)
	errorHandler(err, 202)

	log.Printf("Unmarshal: Key: %v", outArgs)

	for j := 0; j < len(outArgs); j++ {
		file.Content += outArgs[j].Key
	}

	return file, nil
}

//error handling
func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}
