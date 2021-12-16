package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

const (
	network = "tcp"
	address = "localhost:1234"
)

type File struct {
	name    string
	content string
}

type MasterServer struct{}

// Grep /*---------- REMOTE PROCEDURE ---------------------------------------*/
func (m *MasterServer) Grep(payload File, reply *File) error {
	log.Printf("Received: %v", getFileName(payload))
	if getFileName(payload) != "" {
		*reply = File{name: getFileName(payload), content: "Test"}
	}
	return nil
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
	listener, err := net.Listen(network, address)
	errorHandler(err)
	log.Printf("Serving RPC server on port %d", 1234)

	// serve the client TODO multiple clients
	err = http.Serve(listener, nil)
	errorHandler(err)

}

func getFileName(file File) string {
	return file.name
}

func errorHandler(err error) {
	if err != nil {
		log.Fatalf("failure: %v", err)
	}
}
