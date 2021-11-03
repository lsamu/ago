package use

import (
    "context"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

// Timeout 超时
func Timeout(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 是否超时，都会执行，进行收尾
        // wrap the request context with a timeout
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer func() {
            // check if context timeout was reached
            if ctx.Err() == context.DeadlineExceeded {
                // 这里超时，才会执行
                // write response and abort the request
                c.Writer.WriteHeader(http.StatusGatewayTimeout)
                c.Abort()
            }
            // c.JSON(http.StatusOK, "hello")
            cancel()
        }()
        // replace request with context wrapped request
        c.Request = c.Request.WithContext(ctx)
        c.Next() // 实际调用具体的handler处理业务，实际还在这个方法中，所以业务执行结束会回到中间件中执行中间件中的defer函数
    }
}
