package upload

import (
    "context"
    "errors"
    "fmt"
    "github.com/tencentyun/cos-go-sdk-v5"
    "mime/multipart"
    "net/http"
    "net/url"
    "time"
)

type TencentCOS struct {
    Bucket     string
    Region     string
    SecretID   string
    SecretKey  string
    PathPrefix string
    BaseURL    string
}

func NewTencentCOS() (oo TencentCOS, err error) {
    return oo, err
}

func (a *TencentCOS) Upload(file *multipart.FileHeader) (fullPath string, path string, err error) {
    client := a.NewClient()
    f, openError := file.Open()
    if openError != nil {
        return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
    }
    fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
    _, err = client.Object.Put(context.Background(), a.PathPrefix+"/"+fileKey, f, nil)
    if err != nil {
        panic(err)
    }
    return a.BaseURL + "/" + a.PathPrefix + "/" + fileKey, fileKey, nil
}

func (a *TencentCOS) Delete(key string) (err error) {
    client := a.NewClient()
    name := a.PathPrefix + "/" + key
    _, err = client.Object.Delete(context.Background(), name)
    if err != nil {
        return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
    }
    return nil
}

func (a *TencentCOS) NewClient() *cos.Client {
    urlStr, _ := url.Parse("https://" + a.Bucket + ".cos." + a.Region + ".myqcloud.com")
    baseURL := &cos.BaseURL{BucketURL: urlStr}
    client := cos.NewClient(baseURL, &http.Client{
        Transport: &cos.AuthorizationTransport{
            SecretID:  a.SecretID,
            SecretKey: a.SecretKey,
        },
    })
    return client
}
