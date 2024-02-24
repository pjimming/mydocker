package subsystems

const memorySubsystem = "memory"

type MemorySubsystem struct {
}

func (s *MemorySubsystem) CgroupFileName() string {
	return "memory.limit_in_bytes"
}

func (s *MemorySubsystem) Name() string {
	return memorySubsystem
}

func (s *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	return setCgroup(s.Name(), cgroupPath, s.CgroupFileName(), res.MemoryLimit)
}

func (s *MemorySubsystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	return applyCgroup(s.Name(), cgroupPath, pid)
}

func (s *MemorySubsystem) Remove(cgroupPath string) error {
	return removeCgroup(s.Name(), cgroupPath)
}
