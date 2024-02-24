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

// 容器相关目录
const (
	RootUrl         = "/root/"
	lowerDirFormat  = "/root/%s/lower"
	upperDirFormat  = "/root/%s/upper"
	workDirFormat   = "/root/%s/work"
	mergedDirFormat = "/root/%s/merged"
	overlayFSFormat = "lowerdir=%s,upperdir=%s,workdir=%s"
)
