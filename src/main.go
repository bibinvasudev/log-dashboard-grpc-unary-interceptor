// websockets.go
package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "time"

    pb "google.golang.org/grpc/examples/helloworld/helloworld"

    "github.com/gorilla/websocket"

    "golang.org/x/net/context"
    "google.golang.org/grpc"

)

const (
	address     = "localhost:50051"
	defaultName = "world"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func grpcmain() string {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
        return string(r.GetMessage())
}


func main() {
	// Create a simple file server
    fs := http.FileServer(http.Dir("../public"))
    http.Handle("/", fs)
    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

        for {
            // Read message from browser
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }

            // Print the message to the console
            fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
            grpcmessage := grpcmain()
            fmt.Printf("%%%%%%%%%%%%%")
            msg = []byte(grpcmessage)
            fmt.Printf(grpcmessage)
            // Write message back to browser
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    })

    http.ListenAndServe(":8080", nil)
}
