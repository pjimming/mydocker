package main

import (
	"fmt"
	"github.com/pjimming/mydocker/cgroups/subsystems"
	"github.com/pjimming/mydocker/container"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			// 限制进程内存使用量
			Name:  "mem",
			Usage: "memory limit, e.g.: -mem 100m",
		},
		cli.StringFlag{
			// 限制进程cpu使用率
			Name:  "cpushare",
			Usage: "cpu quota, e.g.: -cpushare 100",
		},
		cli.StringFlag{
			// 限制进程cpu使用率
			Name:  "cpuset",
			Usage: "cpuset limit, e.g.: -cpuset 2,4",
		},
		cli.StringFlag{
			// volume
			Name:  "v",
			Usage: "volume, e.g.: -v /etc/conf:/etc/conf",
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

		cmdArray := make([]string, 0)
		for _, cmd := range ctx.Args() {
			cmdArray = append(cmdArray, cmd)
		}
		tty := ctx.Bool("it")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
			CpuShare:    ctx.String("cpushare"),
			CpuSet:      ctx.String("cpuset"),
		}
		logrus.Infof("run cmd = %s", strings.Join(cmdArray, " "))
		volume := ctx.String("v")
		Run(tty, cmdArray, resConf, volume)
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
