package handler

import (
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "log"
    "reflect"
)

const (
    //ApplicationJSON ApplicationJSON
    ApplicationJSON = "application/json"
)

//Parse Parse
func Parse(c *gin.Context, req interface{}) error {
    //log.Printf(c.ContentType())
    var err error
    if c.ContentType() == "" {
        err = c.ShouldBind(req)
        if err != nil {
            log.Printf("格式化有误！%s", err.Error())
            return errors.New(fmt.Sprintf("格式化有误！%s", err.Error()))
        }
    }else {
        //判断是否为json
        err = c.ShouldBindJSON(req)
        if err != nil {
            log.Printf("格式化有误！%+v", err)
            if err.Error() == "EOF" {
                return errors.New("JSON内容不能为空！")
            }
            _, errBody := ioutil.ReadAll(c.Request.Body)
            if errBody != nil {
                log.Printf("序列化参数有误，body解析Error：%s", errBody.Error())
                return errors.New(errBody.Error())
            }
            return errors.New(Translate(err))
        }
    }
    val := reflect.ValueOf(req)
    methodCheck := val.MethodByName("Check")
    if methodCheck.IsValid() {
        var params []reflect.Value
        callParam := methodCheck.Call(params)
        if len(callParam) < 1 {
            panic("param len error.")
        }
        if _, ok := callParam[0].Interface().(error); ok {
            panic("return type error.")
        }
        cusErr := callParam[0].Interface()
        if cusErr != nil {
            return cusErr.(error)
        }
    }
    return nil
}
