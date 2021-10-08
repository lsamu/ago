package use

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/lsamu/ago/lib/logger"
)

// Error Error
func Error() gin.HandlerFunc {
    return func(c *gin.Context){
        defer func() {
            if err := recover(); err != nil {
                DebugStack := ""
                //for _, v := range strings.Split(string(debug.Stack()), "\n") {
                //	DebugStack += v + "<br/>"
                //}
                //保存起来排查日志
                DebugStack += fmt.Sprintf("%s", err) + "<br/>"
                DebugStack += c.Request.Method + "  " + c.Request.Host + c.Request.RequestURI + "<br/>"
                DebugStack += c.Request.UserAgent() + "<br/>"
                DebugStack += c.ClientIP()
                logger.Errorf("系统异常，请联系管理员! %s", DebugStack)
                c.Abort()
                return
            }
        }()
    }
}
