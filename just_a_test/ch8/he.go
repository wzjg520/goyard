package main

import (
	"fmt"
)

func main() {
	var c = make(chan int)
	var a = make(chan string)

	go func() {
		//time.Sleep(100 * time.Second)
		close(c)
		a <- "chan a"

	}()

	for i := 0; i < 1000; i++{
		//time.Sleep(1 * time.Second)
		select {
		case temp := <-c:
			fmt.Println("closed c chan", temp)

		default:
			fmt.Println("defualt")
		}
	}

	//go func() {
	//	loop:
	//	for {
	//		select {
	//		case temp := <-c:
	//			a <- "chan a"
	//		fmt.Println(temp)
	//			fmt.Println(nil)
	//			break loop
	//		default:
	//			fmt.Println("defualt")
	//		}
	//	}
	//
	//
	//}()

	//select {
	//case <- c:
	//	fmt.Println(3)
	//default:
	//	fmt.Println(4)
	//}
	//
	//fmt.Println(<-a)

	//ch := make(chan int)
	//for i := 0; i < 20; i++ {
	//	select {
	//	case x := <-ch:
	//		fmt.Println(x) // "0" "2" "4" "6" "8"
	//	case ch <- i:
	//	}
	//}
}
