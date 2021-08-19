package ding

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/jjonline/go-lib-backend/guzzle"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

// messageURL 钉钉消息api
var (
    messageURL = "https://oapi.dingtalk.com/robot/send"
)

// dingResponse 钉钉响应结构
type dingResponse struct {
    ErrCode int    `json:"errcode"`
    ErrMsg  string `json:"errmsg"`
}

// Ding 钉钉机器人消息发送客户端
//  - ① 每个钉钉机器人每分钟最多发送20条。如果超过20条，会限流10分钟
//  - ② 支持 文本 (text)、链接 (link)、markdown(markdown)、ActionCard、FeedCard消息类型
type Ding struct {
    token  string
    secret string
    client *guzzle.Client
    enable bool
}

// Btn actionCard类型消息多个按钮定义结构
type Btn struct {
    Title     string `json:"title" comment:"按钮名称，必须"`
    ActionURL string `json:"actionURL" comment:"按钮点击后打开的网址，必须"`
}

// Feed feed流消息子结构<多图文>
type Feed struct {
    Title      string `json:"title" comment:"标题，必须"`
    MessageURL string `json:"messageURL" comment:"点击后打开的网址，必须"`
    PicURL     string `json:"picURL" comment:"图片地址，必须"`
}

// New 创建钉钉客户端
//  - token  钉钉access_token，钉钉机器人设置时 Webhook 的URL里的access_token值
//  - secret 钉钉secret，钉钉机器人设置时 启用加签获得以 SEC 开头的秘钥令牌
//  - client 自定义 *http.Client 可自主控制http请求客户端，给 nil 不则使用默认
func New(token, secret string, enable bool, client *http.Client) *Ding {
    return &Ding{
        token:  token,
        secret: secret,
        client: guzzle.New(client),
        enable: enable,
    }
}

func (d *Ding) sign(t int64) string {
    payload := fmt.Sprintf("%d\n%s", t, d.secret)
    h := hmac.New(sha256.New, []byte(d.secret))
    h.Write([]byte(payload))
    data := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(data)
}

func (d *Ding) send(message interface{}) error {
    params := url.Values{}
    params.Set("access_token", d.token)
    if d.secret != "" { // 如果设置密钥,则签名
        t := time.Now().Unix() * 1000
        params.Set("timestamp", strconv.FormatInt(t, 10))
        params.Set("sign", d.sign(t))
    }

    res, err := d.client.PostJSON(context.TODO(), guzzle.ToQueryURL(messageURL, params), message, nil)
    if err != nil {
        return err
    }

    // check response
    body := dingResponse{}
    err = json.Unmarshal(res.Body, &body)
    if err != nil {
        return err
    }
    if body.ErrCode != 0 {
        return fmt.Errorf("%s", body.ErrMsg)
    }

    return nil
}

// Text 发送文本信息
//  - msg       文本消息内容，不宜过长
//  - atMobiles 需要 at 的人的手机号
//  - isAtAll   是否要 at 全员
func (d *Ding) Text(msg string, atMobiles []string, isAtAll bool) error {
    if !d.enable {
        return nil
    }
    message := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": msg,
        },
        "at": map[string]interface{}{
            "atMobiles": atMobiles,
            "isAtAll":   isAtAll,
        },
    }
    return d.send(message)
}

