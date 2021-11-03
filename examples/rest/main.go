package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/lsamu/ago/examples/rest/request"
    "github.com/lsamu/ago/rest"
    "github.com/lsamu/ago/rest/handler"
    "github.com/lsamu/ago/rest/use"
)

func main() {
    server := rest.NewServer(rest.Conf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()
    server.Use(use.Cors("*"))
    server.AddRoute(rest.Route{
        Method: "GET",
        Path:   "/",
        Handler: func(c *gin.Context) {
            c.String(200, "index ago!")
        },
    })
    server.AddRoute(rest.Route{
        Method: "GET",
        Path:   "/hello",
        Handler: func(c *gin.Context) {
            c.String(200, "hello ago!")
        },
    })

    server.AddRoute(rest.Route{
        Method: "GET",
        Path:   "/query",
        Handler: func(c *gin.Context) {
            var err error
            var req request.UserRequest
            err = handler.Parse(c, &req)
            if err != nil {
                handler.JSON(c, handler.CodeErr, err.Error())
                return
            }
            fmt.Printf("%+v", req)
            handler.JSON(c, handler.CodeOK, handler.MsgOK)
        },
    })
    server.Start()
}
