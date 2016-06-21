package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print("commencing countdown.")
	var tick = time.Tick(1 * time.Second)

	for i := 10; i > 0; i-- {
		fmt.Println(i)
		j := <-tick
		fmt.Println(j)
	}

	time.Sleep(100 * time.Second)
}
