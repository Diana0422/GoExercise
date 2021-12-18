package main

import (
	"container/list"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

const (
	network      = "tcp"
	addressLocal = "localhost:5678"
)

type GrepArgs struct {
	Chunk File
	Regex string
}

type File struct {
	Name    string
	Content string
}

type MapResp struct {
	Key   string
	Value bool
}

type Worker int

// Grep /*---------- REMOTE PROCEDURE - MASTER SIDE ---------------------------------------*/
func (w *Worker) Grep(payload []byte, mapRes *list.List) error {

	log.Printf("Received: %v", string(payload))
	var inArgs GrepArgs

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err)

	log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s",
		inArgs.Chunk.Name, inArgs.Chunk.Content, inArgs.Regex)

	mapRes = mapGrep(inArgs.Chunk, inArgs.Regex)

	return nil
}

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	worker := new(Worker)
	// Publish the receiver methods
	err := rpc.Register(worker)
	errorHandler(err)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 5678
	listener, err := net.Listen(network, addressLocal)
	errorHandler(err)
	log.Printf("Serving RPC server on port %d", 5678)

	err = http.Serve(listener, nil)
	errorHandler(err)
}

// MAP -> input (key=chunk, val=regex) => output [(key=str, val=regexIsIn)]
func mapGrep(chunk File, regex string) *list.List {

	res := new(list.List)
	lines := strings.Split(chunk.Content, "\n")

	for _, line := range lines {
		if strings.Contains(regex, line) {
			res.PushBack(MapResp{line, true})
		}
	}

	return res
}

//error handling
func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}
