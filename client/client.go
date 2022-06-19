package main

// Using netcat type tool as a basis for a client starting point
import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // Needs error handling
		log.Println("done")
		done <- struct{}{} // signal goroutine
	}()
	mustCopy(conn, os.Stdin)
	defer conn.Close()
	<-done // Waiting for background goroutine to stop
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
