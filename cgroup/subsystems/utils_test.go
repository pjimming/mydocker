package subsystems

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindCgroupMountPoint(t *testing.T) {
	ast := assert.New(t)

	cpuMountPoint, err := findCgroupMountPoint("cpu")
	ast.Nil(err)
	t.Logf("cpu subsystem mount point %v", cpuMountPoint)

	cpusetMountPoint, err := findCgroupMountPoint("cpuset")
	ast.Nil(err)
	t.Logf("cpuset subsystem mount point %v", cpusetMountPoint)

	memoryMountPoint, err := findCgroupMountPoint("memory")
	ast.Nil(err)
	t.Logf("memory subsystem mount point %v", memoryMountPoint)
}
