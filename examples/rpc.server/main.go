package main

import (
    "github.com/lsamu/ago/rpc"
    "google.golang.org/grpc"
)

func main() {
    server:= rpc.NewServer(rpc.ServerConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    server.AddService(func(s *grpc.Server) {
        //注册服务
    })
    defer server.Stop()
    server.Start()
}
