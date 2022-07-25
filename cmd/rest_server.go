package main

import (
	"github.com/hvuhsg/kiko/communication/rest"
	"github.com/hvuhsg/kiko/execution"
)

func main() {
	engine := execution.NewEngine(20)
	server := rest.NewServer("127.0.0.1", 8080, &engine)
	server.Run()
}
