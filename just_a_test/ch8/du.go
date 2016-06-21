package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"sync"
)

var verbose = flag.Bool("v", true, "show verbose process messages")
var sema = make(chan struct{}, 20)

func main() {
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0 {
		roots = []string{"."}
	}

	for {
		du(roots)
		time.Sleep(10 * time.Second)
	}






}

func du(roots []string) {
	filesizes := make(chan int64)

	var n sync.WaitGroup

	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, filesizes)
	}

	go func() {
		n.Wait()
		close(filesizes)
	}()

	var nfiles, nbytes int64

	var tick <-chan time.Time

	if *verbose {
		tick = time.Tick(50000 * time.Millisecond)
	}

	loop:
	for {
		select {
		case size, ok := <-filesizes:
			if !ok {
				break loop
			}

			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes)
}

func walkDir(dir string, n *sync.WaitGroup, filesizes chan<- int64) {
	list := dirents(dir)
	defer n.Done()

	for _, info := range list {

		if info.IsDir() {
			subdir := filepath.Join(dir, info.Name())
			n.Add(1)
			go walkDir(subdir, n, filesizes)
		} else {
			filesizes <- info.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo {
	sema <- struct {}{}
	defer func() {
		<-sema
	}()
	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "du:%v\n", err)
		return nil
	}



	return entries
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}
