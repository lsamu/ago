package http

import "github.com/go-resty/resty/v2"

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

//post:  map[string]interface  string  Struct []byte
// raw
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

//post:  map[string]interface  string  Struct []byte
// raw
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
