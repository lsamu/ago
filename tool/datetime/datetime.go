package datetime

import (
    "fmt"
    "github.com/urfave/cli/v2"
)

//InitCommand 初始化命令
func InitCommand() *cli.Command {
        return &cli.Command{
            Name:   "datetime",
            Usage:  "ago datetime xxx",
            Action: Generate,
        }
}

//Generate 生成
func Generate(c *cli.Context) (err error) {
    //现在时间戳
    //时间生成时间戳
    //时间戳生成时间
   fmt.Printf("%+v",c.Args())
    return
}

