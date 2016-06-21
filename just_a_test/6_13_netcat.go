package main

import (
)
import (
	"net"
	"log"
	"io"
	"os"
)

var host = "localhost"
var port = "8080"

func main() {
	conn, err := net.Dial("tcp", host + ":" + port)
	if err != nil {
		log.Fatalf("connect to %s failed", conn.RemoteAddr())
	}

	log.Printf("connect to %s", conn.RemoteAddr())

	defer conn.Close()

	var done = make(chan struct {})
	go func() {
		io.Copy(os.Stdout, conn)
		log.Print("done")
		done <- struct {}{}

	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}


func mustCopy(w io.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		log.Fatal(err)
	}

}

