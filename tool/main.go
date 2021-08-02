package main

import (
	"fmt"
	"github.com/lsamu/ago/tool/docker"
	"github.com/lsamu/ago/tool/k8s"
	"github.com/lsamu/ago/tool/rest"
	"github.com/lsamu/ago/tool/rpc"
	"github.com/lsamu/ago/tool/so"
	"github.com/lsamu/ago/tool/sock"
	"github.com/lsamu/ago/tool/tcp"
	"github.com/lsamu/ago/tool/upgrade"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	cmd := cli.NewApp()
	cmd.Commands = []*cli.Command{
		{
			Name:   "rest",
			Usage:  "ago rest create xxx, ago rest update xxx",
			Action: rest.Generate,
		},
		{
			Name:   "rpc",
			Usage:  "ago rpc",
			Action: rpc.Generate,
		},
		{
			Name:   "socket",
			Usage:  "ago socket",
			Action: sock.Generate,
		},
		{
			Name:   "tcp",
			Usage:  "ago tcp",
			Action: tcp.Generate,
		},
		{
			Name:   "so",
			Usage:  "ago so",
			Action: so.Generate,
		},
		{
			Name:   "docker",
			Usage:  "ago docker",
			Action: docker.Generate,
		},
		{
			Name:   "k8s",
			Usage:  "ago k8s",
			Action: k8s.Generate,
		},
		{
			Name:   "upgrade",
			Usage:  "ago upgrade",
			Action: upgrade.Generate,
		},
		{
			Name:   "mysql",
			Usage:  "ago upgrade",
			Action: k8s.Generate,
		},
	}
	err := cmd.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
