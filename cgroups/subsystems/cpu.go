package subsystems

const cpuSubsystem = "cpu"

type CpuSubsystem struct {
}

func (s *CpuSubsystem) CgroupFileName() string {
	return "cpu.shares"
}

func (s *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuShare == "" {
		return nil
	}

	return setCgroup(s.Name(), cgroupPath, s.CgroupFileName(), res.CpuShare)
}

func (s *CpuSubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.CpuShare == "" {
		return nil
	}

	return applyCgroup(s.Name(), cgroupPath, pid)
}

func (s *CpuSubsystem) Remove(cgroupPath string) error {
	return removeCgroup(s.Name(), cgroupPath)
}

func (s *CpuSubsystem) Name() string {
	return cpuSubsystem
}
