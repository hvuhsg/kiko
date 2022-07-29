package main

import (
	"time"

	"github.com/hvuhsg/kiko/communication/rest"
	"github.com/hvuhsg/kiko/execution"
)

func main() {
	engine := execution.NewEngine(20, 0.1, time.Second*1)
	server := rest.NewServer("127.0.0.1", 8080, &engine)
	server.Run()
}
