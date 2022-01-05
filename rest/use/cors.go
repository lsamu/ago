package use

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func Cors(site string) gin.HandlerFunc {
    if site == "" {
        site = "*"
    }
    return cors.New(cors.Config{
        AllowOrigins:     []string{site},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Token"},
        AllowCredentials: true,
        MaxAge:           time.Second * time.Duration(7200),
    })
    // return func(c *gin.Context) {
    //     method := c.Request.Method
    //     if site == "" {
    //         site = "*"
    //     }
    //     c.Header("Access-Control-Allow-Origin", site) // 最好配置成域名
    //     c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
    //     c.Header("Access-Control-Allow-Methods", "POST, GET,PUT,DELETE,OPTIONS")
    //     c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
    //     c.Header("Access-Control-Allow-Credentials", "true")
    //     // 放行所有OPTIONS方法
    //     if method == "OPTIONS" {
    //         c.AbortWithStatus(http.StatusNoContent)
    //     }
    //     // 处理请求
    //     c.Next()
    // }
}
