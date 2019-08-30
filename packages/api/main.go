package main

import (
	"github.com/honerlaw/mentordoc/server/http"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	http.StartServer(waitGroup)

	waitGroup.Wait()
}
