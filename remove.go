package main

import "github.com/pjimming/mydocker/container"

func removeContainer(containerId string) error {
	return container.Remove(containerId)
}
