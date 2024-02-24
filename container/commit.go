package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func Commit(imageName string) error {
	mntPath := "/root/merged"
	imageTar := fmt.Sprintf("/root/%s.tar", imageName)
	logrus.Infof("commit image: %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		return fmt.Errorf("tar folder %s fail, %v", mntPath, err)
	}
	return nil
}
