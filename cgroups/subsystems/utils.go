package subsystems

import (
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

const mountPointIndex = 4

func getCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot, err := findCgroupMountPoint(subsystem)
	if err != nil {
		logrus.Errorf("find cgroups mount point fail, %v", err)
		return "", err
	}

	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	// 指定自动创建时才判断是否存在
	_, err = os.Stat(absPath)
	// 只有不存在才创建
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(absPath, 0755)
		return absPath, err
	}
	// 其他错误或者没有错误都直接返回
	return absPath, err
}

// findCgroupMountPoint 通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
func findCgroupMountPoint(subsystem string) (string, error) {
	// /proc/self/mountinfo 为当前进程的 mountinfo 信息
	// 可以直接通过 cat /proc/self/mountinfo 命令查看
	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// txt 大概是这样的：104 85 0:20 / /sys/fs/cgroups/memory rw,nosuid,nodev,noexec,relatime - cgroups cgroups rw,memory
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		// 对最后一个元素按逗号进行分割，这里的最后一个元素就是 rw,memory
		// 其中的 memory 就表示这是一个 memory subsystem
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// 如果等于指定的 subsystem，那么就返回这个挂载点跟目录，就是第四个元素，
				// 这里就是`/sys/fs/cgroups/memory`,即我们要找的根目录
				return fields[mountPointIndex], nil
			}
		}
	}

	err = scanner.Err()
	return "", err
}