// Markdown 发送markdown信息
//  - title     MD格式消息的标题，会以 ## 即h2形式显示在首行
//  - msg       MD的内容<注意使用MarkDown格式，支持链接、图片>
//  - atMobiles 需要 at 的人的手机号
//  - isAtAll   是否要 at 全员
//  支持的MD语法见：https://developers.dingtalk.com/document/app/custom-robot-access#section-e4x-4y8-9k0
// 		例子1：
//      var account = "acc"
//      var msg = "login use mobile"
//      var time = "2021-07-23 14:48:44"
//      msg = fmt.Sprintf("> Account: %s  \n> Msg: %s  \n> Time:  %s \n", account, msg, time)
//      c.Markdown("login info", msg, nil, false)
// --------------------------------------------------------------------------------------------
//		例子2：
//      msg := []string{
//          "## panic",
//          "> Env: " + conf.Config.Server.Env,
//          "> Code: " + strconv.Itoa(code),
//          "> Msg: " + message,
//          "> ReqID: " + utils.IFaceToString(reqID),
//          "> Url: " + ctx.Request.URL.String(),
//          "> Method: " + ctx.Request.Method,
//          "> Stack：" + stack,
//      }
//      msg := strings.Join(msg, "  \n")
//      c.Markdown("", msg, nil, false)
func (d *Ding) Markdown(title, msg string, atMobiles []string, isAtAll bool) error {
    if !d.enable {
        return nil
    }
    message := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]string{
            "title": "Notify",
            "text":  fmt.Sprintf("## %s  \n%s  \n", title, msg),
        },
        "at": map[string]interface{}{
            "atMobiles": atMobiles,
            "isAtAll":   isAtAll,
        },
    }
    return d.send(message)
}

// Link 发送卡片链接信息
//  - title     卡片消息的标题
//  - msg       卡片消息的正文
//  - msgURL    卡片消息被点击后打开的URL
//  - picURL    卡片消息的封面图<没有图片可传空字符串>
func (d *Ding) Link(title, msg, msgURL, picURL string) error {
    if !d.enable {
        return nil
    }
    message := map[string]interface{}{
        "msgtype": "link",
        "link": map[string]string{
            "title":      title,
            "text":       msg,
            "messageUrl": msgURL,
            "picUrl":     picURL,
        },
    }
    return d.send(message)
}

// ActionCard 发送只有一个按钮的整体跳转类型单个卡片消息
//  - title       卡片消息的标题
//  - msg         卡片消息的正文<支持MarkDown格式>
//  - btnText     卡片消息按钮上的文字
//  - btnURL      卡片消息按钮（其实是整个卡片）被点击后打开的URL
func (d *Ding) ActionCard(title, msg, btnText, btnURL string) error {
    if !d.enable {
        return nil
    }
    message := map[string]interface{}{
        "msgtype": "actionCard",
        "actionCard": map[string]string{
            "title":          title,
            "text":           msg,
            "singleTitle":    btnText,
            "singleURL":      btnURL,
            "btnOrientation": "0", // ActionCard 只有单条消息时，按钮垂直排列or平行排列无意义，给固定值
        },
    }
    return d.send(message)
}

// ActionCardWithMultiBtn 发送有多个按钮的单个卡片消息
//  - title         卡片消息的标题
//  - msg           卡片消息的正文<支持MarkDown格式>
//  - btn           Btn 切片 卡片消息多个按钮定义
//  - isBtnVertical 多个按钮连接是否垂直排列<即多个按钮是否从上至下依次排列> true垂直 false水平
func (d *Ding) ActionCardWithMultiBtn(title, msg string, btn []Btn, isBtnVertical bool) error {
    if !d.enable || btn == nil {
        return nil
    }

    // 按钮排列设置
    btnOrientation := "1"
    if isBtnVertical {
        btnOrientation = "0"
    }

    // set URL
    btnSets := make([]map[string]string, 0)
    for _, val := range btn {
        btnSets = append(btnSets, map[string]string{"title": val.Title, "actionURL": val.ActionURL})
    }

    message := map[string]interface{}{
        "msgtype": "actionCard",
        "actionCard": map[string]interface{}{
            "title":          title,
            "text":           msg,
            "btns":           btnSets,
            "btnOrientation": btnOrientation,
        },
    }
    return d.send(message)
}

// FeedCard 信息流消息<类似微信公众号多多图消息>
func (d *Ding) FeedCard(feed []Feed) error {
    if !d.enable || feed == nil {
        return nil
    }

    message := map[string]interface{}{
        "msgtype": "feedCard",
        "feedCard": map[string][]Feed{
            "links": feed,
        },
    }
    return d.send(message)
}