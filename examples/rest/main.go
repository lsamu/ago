package main

import (
    "github.com/gin-gonic/gin"
    "github.com/lsamu/ago/rest"
)

func main() {
    server:=rest.NewServer(rest.RestConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()

    server.AddRoute(rest.Route{
        Method:  "GET",
        Path:    "/",
        Handler: func(c *gin.Context) {

        },
    })
    server.Start()
}
