package main

import (
	"fmt"
)

func main() {
	a := func() {
		fmt.Print("hello world")
	}

	fmt.Println(string(a))
}
