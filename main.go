package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// we define a server structure were like in every websocket is a map of conections
// where is true for connected or false with disconedted

type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

// we create the server pointing to the server struct to manage memory and we define
// the connections is gonna be like we define before

func CreateServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// we create the server and if is connected put from there is connected
// if is disconected the connection is deleted and is gonna over again
// and read every connection on the readLoop function

func (s *Server) handleConnection(conn *websocket.Conn) {
	fmt.Println("new client connected", conn.RemoteAddr())

	s.mu.Lock()
	s.conns[conn] = true
	s.mu.Unlock()

	s.readLoop(conn)

	s.mu.Lock()
	delete(s.conns, conn)
	s.mu.Unlock()

	fmt.Println("client disconnected", conn.RemoteAddr())
}

// we create a loop to read every message we define the text on bytes and also if
// we have an error reporeted and then we print what message was sended
// and every message is sended to the broadcast

func (s *Server) readLoop(conn *websocket.Conn) {
	defer conn.Close() // Close the connection when this function exits
	text := make([]byte, 1024)
	for {
		n, err := conn.Read(text)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from connection:", err)
			break // Exit the loop on error
		}

		message := text[:n]
		fmt.Println("received message:", string(message))
		s.broadcast(message)
	}
}

// in the broadcast functon we create a for loop in the range of connections
// where is gonna with the propery ws of websocket write in the server what message was send it
// and then broadcast to the server

func (s *Server) broadcast(msg []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			_, err := ws.Write(msg)
			if err != nil {
				fmt.Println("error broadcasting the msg:", err)
			}
		}(ws)
	}
}

// and the in the main we create the server and they respec link and the port of the server
func main() {
	server1 := CreateServer()
	server2 := CreateServer()
	server3 := CreateServer()
	http.Handle("/ws/server1", websocket.Handler(server1.handleConnection))
	http.Handle("/ws/server2", websocket.Handler(server2.handleConnection))
	http.Handle("/ws/server3", websocket.Handler(server3.handleConnection))
	http.ListenAndServe(":8888", nil)
}
