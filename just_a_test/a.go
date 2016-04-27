package main

import (
	"fmt"

)

func main() {
	a := "google"
	b, ok := interface{}(a).(string)
	fmt.Println(b,ok)

}
