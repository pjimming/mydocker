package container

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Remove(id string) error {
	info, err := getInfoById(id)
	if err != nil {
		logrus.Errorf("[Remove][id=%s] get info error, %v", id, err)
		return err
	}

	// 只能删除stopped的容器
	if info.Status != STOP {
		logrus.Errorf("[Remove][id=%s] can not remove, status is %s; pid is %s", id, info.Status, info.Pid)
		return nil
	}

	// 删除宿主机上关于容器的子目录所有文件
	dir := getContainerDir(id)
	if err = os.RemoveAll(dir); err != nil {
		logrus.Errorf("[Remove][id=%s] remove all %s error, %v", id, dir, err)
		return err
	}
	logrus.Infof("remove container [%s] success", id)
	return nil
}
