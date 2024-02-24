package main

import "github.com/pjimming/mydocker/container"

func stopContainer(containerId string) error {
	return container.Stop(containerId)
}
