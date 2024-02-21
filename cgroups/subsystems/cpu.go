package subsystems

import (
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

const cpuSubsystem = "cpu"

type CpuSubsystem struct {
}

func (s *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuShare == "" {
		return nil
	}

	logrus.Debugf("cpushare set, path: %s, limit: %s", cgroupPath, res.CpuShare)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("cpushare subsystem get cgroup path fail, %v", err)
		return err
	}

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
		logrus.Errorf("set cpu shares fail, %v", err)
		return err
	}
	return nil
}

func (s *CpuSubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuShare == "" {
		return nil
	}

	logrus.Debugf("cpu shares set %d, path: %s", pid, res.CpuShare)

	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("cpu shares get cgroups path fail, %v", err)
		return err
	}

	if err = os.WriteFile(path.Join(subsystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		logrus.Errorf("apply %d to cpu tasks fail, %v", pid, err)
		return err
	}
	return nil
}

func (s *CpuSubsystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Errorf("cpu shares subsystem get cgroups path fail, %v", err)
		return err
	}
	logrus.Debugf("remove all %s", subsystemCgroupPath)
	return os.RemoveAll(subsystemCgroupPath)
}

func (s *CpuSubsystem) Name() string {
	return cpuSubsystem
}
