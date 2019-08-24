package main

import (
	"github.com/honerlaw/mentordoc/server"
	"sync"
)

func main() {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	server.StartServer(waitGroup)

	waitGroup.Wait()
}
