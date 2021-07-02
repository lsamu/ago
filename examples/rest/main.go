package main

import "github.com/lsamu/ago/rest"

func main() {
    server:=rest.NewServer(rest.RestConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()
    server.Start()
}
