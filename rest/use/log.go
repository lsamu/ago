package use

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "net/http"
    "time"
)

type SysOperationRecord struct {
    Ip string
    Method string
    Path string
    Agent string
    Body string
    ErrorMessage string
}

// Log 日志
func Log() gin.HandlerFunc {
    return func(c *gin.Context){
        body, _ := ioutil.ReadAll(c.Request.Body)
        record := SysOperationRecord{
            Ip:     c.ClientIP(),
            Method: c.Request.Method,
            Path:   c.Request.URL.Path,
            Agent:  c.Request.UserAgent(),
            Body:   string(body),
        }
        now := time.Now()
        c.Next()
        latency := time.Now().Sub(now)
        status := c.Writer.Status()
        record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
        str := "接收到的请求为" + record.Body + "\n" + "请求方式为" + record.Method + "\n" + "报错信息如下" + record.ErrorMessage + "\n" + "耗时" + latency.String() + "\n"
        if status != 200 {
            subject :=record.Ip + "调用了" + record.Path + "报错了"
            fmt.Println(subject,str)//发送邮件
            //if err := utils.ErrorToEmail(subject, str); err != nil {
            //    utils.Error("ErrorToEmail Failed, err:", zap.Any("err", err))
            //}
            c.AbortWithStatus(http.StatusInternalServerError)
        }
    }
}
