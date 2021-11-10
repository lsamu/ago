package conf

import (
    "github.com/spf13/viper"
    "path"
    "path/filepath"
    "strings"
)

type Conf struct {
    FilePath string //  ./xxx/xxx.toml
}

// InitConfig InitConfig
func InitConfig(conf *Conf)(err error) {
    if conf.FilePath == "" {
        panic("filePath is empty")
    }
    dir, file := filepath.Split(conf.FilePath)
    // fileName := filepath.Base(conf.FilePath)
    fileExt := path.Ext(conf.FilePath)
    viper.SetConfigFile(file)
    viper.SetConfigType(strings.Trim(fileExt, "."))
    viper.AddConfigPath(dir)

    err = viper.ReadInConfig()     // 查找并读取配置文件
    if err != nil {                // 处理读取配置文件的错误
        return
    }
    return
}

// GetConfig GetConfig
func GetConfig(confFile interface{}) (err error) {
    err = viper.Unmarshal(confFile)
    return err
}

// GetConfigByKey GetConfigByKey
func GetConfigByKey(key string, confFile interface{}) (err error) {
    err = viper.UnmarshalKey(key, confFile)
    return err
}