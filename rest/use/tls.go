package use

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/unrolled/secure"
)

// TlsHandler tls传输  domain:xxx.com     使用 router.RunTLS(":8080", "ssl.pem", "ssl.key")
func TlsHandler(domain string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if domain == "" {
            fmt.Println("abort tls")
            c.Abort()
        }
        secureMiddleware := secure.New(secure.Options{
            SSLRedirect: true,
            SSLHost:     domain, // config.SslConfig.Domain,
        })
        err := secureMiddleware.Process(c.Writer, c.Request)
        if err != nil {
            return
        }
        c.Next()
    }
}
