package subsystems

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
)

const memorySubsystem = "memory"

type MemorySubsystem struct {
}

func (s *MemorySubsystem) Name() string {
	return memorySubsystem
}

func (s *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	logrus.Debugf("memory set, path: %s, limit: %s", cgroupPath, res.MemoryLimit)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("memory subsystem get cgroups path fail, %v", err)
		return err
	}

	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err = os.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
		logrus.Errorf("set cgroups memory fail, %v", err)
		return err
	}
	return nil
}

func (s *MemorySubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	logrus.Debugf("memory set pid %d, path: %s", pid, res.CpuShare)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("memory subsystem get cgroups path fail, %v", err)
		return err
	}

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("apply %d to memory tasks fail, %v", pid, err)
		return err
	}
	return nil
}

func (s *MemorySubsystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("memory subsystem get cgroups path fail, %v", err)
		return err
	}
	logrus.Debugf("remove all %s", subsystemCgroupPath)
	return os.RemoveAll(subsystemCgroupPath)
}
