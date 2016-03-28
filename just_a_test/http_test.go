package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://zixun.jia.com")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	str := resp.Body
	fmt.Print(str)

}
