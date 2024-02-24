package main

import (
	"github.com/pjimming/mydocker/container"
)

func execContainer(containerId string, cmdArray []string) error {
	return container.Exec(containerId, cmdArray)
}
