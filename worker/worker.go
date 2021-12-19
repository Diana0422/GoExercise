package main

import (
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
	File  File
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
func (w *Worker) Grep(payload []byte, result *[]byte) error {

	//log.Printf("Received: %v", string(payload))
	var inArgs GrepArgs

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err, 42)

	//log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s", inArgs.File.Name, inArgs.File.Content, inArgs.Regex)

	mapRes := mapGrep(inArgs.File, inArgs.Regex)
	//log.Printf("MapRes: %v", mapRes)

	// Marshaling
	s, err := json.Marshal(&mapRes)
	errorHandler(err, 50)
	log.Printf("Marshaled Data: %s", s)

	*result = s

	return nil
}

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	worker := new(Worker)
	// Publish the receiver methods
	err := rpc.Register(worker)
	errorHandler(err, 57)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 5678
	listener, err := net.Listen(network, addressLocal)
	errorHandler(err, 63)
	log.Printf("Serving RPC server on port %d", 5678)

	err = http.Serve(listener, nil)
	errorHandler(err, 67)
}

// MAP -> input (key=chunk, val=regex) => output [(key=str, val=regexIsIn)]
func mapGrep(chunk File, regex string) []MapResp {

	res := make([]MapResp, 0)

	lines := strings.Split(chunk.Content, "\n")
	for _, line := range lines {
		if strings.Contains(line, regex) {
			res = append(res, MapResp{line, true})
		}
	}

	//log.Println(res)
	return res
}

//error handling
func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}
