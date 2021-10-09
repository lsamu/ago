package main

import sockio "github.com/lsamu/ago/sockio"

func main() {
    server:= sockio.NewServer(sockio.SockConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()
    server.Start()
}
