package guzzle

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ToQueryURL url和查询字符串拼接成完整URL
//  - url   请求URL，例如：https://www.baidu.com 或 https://www.baidu.com/?wd=11
//  - query 需要附加到URL里的查询键值对，没有附加的查询字符串时给 nil ，支持的类型如下：
//    - map[string]string
//    - map[string][]string <==> url.Values
func ToQueryURL(url string, query interface{}) string {
	if query != nil {
		// url里不存在 ? 符号 直接拼接返回
		if !strings.Contains(url, "?") {
			return url + "?" + BuildQuery(query)
		}

		// url里存在 ? 符号，去除右侧 & 符号后，再使用&符拼接返回
		return strings.TrimRight(url, "&") + "&" + BuildQuery(query)
	}
	return url
}

// ToJsonReader 将JSON类型请求的body体参数转换为统一的 io.Reader 类型
//  - body 拟转换为JSON类型请求body体的参数，支持的类型如下：
//    - nil
//    - string 即JSON字面量的字符串
//    - []byte 即JSON字面量的字节流
//    - io.Reader 无需转换，远洋返回
//    - struct等 使用 json.Marshal 转换
func ToJsonReader(body interface{}) io.Reader {
	switch pv := body.(type) {
	case nil:
		return nil
	case io.Reader:
		return pv
	case string:
		return strings.NewReader(pv)
	case []byte:
		return bytes.NewReader(pv)
	default:
		b, _ := json.Marshal(body)
		return bytes.NewReader(b)
	}
}

// ToFormReader 处理参数为Form表单类型
//   - 支持的参数类型如下：
//   - nil
//   - io.Reader
//   - string
//   - []byte
//   - map[string]string
//   - map[string][]string <==> url.Values
func ToFormReader(param interface{}) io.Reader {
	switch pv := param.(type) {
	case nil:
		return nil
	case io.Reader:
		return pv
	case string:
		return strings.NewReader(pv)
	case []byte:
		return bytes.NewReader(pv)
	case map[string]string, map[string][]string, url.Values:
		return strings.NewReader(BuildQuery(pv))
	default:
		return http.NoBody
	}
}

// BuildQuery 处理请求参数为URL里的Query键值对
//   - 支持的能构建的参数类型如下：
//   - map[string]string
//   - map[string][]string <==> url.Values
//   - 除了上述不支持的类型，其他类型将会忽略返回空字符串
func BuildQuery(query interface{}) string {
	switch kv := query.(type) {
	case map[string]string:
		values := make(url.Values)
		for k, v := range kv {
			values.Add(k, v)
		}
		return values.Encode()
	case map[string][]string:
		values := url.Values(kv)
		return values.Encode()
	case url.Values:
		return kv.Encode()
	default:
		return ""
	}
}
