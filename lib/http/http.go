package http

import "github.com/go-resty/resty/v2"

//Get Get
func Get(url string, headers map[string]string) (bs []byte, err error) {
    client := resty.New()
    resp, err := client.R().
        SetHeaders(headers).
        Get(url)
    if err != nil {
        return bs, err
    }
    bs = resp.Body()
    return bs, err
}

//Post Post map[string]interface  string  Struct []byte
func Post(url string, post interface{}, headers map[string]string) (bs []byte, err error) {
    client := resty.New()
    resp, err := client.R().
        SetHeaders(headers).
        SetBody(post).
        Post(url)
    if err != nil {
        return bs, err
    }
    bs = resp.Body()
    return bs, err
}

//Put Put map[string]interface  string  Struct []byte
func Put(url string, post interface{}, headers map[string]string) (bs []byte, err error) {
    client := resty.New()
    resp, err := client.R().
        SetHeaders(headers).
        SetBody(post).
        Put(url)
    if err != nil {
        return bs, err
    }
    bs = resp.Body()
    return bs, err
}

//Delete Delete
func Delete(url string, post interface{}, headers map[string]string) (bs []byte, err error) {
    client := resty.New()
    resp, err := client.R().
        SetHeaders(headers).
        SetBody(post).
        Delete(url)
    if err != nil {
        return bs, err
    }
    bs = resp.Body()
    return bs, err
}
