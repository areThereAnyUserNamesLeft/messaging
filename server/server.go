package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"messaging/message"
	"net"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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
			for client := range clients {
				client <- msg
				log.Println(msg)
			}
		case client := <-entering:
			clients[client] = true
			log.Println(client)
			// will need to build an associative array here to store the clients linked to their UID (Uint64)

		case client := <-leaving:
			delete(clients, client)
			log.Println(client)
			defer close(client)
		}
	}
}

func handleConn(conn net.Conn) error {
	defer conn.Close()
	ch := make(chan string) // outgoing client messages
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return clientWriter(conn, ch)
	})
	who := conn.RemoteAddr().String() // return here change IP to UiD
	fmt.Println(conn.RemoteAddr().String())
	fmt.Println(conn.LocalAddr().String())
	clientList[who].uuid = uuid.NewString() 
	fmt.Println(clientList[who])
	ch <- "Your id is " + clientList[who].uuid
	ch <- "Would you like to set a name for this account so other users know who you are?", 
	log.Println(clientList[who] + " connected ")
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if firstWord(input.Text()) == "relay:" {
			log.Println(
				input.Text(),
			)
			if len(input.Text()) > int(1024*1000) {
				ch <- "That message is too long to send, why not send two?"
			} else {
				messages <- message.RelayMessage(input.Text(), clientList[who], nil)
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

			ch <- message.ListMessage(strings.Join(v, "',  '"), clientList[who])

		} else if firstWord(input.Text()) == "identity:" {
			ch <- message.IdentityMessage(clientList[who], clientList[who])
		} else {
			advice := ` : No messaging protocol used -

Please make sure you prefix your message with one of the following
	List: <your message>
	Identify: <your message>
	Relay: <your message>
`
			ch <- clientList[who] + advice + input.Text()
			ch <- "Don't forget your colon ':' (always good advice!) :)"

		}
	}
	// Needs error handling for input.Err()
	leaving <- ch
	log.Println(clientList[who] + " disconnected")
	return eg.Wait()
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

func clientWriter(conn net.Conn, ch <-chan string) error {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg) // ignore network errors
		if err != nil {
			return err
		}
	}
	return nil
}
