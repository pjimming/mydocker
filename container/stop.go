package container

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// Stop 停止容器
// 1. 找到pid
// 2. 发送SIGTERM信号
// 3. 修改config信息
func Stop(containerId string) error {
	pid, err := getPidById(containerId)
	if err != nil {
		logrus.Errorf("[Stop][id=%s] get pid error, %v", containerId, err)
		return err
	}

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		logrus.Errorf("[Stop][id=%s] atoi error, %v", containerId, err)
		return err
	}

	// 杀死容器进程
	if err = syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		logrus.Errorf("[Stop][id=%s] kill pid = %d error, %v", containerId, pidInt, err)
	}

	info, err := getInfoById(containerId)
	if err != nil {
		logrus.Errorf("[Stop][id=%s] get info error, %v", containerId, err)
		return err
	}

	// 修改容器信息
	info.Status = STOP
	info.Pid = " "
	infoByte, err := json.Marshal(info)
	if err != nil {
		logrus.Errorf("[Stop][id=%s] to json string error, %v", containerId, err)
		return err
	}

	// 覆盖之前的数据
	configFilePath := filepath.Join(getContainerDir(containerId), ConfigName)
	if err = os.WriteFile(configFilePath, infoByte, 0622); err != nil {
		logrus.Errorf("[Stop][id=%s] write %s error, %v", containerId, configFilePath, err)
		return err
	}
	logrus.Infof("[%s] stop container success", containerId)
	return nil
}
