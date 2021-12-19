package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"regexp"
	"strings"
)

const (
	network      = "tcp"
	addressLocal = "localhost:5678"
)

type MapArgs struct {
	File  File
	Regex []string
}

type File struct {
	Name    string
	Content string
}

type MapResp struct {
	Key   string
	Value string
}

type Worker int

// Map /*---------- REMOTE PROCEDURE - MASTER SIDE ---------------------------------------*/
func (w *Worker) Map(payload []byte, result *[]byte) error {

	//log.Printf("Received: %v", string(payload))
	var inArgs MapArgs

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err, 41)
	//log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s", inArgs.File.Name, inArgs.File.Content, inArgs.Regex)

	//map
	var mapRes []MapResp
	chunk := inArgs.File
	regex := inArgs.Regex
	lines := strings.Split(chunk.Content, "\n")
	for _, line := range lines {
		for i := 0; i < len(regex); i++ {
			match, err := regexp.Match(regex[i], []byte(line))
			errorHandler(err, 86)
			if match {
				mapRes = append(mapRes, MapResp{regex[i], line})
			}
		}
	}
	//log.Printf("MapRes: %v", mapRes)

	// Marshalling
	s, err := json.Marshal(&mapRes)
	errorHandler(err, 50)
	log.Printf("Marshaled Data: %s", s)

	//return
	*result = s
	return nil
}

// Reduce -> identity /*---------- REMOTE PROCEDURE - MASTER SIDE --------------*/
func (w *Worker) Reduce(payload []byte, result *[]byte) error {
	*result = payload
	return nil
}

/*------------------ MAIN -------------------------------------------------------*/
func main() {
	worker := new(Worker)
	// Publish the receiver methods
	err := rpc.Register(worker)
	errorHandler(err, 63)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 5678
	listener, err := net.Listen(network, addressLocal)
	errorHandler(err, 69)
	log.Printf("Serving RPC server on port %d", 5678)

	err = http.Serve(listener, nil)
	errorHandler(err, 73)
}

//error handling
func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}
