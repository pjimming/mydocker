package container

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pjimming/mydocker/utils/jsonx"
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
	filePath := fmt.Sprintf("/proc/%s/environ", pid)
	content, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("read %s error, %v", filePath, err)
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

func getImage(imageName string) string {
	return RootUrl + imageName + ".tar"
}
func getUnTar(imageName string) string {
	return RootUrl + imageName + "/"
}

func getRoot(containerId string) string {
	return path.Join(RootUrl, containerId)
}

func getLower(containerId string) string {
	return fmt.Sprintf(lowerDirFormat, containerId)
}

func getUpper(containerId string) string {
	return fmt.Sprintf(upperDirFormat, containerId)
}

func getWorker(containerId string) string {
	return fmt.Sprintf(workDirFormat, containerId)
}

func getMerged(containerId string) string {
	return fmt.Sprintf(mergedDirFormat, containerId)
}

func getOverlayFsDirs(containerId string) string {
	// lowerdir=lower1:lower2:lower3,upperdir=upper,workdir=work
	return fmt.Sprintf(overlayFSFormat,
		getLower(containerId),
		getUpper(containerId),
		getWorker(containerId),
	)
}
