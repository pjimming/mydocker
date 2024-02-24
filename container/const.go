package container

const (
	RUNNING       = "running"
	STOP          = "stopped"
	Exit          = "exited"
	InfoLoc       = "/var/run/mydocker/"
	InfoLocFormat = InfoLoc + "%s/"
	ConfigName    = "config.json"
	IDLength      = 10
	LogFile       = "container.log"
)

// nsenter里的C代码里已经出现mydocker_pid和mydocker_cmd这两个Key,主要是为了控制是否执行C代码里面的setns.
const (
	EnvExecPid = "mydocker_pid"
	EnvExecCmd = "mydocker_cmd"
)

// 容器相关目录
const (
	RootUrl         = "/root/"
	lowerDirFormat  = "/root/%s/lower"
	upperDirFormat  = "/root/%s/upper"
	workDirFormat   = "/root/%s/work"
	mergedDirFormat = "/root/%s/merged"
	overlayFSFormat = "lowerdir=%s,upperdir=%s,workdir=%s"
)
