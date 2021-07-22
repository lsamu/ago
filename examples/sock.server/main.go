package main

import "github.com/lsamu/ago/sock"

func main() {
	server:=sock.NewServer(sock.SockConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
	defer server.Stop()
	server.Start()
}
