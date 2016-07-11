package main

import (
	"net"
	"flag"
	"log"
	"fmt"
	"bufio"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "8080", "port")

type client chan<- string

var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
)

func broadcaster() {

	var clients = make(map[client]bool)

	for {
		select {
		case cli := <-entering:
			clients[cli] = true
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {

	ch := make(chan string)
	go clientWrite(conn, ch)
	who := conn.RemoteAddr().String()

	ch <- "you are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- "hello " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"

	conn.Close()
}

func clientWrite(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}


func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *host + ":" + *port)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	log.Println("listen ", *host, *port)

	go broadcaster()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}
		go handleConn(conn)
	}



}


