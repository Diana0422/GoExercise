package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	debug         = false
	network       = "tcp"
	address       = "localhost:5678"
	mapService    = "Worker.Map"
	reduceService = "Worker.Reduce"
	maxLoad       = 10 //every worker operates on a maximum of 'maxLoad' lines
)

var port string

type MapResp struct {
	Key   string
	Value string
}

type MapRequest struct {
	File  File
	Regex []string
}

type ReduceArgs struct {
	Key    string
	Values []string
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
func (m *MasterServer) Grep(payload []byte, reply *[]byte) error {
	//log.Printf("Received: %v", string(payload))
	var inArgs MapRequest

	// Unmarshalling
	err := json.Unmarshal(payload, &inArgs)
	errorHandler(err, 57)

	if debug {
		log.Printf("Unmarshal: Name: %s, Content: %s, Regex: %s",
			inArgs.File.Name, inArgs.File.Content, inArgs.Regex)
	}

	master := new(MasterClient)
	result, err := master.Grep(inArgs.File, inArgs.Regex)
	errorHandler(err, 64)

	// Marshalling
	s, err := json.Marshal(&result)
	if debug {
		log.Printf("Marshaled Data: %s", s)
	}

	*reply = s
	return err
}

// Grep /*---------- REMOTE PROCEDURE - WORKER SIDE ---------------------------------------*/
func (mc *MasterClient) Grep(srcFile File, regex []string) (*File, error) {
	//MAP PHASE
	log.Println("-->Activate Map Service on workers...")
	mapResp := mapFunction(mc, srcFile, regex)
	log.Print("...Done: All the workers returned from map -->\n\n")

	//SHUFFLE AND SORT PHASE
	log.Println("-->Do Shuffle and sort...")
	mapOutput, err := mergeMapResults(mapResp, mc.numWorkers)
	reduceInput := shuffleAndSort(mapOutput)
	log.Print("...Done: Shuffle and sort -->\n\n")

	//REDUCE PHASE
	log.Println("-->Activate Reduce Service on workers...")
	redResp := reduceFunction(mc, reduceInput)
	log.Print("...Done: All the workers returned from reduce -->\n\n")

	reply, err := mergeFinalResults(redResp, mc.numWorkers)
	if debug {
		log.Printf("Map Data: %s", mapResp)
		log.Printf("Reduced Data: %s", redResp)
		log.Printf("Reply: %s", reply)
	}

	return reply, err
}

func reduceFunction(mc *MasterClient, redIn []ReduceArgs) [][]byte {

	mc.numWorkers = len(redIn)

	//prepare results
	grepChan := make([]*rpc.Call, mc.numWorkers)
	grepResp := make([][]byte, mc.numWorkers)

	//SEND CHUNKS TO REDUCERS
	for i, chunk := range redIn {
		//create a TCP connection to localhost on port 5678
		cli, err := rpc.DialHTTP(network, address)
		errorHandler(err, 110)

		// Marshaling
		rArgs, err := json.Marshal(&chunk)
		errorHandler(err, 114)
		if debug {
			log.Printf("Marshalled Data: %s", rArgs)
		}

		//spawn worker connections
		grepChan[i] = cli.Go(reduceService, rArgs, &grepResp[i], nil)

		log.Printf("Spawned worker connection #%d", i)
	}

	//wait for response
	for i := 0; i < mc.numWorkers; i++ {
		<-grepChan[i].Done
		log.Printf("Worker #%d DONE", i)
	}

	return grepResp
}

func mapFunction(mc *MasterClient, srcFile File, regex []string) [][]byte {
	// chunk the file using getChunks function
	var chunks []File
	chunks = getChunks(srcFile, mc)

	//prepare results
	grepChan := make([]*rpc.Call, mc.numWorkers)
	grepResp := make([][]byte, mc.numWorkers)

	//SEND CHUNKS TO MAPPERS
	for i, chunk := range chunks {
		//create a TCP connection to localhost on port 5678
		cli, err := rpc.DialHTTP(network, address)
		errorHandler(err, 146)

		mArgs := prepareMapArguments(chunk, regex)

		//spawn worker connections
		grepChan[i] = cli.Go(mapService, mArgs, &grepResp[i], nil)

		log.Printf("Spawned worker connection #%d", i)
	}

	//wait for response
	for i := 0; i < mc.numWorkers; i++ {
		<-grepChan[i].Done
		log.Printf("Done: Worker #%d\n", i)
	}

	return grepResp
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
	errorHandler(err, 180)

	for {
	}
}

func serveClients() {
	addr, err := net.ResolveTCPAddr(network, "0.0.0.0:"+port)
	errorHandler(err, 132)

	// Register a HTTP handler
	rpc.HandleHTTP()
	//Listen to TCP connections on port 1234
	listen, err := net.ListenTCP(network, addr)
	errorHandler(err, 194)
	log.Printf("Serving RPC server on address %s , port %s\n", addr, port)

	for {
		// serve the new client
		rpc.Accept(listen)
	}
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
		if i == mc.numWorkers-1 && i != 0 && numLines%maxLoad != 0 {
			linesPerWorker = numLines % maxLoad
		}

		for j := 0; j < linesPerWorker; j++ {
			chunk.Content += lines[currLine] + "\n"
			currLine++
		}
		chunks[i] = *chunk
	}

	return chunks
}

func prepareMapArguments(chunk File, regex []string) interface{} {
	// Arguments
	grepArgs := new(MapRequest)
	grepArgs.Regex = regex
	grepArgs.File = chunk

	// Marshaling
	mArgs, err := json.Marshal(&grepArgs)
	errorHandler(err, 259)
	if debug {
		log.Printf("Marshaled Data: %s", mArgs)
	}

	return mArgs
}

func shuffleAndSort(mapRes []MapResp) []ReduceArgs {
	sort.Slice(mapRes, func(i, j int) bool {
		return mapRes[i].Key < mapRes[j].Key
	})

	var result []ReduceArgs
	var currKey string

	var r *ReduceArgs
	for i, m := range mapRes {
		if currKey != m.Key {
			if i != 0 {
				result = append(result, *r)
			}
			currKey = m.Key
			r = new(ReduceArgs)
			r.Key = currKey
			r.Values = append(r.Values, m.Value)
		} else {
			r.Values = append(r.Values, m.Value)
		}
	}

	result = append(result, *r)
	return result
}

func mergeMapResults(resp [][]byte, dim int) ([]MapResp, error) {

	var mapRes []MapResp

	for i := 0; i < dim; i++ {
		// Unmarshalling
		var temp []MapResp
		err := json.Unmarshal(resp[i], &temp)
		errorHandler(err, 301)
		if debug {
			log.Printf("Received: %s", resp[i])
			log.Printf("Unmarshal: Key: %v", temp)
		}
		mapRes = append(mapRes, temp...)
	}

	return mapRes, nil
}

func mergeFinalResults(resp [][]byte, dim int) (*File, error) {
	file := new(File)
	file.Name = "result.txt"

	for i := 0; i < dim; i++ {
		// Unmarshalling
		var outArgs ReduceArgs

		err := json.Unmarshal(resp[i], &outArgs)
		errorHandler(err, 320)
		if debug {
			log.Printf("Received: %s", resp[i])
			log.Printf("Unmarshal: Key: %v", outArgs)
		}
		for k := 0; k < len(outArgs.Values); k++ {
			file.Content += outArgs.Values[k] + "\n"
		}
	}
	return file, nil
}

//error handling
func errorHandler(err error, line int) {
	if err != nil {
		log.Fatalf("failure at line %d: %v", line, err)
	}
}
