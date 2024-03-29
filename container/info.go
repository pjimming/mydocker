package container

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pjimming/mydocker/utils/jsonx"

	"github.com/sirupsen/logrus"
)

type Info struct {
	Pid         string   `json:"pid"`         // 容器的init进程在宿主机上的 PID
	Id          string   `json:"id"`          // 容器Id
	Name        string   `json:"name"`        // 容器名
	Command     string   `json:"command"`     // 容器内init运行命令
	CreatedTime string   `json:"createTime"`  // 创建时间
	Status      string   `json:"status"`      // 容器的状态
	Volume      string   `json:"volume"`      // 挂载的数据卷
	PortMapping []string `json:"portMapping"` // 端口映射
}

// RecordInfo 记录容器相关信息
func RecordInfo(containerPid int, commandArray []string, containerName, containerId, volume string) error {
	if containerName == "" {
		containerName = containerId
	}
	command := strings.Join(commandArray, "")
	containerInfo := &Info{
		Pid:         strconv.Itoa(containerPid),
		Id:          containerId,
		Name:        containerName,
		Command:     command,
		CreatedTime: time.Now().Format(time.DateTime),
		Status:      RUNNING,
		Volume:      volume,
	}

	infoStr, err := jsonx.ToJsonString(containerInfo)
	if err != nil {
		err = fmt.Errorf("to json string fail, %v", err)
		logrus.Error(err)
		return err
	}

	dirPath := getContainerDir(containerId)
	if err = os.MkdirAll(dirPath, 0622); err != nil {
		err = fmt.Errorf("mkdir all fail, %v", err)
		logrus.Error(err)
		return err
	}

	fileName := path.Join(dirPath, ConfigName)
	file, err := os.Create(fileName)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		err = fmt.Errorf("create file %s fail, %v", fileName, err)
		logrus.Error(err)
		return err
	}

	if _, err = file.WriteString(infoStr); err != nil {
		err = fmt.Errorf("write file %s fail, %v", fileName, err)
		logrus.Error(err)
		return err
	}
	return nil
}

// DeleteInfo 删除容器config信息
func DeleteInfo(containerId string) error {
	if err := os.RemoveAll(getContainerDir(containerId)); err != nil {
		logrus.Errorf("[%s] remove container info fail, %v", containerId, err)
		return err
	}
	return nil
}

// ReadInfo 根据containerId读取信息
func ReadInfo(containerId string) (*Info, error) {
	return getInfoById(containerId)
}
