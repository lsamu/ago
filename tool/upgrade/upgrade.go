package upgrade

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

//Generate Generate
func Generate(c *cli.Context) (err error) {
	//GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get -u github.com/xx/xx
	fmt.Println("upgrade")
	return
}
