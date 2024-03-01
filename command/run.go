package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/pjimming/mydocker/cgroups"
	"github.com/pjimming/mydocker/cgroups/subsystems"
	"github.com/pjimming/mydocker/container"
	"github.com/pjimming/mydocker/utils/randx"
)

var RunCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
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
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment",
		},
	},

	/*
		1. 判断参数是否包含command
		2. 获取用户指定的command
		3. 调用 run function 去准备启动容器
	*/
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}

		imageName := ctx.Args().First()
		var cmdArray []string
		cmdArray = append(cmdArray, ctx.Args().Tail()...)
		tty := ctx.Bool("it")
		detach := ctx.Bool("d")

		// tty 与 detach 不能共存
		if tty && detach {
			return fmt.Errorf("it and d paramter can not both provided")
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("mem"),
			CpuShare:    ctx.String("cpushare"),
			CpuSet:      ctx.String("cpuset"),
		}
		logrus.Infof("run cmd = %s", strings.Join(cmdArray, " "))
		volume := ctx.String("v")
		containerName := ctx.String("name")
		environSlice := ctx.StringSlice("e")
		run(tty, cmdArray, resConf, volume, containerName, imageName, environSlice)
		return nil
	},
}

// run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func run(tty bool, cmd []string, runResConf *subsystems.ResourceConfig, volume, containerName, imageName string, envSlice []string) {
	containerId := randx.RandString(container.IDLength)

	parent, writePipe, err := container.NewParentProcess(tty, volume, containerId, imageName, envSlice)
	if err != nil {
		return
	}
	if err = parent.Start(); err != nil {
		logrus.Errorf("run fail, %v", err)
	}

	// record container info
	if err = container.RecordInfo(parent.Process.Pid, cmd, containerName, containerId, volume); err != nil {
		logrus.Errorf("record container info fail, %v", err)
		return
	}

	// new cgroup manager
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	if err = cgroupManager.Set(runResConf); err != nil {
		logrus.Errorf("set cgroup res fail, %v", err)
	}
	if err = cgroupManager.Apply(parent.Process.Pid, runResConf); err != nil {
		logrus.Errorf("apply %d process cgroup res fail, %v", parent.Process.Pid, err)
	}

	// 在子进程创建后才能通过匹配来发送参数
	sendInitCommand(cmd, writePipe)
	if tty {
		_ = parent.Wait()
		if err = container.DeleteWorkSpace(volume, containerId); err != nil {
			logrus.Errorf("delete work space fail, %v", err)
		}
		_ = container.DeleteInfo(containerId)
		if err = cgroupManager.Destroy(); err != nil {
			logrus.Errorf("cgroup manager destroy fail, %v", err)
		}
	}
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command = %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
