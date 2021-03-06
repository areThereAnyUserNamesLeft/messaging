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
		if firstWord(input.Text()) == "relay:" {
			if len(input.Text()) > int(1024*1000) {
				ch <- "That message is too long to send, why not send two?"
			} else {
				messages <- message.RelayMess(input.Text(), clientList[who])
			}
		} else if firstWord(input.Text()) == "list:" {

			v := make([]string, len(clientList))
			idx := 0
			for _, value := range clientList {
				v[idx] = value
				if idx > 254 {
					break
				}
				idx++
			}

			ch <- message.ListMess(strings.Join(v, "',  '"), clientList[who])

		} else if firstWord(input.Text()) == "identity:" {
			ch <- message.IdentityMess(clientList[who], clientList[who])
		} else {
			advice := " :no messaging protocol used - Please make sure you prefix your message with 'List:', 'Identify': or 'Relay:' then <your message> \n"
			ch <- clientList[who] + advice + input.Text()
			ch <- "Don't forget your colon ':' (always good advice!) :)"

		}
	}
	// Needs error handling for input.Err()
	leaving <- ch
	log.Println(clientList[who] + " disconnected")
	conn.Close()
}

func firstWord(str string) string {
	count := 1
	for i := range str {
		if str[i] == ' ' {
			count -= 1
			if count == 0 {
				return strings.ToLower(str[0:i])
			}
		}
	}
	return strings.ToLower(str)
}

func clientWriter(conn net.Conn, ch <-chan string) {

	for msg := range ch {
		fmt.Fprintln(conn, msg) // ignore network errors
	}

}
