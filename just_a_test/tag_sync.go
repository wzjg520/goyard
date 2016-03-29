package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func help() {
	fmt.Println(`
		Commands:
			getTagArticles(get)
			parseData(parse)
	`)
}

func GetTagArticles(args []string) int {
	fmt.Println("GET TAG ARTICLES")
	return 0
}

func ParseData(args []string) int {
	fmt.Println("PARSE DATA")
	return 0
}

func getCommandHandles() map[string]func(args []string) int {
	return map[string]func([]string) int{
		"getTagArticles": GetTagArticles,
		"get":            GetTagArticles,
		"parseData":      ParseData,
		"parse":          ParseData,
	}
}

func main() {
	fmt.Println("Tag Synchronization Solution")
	help()
	r := bufio.NewReader(os.Stdin)
	handlers := getCommandHandles()
	for {
		fmt.Print("Command>")
		b, _, _ := r.ReadLine()
		line := string(b)
		token := strings.Split(line, " ")
		if handler, ok := handlers[token[0]]; ok {
			ret := handler(token)
			if ret != 0 {
				break
			}
		} else {
			fmt.Println("Unknown Command: ", token[0])
		}
	}

}
