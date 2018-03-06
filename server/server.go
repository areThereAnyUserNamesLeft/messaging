package main

import (
	"bufio"
	"fmt"
	"log"
	"messaging/message"
	"net"
	"strconv"
	"strings"
)

var clientList map[string]string = make(map[string]string)

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
				log.Println(msg)
			}
		case cli := <-entering:
			clients[cli] = true
			log.Println(cli)
			// will need to build an associative array here to store the clients linked to their UID (Uint64)

		case cli := <-leaving:
			delete(clients, cli)
			log.Println(cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String() // return here change IP to UiD
	fmt.Println(conn.RemoteAddr().String())
	fmt.Println(conn.LocalAddr().String())
	clientList[who] = strconv.Itoa(int(message.Uint64())) // No point saving as a uint64 if I only need it as a str
	fmt.Println(clientList[who])
	ch <- "You are " + who
	log.Println(clientList[who] + " connected ")
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		log.Println(clientList[who] + " :sent " + input.Text())
		if strings.Index(input.Text(), "Relay") != -1 {
			messages <- message.RelayMess(input.Text(), clientList[who])
		}
		if strings.Index(input.Text(), "List") != -1 {
			v := make([]string, len(clientList))
			idx := 0
			for _, value := range clientList {
				v[idx] = value
				idx++
			}

			ch <- message.ListMess(strings.Join(v, "\n"), clientList[who])
		}
		if strings.Index(input.Text(), "Identity") != -1 {
			ch <- message.ListMess(clientList[who], clientList[who])
		}
	}
	// Needs error handling for input.Err()
	leaving <- ch
	log.Println(clientList[who] + " disconnected")
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {

	for msg := range ch {
		fmt.Fprintln(conn, msg) // ignore network errors
	}

}
