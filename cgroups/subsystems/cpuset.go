package subsystems

import (
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

const cpusetSubsystem = "cpuset"

type CpusetSubsystem struct {
}

func (s *CpusetSubsystem) Name() string {
	return cpusetSubsystem
}

func (s *CpusetSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	logrus.Infof("cpu set, path: %s, limit: %s", cgroupPath, res.CpuSet)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("cpuset subsystem get cgroup path fail, %v", err)
		return err
	}

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "cpuset.cpus"), []byte(res.CpuShare), 0644); err != nil {
		logrus.Errorf("set cpuset fail, %v", err)
		return err
	}
	return nil
}

func (s *CpusetSubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	logrus.Infof("cpuset set %d, path: %s", pid, res.CpuSet)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("cpuset get cgroups path fail, %v", err)
		return err
	}

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("apply %d to cpuset tasks fail, %v", pid, err)
		return err
	}
	return nil
}

func (s *CpusetSubsystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("cpuset subsystem get cgroups path fail, %v", err)
		return err
	}
	logrus.Infof("remove all %s", subsystemCgroupPath)
	return os.RemoveAll(subsystemCgroupPath)
}
