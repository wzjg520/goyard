package main

import (
	"net"
	"log"
	"time"
	"fmt"
	"strings"
	"bufio"
)

var host = "localhost"
var port = "8080"

func main() {

	l, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatalf("Listening %s:%s failed error:%s", host, port, err)
	}

	log.Printf("Listening %s:%s", host, port)
	for {
		conn, err := l.Accept()

		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)

	}
}

func handleConn(conn net.Conn) error {
	defer conn.Close()

	log.Printf("Accept from %s", conn.RemoteAddr())

	input := bufio.NewScanner(conn)

	for input.Scan() {
		go echo(conn, input.Text(), 1 * time.Second)
	}

	return nil
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}
