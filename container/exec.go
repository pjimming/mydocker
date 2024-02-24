package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func Exec(containerId string, cmdArray []string) error {
	pid, err := getPidById(containerId)
	if err != nil {
		logrus.Errorf("[Exec] %s get pid fail, %v", containerId, err)
		return err
	}

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmdStr := strings.Join(cmdArray, " ")
	logrus.Infof("[Exec] container id: %s; container pid: %s; command: %s", containerId, pid, cmdStr)
	_ = os.Setenv(EnvExecPid, pid)
	_ = os.Setenv(EnvExecCmd, cmdStr)
	// 把指定PID进程的环境变量传递给新启动的进程，实现通过exec命令也能查询到容器的环境变量
	envs, err := getEnvsByPid(pid)
	if err != nil {
		logrus.Errorf("get envs by pid: [%s] error, %v", pid, err)
		return err
	}
	cmd.Env = append(os.Environ(), envs...)
	if err = cmd.Run(); err != nil {
		logrus.Errorf("[Exec] exec container %s error, %v", containerId, err)
		return err
	}
	return nil
}
