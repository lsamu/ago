package upload

import (
    "errors"
    "github.com/lsamu/ago/lib/secret"
    "io"
    "mime/multipart"
    "os"
    "path"
    "strings"
    "time"
)

//Local Local
type Local struct {
    Path string
}

//NewLocal NewLocal
func NewLocal() (oo Local, err error) {
    return oo, err
}

//Upload Upload
func (a *Local) Upload(file *multipart.FileHeader) (fullPath string, path1 string, err error) {
    // 读取文件后缀
    ext := path.Ext(file.Filename)
    // 读取文件名并加密
    name := strings.TrimSuffix(file.Filename, ext)
    name = secret.MD5(name)
    // 拼接新文件名
    filename := name + "_" + time.Now().Format("20060102150405") + ext
    // 尝试创建此路径
    mkdirErr := os.MkdirAll(a.Path, os.ModePerm)
    if mkdirErr != nil {
        return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
    }
    // 拼接路径和文件名
    p := a.Path + "/" + filename
    f, openError := file.Open() // 读取文件
    if openError != nil {
        return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
    }
    defer f.Close() // 创建文件 defer 关闭
    out, createErr := os.Create(p)
    if createErr != nil {
        return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
    }
    defer out.Close()             // 创建文件 defer 关闭
    _, copyErr := io.Copy(out, f) // 传输（拷贝）文件
    if copyErr != nil {
        return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
    }
    return p, filename, nil
}

//Delete Delete
func (a *Local) Delete(key string) (err error) {
    p := a.Path + "/" + key
    if strings.Contains(p, a.Path) {
        if err := os.Remove(p); err != nil {
            return errors.New("本地文件删除失败, err:" + err.Error())
        }
    }
    return nil
}
