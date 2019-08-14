package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/honerlaw/mentordoc/server"
	"sync"
)

func main() {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	server.StartServer(waitGroup)

	waitGroup.Wait()
}
