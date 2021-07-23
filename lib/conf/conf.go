package conf

import (
    "encoding/json"
    "github.com/BurntSushi/toml"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
    "strings"
)

// Load 加载文件
func Load(path string, v interface{}) error {
    root := os.Getenv("SERVER_ROOT")
    if "" == root {
        root = "."
    }
    confpath := root + "/" + path
    content, err := ioutil.ReadFile(confpath)
    if nil != err {
        panic("Can not found conf file")
    }

    if strings.HasSuffix(confpath, ".json") {
        err = json.Unmarshal(content, v)
    }
    if strings.HasSuffix(confpath, ".yml") {
        err = yaml.Unmarshal(content, v)
    }
    if strings.HasSuffix(confpath, ".toml") {
        err = toml.Unmarshal(content, v)
    }
    return err
}