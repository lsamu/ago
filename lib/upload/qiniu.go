package upload

import (
    "context"
    "errors"
    "fmt"
    "github.com/qiniu/api.v7/v7/auth/qbox"
    "github.com/qiniu/api.v7/v7/storage"
    "mime/multipart"
    "time"
)

//QiNiu QiNiu
type QiNiu struct {
    Bucket        string
    AccessKey     string
    SecretKey     string
    ImgPath       string
    UseHTTPS      bool
    UseCdnDomains bool
    Zone          string
}

//NewQiNiu NewQiNiu
func NewQiNiu() (oo QiNiu, err error) {
    return oo, err
}


//Upload Upload
func (a *QiNiu) Upload(file *multipart.FileHeader) (fullPath string, path string, err error) {
    putPolicy := storage.PutPolicy{Scope: a.Bucket}
    mac := qbox.NewMac(a.AccessKey, a.SecretKey)
    upToken := putPolicy.UploadToken(mac)
    cfg := a.qiniuConfig()
    formUploader := storage.NewFormUploader(cfg)
    ret := storage.PutRet{}
    putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}
    f, openError := file.Open()
    if openError != nil {
        return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
    }
    fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
    putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
    if putErr != nil {
        return "", "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
    }
    return a.ImgPath + "/" + ret.Key, ret.Key, nil
}

//Delete Delete
func (a *QiNiu) Delete(key string) (err error) {
    mac := qbox.NewMac(a.AccessKey, a.SecretKey)
    cfg := a.qiniuConfig()
    bucketManager := storage.NewBucketManager(mac, cfg)
    if err := bucketManager.Delete(a.Bucket, key); err != nil {
        return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
    }
    return nil
}

func (a *QiNiu) qiniuConfig() *storage.Config {
    cfg := storage.Config{
        UseHTTPS:      a.UseHTTPS,
        UseCdnDomains: a.UseCdnDomains,
    }
    switch a.Zone { // 根据配置文件进行初始化空间对应的机房
    case "ZoneHuadong":
        cfg.Zone = &storage.ZoneHuadong
    case "ZoneHuabei":
        cfg.Zone = &storage.ZoneHuabei
    case "ZoneHuanan":
        cfg.Zone = &storage.ZoneHuanan
    case "ZoneBeimei":
        cfg.Zone = &storage.ZoneBeimei
    case "ZoneXinjiapo":
        cfg.Zone = &storage.ZoneXinjiapo
    }
    return &cfg
}
