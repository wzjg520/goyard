package main

import (
	//	"bufio"
	"fmt"
	"net/http"
	"os"

	"go/types"

	"golang.org/x/net/html"
)

type depth struct {
	nowDepth int
	flag     bool
}

var mdepth depth

func forEachNode(n *html.Node, start, end func(n *html.Node)) {
	var isSpace = false
	if start != nil {
		start(n, *isSpace)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, start, end)
	}

	if end != nil {
		end(n, *isSpace)
	}

}

func startElement(n *html.Node, isSpace *isSpace) {

	if n.Type == html.ElementNode {

		var attrs string
		var enter = ""
		mdepth.flag = false

		for _, v := range n.Attr {
			attrs += v.Key + "=\"" + v.Val + "\" "
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			//			fmt.Println(c.Type)
			//			fmt.Println(c.Data)
			if c.Type == html.ElementNode {
				enter = "\n"
				mdepth.flag = true
				isSpace = true
				break
			}
		}

		if n.Data == "meta" || n.Data == "link" {
			fmt.Printf("%*s<%s %s/>%s", mdepth.nowDepth*4, "", n.Data, attrs, enter)
		} else {
			fmt.Printf("%*s<%s %s>%s", mdepth.nowDepth*4, "", n.Data, attrs, enter)
		}

		mdepth.nowDepth++
	}
}

func endElement(n *html.Node, isSpace *isSpace) {
	if n.Type == html.ElementNode {
		mdepth.nowDepth--
		fmt.Println(isSpace)
		if n.Data == "meta" || n.Data == "link" {
			//fmt.Printf("%*s<%s/>\n", depth*4, "", n.Data)
		} else if mdepth.flag == true {
			//fmt.Println(strings.EqualFold(n.Data, "mata"))
			fmt.Printf("%*s<%s a/>\n", mdepth.nowDepth*4, "", n.Data)
		} else {
			fmt.Printf("<%s/>\n", n.Data)
		}

	}
}

func visit(stack []string, links []string, n *html.Node) ([]string, []string) {
	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
				}
			}
		}

		stack = append(stack, n.Data)
		//		fmt.Println(stack)

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links, stack = visit(stack, links, c)
	}

	return links, stack
}

func grab(url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("geting %s:%s...", url, resp.Status)
	}

	defer resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "get url:%v\n", err)
		os.Exit(1)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks:%v\n", err)
		os.Exit(1)
	}

	forEachNode(doc, startElement, endElement)
	return nil
	//	links, stacks := visit(nil, nil, doc)

	//	return links, stacks, nil

	//	for _, stack := range stacks {
	//		//		fmt.Println(stack)
	//	}

}

func main() {

	//	r := bufio.NewReader(os.Stdin)

	//	for {
	//		fmt.Print("commend>")
	//		b, _, _ := r.ReadLine()
	//		url := string(b)
	//		grab(url)
	//		//		links, _, _ := grab(url)
	//		//		for _, link := range links {
	//		//			fmt.Println(link)
	//		//		}

	//	}
	url := os.Args[1]
	//	fmt.Print(url)
	err := grab(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
