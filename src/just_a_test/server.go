package main

import (
	"net"
	"flag"
	"log"
	"time"
	"strings"
	"strconv"
	"bytes"
	"fmt"
	"io"
)


var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "8080", "port")

func main() {

	flag.Parse()

	l, err := net.Listen("tcp", *host + ":" + *port)

	defer l.Close()

	if err != nil {
		log.Printf("Error listening: ", err)
	}

	log.Printf("Listening: %s:%s", *host, *port)
	serverGo(l)
}


func serverGo(l net.Listener) {
	for {
		conn, err := l.Accept()

		if err != nil {
			log.Printf("Error accepting:" , err)
			continue
		}

		log.Printf("Accept message from %s ->> %s", conn.RemoteAddr(), conn.LocalAddr())
		go requestHandle(conn)
	}
}

func requestHandle(conn net.Conn) {



	defer conn.Close()

	for {
		conn.SetDeadline(time.Now().Add(1000 * time.Second))
		msg := read(conn)
		if msg == "" {
			log.Println(string(msg), "nothing")
		} else if strings.TrimSpace(msg) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10)

			fmt.Fprint(conn, daytime + "recived " + msg)
		} else if msg == "close" {
			break
		} else {
			fmt.Fprint(conn, "recived " + msg)
		}
		log.Println("server:", msg)
	}

}

func read(conn net.Conn) string {
	readByte := make([]byte, 10)
	var buffer bytes.Buffer

	for {
		readLen, err := conn.Read(readByte)

		buffer.Write(readByte)

		if err == io.EOF {
			conn.Close()
			return "close"
		}

		if readLen < 10 {
			break
		}

		readByte = make([]byte, 10)
	}

	msg := buffer.String()


	return msg

}
