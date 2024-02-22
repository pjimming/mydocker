package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

/*
容器的文件系统相关操作
*/

// NewWorkSpace create an overlays filesystem as container root workspace
func NewWorkSpace(rootPath, mntPath string) {
	createLower(rootPath)
	createDirs(rootPath)
	mountOverlayFs(rootPath, mntPath)
}

// createLower 把busybox作为overlays的lower层
func createLower(rootPath string) {
	busyboxPath := filepath.Join(rootPath, "busybox/")
	busyboxTarPath := filepath.Join(rootPath, "busybox.tar")
	// 检查是否已经存在busybox文件夹
	exist, err := pathExists(busyboxPath)
	if err != nil {
		logrus.Errorf("check %s exist fail, %v", busyboxPath, err)
	}
	if !exist {
		if err = os.Mkdir(busyboxPath, 0777); err != nil {
			logrus.Errorf("mkdir %s fail, %v", busyboxPath, err)
		}
		if _, err = exec.Command("tar", "-xvf", busyboxTarPath, "-C", busyboxPath).CombinedOutput(); err != nil {
			logrus.Errorf("unTar dir %s fail, %v", busyboxTarPath, err)
		}
	}
}

// createDirs 创建overlay fs需要的的upper、worker目录
func createDirs(rootPath string) {
	upperPath := filepath.Join(rootPath, "upper/")
	if err := os.Mkdir(upperPath, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", upperPath, err)
	}
	workPath := filepath.Join(rootPath, "work/")
	if err := os.Mkdir(workPath, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", workPath, err)
	}
}

func mountOverlayFs(rootPath, mntPath string) {
	// mount -t overlay overlay -o lowerdir=lower1:lower2:lower3,upperdir=upper,workdir=work merged
	// 创建对应的挂载路径
	if err := os.Mkdir(mntPath, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", mntPath, err)
	}
	// 拼接参数
	// lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/merged
	dirs := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		filepath.Join(rootPath, "busybox/"),
		filepath.Join(rootPath, "upper/"),
		filepath.Join(rootPath, "work/"),
	)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntPath)
	logrus.Infof(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount overlay fail, %v", err)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWorkSpace(rootPath, mntPath string) {
	umountOverlayFs(mntPath)
	deleteDirs(rootPath)
}

func umountOverlayFs(mntPath string) {
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logrus.Infof(cmd.String())
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount command run fail, %v", err)
	}
	if err := os.RemoveAll(mntPath); err != nil {
		logrus.Errorf("remove dir %s fail, %v", mntPath, err)
	}
}

func deleteDirs(rootPath string) {
	upperPath := filepath.Join(rootPath, "upper")
	if err := os.RemoveAll(upperPath); err != nil {
		logrus.Errorf("remove dir %s fail, %v", upperPath, err)
	}
	logrus.Infof("rm dir %s", upperPath)
	workPath := filepath.Join(rootPath, "work")
	if err := os.RemoveAll(workPath); err != nil {
		logrus.Errorf("remove dir %s fail, %v", workPath, err)
	}
	logrus.Infof("rm dir %s", workPath)
}
