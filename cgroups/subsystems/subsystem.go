package subsystems

// ResourceConfig 用于传递资源限制配置的结构体，包含内存限制，CPU 时间片权重，CPU核心数
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// Subsystem 接口，每个Subsystem可以实现下面的4个接口，
// 这里将cgroup抽象成了path,原因是cgroup在hierarchy的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	// Name 返回子系统的名称，如cpu、memory
	Name() string
	// Set 设置某个cgroup在这个子系统的资源限制
	Set(cgroupPath string, res *ResourceConfig) error
	// Apply 将进程添加到某个cgroup中
	Apply(cgroupPath string, pid int, res *ResourceConfig) error
	// Remove 移除某个cgroup
	Remove(cgroupPath string) error
	// CgroupFileName 返回控制子系统资源的cgroup文件名
	CgroupFileName() string
}

var Ins = []Subsystem{
	&MemorySubsystem{},
	&CpuSubsystem{},
	&CpusetSubsystem{},
}
