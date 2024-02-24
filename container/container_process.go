package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"syscall"
)

// NewParentProcess 启动一个新进程
/*
这里是父进程，也就是当前进程执行的内容。
1.这里的/proc/self/exe调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境。
4.如果用户指定了-it参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool, volume, containerId string) (*exec.Cmd, *os.File, error) {
	// 创建匿名管道用于传递参数，将readPipe作为子进程的ExtraFiles，子进程从readPipe中读取参数
	// 父进程中则通过writePipe将参数写入管道
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		logrus.Errorf("new pipe fail, %v", err)
		return nil, nil, err
	}

	args := []string{"init"}
	cmd := exec.Command("/proc/self/exe", args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		containerDir := getContainerDir(containerId)
		if err = os.MkdirAll(containerDir, 0622); err != nil {
			logrus.Errorf("[NewParentProcess] mkdir %s all fail, %v", containerDir, err)
			return nil, nil, err
		}
		stdLogFilePath := path.Join(containerDir, LogFile)
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			logrus.Errorf("[NewParentProcess] create %s error, %v", stdLogFilePath, err)
			return nil, nil, err
		}
		cmd.Stdout = stdLogFile
	}

	cmd.ExtraFiles = []*os.File{readPipe}
	mntPath := "/root/merged/"
	rootPath := "/root/"
	NewWorkSpace(rootPath, mntPath, volume)
	cmd.Dir = mntPath

	return cmd, writePipe, nil
}
