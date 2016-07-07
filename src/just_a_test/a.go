package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signRev := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	signal.Notify(signRev, sigs...)
	for sig := range signRev {
		fmt.Printf("Received a singal: %s\n", sig)
	}

}
