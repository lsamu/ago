package http_test

import (
    "io/ioutil"
    "net/http/httptest"
    "strings"
)

// HTTPTest 接口测试
type HTTPTest struct {

}

//NewHTTPTest NewHTTPTest
func NewHTTPTest() *HTTPTest {
    return &HTTPTest{
    }
}

// Get 根据特定请求uri，发起get请求返回响应
func (h *HTTPTest) Get(uri string) []byte {
    // 构造get请求
    _ = httptest.NewRequest("GET", uri, nil)
    // 初始化响应
    w := httptest.NewRecorder()

    // 调用相应的handler接口
    //h.router.ServeHTTP(w, req)
    // 提取响应
    result := w.Result()
    defer result.Body.Close()
    // 读取响应body
    body, _ := ioutil.ReadAll(result.Body)
    return body
}

// PostForm 根据特定请求uri和参数param，以表单形式传递参数，发起post请求返回响应
func (h *HTTPTest) PostForm(uri string, request string) []byte {
    // 构造post请求
    req := httptest.NewRequest("POST", uri, strings.NewReader(request))
    req.Header.Set("Content-Type", "application/json")

    // 初始化响应
    w := httptest.NewRecorder()

    // 调用相应handler接口
    //h.router.ServeHTTP(w, req)

    // 提取响应
    result := w.Result()
    defer result.Body.Close()

    // 读取响应body
    body, _ := ioutil.ReadAll(result.Body)
    return body
}
