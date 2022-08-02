package main

import (
	"flag"
	"time"

	"github.com/hvuhsg/kiko/communication/rest"
	"github.com/hvuhsg/kiko/execution"
)

var (
	port = flag.Int("port", 8080, "The server port")
	host = flag.String("host", "127.0.0.1", "The server host")
)

func main() {
	flag.Parse()

	engine := execution.NewEngine(5, 0.05, time.Millisecond*20)
	go engine.Optimize()

	server := rest.NewServer(*host, *port, &engine)
	server.Run()
}
