package container

import (
	"fmt"
	"github.com/pjimming/mydocker/utils/jsonx"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// getContainerDir 获取容器记录在宿主机上的dir
func getContainerDir(containerId string) string {
	return fmt.Sprintf(InfoLocFormat, containerId)
}

// 根据containerId获取容器的pid
func getPidById(id string) (string, error) {
	info, err := getInfoById(id)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}

// 根据pid获取进程的env
func getEnvsByPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	content, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("read %s error, %v", path, err)
		return nil, err
	}

	envs := strings.Split(string(content), "\u0000")
	return envs, nil
}

// 根据 containerId 获取 Info
func getInfoById(id string) (*Info, error) {
	dir := getContainerDir(id)
	configFilePath := filepath.Join(dir, ConfigName)

	info := new(Info)
	if err := jsonx.ReadJsonFile(configFilePath, info); err != nil {
		logrus.Errorf("read json file %s error, %v", configFilePath, err)
		return nil, err
	}
	return info, nil
}
