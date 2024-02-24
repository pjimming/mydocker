package main

import "github.com/pjimming/mydocker/container"

func logContainer(containerId string) error {
	return container.LogContainer(containerId)
}
