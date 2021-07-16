package main

import (
	"fmt"
	"github.com/lsamu/ago/tool/upgrade"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	cmd := cli.NewApp()
	cmd.Commands = []*cli.Command{
		{
			Name:   "rest",
			Usage:  "ago rest",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "rpc",
			Usage:  "ago rpc",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "socket",
			Usage:  "ago socket",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "tcp",
			Usage:  "ago tcp",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "docker",
			Usage:  "ago docker",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "k8s",
			Usage:  "ago k8s",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "upgrade",
			Usage:  "ago upgrade",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "mysql",
			Usage:  "ago upgrade",
			Action: upgrade.Upgrade,
		},
		{
			Name:   "migration",
			Usage:  "ago upgrade",
			Action: upgrade.Upgrade,
		},
	}
	err := cmd.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
