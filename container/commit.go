package container

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Commit(containerId, imageName string) error {
	mntPath := getMerged(containerId)
	imageTar := getImage(imageName)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder %s fail, %v", mntPath, err)
		return err
	}
	logrus.Infof("commit %s container success, image: %s", containerId, imageTar)
	return nil
}
