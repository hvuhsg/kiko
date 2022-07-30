/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"time"

	gserver "github.com/hvuhsg/kiko/communication/grpc"
	"github.com/hvuhsg/kiko/execution"
)

var (
	port = flag.Int("port", 50051, "The server port")
	host = flag.String("host", "127.0.0.1", "The server host")
)

func main() {
	flag.Parse()

	engine := execution.NewEngine(2, 0.05, time.Millisecond*100)
	go engine.Optimize()

	grpcServer := gserver.NewServer(&engine, *host, *port)
	grpcServer.Run()
}
