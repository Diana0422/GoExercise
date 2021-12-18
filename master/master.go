package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

const (
	network      = "tcp"
	addressLocal = "localhost:1234"
	address      = "localhost:5678"
	service2     = "Worker.Grep"
	maxLoad      = 10 //every worker operates on a maximum of 'maxLoad' lines
)

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
	log.Printf("Received: %v", string(payload))
	var inArgs GrepRequest

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err)

	log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s",
		inArgs.File.Name, inArgs.File.Content, inArgs.Regex)

	master := new(MasterClient)
	reply, err = master.Grep(inArgs.File, inArgs.Regex)

	return err
}

// Grep /*---------- REMOTE PROCEDURE - WORKER SIDE ---------------------------------------*/
func (mc *MasterClient) Grep(srcFile File, regex string) (*File, error) {
	// chunk the file using getChunks function
	var chunks []File
	chunks = getChunks(srcFile, mc)
	log.Println(chunks) // just for now

	//prepare results
	grepChan := make([]*rpc.Call, mc.numWorkers)
	grepResp := make([]GrepResp, mc.numWorkers)

	//SEND CHUNKS TO WORKERS
	for i, chunk := range chunks {
		//create a TCP connection to localhost on port 5678
		cli, err := rpc.DialHTTP(network, address)
		errorHandler(err)

		mArgs := prepareArguments(chunk, regex)

		//spawn worker connections
		grepChan[i] = cli.Go(service2, mArgs, grepResp[i], nil)
	}

	//wait for response
	for i := 0; i < mc.numWorkers; i++ {
		<-grepChan[i].Done
	}

	//merge results
	reply, err := mergeMapResults(grepResp)
	return reply, err
}

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	master := new(MasterServer)
	// Publish the receiver methods
	err := rpc.Register(master)
	errorHandler(err)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 1234
	listener, err := net.Listen(network, addressLocal)
	errorHandler(err)
	log.Printf("Serving RPC server on port %d", 1234)

	// serve the client TODO multiple clients
	err = http.Serve(listener, nil)
	errorHandler(err)
}

/*------------------ LOCAL FUNCTIONS -------------------------------------------------------*/
func getChunks(srcFile File, mc *MasterClient) []File {

	//Separate the content of the original file in lines
	lines := strings.Split(srcFile.Content, "\n")
	numLines := len(lines)

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

	//create and populate chunk buffer
	chunks := make([]File, mc.numWorkers)
	currLine := 0
	for i := 0; i < mc.numWorkers; i++ {
		//create a new chunk
		chunk := new(File)
		chunk.Name = "chunk" + strconv.Itoa(i) + ".txt"

		//write 'linesPerWorker' lines from src to chunk
		if i == mc.numWorkers-1 && i != 0 {
			linesPerWorker = numLines % maxLoad
		}

		for j := 0; j < linesPerWorker; j++ {
			chunk.Content += lines[currLine]
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
	errorHandler(err)
	log.Printf("Marshaled Data: %s", mArgs)

	return mArgs
}

func mergeMapResults(resp []GrepResp) (*File, error) {
	file := new(File)
	file.Name = "result.txt"

	for _, result := range resp {
		file.Content += result.Key
	}

	return file, nil
}

//error handling
func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}
