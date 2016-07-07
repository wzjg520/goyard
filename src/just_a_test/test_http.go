package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://zixun.jia.com")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

}
