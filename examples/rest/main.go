package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/lsamu/ago/rest"
    "github.com/lsamu/ago/rest/handler"
    "net/http"
)

func main() {
    server := rest.NewServer(rest.RestConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()

    server.Use(func(c *gin.Context) {
        method := c.Request.Method
        c.Header("Access-Control-Allow-Origin", "*") //最好配置成域名
        c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
        c.Header("Access-Control-Allow-Methods", "POST, GET,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
        //放行所有OPTIONS方法
        if method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
        }
        // 处理请求
        c.Next()
    })

    //server.Use(func(c *gin.Context) {
    //    c.String(200,"222")
    //})

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
            type request struct {
                Type string `json:"type" form:"type"`
            }
            var err error
            req := new(request)
            err = handler.Parse(c, req)
            if err != nil {
                handler.JSON(c, handler.CodeErr, err.Error())
                return
            }
            fmt.Println("%+v", req)
            handler.JSON(c, handler.CodeOK, handler.MsgOK)
        },
    })
    server.Start()
}
