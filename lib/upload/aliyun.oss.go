package upload

import (
    "errors"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "mime/multipart"
    "path/filepath"
    "time"
)

type AliYunOSS struct {
    bucket  *oss.Bucket
    baseUrl string
}

func NewAliYunOSS(Endpoint, AccessKeyId, AccessKeySecret, BucketName, baseUrl string) (oo AliYunOSS, err error) {
    client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
    if err != nil {
        return oo, err
    }
    bucket, err := client.Bucket(BucketName)
    if err != nil {
        return oo, err
    }
    oo = AliYunOSS{bucket: bucket, baseUrl: baseUrl}
    return
}

func (a *AliYunOSS) Upload(file *multipart.FileHeader) (fullPath string, path string, err error) {
    // 读取本地文件。
    f, openError := file.Open()
    if openError != nil {
        return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
    }
    //上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
    yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
    // 上传文件流。
    err = a.bucket.PutObject(yunFileTmpPath, f)
    if err != nil {
        return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
    }
    return a.baseUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (a *AliYunOSS) Delete(key string) (err error) {
    // 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
    // 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
    err = a.bucket.DeleteObject(key)
    if err != nil {
        return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
    }
    return
}
