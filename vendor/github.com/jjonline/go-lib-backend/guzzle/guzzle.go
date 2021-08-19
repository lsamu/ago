package guzzle

import (
	"context"
	"errors"
	"io"
	"net/http"
)

// ErrResponseNotOK 当请求响应码非200时返回的错误
//   - 调用方只关注响应码为200的场景时，直接判断err是否为nil即可
//     result, err := client.JSON(xx,xx,xx)
//     if err != nil {
//     	 return
//     }
//     // your code // http响应码为200时的逻辑
//	 ------------------------------------------------------------------------
//   - 调用方若需处理非200时返回值，如下处理：
//     if err != nil && errors.Is(err, guzzle.ErrResponseNotOK) {
//       // http响应码非200，此时result也是有值的
//     }
var ErrResponseNotOK = errors.New("failed response status code is not equal 200")

// defaultUserAgent 默认UA头，调用方法时可覆盖
var defaultUserAgent = "guzzle/go (module github.com/jjonline/go-lib-backend/guzzle)"

// Result 响应封装
type Result struct {
	StatusCode    int         // 响应码
	ContentLength int64       // 响应长度
	Header        http.Header // 响应头
	Body          []byte      // 读取出来的响应body体字节内容
}

// Client http客户端相关方法封装
type Client struct {
	client *http.Client
}

// New 创建一个http客户端实例对象
//   - client *http.Client 可以自定义http请求的相关参数例如请求超时控制，使用默认则传 nil
func New(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		client: client,
	}
}

// NewRequest 新建http请求，链式初始化请求，需链式 Do 方法才实际执行<可灵活自定义以实现诸如 http.MethodOptions 类型请求>
//   - method 请求方法：GET、POST等，使用 http.MethodGet http.MethodPost 等常量
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体 io.Reader 类型
func (c *Client) NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置请求context
	req = req.WithContext(ctx)

	return req, nil
}

// Do 处理请求：用于链式调用
func (c *Client) Do(req *http.Request) (result Result, err error) {
	// set default user-agent if you do not set
	// header key is case-insensitive
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return result, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	row, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	// 非200时返回错误同时结果集仍然返回内容，以方便调用方需要处理状态码非200的场景
	if res.StatusCode != http.StatusOK {
		err = ErrResponseNotOK
	}

	// set result
	result.Body = row
	result.StatusCode = res.StatusCode
	result.Header = res.Header
	result.ContentLength = res.ContentLength

	return result, err
}

// Request 执行请求：实际执行请求<可灵活自定义以实现诸如 http.MethodOptions 类型请求>
//   - method 请求方法：GET、POST等，使用 http.MethodGet http.MethodPost 等常量
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体 io.Reader 类型
//   - head   请求header部分
func (c *Client) Request(ctx context.Context, method, url string, body io.Reader, head map[string]string) (Result, error) {
	req, err := c.NewRequest(ctx, method, url, body)
	if err != nil {
		return Result{}, err
	}
	for key, val := range head {
		req.Header.Add(key, val)
	}
	return c.Do(req)
}

// Get 执行 get 请求
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - query  GET请求URl中的Query键值对，支持类型：map[string]string、map[string][]string<等价于 url.Values>，无则给 nil
//   - head   请求header部分键值对，无则给 nil
func (c *Client) Get(ctx context.Context, url string, query interface{}, head map[string]string) (Result, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, ToQueryURL(url, query), nil)
	if err != nil {
		return Result{}, err
	}
	for key, val := range head {
		req.Header.Add(key, val)
	}
	return c.Do(req)
}

// Delete 执行 delete 请求
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - query  GET请求URl中的Query键值对，支持类型：map[string]string、map[string][]string<等价于 url.Values>，无则给 nil
//   - head   请求header部分键值对，无则给 nil
func (c *Client) Delete(ctx context.Context, url string, query interface{}, head map[string]string) (Result, error) {
	req, err := c.NewRequest(ctx, http.MethodDelete, ToQueryURL(url, query), nil)
	if err != nil {
		return Result{}, err
	}

	for key, val := range head {
		req.Header.Add(key, val)
	}

	return c.Do(req)
}

// JSON 执行 post/put/patch/delete 请求，采用 json 格式<比较底层的方法>
//   - method 请求方法：GET、POST等，使用 http.MethodGet http.MethodPost 等常量
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体 io.Reader 类型
//   - head   请求header部分键值对
func (c *Client) JSON(ctx context.Context, method, url string, body io.Reader, head map[string]string) (Result, error) {
	req, err := c.NewRequest(ctx, method, url, body)
	if err != nil {
		return Result{}, err
	}

	for key, val := range head {
		req.Header.Add(key, val)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.Do(req)
}

// Form 执行 post 请求，采用 form 表单格式<比较底层的方法>
//   - method 请求方法：GET、POST等，使用 http.MethodGet http.MethodPost 等常量
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体 io.Reader 类型
//   - head   请求header部分键值对
func (c *Client) Form(ctx context.Context, method, url string, body io.Reader, head map[string]string) (Result, error) {
	req, err := c.NewRequest(ctx, method, url, body)
	if err != nil {
		return Result{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, val := range head {
		req.Header.Add(key, val)
	}

	return c.Do(req)
}

// PostJSON 执行 post 请求，采用 json 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组、结构体等，最终会转换为 io.Reader 类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PostJSON(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.JSON(ctx, http.MethodPost, url, ToJsonReader(body), head)
}

// PutJSON 执行 put 请求，采用 json 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组、结构体等；请传 guzzle.ToJsonReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PutJSON(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.JSON(ctx, http.MethodPut, url, ToJsonReader(body), head)
}

// PatchJSON 执行 patch 请求，采用 json 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组、结构体等；请传 guzzle.ToJsonReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PatchJSON(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.JSON(ctx, http.MethodPatch, url, ToJsonReader(body), head)
}

// DeleteJSON 执行 delete 请求，采用 json 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组、结构体等；请传 guzzle.ToJsonReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) DeleteJSON(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.JSON(ctx, http.MethodDelete, url, ToJsonReader(body), head)
}

// PostForm 执行 post 请求，采用 form 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组等；请传 guzzle.ToFormReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PostForm(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.Form(ctx, http.MethodPost, url, ToFormReader(body), head)
}

// PutForm 执行 put 请求，采用 form 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组等；请传 guzzle.ToFormReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PutForm(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.Form(ctx, http.MethodPut, url, ToFormReader(body), head)
}

// PatchForm 执行 patch 请求，采用 form 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组等；请传 guzzle.ToFormReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) PatchForm(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.Form(ctx, http.MethodPatch, url, ToFormReader(body), head)
}

// DeleteForm 执行 delete 请求，采用 form 格式
//   - url    请求完整URL<可使用 guzzle.ToQueryURL 构造url里的query查询串>
//   - body   请求body体，支持：字符串、字节数组等；请传 guzzle.ToFormReader 支持的参数类型
//   - head   请求header部分键值对，无传nil
func (c *Client) DeleteForm(ctx context.Context, url string, body interface{}, head map[string]string) (Result, error) {
	return c.Form(ctx, http.MethodDelete, url, ToFormReader(body), head)
}
