package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const (
	network    = "tcp"
	address    = "localhost:1234"
	numWorkers = 5
)

type File struct {
	Name    string
	Content string
}

type MasterServer struct {
	//Workers    []Worker
	numWorkers int
	filepath   string
}

// Grep /*---------- REMOTE PROCEDURE - CLIENT SIDE ---------------------------------------*/
func (m *MasterServer) Grep(payload []byte, reply *File) error {
	log.Printf("Received: %v", payload)
	var file File

	// Unmarshalling
	err := json.Unmarshal(payload, &file)
	errorHandler(err)
	log.Printf("Unmarshal: Name: %s, Content: %s", file.Name, file.Content)

	// chunk the file using getChunks function
	var chunks []string
	chunks = getChunks(file.Name)

	log.Println(chunks) // just for now
	*reply = file       // just for now
	// TODO spawn workers (n_workers)
	// TODO call worker's remote procedure to handle the mapping
	return nil
}

// Grep /*---------- REMOTE PROCEDURE - WORKER SIDE ---------------------------------------*/

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	master := new(MasterServer)
	// Publish the receiver methods
	err := rpc.Register(master)
	errorHandler(err)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 1234
	listener, err := net.Listen(network, address)
	errorHandler(err)
	log.Printf("Serving RPC server on port %d", 1234)

	// serve the client TODO multiple clients
	err = http.Serve(listener, nil)
	errorHandler(err)

}

/*------------------ LOCAL FUNCTIONS -------------------------------------------------------*/
func getChunks(srcName string) []string {
	//retrieve number of lines per worker
	srcFile, err := os.Open(srcName)
	errorHandler(err)

	numLines := lineCounter(srcFile)
	linesPerWorker := numLines / numWorkers

	//create chunk buffer
	chunkNames := make([]string, numWorkers)
	scanner := bufio.NewScanner(srcFile)

	for i := 0; i < numWorkers; i++ {
		//create and open with append mode a new chunk
		name := "chunk" + string(i) + ".txt"
		file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		errorHandler(err)

		//write 'linesPerWorker' lines from src to chunk
		count := 0
		for scanner.Scan() {
			_, err = fmt.Fprintln(file, scanner)
			errorHandler(err)
			count++

			//chunk complete: append
			if count == linesPerWorker {
				file.Close()
				chunkNames = append(chunkNames, name)
				break
			}

			//scanner error
			err := scanner.Err()
			errorHandler(err)
		}
	}

	return chunkNames
}

//count number of lines in a file
func lineCounter(file *os.File) int {
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	err := scanner.Err()
	errorHandler(err)

	return count
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
