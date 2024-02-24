package subsystems

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const mountPointIndex = 4

func getCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot, err := findCgroupMountPoint(subsystem)
	if err != nil {
		logrus.Errorf("find %s cgroups mount point fail, %v", subsystem, err)
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

func removeCgroup(subsystem, cgroupPath string) error {
	subsystemCgroupPath, err := getCgroupPath(subsystem, cgroupPath, false)
	if err != nil {
		logrus.Errorf("%s subsystem get cgroups path fail, %v", subsystem, err)
		return err
	}
	logrus.Infof("remove all %s", subsystemCgroupPath)
	return os.RemoveAll(subsystemCgroupPath)
}

func setCgroup(subsystem, cgroupPath, cgroupFileName, limit string) error {
	subsystemCgroupPath, err := getCgroupPath(subsystem, cgroupPath, true)
	if err != nil {
		logrus.Errorf("%s subsystem get cgroup path fail, %v", subsystem, err)
		return err
	}

	logrus.Infof("%s set cgroup %s, limit: %s",
		subsystem,
		filepath.Join(subsystemCgroupPath, cgroupFileName),
		limit,
	)

	if err = os.WriteFile(path.Join(subsystemCgroupPath, cgroupFileName), []byte(limit), 0644); err != nil {
		logrus.Errorf("set %s fail, %v", subsystem, err)
		return err
	}
	return nil
}

func applyCgroup(subsystem, cgroupPath string, pid int) error {
	subsystemCgroupPath, err := getCgroupPath(subsystem, cgroupPath, false)
	if err != nil {
		logrus.Errorf("%s get cgroups path fail, %v", subsystem, err)
		return err
	}

	logrus.Infof("%s apply cgroup %s, pid: %d",
		subsystem,
		subsystemCgroupPath,
		pid,
	)

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("apply %d to cpu tasks fail, %v", pid, err)
		return err
	}
	return nil
}
