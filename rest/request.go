package rest

import (
	"errors"
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
func Parse(c *gin.Context, r interface{}) error {
	//判断是否为json
	err := c.ShouldBindJSON(r)
	if err != nil {
		_, errBody := ioutil.ReadAll(c.Request.Body)
		if errBody != nil {
			log.Printf("序列化参数有误，body解析Error：", errBody.Error())
			return errors.New(errBody.Error())
		}
		return errors.New(Translate(err))
	}
	val := reflect.ValueOf(r)
	vv0 := val.MethodByName("Check")
	if vv0.IsValid() {
		var params []reflect.Value
		aa := vv0.Call(params)
		if len(aa) < 1 {
			panic("param len error.")
		}
		if _, ok := aa[0].Interface().(error); ok {
			panic("return type error.")
		}
		cusErr := aa[0].Interface()
		if cusErr != nil {
			return cusErr.(error)
		}
	}
	return nil
}
