package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
容器的文件系统相关操作
*/

// NewWorkSpace create an overlays filesystem as container root workspace
/*
1）创建lower层
2）创建upper、worker层
3）创建merged目录并挂载overlayFS
4）如果有指定volume则挂载volume
*/
func NewWorkSpace(volume, imageName, containerId string) error {
	if err := createLower(imageName, containerId); err != nil {
		logrus.Errorf("[NewWorkSpace][image:%s][containerId:%s] create lower error, %v", imageName, containerId, err)
		return err
	}
	if err := createUpperAndWorker(containerId); err != nil {
		logrus.Errorf("[NewWorkSpace][containerId:%s] create upper and worker error, %v", containerId, err)
		return err
	}
	if err := mountOverlayFs(containerId); err != nil {
		logrus.Errorf("[NewWorkSpace][containerId:%s] mount overlayFs error, %v", containerId, err)
		return err
	}

	if volume != "" {
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			logrus.Errorf("volume extract fail: %v", err)
			return err
		}
		if err = mountVolume(containerId, hostPath, containerPath); err != nil {
			logrus.Errorf("[NewWorkSpace][ContainerId:%s] mount volume error, %v", containerId, err)
			return err
		}
	}
	return nil
}

// createLower 把busybox作为overlays的lower层
func createLower(imageName, containerId string) error {
	imagePath := getImage(imageName)
	lower := getLower(containerId)

	if err := os.MkdirAll(lower, 0622); err != nil {
		logrus.Errorf("mkdir all %s error, %v", lower, err)
		return err
	}
	if _, err := exec.Command("tar", "-xvf", imagePath, "-C", lower).CombinedOutput(); err != nil {
		logrus.Errorf("[createLower][tar -xvf %s -C %s] error, %v", imagePath, lower, err)
		return err
	}
	return nil
}

// createUpperAndWorker 创建overlay fs需要的的upper、worker目录
func createUpperAndWorker(containerId string) error {
	upper := getUpper(containerId)
	if err := os.MkdirAll(upper, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", upper, err)
		return err
	}
	worker := getWorker(containerId)
	if err := os.Mkdir(worker, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", worker, err)
		return err
	}
	return nil
}

func mountOverlayFs(containerId string) error {
	// mount -t overlay overlay -o lowerdir=lower1:lower2:lower3,upperdir=upper,workdir=work merged
	// 创建对应的挂载路径
	mntPath := getMerged(containerId)
	if err := os.MkdirAll(mntPath, 0777); err != nil {
		logrus.Errorf("mkdir %s fail, %v", mntPath, err)
		return err
	}
	// 拼接参数
	// lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/merged
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", getOverlayFsDirs(containerId), mntPath)
	logrus.Infof(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("mount overlay fail, %v", err)
		return err
	}
	return nil
}

func volumeExtract(volume string) (sourcePath, destinationPath string, err error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invaild volume [%s]", volume)
	}

	sourcePath, destinationPath = parts[0], parts[1]
	if sourcePath == "" || destinationPath == "" {
		return "", "", fmt.Errorf("invaild volume [%s]", volume)
	}
	return
}

// 使用 bind mount 挂载 volume
func mountVolume(containerId, hostPath, containerPath string) error {
	// 宿主机目录
	if err := os.MkdirAll(hostPath, 0777); err != nil {
		logrus.Errorf("[mountVolume] mkdir %s fail, %v", hostPath, err)
		return err
	}
	// 容器目录
	mntPath := getMerged(containerId)
	containerVolumePath := filepath.Join(mntPath, containerPath)
	if err := os.MkdirAll(containerVolumePath, 0777); err != nil {
		logrus.Errorf("[mountVolume] mkdir %s fail, %v", containerVolumePath, err)
		return err
	}
	// 通过bind mount 将宿主机目录挂载到容器目录
	// mount -o bind /hostPath /containerVolumePath
	cmd := exec.Command("mount", "-o", "bind", hostPath, containerVolumePath)

	logrus.Infof("[mountVolume] cmd = %s", cmd.String())
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Errorf("mount volume failed. %v", err)
		return err
	}
	return nil
}

// DeleteWorkSpace 删除overlayFs当容器退出
/*
和创建相反
1）有volume则卸载volume
2）卸载merged目录
3）卸载upper、worker层
4）移除该容器的overlayFs目录
*/
func DeleteWorkSpace(volume, containerId string) error {
	logrus.Infof("[DeleteWorkSpace] volume:%s; containerId:%s", volume, containerId)
	// 1. umount volume
	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			logrus.Errorf("[DeleteWorkSpace] volume %s extract fail, %v", volume, err)
			return err
		}
		if err = umountVolume(containerId, containerPath); err != nil {
			logrus.Errorf("[DeleteWorkSpace] umount volume fail, %v", err)
			return err
		}
	}

	if err := umountOverlayFs(containerId); err != nil {
		logrus.Errorf("[DeleteWorkSpace] umount overlayFs error, %v", err)
		return err
	}
	if err := deleteDirs(containerId); err != nil {
		logrus.Errorf("[DeleteWorkSpace] deleteDirs error, %v", err)
		return err
	}
	return nil
}

func umountVolume(containerId, containerPath string) error {
	// mntPath 为容器在宿主机上的挂载点，例如 /root/merged
	// containerPath 为 volume 在容器中对应的目录，例如 /root/tmp
	// containerPathInHost 则是容器中目录在宿主机上的具体位置，例如 /root/{containerId}/merged/root/tmp
	containerPathInHost := path.Join(getMerged(containerId), containerPath)
	cmd := exec.Command("umount", containerPathInHost)

	logrus.Infof("[umountVolume] cmd = %s", cmd.String())
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Errorf("umount volume failed. %v", err)
		return err
	}
	logrus.Infof("umount volume %s success", containerPathInHost)
	return nil
}

func umountOverlayFs(containerId string) error {
	mntPath := getMerged(containerId)
	cmd := exec.Command("umount", mntPath)

	logrus.Info(cmd.String())
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Errorf("umount command run fail, %v", err)
		return err
	}
	logrus.Infof("umount overlayFs %s success", mntPath)
	return nil
}

func deleteDirs(containerId string) error {
	rootDir := getRoot(containerId)
	if err := os.RemoveAll(rootDir); err != nil {
		logrus.Errorf("[deleteDirs] rm -rf %s error, %v", rootDir, err)
		return err
	}
	logrus.Infof("rm -rf %s success", rootDir)
	return nil
}
