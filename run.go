package main

import (
	"github.com/pjimming/mydocker/cgroups"
	"github.com/pjimming/mydocker/cgroups/subsystems"
	"github.com/pjimming/mydocker/container"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始前面创建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, cmd []string, runResConf *subsystems.ResourceConfig, volume string) {
	parent, writePipe, err := container.NewParentProcess(tty, volume)
	if err != nil {
		return
	}
	if err = parent.Start(); err != nil {
		logrus.Errorf("run fail, %v", err)
	}

	// new cgroup manager
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer func() {
		if err = cgroupManager.Destroy(); err != nil {
			logrus.Errorf("cgroup manager destroy fail, %v", err)
		}
	}()

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
		container.DeleteWorkSpace("/root", "/root/merged", volume)
	}
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command = %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
