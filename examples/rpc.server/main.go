package main

import (
    "github.com/lsamu/ago/arpc"
    "google.golang.org/grpc"
)

func main() {
    server:=arpc.NewServer(arpc.ServerConf{
        Host: "0.0.0.0",
        Port: 8888,
    }, func(server *grpc.Server) {

    })
    defer server.Stop()
    server.Start()
}
