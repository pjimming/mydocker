package command

import (
	"fmt"
	"github.com/pjimming/mydocker/container"
	"github.com/urfave/cli"
)

var StopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := ctx.Args().Get(0)
		return stopContainer(containerId)
	},
}

func stopContainer(containerId string) error {
	return container.Stop(containerId)
}
