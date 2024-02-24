package subsystems

const cpusetSubsystem = "cpuset"

type CpusetSubsystem struct {
}

func (s *CpusetSubsystem) CgroupFileName() string {
	return "cpuset.cpus"
}

func (s *CpusetSubsystem) Name() string {
	return cpusetSubsystem
}

func (s *CpusetSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	return setCgroup(s.Name(), cgroupPath, s.CgroupFileName(), res.CpuSet)
}

func (s *CpusetSubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	return applyCgroup(s.Name(), cgroupPath, pid)
}

func (s *CpusetSubsystem) Remove(cgroupPath string) error {
	return removeCgroup(s.Name(), cgroupPath)
}
