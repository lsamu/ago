# guzzle

## 一、包功能说明

> 类似php客户端 [guzzle](https://github.com/guzzle/guzzle) ，http请求库简单封装。

## 二、功能说明

目前可以发送get、post、put、delete、patch，以及使用body发送json数据体的请求

下面给个Get请求示例：

````
// init use default http.Client
client := guzzle.New(nil)

// get
res, err := client.Get("https://dev.dev", url.Values{}, map[string]string{"header-name": "header-value"})
````
