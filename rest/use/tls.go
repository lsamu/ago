package use

import (
    "github.com/gin-gonic/gin"
    "github.com/unrolled/secure"
)

func TlsHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        secureMiddleware := secure.New(secure.Options{
            SSLRedirect: true,
            SSLHost:     "xxx.com", // config.SslConfig.Domain,
        })
        err := secureMiddleware.Process(c.Writer, c.Request)
        if err != nil {
            return
        }
        c.Next()
    }
}
