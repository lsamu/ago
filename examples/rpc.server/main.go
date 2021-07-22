package main

import (
    "github.com/lsamu/ago/rpca"
    "google.golang.org/grpc"
)

func main() {
    server:= rpca.NewServer(rpca.ServerConf{
        Host: "0.0.0.0",
        Port: 8888,
    }, func(server *grpc.Server) {

    })
    defer server.Stop()
    server.Start()
}
