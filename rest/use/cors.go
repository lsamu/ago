package use

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func Cors(site string) gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method
        if site == "" {
            site = "*"
        }
        c.Header("Access-Control-Allow-Origin", site) //最好配置成域名
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
    }
}
