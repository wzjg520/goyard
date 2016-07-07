package main

import (
	"net"
	"log"
	"flag"
	"fmt"
	"bufio"
	"os"
	"strings"
	"time"
	"io"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "9000", "port")


func main() {
	flag.Parse()

	var conn, err = net.Dial("tcp", *host + ":" + *port)
	conn.SetDeadline(time.Now().Add(1000 * time.Second))

	if err != nil {
		log.Fatalf("content to %s:%s falsed", *host, *port)
	}

	log.Printf("content to %s:%s.", *host, *port)

	r := bufio.NewReader(os.Stdin)

	go func() {


		_, err := io.Copy(os.Stdout, conn)

		if err != nil {
			log.Fatalln(os.Stderr, "Fatal error:", err.Error())
			os.Exit(1)
		}

	}()

	for {
		fmt.Print("Commend>")
		b, _, _ := r.ReadLine()
		line := string(b)
		var token = strings.Split(line, " ")
		requestServer(conn, token[0])
	}



}

func requestServer(conn net.Conn, msg string) {

	if msg == "exit" {
		conn.Close()
		os.Exit(1)
	}

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}



}



