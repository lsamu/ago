package handler

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Response struct {
    Code int         `json:"code"`
    Data interface{} `json:"data"`
    Msg  string      `json:"msg"`
}

const (
    ERROR   = -999
    SUCCESS = 0
)

const (
    CodeOK       = 0
    CodeErr      = 500 //内部错误
    CodeErr301   = 301 //未授权
    CodeErr401   = 401 //拒绝访问
    CodeErr404   = 404 //没有找到
    CodeParamErr = 501 //入参错误
    CodeSQLErr   = 502 //sql 错误
    MsgOK        = "操作成功"
    MsgErr       = "操作失败"
    MsgErr301    = "未授权！"
    MsgErr401    = "拒绝访问！"
    MsgErr404    = "没有找到！"
)

type CustomError struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

func JSON(c *gin.Context, code int, msg string, data ...interface{}) {
    //c.JSON(http.StatusOK, Response{
    //    code,
    //    data,
    //    msg,
    //})
    //c.Abort()
    result := new(Response)
    result.Code = code
    result.Msg = msg
    if len(data) > 0 && data[0] != nil {
        result.Data = data[0]
    } else {
        result.Data = ""
    }
    c.JSON(http.StatusOK, result)
    c.Abort()
}
