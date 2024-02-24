package command

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/pjimming/mydocker/container"
)

// InitCommand 内部方法，没有暴露给外部使用
var InitCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	/*
		1. 获取传递过来的 command 参数
		2. 执行容器初始化操作
	*/
	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		logrus.Infof("[InitCommand] init command %s", cmd)
		return container.RunContainerInitProcess()
	},
}
