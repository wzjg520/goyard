package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//	resp, err := http.Get("http://zixun.jia.com")
	//	if err != nil {
	//		return
	//	}
	//	defer resp.Body.Close()
	//	str := resp.Body
	//	fmt.Print(str)

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://zixun.jia.com")
		fmt.Fprintf(w, "hello %q", html.EscapeString(r.URL.Path))
		result, err := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "%s", html.UnescapeString(string(result)))
		if err != nil {
			return
		}
		defer resp.Body.Close()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
