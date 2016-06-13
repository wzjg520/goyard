package main

import (
	"fmt"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()
	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}

}

// ------------------------
//import (
//	"io"
//	"log"
//	"net"
//	"time"
//)

//func main() {
//	listener, err := net.Listen("tcp", "localhost:8000")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			log.Print(err)
//			continue
//		}
//		go handleConn(conn)
//	}
//}

//func handleConn(c net.Conn) {
//	defer c.Close()
//	for {
//		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
//		if err != nil {
//			return
//		}
//		time.Sleep(1 * time.Second)
//	}
//}

// -------------------------------------------------------
//import (
//	"fmt"
//	"time"
//)

//func main() {
//	go spinner(100 * time.Millisecond)
//	const n = 45
//	fibN := fib(n)
//	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
//}

//func spinner(delay time.Duration) {
//	for {
//		for _, r := range `-\|/` {
//			fmt.Printf("\r%c", r)
//			time.Sleep(delay)
//		}
//	}
//}

//func fib(x int) int {
//	if x < 2 {
//		return x
//	}
//	return fib(x-1) + fib(x-2)
//}

// ----------
