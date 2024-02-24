package cgroups

import (
	"fmt"
	"github.com/pjimming/mydocker/cgroups/subsystems"
	"strings"
)

type CgroupManager struct {
	Path string
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{Path: path}
}

func (c *CgroupManager) Apply(pid int, res *subsystems.ResourceConfig) error {
	var errMsg []string
	for _, subsys := range subsystems.Ins {
		if err := subsys.Apply(c.Path, pid, res); err != nil {
			errMsg = append(errMsg, err.Error())
		}
	}

	if len(errMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsg, ", "))
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	var errMsg []string
	for _, subsys := range subsystems.Ins {
		if err := subsys.Set(c.Path, res); err != nil {
			errMsg = append(errMsg, err.Error())
		}
	}

	if len(errMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsg, ", "))
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	var errMsg []string
	for _, subsys := range subsystems.Ins {
		if err := subsys.Remove(c.Path); err != nil {
			errMsg = append(errMsg, err.Error())
		}
	}

	if len(errMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsg, ", "))
	}
	return nil
}
