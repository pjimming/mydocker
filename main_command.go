package main

import (
	"fmt"
	"github.com/pjimming/mydocker/container"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroup limit
			mydocker run -it [command]`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},

	/*
		1. 判断参数是否包含command
		2. 获取用户指定的command
		3. 调用 Run function 去准备启动容器
	*/
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := ctx.Args().Get(0)
		tty := ctx.Bool("it")
		logrus.Infof("run cmd = %s", cmd)
		Run(tty, []string{cmd})
		return nil
	},
}

// 内部方法，没有暴露给外部使用
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	/*
		1. 获取传递过来的 command 参数
		2. 执行容器初始化操作
	*/
	Action: func(ctx *cli.Context) error {
		logrus.Infof("[initCommand] init come on")
		cmd := ctx.Args().Get(0)
		logrus.Infof("[initCommand] init command %s", cmd)
		return container.RunContainerInitProcess()
	},
}
