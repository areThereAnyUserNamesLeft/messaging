package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)

	}
	go broadcast() // Sends messages to all users
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string // outbound message channel
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // All messages from clients
)

func broadcast() {
	clients := make(map[client]bool) // all incoming clients
	for {
		select {
		case msg := <-messages:
			// Send incoming messages to all clients outgoing message channel
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
			// will need to build an associative array here to store the clients linked to their UID (Uint64)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String() // return here change IP to UiD
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// Needs error handling for input.Err()
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // ignore network errors
	}
}
