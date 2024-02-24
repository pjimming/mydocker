package command

import (
	"fmt"
	"github.com/pjimming/mydocker/container"
	"github.com/urfave/cli"
)

var RemoveCommand = cli.Command{
	Name:  "rm",
	Usage: "remove a stopped container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container id")
		}
		containerId := ctx.Args().Get(0)
		return removeContainer(containerId)
	},
}

func removeContainer(containerId string) error {
	return container.Remove(containerId)
}
