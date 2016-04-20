package main

import (
	"fmt"
	"log"
	"net/http"
)

const test_string = "hello world"

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.Header().Set("Content-Length", fmt.Sprint(len(test_string)))
	w.Write([]byte(test_string))
}

func main() {
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}
